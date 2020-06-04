package ingester

import (
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/timdrysdale/chmsg"
	"github.com/timdrysdale/gradex-cli/count"
	"github.com/timdrysdale/gradex-cli/extract"
	img "github.com/timdrysdale/gradex-cli/image"
	"github.com/timdrysdale/gradex-cli/optical"
	"github.com/timdrysdale/gradex-cli/pagedata"
	"github.com/timdrysdale/gradex-cli/parsesvg"
	"github.com/timdrysdale/gradex-cli/util"
	"github.com/timdrysdale/unipdf/v3/creator"
	"github.com/timdrysdale/unipdf/v3/model/optimize"
)

func TestOpticalBoxReading(t *testing.T) {

	util.EnsureDir("tmp-flatten")
	pdfPath := "./test-flatten/Practice-B999999-maTDD-marked-comments.pdf"
	jpgPath := "./tmp-flatten/Practice-B999999-maTDD-marked-comments.jpg"

	err := img.ConvertPDFToJPEGs(pdfPath, "./tmp-flatten", jpgPath)

	assert.NoError(t, err)

	widthPx, heightPx, err := optical.GetImageDimension(jpgPath)

	assert.NoError(t, err)

	fieldsMapByPage, err := extract.ExtractTextFieldsStructFromPDF(pdfPath)

	textfields := make(map[string]extract.TextField)

	//get the first one
	for _, v := range fieldsMapByPage {
		textfields = v
		break
	}

	assert.NoError(t, err)
	expandPx := -5
	boxes, err := parsesvg.GetImageBoxesForTextFields(textfields, heightPx, widthPx, true, expandPx)

	results, err := optical.CheckBoxFile(jpgPath, boxes)
	assert.NoError(t, err)

	resultMap := make(map[string]bool)

	for i, result := range results {
		resultMap[boxes[i].ID] = result
	}

	expectedMap := map[string]bool{
		"page-bad":    false,
		"page-ok":     true,
		"q1-mark":     true,
		"q1-number":   true,
		"q1-section":  true,
		"q2-mark":     false,
		"q2-number":   false,
		"q2-section":  false,
		"subtotal-00": true,
		"subtotal-01": false,
		"subtotal-02": false,
		"subtotal-03": false,
		"subtotal-04": true,
		"subtotal-05": false,
		"subtotal-06": false,
		"subtotal-07": false,
		"subtotal-08": false,
		"subtotal-09": true,
		"subtotal-10": false,
		"subtotal-11": false,
		"subtotal-12": false,
		"subtotal-13": false,
		"subtotal-14": false,
		"subtotal-15": false,
		"subtotal-16": false,
	}

	assert.Equal(t, expectedMap, resultMap)

	reader, err := os.Open(jpgPath)
	assert.NoError(t, err)

	defer reader.Close()

	testImage, _, err := image.Decode(reader)
	assert.NoError(t, err)

	for idx := 0; idx < len(boxes); idx = idx + 1 {

		checkImage := testImage.(optical.SubImager).SubImage(boxes[idx].Bounds)

		actualImagePath := filepath.Join("./tmp-flatten/", boxes[idx].ID+".jpg")
		expectedImagePath := filepath.Join("./expected/visual", boxes[idx].ID+".jpg")
		of, err := os.Create(actualImagePath)

		if err != nil {
			t.Errorf("problem saving checkbox image to file %v\n", err)
		}

		err = jpeg.Encode(of, checkImage, nil)

		if err != nil {
			t.Errorf("writing file %v\n", err)
		}

		of.Close()
		imgPath1 := expectedImagePath
		imgPath2 := actualImagePath

		out, err := exec.Command("compare", "-metric", "ae", imgPath1, imgPath2, "null:").CombinedOutput()
		assert.NoError(t, err)
		if err != nil {

			diffPath := filepath.Join(filepath.Dir(imgPath2),
				strings.TrimSuffix(filepath.Base(imgPath2), filepath.Ext(imgPath2))+
					"-diff"+filepath.Ext(imgPath2))
			cmd := exec.Command("compare", imgPath1, imgPath2, diffPath)
			cmd.Run()

			fmt.Printf("Images differ on page %d by %s (metric ae)\n see %s\n", idx, out, diffPath)
		}

	}

}

func TestFlattenProcessedMarked(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	// process a marked paper with keyed entries
	// check that keyed entries are picked up
	mch := make(chan chmsg.MessageInfo)

	logFile := "./flatten-process-testing.log"

	f, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	assert.NoError(t, err)

	defer f.Close()

	logger := zerolog.New(f).With().Timestamp().Logger()

	g, err := New("./tmp-delete-me", mch, &logger)

	assert.NoError(t, err)

	assert.Equal(t, "./tmp-delete-me", g.Root())

	os.RemoveAll("./tmp-delete-me")

	g.EnsureDirectoryStructure()

	templateFiles, err := g.GetFileList("./test-fs/etc/overlay/template")
	assert.NoError(t, err)

	for _, file := range templateFiles {
		destination := filepath.Join(g.OverlayTemplate(), filepath.Base(file))
		err := Copy(file, destination)
		assert.NoError(t, err)
	}

	exam := "Practice"
	stage := "marked"

	err = g.SetupExamDirs(exam)

	assert.NoError(t, err)

	source := "./test-flatten/Practice-B999999-maTDD-marked-comments.pdf"

	destinationDir := g.Ingest()

	assert.NoError(t, err)

	err = g.CopyToDir(source, destinationDir)

	assert.NoError(t, err)

	//destinationDir, err := g.FlattenProcessedPapersFromDir(exam, stage)
	err = g.StageFromIngest()
	assert.NoError(t, err)

	err = g.FlattenProcessedPapers(exam, stage)

	assert.NoError(t, err)

	// pagedata check
	flattenedPath := "tmp-delete-me/usr/exam/Practice/23-marker-flattened/Practice-B999999-maTDD-marked-comments.pdf"
	pdMap, err := pagedata.UnMarshalAllFromFile(flattenedPath)
	if err != nil {
		t.Error(err)
	}

	// check question values, and comments on page 1 in currentData
	// check previousData has one set, with null values for comments and values

	assert.Equal(t, "B999999", pdMap[1].Current.Item.Who)
	assert.Equal(t, "further-processing", pdMap[1].Current.Process.ToDo)
	assert.Equal(t, "ingester", pdMap[1].Current.Process.For)

	numPages, err := count.Pages(flattenedPath)

	assert.NoError(t, err)

	actualFields := make(map[int]map[string]string)

	for i := 1; i <= numPages; i++ {
		fieldsForPage := make(map[string]string)

		for _, item := range pdMap[i].Current.Data {
			fieldsForPage[item.Key] = item.Value
		}

		actualFields[i] = fieldsForPage
	}

	expectedFields := make(map[int]map[string]string)

	expectedFields[1] = map[string]string{
		"tf-page-ok":             "X",
		"tf-q1-mark":             "6/12",
		"tf-q1-number":           "1",
		"tf-q1-section":          "A",
		"tf-subtotal-00":         "1/2",
		"tf-subtotal-04":         "2/4",
		"tf-subtotal-09":         "3/6",
		"tf-page-ok-optical":     markDetected,
		"tf-q1-mark-optical":     markDetected,
		"tf-q1-number-optical":   markDetected,
		"tf-q1-section-optical":  markDetected,
		"tf-subtotal-00-optical": markDetected,
		"tf-subtotal-04-optical": markDetected,
		"tf-subtotal-09-optical": markDetected,
	}

	expectedFields[2] = map[string]string{
		"tf-page-bad":         "X",
		"tf-page-bad-optical": markDetected,
	}

	expectedFields[3] = map[string]string{
		"tf-page-ok":             "x",
		"tf-q1-mark":             "17",
		"tf-q1-number":           "1",
		"tf-q1-section":          "B",
		"tf-subtotal-01":         "2",
		"tf-subtotal-03":         "2",
		"tf-subtotal-06":         "1",
		"tf-subtotal-08":         "2",
		"tf-subtotal-10":         "2",
		"tf-subtotal-11":         "3",
		"tf-subtotal-14":         "5",
		"tf-page-ok-optical":     markDetected,
		"tf-q1-mark-optical":     markDetected,
		"tf-q1-number-optical":   markDetected,
		"tf-q1-section-optical":  markDetected,
		"tf-subtotal-01-optical": markDetected,
		"tf-subtotal-03-optical": markDetected,
		"tf-subtotal-06-optical": markDetected,
		"tf-subtotal-08-optical": markDetected,
		"tf-subtotal-10-optical": markDetected,
		"tf-subtotal-11-optical": markDetected,
		"tf-subtotal-14-optical": markDetected,
	}

	// checking by actual field value allows check for false positive optical marks
	// without specifying all the null fields in the expected values - it is
	// assumed automatically
	for page, fields := range actualFields {
		for k, v := range fields {
			expectedValue := "" //assume MUST BE empty field if not specified
			if _, ok := expectedFields[page][k]; ok {
				expectedValue = expectedFields[page][k]
			}
			assert.Equal(t, expectedValue, v)
			if expectedValue != v {
				fmt.Println(k)
			}

		}
	}

	// visual check (comments, in particular, as well as flattening of typed values)
	actualPdf := "./tmp-delete-me/usr/exam/Practice/23-marker-flattened/Practice-B999999-maTDD-marked-comments.pdf"
	expectedPdf := "./expected/visual/Practice-B999999-maTDD-marked-comments.pdf"

	_, err = os.Stat(actualPdf)
	assert.NoError(t, err)
	_, err = os.Stat(expectedPdf)
	assert.NoError(t, err)
	result, err := visuallyIdenticalMultiPagePDF(actualPdf, expectedPdf)
	assert.NoError(t, err)
	assert.True(t, result)
	if !result {
		fmt.Println(actualPdf)
	}

	os.RemoveAll("./tmp-delete-me")

}

func TestFlattenProcessedStylus(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	// process a marked paper with keyed entries
	// check that keyed entries are picked up
	mch := make(chan chmsg.MessageInfo)

	logFile := "./flatten-process-testing.log"

	f, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	assert.NoError(t, err)

	defer f.Close()

	logger := zerolog.New(f).With().Timestamp().Logger()

	g, err := New("./tmp-delete-me", mch, &logger)

	assert.NoError(t, err)

	assert.Equal(t, "./tmp-delete-me", g.Root())

	os.RemoveAll("./tmp-delete-me")

	g.EnsureDirectoryStructure()

	templateFiles, err := g.GetFileList("./test-fs/etc/overlay/template")
	assert.NoError(t, err)

	for _, file := range templateFiles {
		destination := filepath.Join(g.OverlayTemplate(), filepath.Base(file))
		err := Copy(file, destination)
		assert.NoError(t, err)
	}

	exam := "Practice"
	stage := "marked"

	err = g.SetupExamDirs(exam)

	assert.NoError(t, err)

	source := "./test-flatten/Practice-B999999-maTDD-marked-stylus.pdf"

	destinationDir := g.Ingest()

	assert.NoError(t, err)

	err = g.CopyToDir(source, destinationDir)

	assert.NoError(t, err)

	//destinationDir, err := g.FlattenProcessedPapersFromDir(exam, stage)
	err = g.StageFromIngest()
	assert.NoError(t, err)

	err = g.FlattenProcessedPapers(exam, stage)
	assert.NoError(t, err)

	// pagedata check
	flattenedPath := "tmp-delete-me/usr/exam/Practice/23-marker-flattened/Practice-B999999-maTDD-marked-stylus.pdf"
	pdMap, err := pagedata.UnMarshalAllFromFile(flattenedPath)
	if err != nil {
		t.Error(err)
	}

	//parsesvg.PrettyPrintStruct(pdMap)

	// check question values, and comments on page 1 in currentData
	// check previousData has one set, with null values for comments and values

	assert.Equal(t, "B999999", pdMap[1].Current.Item.Who)
	assert.Equal(t, "further-processing", pdMap[1].Current.Process.ToDo)
	assert.Equal(t, "ingester", pdMap[1].Current.Process.For)

	numPages, err := count.Pages(flattenedPath)

	assert.NoError(t, err)

	actualFields := make(map[int]map[string]string)

	for i := 1; i <= numPages; i++ {
		fieldsForPage := make(map[string]string)

		for _, item := range pdMap[i].Current.Data {
			fieldsForPage[item.Key] = item.Value
		}

		actualFields[i] = fieldsForPage
	}

	expectedFields := make(map[int]map[string]string)

	expectedFields[1] = map[string]string{
		"tf-page-ok-optical":     markDetected,
		"tf-q1-mark-optical":     markDetected,
		"tf-q1-number-optical":   markDetected,
		"tf-q1-section-optical":  markDetected,
		"tf-subtotal-00-optical": markDetected,
		"tf-subtotal-03-optical": markDetected,
		"tf-subtotal-05-optical": markDetected,
	}

	expectedFields[2] = map[string]string{
		"tf-page-bad-optical": markDetected,
	}

	expectedFields[3] = map[string]string{
		"tf-page-ok-optical":     markDetected,
		"tf-q1-mark-optical":     markDetected,
		"tf-q1-number-optical":   markDetected,
		"tf-q1-section-optical":  markDetected,
		"tf-subtotal-01-optical": markDetected,
		"tf-subtotal-03-optical": markDetected,
		"tf-subtotal-04-optical": markDetected,
		"tf-subtotal-06-optical": markDetected,
		"tf-subtotal-07-optical": markDetected,
		"tf-subtotal-10-optical": markDetected,
		"tf-subtotal-13-optical": markDetected,
	}

	// checking by actual field value allows check for false positive optical marks
	// without specifying all the null fields in the expected values - it is
	// assumed automatically
	for page, fields := range actualFields {
		for k, v := range fields {
			expectedValue := "" //assume MUST BE empty field if not specified
			if _, ok := expectedFields[page][k]; ok {
				expectedValue = expectedFields[page][k]
			}
			assert.Equal(t, expectedValue, v)
			if expectedValue != v {
				fmt.Println(k)
			}

		}
	}

	// visual check (comments, in particular, as well as flattening of typed values)
	actualPdf := "./tmp-delete-me/usr/exam/Practice/23-marker-flattened/Practice-B999999-maTDD-marked-stylus.pdf"
	expectedPdf := "./expected/visual/Practice-B999999-maTDD-marked-stylus.pdf"

	_, err = os.Stat(actualPdf)
	assert.NoError(t, err)
	_, err = os.Stat(expectedPdf)
	assert.NoError(t, err)
	result, err := visuallyIdenticalMultiPagePDF(actualPdf, expectedPdf)
	assert.NoError(t, err)
	assert.True(t, result)
	if !result {
		fmt.Println(actualPdf)
	}

	os.RemoveAll("./tmp-delete-me")

}

func TestFlattenProcessedMarkedAncestor(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	// process a marked paper with keyed entries
	// check that keyed entries are picked up
	mch := make(chan chmsg.MessageInfo)

	logFile := "./flatten-process-testing.log"

	f, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	assert.NoError(t, err)

	defer f.Close()

	logger := zerolog.New(f).With().Timestamp().Logger()

	g, err := New("./tmp-delete-me", mch, &logger)

	assert.NoError(t, err)

	assert.Equal(t, "./tmp-delete-me", g.Root())

	os.RemoveAll("./tmp-delete-me")

	g.EnsureDirectoryStructure()

	templateFiles, err := g.GetFileList("./test-fs/etc/overlay/template")
	assert.NoError(t, err)

	for _, file := range templateFiles {
		destination := filepath.Join(g.OverlayTemplate(), filepath.Base(file))
		err := Copy(file, destination)
		assert.NoError(t, err)
	}

	exam := "Practice"
	stage := "marked"

	err = g.SetupExamDirs(exam)

	assert.NoError(t, err)

	source := "./test-flatten/Practice-B999999-maTDD-marked-comments.pdf"

	destinationDir := g.GetExamDirNamed(exam, markerBack, "XX")

	assert.NoError(t, err)

	err = g.CopyToDir(source, destinationDir)

	assert.NoError(t, err)

	longname := "Practice Examination - A Huge Long Name"

	pds := []pagedata.PageData{
		pagedata.PageData{
			Current: pagedata.PageDetail{
				Item: pagedata.ItemDetail{
					Who:  "B999999",
					What: longname,
				},
				UUID:    "ancestorp1",
				Follows: "",
			},
		},
		pagedata.PageData{
			Current: pagedata.PageDetail{
				Item: pagedata.ItemDetail{
					Who:  "B999999",
					What: longname,
				},
				UUID:    "ancestorp2",
				Follows: "",
			},
		},
		pagedata.PageData{
			Current: pagedata.PageDetail{
				Item: pagedata.ItemDetail{
					Who:  "B999999",
					What: longname,
				},
				UUID:    "ancestorp3",
				Follows: "",
			},
		},
	}

	ancestorPath := filepath.Join(g.GetExamDir(exam, anonPapers), "Practice-B999999-ancestor.pdf")

	err = createPDF(ancestorPath, pds)

	assert.NoError(t, err)

	g.SetChangeAncestor(true)
	err = g.FlattenProcessedPapers(exam, stage)

	assert.NoError(t, err)

	// pagedata check
	flattenedPath := "tmp-delete-me/usr/exam/Practice/23-marker-flattened/Practice-B999999-maTDD-marked-comments.pdf"
	pdMap, err := pagedata.UnMarshalAllFromFile(flattenedPath)
	if err != nil {
		t.Error(err)
	}

	// check question values, and comments on page 1 in currentData
	// check previousData has one set, with null values for comments and values

	assert.Equal(t, "B999999", pdMap[1].Current.Item.Who)
	assert.Equal(t, "further-processing", pdMap[1].Current.Process.ToDo)
	assert.Equal(t, "ingester", pdMap[1].Current.Process.For)

	numPages, err := count.Pages(flattenedPath)

	assert.NoError(t, err)
	assert.Equal(t, 3, numPages)

	_, err = g.ReportOnProcessedDir(exam, g.GetExamDir(exam, markerFlattened), true, false)
	assert.NoError(t, err)

	pdA, err := pagedata.UnMarshalAllFromFile(ancestorPath)
	assert.NoError(t, err)

	pdN, err := pagedata.UnMarshalAllFromFile(flattenedPath)
	assert.NoError(t, err)

	assert.Equal(t, pdA[1].Current.Item.What, pdN[1].Current.Item.What)

	aA, err := GetPageSummaryMap(pdA)
	assert.NoError(t, err)

	aN, err := GetPageSummaryMap(pdN)
	assert.NoError(t, err)

	assert.Equal(t, true, aN[1].IsLinked)

	//these are the root UUID in the ancestor and the new file...
	assert.Equal(t, aA[1].FirstLink, aN[1].FirstLink)
	assert.Equal(t, aA[2].FirstLink, aN[2].FirstLink)
	assert.Equal(t, aA[3].FirstLink, aN[3].FirstLink)

	assert.Equal(t, "TDD", aN[1].WasFor)

	os.RemoveAll("./tmp-delete-me")
}

func createPDF(path string, pds []pagedata.PageData) error {

	c := creator.New()
	c.SetPageMargins(0, 0, 0, 0) // we're not printing so use the whole page
	c.SetPageSize(creator.PageSizeA4)

	for idx, pd := range pds {

		c.NewPage()
		p := c.NewParagraph(fmt.Sprintf("Ancestor Page %d", idx))
		p.SetFontSize(12)
		p.SetPos(200, 10)
		c.Draw(p)

		pagedata.MarshalOneToCreator(c, &pd)
	}

	c.SetOptimizer(optimize.New(optimize.Options{
		CombineDuplicateDirectObjects:   true,
		CombineIdenticalIndirectObjects: true,
		CombineDuplicateStreams:         true,
		CompressStreams:                 true,
		UseObjectStreams:                true,
		ImageQuality:                    90,
		ImageUpperPPI:                   150,
	}))

	of, err := os.Create(path)
	if err != nil {
		return err
	}

	defer of.Close()

	c.Write(of)
	return nil
}
