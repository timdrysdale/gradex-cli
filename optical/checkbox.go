package optical

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

type Box struct {
	Vanilla bool
	Bounds  image.Rectangle
	ID      string
}

func CheckBoxFile(inputPath string, boxes []Box) ([]bool, error) {

	var results []bool

	reader, err := os.Open(inputPath)

	if err != nil {
		return []bool{}, err
	}

	defer reader.Close()

	wholeImage, _, err := image.Decode(reader)

	if err != nil {
		return []bool{}, err
	}

	for idx := 0; idx < len(boxes); idx = idx + 1 {
		result := CheckBox(wholeImage, boxes[idx])
		results = append(results, result)
	}

	return results, nil
}

// future aim is to read handwriting - currently behaves same as checkbox
func DataBoxFile(inputPath string, boxes []Box) ([]bool, []image.Image, error) {

	var (
		results []bool
		images  []image.Image
	)

	reader, err := os.Open(inputPath)

	if err != nil {
		return []bool{}, []image.Image{}, err
	}

	defer reader.Close()

	wholeImage, _, err := image.Decode(reader)

	if err != nil {
		return []bool{}, []image.Image{}, err
	}

	for idx := 0; idx < len(boxes); idx = idx + 1 {
		result, img := DataBox(wholeImage, boxes[idx])
		results = append(results, result)
		images = append(images, img)
	}

	return results, images, nil
}

//see https://stackoverflow.com/questions/16072910/trouble-getting-a-subimage-of-an-image-in-go
type SubImager interface {
	SubImage(r image.Rectangle) image.Image
}

func CheckBoxDebug(im image.Image, box Box) (bool, image.Image, float64) {

	checkImage := im.(SubImager).SubImage(box.Bounds)
	cum := uint32(0)
	bounds := checkImage.Bounds()

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := checkImage.At(x, y).RGBA()
			cum = cum + r>>11 + g>>11 + b>>11
		}
	}

	colourCount := 3
	pixelCount := colourCount * (bounds.Max.X - bounds.Min.X) * (bounds.Max.Y - bounds.Min.Y)
	averagePixelValue := float64(cum / uint32(pixelCount))

	if box.Vanilla {
		return averagePixelValue < 30.0, checkImage, averagePixelValue
	} else {
		return averagePixelValue > 1.0, checkImage, averagePixelValue
	}
}

func CheckBox(im image.Image, box Box) bool {

	checkImage := im.(SubImager).SubImage(box.Bounds)
	cum := uint32(0)
	bounds := checkImage.Bounds()

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := checkImage.At(x, y).RGBA()
			cum = cum + r>>11 + g>>11 + b>>11
		}
	}

	colourCount := 3
	pixelCount := colourCount * (bounds.Max.X - bounds.Min.X) * (bounds.Max.Y - bounds.Min.Y)
	averagePixelValue := float64(cum / uint32(pixelCount))

	if box.Vanilla {
		return averagePixelValue < 30.0
	} else {
		return averagePixelValue > 1.0
	}
}

// This will do handwriting recognition in future (wishlist!)
func DataBox(im image.Image, box Box) (bool, image.Image) {

	checkImage := im.(SubImager).SubImage(box.Bounds)
	cum := uint32(0)
	bounds := checkImage.Bounds()

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := checkImage.At(x, y).RGBA()
			cum = cum + r>>11 + g>>11 + b>>11
		}
	}

	colourCount := 3
	pixelCount := colourCount * (bounds.Max.X - bounds.Min.X) * (bounds.Max.Y - bounds.Min.Y)
	averagePixelValue := float64(cum / uint32(pixelCount))

	if box.Vanilla {
		return averagePixelValue < 30.0, checkImage
	} else {
		return averagePixelValue > 1.0, checkImage
	}

}
