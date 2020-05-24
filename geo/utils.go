package geo

import (
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
