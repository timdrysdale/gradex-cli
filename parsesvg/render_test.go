package parsesvg

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/timdrysdale/gradex-cli/comment"
)

func TestRenderImagePrefillBackwardsCompatibility(t *testing.T) {

	var comments = make(map[int][]comment.Comment)

	comments[0] = []comment.Comment{c00}

	comments[1] = []comment.Comment{c10, c11}

	comments[2] = []comment.Comment{c20, c21}

	svgLayoutPath := "./test/layout-312pt.svg"

	pdfOutputPath := "./test/render-mark-spread-commments-backwards-compatibility.pdf"

	expectedPdf := "./expected/render-mark-spread-commments-backwards-compatibility.pdf"

	diffPdf := "./test/render-mark-spread-commments-backwards-compatibility-diff.pdf"

	previousImagePath := "./test/script.jpg"

	spreadName := "mark"

	pageNumber := int(1)

	contents := SpreadContents{
		SvgLayoutPath:         svgLayoutPath,
		SpreadName:            spreadName,
		PreviousImagePath:     previousImagePath,
		PageNumber:            pageNumber,
		PdfOutputPath:         pdfOutputPath,
		Comments:              comments,
		TemplatePathsRelative: true,
	}

	err := RenderSpreadExtra(contents)
	if err != nil {
		t.Error(err)
	}
	result, err := visuallyIdenticalPDF(pdfOutputPath, expectedPdf, diffPdf)
	assert.NoError(t, err)
	assert.True(t, result)

}

func TestRenderImagePrefillReplaceImage(t *testing.T) {

	var comments = make(map[int][]comment.Comment)

	comments[0] = []comment.Comment{c00}

	comments[1] = []comment.Comment{c10, c11}

	comments[2] = []comment.Comment{c20, c21}

	svgLayoutPath := "./test/layout-312pt.svg"

	pdfOutputPath := "./test/render-mark-spread-commments-prefill-header.pdf"
	expectedPdf := "./expected/render-mark-spread-commments-prefill-header.pdf"
	diffPdf := "./test/render-mark-spread-commments-prefill-header-diff.pdf"

	previousImagePath := "./img/a4-page.jpg"

	prefillImagePaths := make(map[string]string)

	// does not use extension - e
	prefillImagePaths["mark-header"] = "./test/prefill-header"

	spreadName := "mark"

	pageNumber := int(1)

	contents := SpreadContents{
		SvgLayoutPath:         svgLayoutPath,
		SpreadName:            spreadName,
		PreviousImagePath:     previousImagePath,
		PageNumber:            pageNumber,
		PdfOutputPath:         pdfOutputPath,
		Comments:              comments,
		PrefillImagePaths:     prefillImagePaths,
		TemplatePathsRelative: true,
	}

	err := RenderSpreadExtra(contents)
	if err != nil {
		t.Error(err)
	}
	result, err := visuallyIdenticalPDF(pdfOutputPath, expectedPdf, diffPdf)
	assert.NoError(t, err)
	assert.True(t, result)
}
