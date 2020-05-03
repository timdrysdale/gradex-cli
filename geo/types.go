package geo

type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type Dim struct {
	Width        float64 `json:"width"`
	Height       float64 `json:"height"`
	DynamicWidth bool    `json:"dynamicWidth"` //dynamicW
}

type Rect struct {
	Corner Point `json:"corner"`
	Dim    Dim   `json:"dim"`
}

const (
	AnchorsLayer      = "anchors"
	AnchorReference   = "ref-anchor"
	ChromeLayer       = "chrome"
	TextFieldsLayer   = "textfields"
	TextPrefillsLayer = "textprefills"
	PagesLayer        = "pages"
	ImagesLayer       = "images"
	Translate         = "translate"
	SVGElement        = "svg-"
	JPGElement        = "jpg-"
	Previous          = "previous"
	PXIN              = float64(96)                //pixels per inch
	PPIN              = float64(72)                //points per inch
	PPMM              = float64(PPIN * 1.0 / 25.4) //points per millimetre
	PPPX              = float64(PPIN / PXIN)       // points per pixel

)
