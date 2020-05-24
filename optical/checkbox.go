package optical

import (
	"errors"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"os"
)

type Box struct {
	Vanilla bool
	Bounds  image.Rectangle
	ID      string
}

// https://gist.github.com/sergiotapia/7882944
func GetImageDimension(imagePath string) (int, int, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return 0, 0, err
	}

	image, _, err := image.DecodeConfig(file)
	if err != nil {
		return 0, 0, err
	}
	return image.Width, image.Height, nil
}

func ExpandBound(box *Box, extra int) error {

	if box == nil {
		return errors.New("nil pointer to box")
	}

	(*box).Bounds.Min.X = (*box).Bounds.Min.X - extra
	(*box).Bounds.Min.Y = (*box).Bounds.Min.Y - extra
	(*box).Bounds.Max.X = (*box).Bounds.Max.X + extra
	(*box).Bounds.Max.Y = (*box).Bounds.Max.Y + extra

	return nil

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

	vanillaThreshold := uint32(math.Round(0.99 * 3 * 65535))
	chocolateThreshold := uint32(math.Round(0.01 * 3 * 65535))

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := checkImage.At(x, y).RGBA()
			luminosity := r + g + b

			if box.Vanilla && luminosity < vanillaThreshold {
				cum++ // count non-vanilla pixels, on vanilla background)
			}
			if !box.Vanilla && luminosity > chocolateThreshold {
				cum++ // count non-chocolate pixels, on chocolate background
			}
		}
	}

	colourCount := 3
	pixelCount := colourCount * (bounds.Max.X - bounds.Min.X) * (bounds.Max.Y - bounds.Min.Y)
	markedPixelFraction := float64(cum) / float64(pixelCount)
	thresh := 0.02
	return markedPixelFraction > thresh, checkImage, markedPixelFraction
}

func CheckBox(im image.Image, box Box) bool {

	result, _, _ := CheckBoxDebug(im, box)

	return result
}

// This will do handwriting recognition in future (wishlist!)
func DataBox(im image.Image, box Box) (bool, image.Image) {

	result, img, _ := CheckBoxDebug(im, box)

	return result, img

}
