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

	xa := math.Round(rect[0])
	xb := math.Round(rect[2])
	ya := math.Round(rect[1])
	yb := math.Round(rect[3])

	minX := int(math.Min(xa, xb))
	minY := int(math.Min(ya, yb))
	maxX := int(math.Max(xa, xb))
	maxY := int(math.Max(ya, yb))

	ir = image.Rectangle{
		Min: image.Point{
			X: minX,
			Y: minY,
		},

		Max: image.Point{
			X: maxX,
			Y: maxY,
		},
	}

	return ir, nil

}
