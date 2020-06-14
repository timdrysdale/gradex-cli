package parsesvg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/mattetti/filebuffer"
	"github.com/stretchr/testify/assert"
	"github.com/timdrysdale/gradex-cli/comment"
	"github.com/timdrysdale/gradex-cli/geo"
	"github.com/timdrysdale/gradex-cli/image"

	"github.com/timdrysdale/gradex-cli/pagedata"
	"github.com/timdrysdale/unipdf/v3/annotator"
	"github.com/timdrysdale/unipdf/v3/creator"
	"github.com/timdrysdale/unipdf/v3/model"
	"github.com/timdrysdale/unipdf/v3/model/optimize"
)

func visuallyIdenticalPDF(pdf1, pdf2, diff string) (bool, error) {

	return image.VisuallyIdenticalMultiPagePDF(pdf1, pdf2)
}

func TestRenderSpreadTextFieldPrefill(t *testing.T) {

	svgLayoutPath := "./test/layout-a4-prefill.svg"

	pdfOutputPath := "./test/render-mark-spread-textfield-prefill.pdf"
	expectedPdf := "./expected/render-mark-spread-textfield-prefill.pdf"
	diffPdf := "./test/render-mark-spread-textfield-prefill-diff.pdf"

	previousImagePath := "./test/a4-three-square.jpg"

	spreadName := "mark"

	pageNumber := int(0)

	svgBytes, err := ioutil.ReadFile(svgLayoutPath)

	if err != nil {
		t.Error(err)
	}

	layout, err := DefineLayoutFromSVG(svgBytes)
	if err != nil {
		t.Errorf("Error defining layout %v", err)
	}

	if false {
		PrettyPrintStruct(layout)
	}

	textprefills := DocPrefills{}

	textprefills[0] = make(map[string]string)

	textprefills[0]["top-box"] = "THERE"
	textprefills[0]["top-box-00"] = "HELLO"

	textfieldvalues := DocPrefills{}
	textfieldvalues[0] = make(map[string]string)
	textfieldvalues[0]["bottom-box"] = "IPSUM"
	textfieldvalues[0]["bottom-box-00"] = "LOREM"

	contents := SpreadContents{
		SvgLayoutPath:     svgLayoutPath,
		SpreadName:        spreadName,
		PreviousImagePath: previousImagePath,
		PageNumber:        pageNumber,
		PdfOutputPath:     pdfOutputPath,
		Prefills:          textprefills,
		TextFieldValues:   textfieldvalues,
	}

	err = RenderSpreadExtra(contents)

	if err != nil {
		t.Error(err)
	}
	result, err := visuallyIdenticalPDF(pdfOutputPath, expectedPdf, diffPdf)
	assert.NoError(t, err)
	assert.True(t, result)
}

func TestRenderComboBox(t *testing.T) {

	svgLayoutPath := "./test/layout-a4-combo.svg"

	pdfOutputPath := "./test/render-combo.pdf"
	diffPdf := "./test/render-combo-diff.pdf"
	expectedPdf := "./expected/render-combo.pdf"

	previousImagePath := ""

	spreadName := "mark"

	pageNumber := int(1)

	combos := DocComboBoxes{}

	combos[1] = PageComboBoxes{}

	combos[1]["question"] = ComboOptions{
		Options: []string{"Guten Tag", "Bonjour", "Hola", "Gday"},
	}

	contents := SpreadContents{
		SvgLayoutPath:     svgLayoutPath,
		SpreadName:        spreadName,
		PreviousImagePath: previousImagePath,
		PageNumber:        pageNumber,
		PdfOutputPath:     pdfOutputPath,
		ComboBoxes:        combos,
	}

	err := RenderSpreadExtra(contents)
	assert.NoError(t, err)

	result, err := visuallyIdenticalPDF(pdfOutputPath, expectedPdf, diffPdf)
	assert.NoError(t, err)
	assert.True(t, result)

}

func TestRenderSpreadMarkPrefill(t *testing.T) {

	svgLayoutPath := "./test/layout-a4-prefill.svg"

	pdfOutputPath := "./test/render-mark-spread-prefill.pdf"
	expectedPdf := "./expected/render-mark-spread-prefill.pdf"
	diffPdf := "./test/render-mark-spread-prefill.pdf"

	previousImagePath := "./test/a4-three-square.jpg"

	spreadName := "mark"

	pageNumber := int(0)

	svgBytes, err := ioutil.ReadFile(svgLayoutPath)

	if err != nil {
		t.Error(err)
	}

	layout, err := DefineLayoutFromSVG(svgBytes)
	if err != nil {
		t.Errorf("Error defining layout %v", err)
	}

	if false {
		PrettyPrintStruct(layout)
	}

	textprefills := DocPrefills{}

	textprefills[0] = make(map[string]string)

	textprefills[0]["top-box"] = "THERE"
	textprefills[0]["top-box-00"] = "HELLO"

	contents := SpreadContents{
		SvgLayoutPath:     svgLayoutPath,
		SpreadName:        spreadName,
		PreviousImagePath: previousImagePath,
		PageNumber:        pageNumber,
		PdfOutputPath:     pdfOutputPath,
		Prefills:          textprefills,
	}

	err = RenderSpreadExtra(contents)

	if err != nil {
		t.Error(err)
	}
	result, err := visuallyIdenticalPDF(pdfOutputPath, expectedPdf, diffPdf)
	assert.NoError(t, err)
	assert.True(t, result)
}

const expectedLayoutJSON = `{"anchor":{"x":1.2588355559055121e-15,"y":-0.0003496930299212599},"dim":{"width":901.4173228346458,"height":884.4094488188978,"dynamicWidth":false},"id":"a4-portrait-layout","anchors":{"img-previous-mark":{"x":0,"y":42.51951212598426},"mark-header":{"x":6.294177637795276e-16,"y":1.062992040944882},"svg-check-flow":{"x":7.086614173228347,"y":1.062992040944882},"svg-mark-flow":{"x":655.9848283464568,"y":1.0628233228346458},"svg-mark-ladder":{"x":600.4855842519686,"y":1.0628233228346458},"svg-moderate-active":{"x":762.7586173228348,"y":1.0628233228346458},"svg-moderate-inactive":{"x":763.2376157480315,"y":1.062905102362205}},"pageDims":{"check":{"width":111.55415811023623,"height":883.3464566929134,"dynamicWidth":false},"mark":{"width":763.2376157480315,"height":883.3464566929134,"dynamicWidth":false},"moderate-active":{"width":899.7675590551182,"height":883.3464566929134,"dynamicWidth":false},"moderate-inactive":{"width":786.7112314960631,"height":883.3464566929134,"dynamicWidth":false},"width-moderate":{"width":1.417039398425197,"height":881.5748031496064,"dynamicWidth":true}},"filenames":{"mark-header":"./test/ladders-a4-portrait-header","svg-check-flow":"./test/sidebar-312pt-check-flow","svg-mark-flow":"./test/sidebar-312pt-mark-flow","svg-mark-ladder":"./test/sidebar-312pt-mark-ladder","svg-moderate-active":"./test/sidebar-312pt-moderate-flow-comment-active","svg-moderate-inactive":"./test/sidebar-312pt-moderate-inactive"},"ImageDims":{"mark-header":{"width":592.4409448818898,"height":39.68503937007874,"dynamicWidth":false},"previous-check":{"width":1.417039398425197,"height":881.5748031496064,"dynamicWidth":true},"previous-mark":{"width":595.2755905511812,"height":839.0551181102363,"dynamicWidth":false},"previous-moderate":{"width":763.2376157480315,"height":881.5748031496064,"dynamicWidth":false}}}`

var c00 = comment.Comment{Pos: geo.Point{X: 117.819, Y: 681.924}, Text: "This is a comment on page 1 - wrong page!"}
var c10 = comment.Comment{Pos: geo.Point{X: 326.501, Y: 593.954}, Text: "this is a comment on page 2", Page: 1}
var c11 = comment.Comment{Pos: geo.Point{X: 141.883, Y: 685.869}, Text: "this is a second comment on page 2", Page: 1}
var c20 = comment.Comment{Pos: geo.Point{X: 387.252, Y: 696.52}, Text: "this is a comment on page 3 - wrong page!", Page: 2}
var c21 = comment.Comment{Pos: geo.Point{X: 184.487, Y: 659.439}, Text: "this is a second comment on page 3 - wrong page!", Page: 2}

func TestDefineLayoutFromSvg(t *testing.T) {
	svgFilename := "./test/layout-312pt-static-mark-dynamic-moderate-static-check-v2.svg"
	svgBytes, err := ioutil.ReadFile(svgFilename)

	if err != nil {
		t.Error(err)
	}

	got, err := DefineLayoutFromSVG(svgBytes)
	if err != nil {
		t.Errorf("Error defining layout %v", err)
	}

	want := Layout{}

	_ = json.Unmarshal([]byte(expectedLayoutJSON), &want)

	if !reflect.DeepEqual(want.Anchor, got.Anchor) {
		t.Errorf("Anchor is different\n%v\n%v", want.Anchor, got.Anchor)
	}
	if !reflect.DeepEqual(want.Dim, got.Dim) {
		t.Errorf("Dim is different\n%v\n%v", want.Dim, got.Dim)
	}
	if !reflect.DeepEqual(want.ID, got.ID) {
		t.Errorf("ID is different\n%v\n%v", want.ID, got.ID)
	}
	if !reflect.DeepEqual(want.Anchors, got.Anchors) {
		t.Errorf("Anchors are different\n%v\n%v", want.Anchors, got.Anchors)
	}

	if !reflect.DeepEqual(want.PageDims, got.PageDims) {
		t.Errorf("PageDims are different\n%v\n%v", want.PageDims, got.PageDims)
	}
	if !reflect.DeepEqual(want.ImageDims, got.ImageDims) {
		t.Errorf("ImageDims are different\n%v\n%v", want.ImageDims, got.ImageDims)
	}
	if !reflect.DeepEqual(want.Filenames, got.Filenames) {
		t.Errorf("Filenames are different\n%v\n%v", want.Filenames, got.Filenames)
	}

}

func testPrintLayout(t *testing.T) {
	// helper for writing the tests on this file - not actually a test
	svgFilename := "./test/layout-312pt-static-mark-dynamic-moderate-static-check-v2.svg"
	svgBytes, err := ioutil.ReadFile(svgFilename)

	if err != nil {
		t.Error(err)
	}

	got, err := DefineLayoutFromSVG(svgBytes)
	if err != nil {
		t.Errorf("Error defining layout %v", err)
	}
	PrintLayout(got)

}

func testPrintSpreadsFromLayout(t *testing.T) {
	svgFilename := "./test/layout-312pt-static-mark-dynamic-moderate-static-check-v2.svg"
	svgBytes, err := ioutil.ReadFile(svgFilename)

	if err != nil {
		t.Error(err)
	}

	layout, err := DefineLayoutFromSVG(svgBytes)
	if err != nil {
		t.Errorf("Error defining layout %v", err)
	}

	spread := Spread{}

	spread.Name = "mark"

	foundPage := false
	for k, v := range layout.PageDims {
		if strings.Contains(k, spread.Name) {
			spread.Dim = v
			foundPage = true
		}
	}

	if !foundPage {
		fmt.Printf("no page info for this spread %s\n", spread.Name)
		return
	}

	/* TODO - CUT THIS STALE CODE?

	 ImageDims := geo.Dim{}

		for k, v := range layout.ImageDims {
			if strings.Contains(k, spread.Name) {
				ImageDims = v
			}
		}*/

	// find svg & img elements for this name
	var svgFilenames, imgFilenames []string

	for k, _ := range layout.Filenames {
		if strings.Contains(k, spread.Name) {

			// assume jpg- or no prefix is image; svg- is ladder (image plus acroforms)
			if strings.HasPrefix(k, geo.SVGElement) {
				svgFilenames = append(svgFilenames, k) //we'll get the contents later
			} else {
				imgFilenames = append(imgFilenames, k)
			}
		}
	}

	// get all the textfields (and put image of associated chrome into images list)
	// note that if page dynamic, textfields are ALL dynamically shifting wrt to dynamic page edge,
	// no matter what side of the previous image edge they are. This means we only need one set of dims
	// the layout engine will just add the amount of the previous image's size in the dynamic dimension
	// We need to add the anchor position to the textfield positions (which are relative to that anchor)

	//	TranslatePosition()

	for _, svgname := range svgFilenames {

		offset := geo.Point{}

		if thisAnchor, ok := layout.Anchors[svgname]; !ok {
			//default to layout anchor if not in the list
			offset = geo.Point{X: 0, Y: 0}
			//fmt.Printf("didn't find anchor for %s\n", svgname)
		} else {

			//fmt.Printf("%s@%v ref@%v\n", svgname, thisAnchor, layout.Anchor)
			offset = DiffPosition(layout.Anchor, thisAnchor)
			//fmt.Printf("Offset %s %v\n", svgname, offset)
		}

		svgfilename := fmt.Sprintf("%s.svg", layout.Filenames[svgname])
		imgfilename := fmt.Sprintf("%s.jpg", layout.Filenames[svgname]) //fixed by pdf library (I think)

		svgBytes, err := ioutil.ReadFile(svgfilename)
		if err != nil {
			t.Errorf("Error opening svg file %s", svgfilename)
		}

		ladder, err := DefineLadderFromSVG(svgBytes)
		if err != nil {
			t.Errorf("Error defining ladder %v", err)
		}

		if ladder == nil {
			continue //throw error?
		}
		spread.Ladders = append(spread.Ladders, *ladder)

		// append chrome image to the images list
		image := ImageInsert{
			Filename: imgfilename,
			Corner:   TranslatePosition(ladder.Anchor, offset),
			Dim:      ladder.Dim,
		}

		spread.Images = append(spread.Images, image) //add chrome to list of images to include

		//append TextFields to the Textfield list

		for _, tf := range ladder.TextFields {

			//shift the text field and add it to the list
			//let engine take care of mangling name to suit page

			tf.Rect.Corner = TranslatePosition(tf.Rect.Corner, offset)
			spread.TextFields = append(spread.TextFields, tf)
		}

	}

	//get all images, other than image-previous, that comes separately via own arg
	//Since these are the images for the textfield chrome, it's the same story - page layout engine will sort.
	//note that we haven't got previous image, so just send filename as 'previous-image' and let engine work it out
	//note that _all_ non-svg images need an image dims box....else their size will depend on their quality (dpi)

	for _, imgname := range imgFilenames {

		if _, ok := layout.ImageDims[imgname]; !ok {
			t.Errorf("No size for image %s (must be provided in layout)\n", imgname)
		}

		offset := geo.Point{}

		if thisAnchor, ok := layout.Anchors[imgname]; !ok {
			//default to layout anchor if not in the list
			offset = layout.Anchor
		} else {

			offset = DiffPosition(layout.Anchor, thisAnchor)
			//fmt.Printf("Previous: %v %v\n", layout.Anchor, offset)
		}

		imgfilename := imgname //in case not specified, e.g. previous image

		if filename, ok := layout.Filenames[imgname]; ok {
			imgfilename = fmt.Sprintf("%s.jpg", filename)
		}
		// append chrome image to the images list

		image := ImageInsert{
			Filename: imgfilename,
			Corner:   offset,
			Dim:      layout.ImageDims[imgname], //TODO need to get dim from images layer
		}
		spread.Images = append(spread.Images, image) //add chrome to list of images to include

	}

	offset := DiffPosition(layout.Anchors["img-previous-mark"], layout.Anchor) //TODO change example layout & parserto image-previous-mark
	//fmt.Printf("image-mark: %v %v\n", layout.Anchor, offset)
	previousImage := ImageInsert{
		Filename: "./test/script.jpg",
		Corner:   offset,                            //geo.Point{X: 0, Y: 61.5},               //offset,
		Dim:      layout.ImageDims["previous-mark"], //
	}

	// draw images
	c := creator.New()
	c.SetPageMargins(0, 0, 0, 0) // we're not printing

	img, err := c.NewImageFromFile(previousImage.Filename)

	if err != nil {
		t.Errorf("Error opening image file: %s", err)
	}

	//see timdrysdale/pagescale if confused
	if spread.Dim.DynamicWidth {
		img.ScaleToHeight(spread.Dim.Height)
		spread.Dim.Width = img.Width()
	} else {
		imgScaledWidth := img.Width() * spread.Dim.Height / img.Height()

		if imgScaledWidth > spread.Dim.Width {
			// oops, we're too big, so scale using width instead
			img.ScaleToWidth(spread.Dim.Width)
		} else {
			img.ScaleToHeight(spread.Dim.Height)
		}

	}

	img.SetPos(previousImage.Corner.X, previousImage.Corner.Y)

	c.SetPageSize(creator.PageSize{spread.Dim.Width, spread.Dim.Height})
	c.NewPage()
	c.Draw(img)

	for _, v := range spread.Images {
		//fmt.Printf("Printing image %s to pdf\n", v.Filename)
		img, err := c.NewImageFromFile(v.Filename)

		if err != nil {
			t.Errorf("Error opening image file: %s", err)
		}
		// all these images are static
		img.SetWidth(v.Dim.Width)
		img.SetHeight(v.Dim.Height)
		img.SetPos(v.Corner.X, v.Corner.Y) //TODO check this has correct sense for non-zero offsets
		//fmt.Printf("Setting position to (%f, %f)\n------------------\n", v.Corner.X, v.Corner.Y)
		// create new page with image

		c.Draw(img)
	}

	// write to memory
	var buf bytes.Buffer

	err = c.Write(&buf)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// convert buffer to readseeker
	var bufslice []byte
	fbuf := filebuffer.New(bufslice)
	fbuf.Write(buf.Bytes())

	// read in from memory
	pdfReader, err := model.NewPdfReader(fbuf)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	pdfWriter := model.NewPdfWriter()

	page, err := pdfReader.GetPage(1)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	form := model.NewPdfAcroForm()

	for _, tf := range spread.TextFields {

		tfopt := annotator.TextFieldOptions{Value: tf.Prefill} //TODO - MaxLen?!
		name := fmt.Sprintf("Page-00-%s", tf.ID)
		textf, err := annotator.NewTextField(page, name, formRect(tf.Rect, layout.Dim), tfopt)
		if err != nil {
			panic(err)
		}
		*form.Fields = append(*form.Fields, textf.PdfField)
		page.AddAnnotation(textf.Annotations[0].PdfAnnotation)
	}

	err = pdfWriter.SetForms(form)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	err = pdfWriter.AddPage(page)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	of, err := os.Create("./test/mark-spread.pdf")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	defer of.Close()

	pdfWriter.SetOptimizer(optimize.New(optimize.Options{
		CombineDuplicateDirectObjects:   true,
		CombineIdenticalIndirectObjects: true,
		CombineDuplicateStreams:         true,
		CompressStreams:                 true,
		UseObjectStreams:                true,
		ImageQuality:                    80,
		ImageUpperPPI:                   100,
	}))

	pdfWriter.Write(of)
}

// gs -dNOPAUSE -sDEVICE=jpeg -sOutputFile=mark-spread-gs.jpg -dJPEGQ=95 -r300 -q mark-spread.pdf -c quit
func TestRenderSpreadMark(t *testing.T) {

	svgLayoutPath := "./test/layout-312pt.svg"

	pdfOutputPath := "./test/render-mark-spread.pdf"

	expectedPdf := "./expected/render-mark-spread.pdf"

	diffPdf := "./test/render-mark-spread-diff.pdf"

	previousImagePath := "./test/script.jpg"

	spreadName := "mark"

	pageNumber := int(16)
	contents := SpreadContents{
		SvgLayoutPath:         svgLayoutPath,
		SpreadName:            spreadName,
		PreviousImagePath:     previousImagePath,
		PageNumber:            pageNumber,
		PdfOutputPath:         pdfOutputPath,
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

func TestRenderSpreadMarkOldAndNewComments(t *testing.T) {

	var comments = make(map[int][]comment.Comment)

	comments[0] = []comment.Comment{c00}

	comments[1] = []comment.Comment{c10, c11}

	comments[2] = []comment.Comment{c20, c21}

	svgLayoutPath := "./test/layout-312pt.svg"

	pdfOutputPath := "./test/render-mark-spread-old-and-new-commments.pdf"

	previousImagePath := "./test/script.jpg"

	spreadName := "mark"

	pageNumber := int(1)

	thisPageData := pagedata.PageData{
		Current: pagedata.PageDetail{
			Comments: []comment.Comment{
				comment.Comment{
					Pos:   geo.Point{X: 120, Y: 300},
					Text:  "Old comment!",
					Page:  1,
					Label: "-ABC",
				},
				comment.Comment{
					Pos:   geo.Point{X: 240, Y: 200},
					Text:  "Another Old comment",
					Page:  1,
					Label: "-ABC",
				},
			},
		},
		Previous: []pagedata.PageDetail{
			pagedata.PageDetail{
				Process: pagedata.ProcessDetail{
					For: "ABC",
				},
			},
		},
	}

	contents := SpreadContents{
		SvgLayoutPath:         svgLayoutPath,
		SpreadName:            spreadName,
		PreviousImagePath:     previousImagePath,
		PageNumber:            pageNumber,
		PdfOutputPath:         pdfOutputPath,
		Comments:              comments,
		PageData:              thisPageData,
		TemplatePathsRelative: true,
	}

	err := RenderSpreadExtra(contents)
	assert.NoError(t, err)

	result, err := visuallyIdenticalPDF(pdfOutputPath, "./expected/render-mark-spread-old-and-new-commments.pdf", "./test/render-mark-spread-old-and-new-commments-diff.pdf")

	assert.NoError(t, err)
	assert.True(t, result)

}

func TestRenderSpreadMarkComment(t *testing.T) {

	var comments = make(map[int][]comment.Comment)

	comments[0] = []comment.Comment{c00}

	comments[1] = []comment.Comment{c10, c11}

	comments[2] = []comment.Comment{c20, c21}

	svgLayoutPath := "./test/layout-312pt.svg"

	pdfOutputPath := "./test/render-mark-spread-commments.pdf"

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
	assert.NoError(t, err)

	result, err := visuallyIdenticalPDF(pdfOutputPath, "./expected/render-mark-spread-commments.pdf",
		"./test/render-mark-spread-commments-diff.pdf")

	assert.NoError(t, err)
	assert.True(t, result)

}

func TestRenderSpreadModerate(t *testing.T) {

	svgLayoutPath := "./test/layout-312pt.svg"

	pdfOutputPath := "./test/render-moderate-active-spread.pdf"
	expectedPdf := "./expected/render-moderate-active-spread.pdf"
	diffPdf := "./test/render-moderate-active-spread-diff.pdf"

	previousImagePath := "./test/mark-spread-gs.jpg"

	spreadName := "moderate-active"

	pageNumber := int(16)

	contents := SpreadContents{
		SvgLayoutPath:         svgLayoutPath,
		SpreadName:            spreadName,
		PreviousImagePath:     previousImagePath,
		PageNumber:            pageNumber,
		PdfOutputPath:         pdfOutputPath,
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

func TestRenderSpreadCheck(t *testing.T) {

	svgLayoutPath := "./test/layout-312pt.svg"

	pdfOutputPath := "./test/render-check-spread.pdf"
	expectedPdf := "./expected/render-check-spread.pdf"
	diffPdf := "./test/render-check-spread-diff.pdf"

	previousImagePath := "./test/moderate-active-gs.jpg"

	spreadName := "check"

	pageNumber := int(16)
	contents := SpreadContents{
		SvgLayoutPath:         svgLayoutPath,
		SpreadName:            spreadName,
		PreviousImagePath:     previousImagePath,
		PageNumber:            pageNumber,
		PdfOutputPath:         pdfOutputPath,
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

func TestRenderSpreadCheckAfterInactive(t *testing.T) {

	svgLayoutPath := "./test/layout-312pt.svg"

	pdfOutputPath := "./test/render-check-after-inactive-spread.pdf"
	expectedPdf := "./expected/render-check-after-inactive-spread.pdf"
	diffPdf := "./test/render-check-after-inactive-spread-diff.pdf"

	previousImagePath := "./test/moderate-inactive-spread-gs.jpg"

	spreadName := "check"

	pageNumber := int(16)
	contents := SpreadContents{
		SvgLayoutPath:         svgLayoutPath,
		SpreadName:            spreadName,
		PreviousImagePath:     previousImagePath,
		PageNumber:            pageNumber,
		PdfOutputPath:         pdfOutputPath,
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

func testPrettyPrintLayout(t *testing.T) {
	// helper for writing the tests on this file - not actually a test
	svgFilename := "./test/layout-312pt-static-mark-dynamic-moderate-comment-static-check.svg"
	//svgFilename := "./test/layout-312pt-static-mark-dynamic-moderate-static-check-v2.svg"
	svgBytes, err := ioutil.ReadFile(svgFilename)

	if err != nil {
		t.Error(err)
	}

	got, err := DefineLayoutFromSVG(svgBytes)
	if err != nil {
		t.Errorf("Error defining layout %v", err)
	}
	PrettyPrintLayout(got)

}
