package ingester

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/timdrysdale/chmsg"
	"github.com/timdrysdale/gradex-cli/count"
	"github.com/timdrysdale/gradex-cli/pagedata"
)

func TestFlattenProcessedMarked(t *testing.T) {

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

	err = g.SetupExamPaths(exam)

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
	flattenedPath := "tmp-delete-me/usr/exam/Practice/23-marked-flattened/Practice-B999999-maTDD-marked-comments.pdf"
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

	/*
		for page, fields := range expectedFields {
			for k, v := range fields {
				assert.Equal(t, v, actualFields[page][k])
				if v != actualFields[page][k] {
					fmt.Println(k)
				}
			}
		}
	*/
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
	actualPdf := "./tmp-delete-me/usr/exam/Practice/23-marked-flattened/Practice-B999999-maTDD-marked-comments.pdf"
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

func testFlattenProcessedMarkedStylus(t *testing.T) {

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

	err = g.SetupExamPaths(exam)

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
	flattenedPath := "tmp-delete-me/usr/exam/Practice/23-marked-flattened/Practice-B999999-maTDD-marked-stylus.pdf"
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

	//parsesvg.PrettyPrintStruct(actualFields)

	expectedFields := make(map[int]map[string]string)

	expectedFields[1] = map[string]string{
		"page-ok": "X",
	}

	expectedFields[2] = map[string]string{
		"page-bad": "X",
	}

	expectedFields[3] = map[string]string{
		"page-ok": "x",
	}

	for page, fields := range expectedFields {
		for k, v := range fields {
			assert.Equal(t, v, actualFields[page][k])
		}
	}

	// visual check (comments, in particular, as well as flattening of typed values)
	actualPdf := "./tmp-delete-me/usr/exam/Practice/23-marked-flattened/Practice-B999999-maTDD-marked-stylus.pdf"
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
