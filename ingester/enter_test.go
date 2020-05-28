package ingester

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
