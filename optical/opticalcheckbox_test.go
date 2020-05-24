package optical

import (
	"fmt"
	"image"
	jpeg "image/jpeg"
	_ "image/png"
	"log"
	"os"
	"testing"
)

var testCheckBoxes = []Box{
	{Vanilla: true, Bounds: image.Rect(10, 10, 40, 40)},
	{Vanilla: true, Bounds: image.Rect(10, 60, 40, 90)},
	{Vanilla: true, Bounds: image.Rect(60, 10, 90, 40)},
	{Vanilla: true, Bounds: image.Rect(60, 60, 90, 90)},
	{Vanilla: true, Bounds: image.Rect(110, 10, 140, 40)},
	{Vanilla: true, Bounds: image.Rect(110, 60, 140, 90)},
	{Vanilla: true, Bounds: image.Rect(160, 10, 190, 40)},
	{Vanilla: true, Bounds: image.Rect(160, 60, 190, 90)},

	{Vanilla: false, Bounds: image.Rect(10, 10, 40, 40)},
	{Vanilla: false, Bounds: image.Rect(10, 60, 40, 90)},
	{Vanilla: false, Bounds: image.Rect(60, 10, 90, 40)},
	{Vanilla: false, Bounds: image.Rect(60, 60, 90, 90)},
	{Vanilla: false, Bounds: image.Rect(110, 10, 140, 40)},
	{Vanilla: false, Bounds: image.Rect(110, 60, 140, 90)},
	{Vanilla: false, Bounds: image.Rect(160, 10, 190, 40)},
	{Vanilla: false, Bounds: image.Rect(160, 60, 190, 90)},
}

var expectedBox = []bool{
	true,
	false,
	true,
	true,
	true,
	true,
	false,
	true,
	true,
	true,
	true,
	true,
	true,
	false,
	true,
	true,
}

func TestCheckBoxDebug(t *testing.T) {

	reader, err := os.Open("./img/test.png")
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	testImage, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}

	for idx := 0; idx < len(testCheckBoxes); idx = idx + 1 {
		actual, checkImage, avgCount := CheckBoxDebug(testImage, testCheckBoxes[idx])
		wanted := expectedBox[idx]
		if actual != wanted {

			of, err := os.Create("failedTestCheckBox.jpg")

			if err != nil {
				t.Errorf("problem saving failed checkbox image to file %v\n", err)
			}

			defer of.Close()

			err = jpeg.Encode(of, checkImage, nil)

			if err != nil {
				t.Errorf("writing file %v\n", err)
			}

			t.Errorf("Unexpected result for checkbox %d; got %v wanted %v; avg pixel value was %f; see failed.jpg\n",
				idx, actual, wanted, avgCount)
		}

	}

}

func TestCheckBox(t *testing.T) {

	reader, err := os.Open("./img/test.png")
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	testImage, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}

	for idx := 0; idx < len(testCheckBoxes); idx = idx + 1 {
		actual := CheckBox(testImage, testCheckBoxes[idx])
		wanted := expectedBox[idx]
		if actual != wanted {

			t.Errorf("Unexpected result for checkbox %d; got %v wanted %v\n",
				idx, actual, wanted)
		}

	}

}

func TestDataBox(t *testing.T) {

	reader, err := os.Open("./img/test.png")
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	testImage, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}

	for idx := 0; idx < len(testCheckBoxes); idx = idx + 1 {
		actual, _ := DataBox(testImage, testCheckBoxes[idx])
		wanted := expectedBox[idx]
		if actual != wanted {
			t.Errorf("Unexpected result for checkbox %d; got %v wanted %v\n",
				idx, actual, wanted)
		}

	}

}

func TestCheckBoxFile(t *testing.T) {

	results, err := CheckBoxFile("./img/test.png", testCheckBoxes)

	if err != nil {
		t.Errorf("error %v\n", err)
	}

	for idx, result := range results {

		if result != expectedBox[idx] {
			t.Errorf("Unexpected result for checkbox %d; got %v wanted %v\n",
				idx, result, expectedBox[idx])
		}

	}
}

func TestDataBoxFile(t *testing.T) {

	results, images, err := DataBoxFile("./img/test.png", testCheckBoxes)

	if err != nil {
		t.Errorf("error %v\n", err)
	}

	for idx, result := range results {

		of, err := os.Create(fmt.Sprintf("dataBox%d.jpg", idx))

		if err != nil {
			t.Errorf("problem saving failed checkbox image to file %v\n", err)
		}

		defer of.Close()

		err = jpeg.Encode(of, images[idx], nil)

		if result != expectedBox[idx] {
			t.Errorf("Unexpected result for checkbox %d; got %v wanted %v\n",
				idx, result, expectedBox[idx])
		}

	}
}
