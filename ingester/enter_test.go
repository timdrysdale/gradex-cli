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

func TestSelectForEnter(t *testing.T) {

	file1 := "./ingester/test-flatten/Practice-B999999-maTDD-marked-comments.pdf"
	file2 := "./ingester/test-flatten/Practice-B999999-maTDD-marked-stylus.pdf"

	files := []string{
		file1,
		file2,
	}

	pdfFiles := make(map[string]bool)

	for _, file := range files {

		if IsPDF(file) {
			pdfFiles[file] = false
		}

	}

	pdByFile := make(map[string]map[int]pagedata.PageData)

	file1pd := pagedata.PageData{
		Current: pagedata.PageDetail{
			Data: []pagedata.Field{
				pagedata.Field{
					Key:   "tf-foo-optical",
					Value: markDetected,
				},
				pagedata.Field{
					Key:   "tf-foo",
					Value: "x",
				},
			},
		},
	}
	file2pd := pagedata.PageData{
		Current: pagedata.PageDetail{
			Data: []pagedata.Field{
				pagedata.Field{
					Key:   "tf-foo-optical",
					Value: markDetected,
				},
				pagedata.Field{
					Key:   "tf-foo",
					Value: "",
				},
			},
		},
	}

	file1pdMap := make(map[int]pagedata.PageData)
	file1pdMap[1] = file1pd
	file1pdMap[2] = file1pd
	file2pdMap := make(map[int]pagedata.PageData)
	file2pdMap[1] = file2pd
	file2pdMap[2] = file2pd

	pdByFile[file1] = file1pdMap
	pdByFile[file2] = file2pdMap

	//parsesvg.PrettyPrintStruct(pdfFiles)

	selectByOpticalOnly(&pdfFiles, pdByFile)

	//parsesvg.PrettyPrintStruct(pdfFiles)

	assert.False(t, pdfFiles[file1])
	assert.True(t, pdfFiles[file2])

}

func TestMergeAfterEnter(t *testing.T) {

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
	stage := "entered"

	err = g.SetupExamDirs(exam)

	assert.NoError(t, err)

	source := "./test-enter/Practice-B999999-maTDD-marked-stylus-enXd.pdf"
	destinationDir := g.Ingest()
	err = g.CopyToDir(source, destinationDir)
	assert.NoError(t, err)

	//destinationDir, err := g.FlattenProcessedPapersFromDir(exam, stage)
	err = g.StageFromIngest()
	assert.NoError(t, err)

	err = g.FlattenProcessedPapers(exam, stage)
	assert.NoError(t, err)

	// pagedata check
	flattenedPath := "tmp-delete-me/usr/exam/Practice/43-enter-flattened/Practice-B999999-maTDD-marked-stylus-enXd.pdf"
	pdMap, err := pagedata.UnMarshalAllFromFile(flattenedPath)
	assert.NoError(t, err)

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
		"tf-page-ok":            "x",
		"tf-q1-mark":            "6/12",
		"tf-q1-number":          "1",
		"tf-q1-section":         "A",
		"tf-page-ok-optical":    markDetected,
		"tf-q1-mark-optical":    markDetected,
		"tf-q1-number-optical":  markDetected,
		"tf-q1-section-optical": markDetected,
	}

	expectedFields[2] = map[string]string{} //skipped page on purpose

	expectedFields[3] = map[string]string{
		"tf-page-ok":            "x",
		"tf-q1-mark":            "17",
		"tf-q1-number":          "1",
		"tf-q1-section":         "B",
		"tf-page-ok-optical":    markDetected,
		"tf-q1-mark-optical":    markDetected,
		"tf-q1-number-optical":  markDetected,
		"tf-q1-section-optical": markDetected,
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

	os.RemoveAll("./tmp-delete-me")

}
