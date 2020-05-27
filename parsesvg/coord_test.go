package parsesvg

import (
	"image"
	"image/jpeg"
	"log"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/timdrysdale/gradex-cli/optical"
)

func TestGetTextFieldSpread(t *testing.T) {

	svgLayoutPath := "./test/layout-312pt.svg"
	spreadName := "mark"

	spread, err := GetTextFieldSpread(svgLayoutPath, spreadName)
	assert.NoError(t, err)

	tfMap := make(map[string]TextField)

	for _, tf := range spread.TextFields {

		tfMap[tf.ID] = tf
	}

	// remember not to use the similarish-but-wrong figures stored in the Ladder
	width := 178.58267716535434
	badX := 128.6072223023622
	badY := 119.15010310796221

	goodX := 127.05220371968505
	goodY := 564.4731701945764

	assert.Equal(t, goodX, tfMap["page-ok"].Rect.Corner.X)
	assert.Equal(t, goodY, tfMap["page-ok"].Rect.Corner.Y)
	assert.Equal(t, badX, tfMap["page-bad"].Rect.Corner.X)
	assert.Equal(t, badY, tfMap["page-bad"].Rect.Corner.Y)
	assert.Equal(t, width, spread.Dim.Width)

	err = SwapTextFieldXCoords(&spread)
	tfMap = make(map[string]TextField)

	for _, tf := range spread.TextFields {
		tfMap[tf.ID] = tf
	}

	swapBadX := width - badX
	swapGoodX := width - goodX

	assert.Equal(t, swapGoodX, tfMap["page-ok"].Rect.Corner.X)
	assert.Equal(t, goodY, tfMap["page-ok"].Rect.Corner.Y)
	assert.Equal(t, swapBadX, tfMap["page-bad"].Rect.Corner.X)
	assert.Equal(t, badY, tfMap["page-bad"].Rect.Corner.Y)
	assert.Equal(t, width, spread.Dim.Width)

	assert.NoError(t, err)

	textfields, err := GetTextFieldsByTopRight(svgLayoutPath, spreadName)

	assert.NoError(t, err)

	tfMap = make(map[string]TextField)

	for _, tf := range textfields {
		tfMap[tf.ID] = tf
	}

	assert.Equal(t, swapGoodX, tfMap["page-ok"].Rect.Corner.X)
	assert.Equal(t, goodY, tfMap["page-ok"].Rect.Corner.Y)
	assert.Equal(t, swapBadX, tfMap["page-bad"].Rect.Corner.X)
	assert.Equal(t, badY, tfMap["page-bad"].Rect.Corner.Y)
	assert.Equal(t, width, spread.Dim.Width)

}

func TestGetImageBoxes(t *testing.T) {

	svgLayoutPath := "./test/layout-312pt.svg"
	spreadName := "mark"

	widthPx := 1883
	heightPx := 2150

	vanilla := true

	boxes, err := GetImageBoxesForTextFields(svgLayoutPath, spreadName, widthPx, heightPx, vanilla, -2)
	assert.NoError(t, err)

	boxMap := make(map[string]optical.Box)

	for _, bx := range boxes {
		boxMap[bx.ID] = bx
	}

	// TODO verify these box values
	assert.Equal(t, 1764, boxMap["page-bad"].Bounds.Min.X)
	assert.Equal(t, 292, boxMap["page-bad"].Bounds.Min.Y)
	assert.Equal(t, 1816, boxMap["page-bad"].Bounds.Max.X)
	assert.Equal(t, 342, boxMap["page-bad"].Bounds.Max.Y)

	assert.Equal(t, 1760, boxMap["page-ok"].Bounds.Min.X)
	assert.Equal(t, 1374, boxMap["page-ok"].Bounds.Min.Y)
	assert.Equal(t, 1812, boxMap["page-ok"].Bounds.Max.X)
	assert.Equal(t, 1425, boxMap["page-ok"].Bounds.Max.Y)

	// Test images created manually from marked exam
	//$ gs -dNOPAUSE -sDEVICE=jpeg -sOutputFile=stylus-%d.jpg -dJPEGQ=90 -r175 -q Practice-B999999-maTDD-marked-stylus.pdf -c quit

	// get some box values....

	pageBoxes := []optical.Box{
		boxMap["page-bad"],
		boxMap["page-ok"],
	}

	reader, err := os.Open("./img/stylus-1.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	testImage, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}

	actual, checkImage, _ := optical.CheckBoxDebug(testImage, pageBoxes[0])

	assert.Equal(t, false, actual)

	actualBadPath := "./img/stylus-1-bad.jpg"
	of, err := os.Create(actualBadPath)
	assert.NoError(t, err)

	defer of.Close()

	err = jpeg.Encode(of, checkImage, nil)
	assert.NoError(t, err)

	actual, checkImage, _ = optical.CheckBoxDebug(testImage, pageBoxes[1])

	assert.Equal(t, true, actual)

	actualOkPath := "./img/stylus-1-ok.jpg"

	of, err = os.Create(actualOkPath)

	assert.NoError(t, err)

	defer of.Close()

	err = jpeg.Encode(of, checkImage, nil)
	assert.NoError(t, err)

	// do visual checks on boxes

	expectedBadPath := "./img/stylus-1-bad-expected.jpg"
	expectedOkPath := "./img/stylus-1-ok-expected.jpg"

	_, err = exec.Command("compare", "-metric", "ae", actualBadPath, expectedBadPath, "null:").CombinedOutput()
	assert.NoError(t, err) //throws error if not same

	_, err = exec.Command("compare", "-metric", "ae", actualOkPath, expectedOkPath, "null:").CombinedOutput()
	assert.NoError(t, err) //throws error if not same

	// Try with the production version
	results, err := optical.CheckBoxFile("./img/stylus-1.jpg", pageBoxes)

	assert.NoError(t, err)

	assert.Equal(t, []bool{false, true}, results)

	//PrettyPrintStruct(boxes)
}

//BenchmarkGetImageBoxes-32    	      73	  35286670 ns/op
// 35.286670 ms/op
func BenchmarkGetImageBoxes(b *testing.B) {

	svgLayoutPath := "./test/layout-312pt.svg"
	spreadName := "mark"
	widthPx := 1883
	heightPx := 2150
	vanilla := true

	// run the function b.N times
	for n := 0; n < b.N; n++ {

		_, err := GetImageBoxesForTextFields(svgLayoutPath, spreadName, widthPx, heightPx, vanilla, -2)

		assert.NoError(b, err)

	}

}

//BenchmarkDoCheckBoxesWholePage-32    	      12	  92800839 ns/op
// 92.800839 ms
// for reading in the file and checking all the boxes (17 subtotal, 2 page, 6 question)
func BenchmarkDoCheckBoxesWholePage(b *testing.B) {

	svgLayoutPath := "./test/layout-312pt.svg"
	spreadName := "mark"
	widthPx := 1883
	heightPx := 2150
	vanilla := true
	boxes, err := GetImageBoxesForTextFields(svgLayoutPath, spreadName, widthPx, heightPx, vanilla, -2)
	assert.NoError(b, err)

	// run the function b.N times
	for n := 0; n < b.N; n++ {

		_, err := optical.CheckBoxFile("./img/stylus-1.jpg", boxes)

		assert.NoError(b, err)

	}

}

// BenchmarkDoCheckBoxesOneBox-32    	      14	  79123245 ns/op
// 79.123245 ms
// for reading in the file and checking one box
// Hence most of the time is in loading the file
// So may as well optically read all the boxes on the page in one go
func BenchmarkDoCheckBoxesOneBox(b *testing.B) {

	svgLayoutPath := "./test/layout-312pt.svg"
	spreadName := "mark"
	widthPx := 1883
	heightPx := 2150
	vanilla := true
	boxes, err := GetImageBoxesForTextFields(svgLayoutPath, spreadName, widthPx, heightPx, vanilla, -2)
	assert.NoError(b, err)

	// run the function b.N times
	for n := 0; n < b.N; n++ {

		_, err := optical.CheckBoxFile("./img/stylus-1.jpg", []optical.Box{boxes[0]})

		assert.NoError(b, err)

	}

}