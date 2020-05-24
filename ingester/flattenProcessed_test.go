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
	"github.com/timdrysdale/gradex-cli/parsesvg"
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

	parsesvg.PrettyPrintStruct(pdMap)

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

		"page-bad":    "",
		"page-ok":     "X",
		"q1-mark":     "6/12",
		"q1-number":   "1",
		"q1-section":  "A",
		"subtotal-00": "1/2",
		"subtotal-04": "2/4",
		"subtotal-09": "3/6",
	}

	expectedFields[2] = map[string]string{
		"page-bad": "X",
	}

	expectedFields[3] = map[string]string{
		"page-ok":     "x",
		"q1-mark":     "17",
		"q1-number":   "1",
		"q1-section":  "B",
		"subtotal-01": "2",
		"subtotal-03": "2",
		"subtotal-06": "1",
		"subtotal-08": "2",
		"subtotal-10": "2",
		"subtotal-11": "3",
		"subtotal-14": "5",
	}

	for page, fields := range expectedFields {
		for k, v := range fields {
			assert.Equal(t, v, actualFields[page][k])
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
