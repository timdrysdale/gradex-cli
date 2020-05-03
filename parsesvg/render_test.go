package parsesvg

import (
	"testing"

	"github.com/timdrysdale/pdfcomment"
)

func TestRenderImagePrefillBackwardsCompatibility(t *testing.T) {

	var comments = make(map[int][]pdfcomment.Comment)

	comments[0] = []pdfcomment.Comment{c00}

	comments[1] = []pdfcomment.Comment{c10, c11}

	comments[2] = []pdfcomment.Comment{c20, c21}

	svgLayoutPath := "./test/layout-312pt-static-mark-dynamic-moderate-comment-static-check.svg"

	pdfOutputPath := "./test/render-mark-spread-commments-backwards-compatibility.pdf"

	previousImagePath := "./test/script.jpg"

	spreadName := "mark"

	pageNumber := int(1)

	contents := SpreadContents{
		SvgLayoutPath:     svgLayoutPath,
		SpreadName:        spreadName,
		PreviousImagePath: previousImagePath,
		PageNumber:        pageNumber,
		PdfOutputPath:     pdfOutputPath,
		Comments:          comments,
	}

	err := RenderSpreadExtra(contents)
	if err != nil {
		t.Error(err)
	}

}
func TestRenderImagePrefillNoPreviousImage(t *testing.T) {

	var comments = make(map[int][]pdfcomment.Comment)

	comments[0] = []pdfcomment.Comment{c00}

	comments[1] = []pdfcomment.Comment{c10, c11}

	comments[2] = []pdfcomment.Comment{c20, c21}

	svgLayoutPath := "./test/layout-312pt-static-mark-dynamic-moderate-comment-static-check.svg"

	pdfOutputPath := "./test/render-mark-spread-commments-no-previous-image.pdf"

	previousImagePath := ""

	spreadName := "mark"

	pageNumber := int(1)

	contents := SpreadContents{
		SvgLayoutPath:     svgLayoutPath,
		SpreadName:        spreadName,
		PreviousImagePath: previousImagePath,
		PageNumber:        pageNumber,
		PdfOutputPath:     pdfOutputPath,
		Comments:          comments,
	}

	err := RenderSpreadExtra(contents)
	if err != nil {
		t.Error(err)
	}

}

func TestRenderImagePrefillReplaceImage(t *testing.T) {

	var comments = make(map[int][]pdfcomment.Comment)

	comments[0] = []pdfcomment.Comment{c00}

	comments[1] = []pdfcomment.Comment{c10, c11}

	comments[2] = []pdfcomment.Comment{c20, c21}

	svgLayoutPath := "./test/layout-312pt-static-mark-dynamic-moderate-comment-static-check.svg"

	pdfOutputPath := "./test/render-mark-spread-commments-prefill-header.pdf"

	previousImagePath := ""

	prefillImagePaths := make(map[string]string)

	// does not use extension - e
	prefillImagePaths["mark-header"] = "./test/prefill-header"

	spreadName := "mark"

	pageNumber := int(1)

	contents := SpreadContents{
		SvgLayoutPath:     svgLayoutPath,
		SpreadName:        spreadName,
		PreviousImagePath: previousImagePath,
		PageNumber:        pageNumber,
		PdfOutputPath:     pdfOutputPath,
		Comments:          comments,
		PrefillImagePaths: prefillImagePaths,
	}

	err := RenderSpreadExtra(contents)
	if err != nil {
		t.Error(err)
	}

}
