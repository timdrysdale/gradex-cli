package geo

import (
	"errors"
	"image"
	"math"
)

// For converting SVG to IMG coordinates (for optical checkboxes)
func ConvertToImageRectangle(rect Rect) image.Rectangle {

	ir := image.Rectangle{
		Min: image.Point{
			X: int(math.Round(rect.Corner.X)),
			Y: int(math.Round(rect.Corner.Y)),
		},

		Max: image.Point{
			X: int(math.Round(rect.Corner.X + rect.Dim.Width)),
			Y: int(math.Round(rect.Corner.Y + rect.Dim.Height)),
		},
	}

	return ir

}

// []float64{rect.Corner.X, dim.Height - rect.Corner.Y, (rect.Corner.X + rect.Dim.Width), dim.Height - (rect.Corner.Y + rect.Dim.Height)}

func ConvertPDFRectToImageRectangle(rect []float64) (image.Rectangle, error) {

	ir := image.Rectangle{}

	if len(rect) != 4 {
		return ir, errors.New("expected four elements in input array")
	}

	ir = image.Rectangle{
		Min: image.Point{
			X: int(math.Round(rect[0])),
			Y: int(math.Round(rect[1])),
		},

		Max: image.Point{
			X: int(math.Round(rect[2])),
			Y: int(math.Round(rect[3])),
		},
	}

	return ir, nil

}
