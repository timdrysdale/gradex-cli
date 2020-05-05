package parsesvg

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strconv"

	"strings"

	"github.com/timdrysdale/gradex-cli/geo"
)

func ParseSvg(input []byte) *Csvg__svg {

	var svg Csvg__svg

	xml.Unmarshal(input, &svg)

	return &svg
}

// return the translation applied to the SVG object
// does not mangle the Y dimension
func getTranslate(transform string) (float64, float64) {

	if len(transform) <= 0 {
		return 0.0, 0.0
	}

	if !strings.Contains(transform, geo.Translate) {
		return 0.0, 0.0
	}

	openBracket := strings.Index(transform, "(")
	comma := strings.Index(transform, ",")
	closeBracket := strings.Index(transform, ")")

	if openBracket < 0 || comma < 0 || closeBracket < 0 {
		return 0.0, 0.0
	}

	if openBracket == comma || comma == closeBracket {
		return 0.0, 0.0
	}

	dx, err := strconv.ParseFloat(transform[openBracket+1:comma], 64)
	if err != nil {
		return 0.0, 0.0
	}
	dy, err := strconv.ParseFloat(transform[comma+1:closeBracket], 64)
	if err != nil {
		return 0.0, 0.0
	}

	return dx, dy

}

func scanUnitStringToPP(str string) (float64, error) {

	str = strings.TrimSpace(str)
	length := len(str)
	units := str[length-2 : length]
	value, err := strconv.ParseFloat(str[0:length-2], 64)
	if err != nil {
		return 0, errors.New(fmt.Sprintf("Couldn't parse  %s when split into value %s with units %s", str, str[0:length-2], units))
	}

	switch units {
	case "mm":
		return value * geo.PPMM, nil
	case "px":
		return value * geo.PPPX, nil
	case "pt":
		return value, nil //TODO check pt doesn't somehow default to not present
	case "in":
		return value * geo.PPIN, nil
	}

	return 0, errors.New(fmt.Sprintf("didn't understand the units %s", units))

}

func getLadderDim(svg *Csvg__svg) (geo.Dim, error) {
	dim := geo.Dim{}

	if svg == nil {
		return dim, errors.New("nil pointer to svg")
	}

	w, err := scanUnitStringToPP(svg.Width)
	if err != nil {
		return dim, err
	}
	h, err := scanUnitStringToPP(svg.Height)
	if err != nil {
		return dim, err
	}

	return geo.Dim{Width: w, Height: h, DynamicWidth: false}, nil

}

func DefineLadderFromSVG(input []byte) (*Ladder, error) {

	var svg Csvg__svg
	ladder := &Ladder{}

	err := xml.Unmarshal(input, &svg)

	if err != nil {
		return nil, err
	}

	if svg.Cmetadata__svg.CRDF__rdf != nil {
		if svg.Cmetadata__svg.CRDF__rdf.CWork__cc != nil {
			if svg.Cmetadata__svg.CRDF__rdf.CWork__cc.Ctitle__dc != nil {
				ladder.ID = svg.Cmetadata__svg.CRDF__rdf.CWork__cc.Ctitle__dc.String
			}
		}
	}

	//fmt.Printf("-------(%s)-----------\n", ladder.ID)
	ladder.Anchor = geo.Point{X: 0, Y: 0}

	ladderDim, err := getLadderDim(&svg)
	if err != nil {
		return nil, err
	}

	ladder.Dim = ladderDim
	// Note to future self - had a global dx, dy here
	// but doing translate parsing properly should avoid need for that...
	// look for reference anchor position
	for _, g := range svg.Cg__svg {

		if g.AttrInkscapeSpacelabel == geo.AnchorsLayer {

			gdx, gdy := getTranslate(g.Transform)

			for _, r := range g.Cpath__svg {
				if r.Title != nil {
					if r.Title.String == geo.AnchorReference { // was force true?
						x, err := strconv.ParseFloat(r.Cx, 64)
						if err != nil {
							return nil, err
						}
						y, err := strconv.ParseFloat(r.Cy, 64)
						if err != nil {
							return nil, err
						}
						rdx, rdy := getTranslate(r.Transform)
						newX := x + gdx + rdx
						newY := y + gdy + rdy
						ladder.Anchor = geo.Point{X: newX, Y: newY}
					}
				}
			}
		}

	}

	// look for textFields
	for _, g := range svg.Cg__svg {
		if g.AttrInkscapeSpacelabel == geo.TextFieldsLayer {
			gdx, gdy := getTranslate(g.Transform)
			for _, r := range g.Crect__svg {
				tf := TextField{}
				if r.Title != nil { //avoid seg fault, obvs
					tf.ID = r.Title.String
				}

				tf.TabSequence = getTabSequence(r)

				if r.Desc != nil {
					tf.Prefill = r.Desc.String
				}
				w, err := strconv.ParseFloat(r.Width, 64)
				if err != nil {
					return nil, err
				}
				h, err := strconv.ParseFloat(r.Height, 64)
				if err != nil {
					return nil, err
				}

				tf.Rect.Dim.Width = w
				tf.Rect.Dim.Height = h
				tf.Rect.Dim.DynamicWidth = false
				rdx, rdy := getTranslate(r.Transform)
				x, err := strconv.ParseFloat(r.Rx, 64)
				if err != nil {
					return nil, err
				}
				y, err := strconv.ParseFloat(r.Ry, 64)
				if err != nil {
					return nil, err
				}
				tf.Rect.Corner.X = x + rdx + gdx
				tf.Rect.Corner.Y = y + rdy + gdy

				ladder.TextFields = append(ladder.TextFields, tf)
			}
		}

	}
	// sort textfields based on tab order

	sort.Slice(ladder.TextFields, func(i, j int) bool {
		return ladder.TextFields[i].TabSequence < ladder.TextFields[j].TabSequence
	})

	// look for prefill textboxes (not editable in pdf)

	for _, g := range svg.Cg__svg {
		gdx, gdy := getTranslate(g.Transform)
		if g.AttrInkscapeSpacelabel == geo.TextPrefillsLayer {

			for _, r := range g.Crect__svg {
				tp := TextPrefill{}
				if r.Title != nil { //avoid seg fault, obvs
					tp.ID = r.Title.String
				}

				if r.Desc != nil {
					tp.Properties = r.Desc.String
				}
				w, err := strconv.ParseFloat(r.Width, 64)
				if err != nil {
					return nil, err
				}
				h, err := strconv.ParseFloat(r.Height, 64)
				if err != nil {
					return nil, err
				}

				tp.Rect.Dim.Width = w
				tp.Rect.Dim.Height = h
				tp.Rect.Dim.DynamicWidth = false

				x, err := strconv.ParseFloat(r.Rx, 64)
				if err != nil {
					return nil, err
				}
				y, err := strconv.ParseFloat(r.Ry, 64)
				if err != nil {
					return nil, err
				}
				rdx, rdy := getTranslate(r.Transform)
				tp.Rect.Corner.X = x + rdx + gdx
				tp.Rect.Corner.Y = y + rdy + gdy

				err = UnmarshalTextPrefill(&tp)
				if err != nil {
					return nil, err
				}
				ladder.TextPrefills = append(ladder.TextPrefills, tp)
			}
		}

	}

	for _, g := range svg.Cg__svg {
		gdx, gdy := getTranslate(g.Transform)
		if g.AttrInkscapeSpacelabel == geo.ComboBoxesLayer {

			for _, r := range g.Crect__svg {
				cb := ComboBox{}
				if r.Title != nil { //avoid seg fault, obvs
					cb.ID = r.Title.String
				}

				if r.Desc != nil {
					cb.Properties = r.Desc.String
				}
				w, err := strconv.ParseFloat(r.Width, 64)
				if err != nil {
					return nil, err
				}
				h, err := strconv.ParseFloat(r.Height, 64)
				if err != nil {
					return nil, err
				}

				cb.Rect.Dim.Width = w
				cb.Rect.Dim.Height = h
				cb.Rect.Dim.DynamicWidth = false

				x, err := strconv.ParseFloat(r.Rx, 64)
				if err != nil {
					return nil, err
				}
				y, err := strconv.ParseFloat(r.Ry, 64)
				if err != nil {
					return nil, err
				}
				rdx, rdy := getTranslate(r.Transform)
				cb.Rect.Corner.X = x + rdx + gdx
				cb.Rect.Corner.Y = y + rdy + gdy

				err = UnmarshalComboBox(&cb)
				if err != nil {
					return nil, err
				}
				ladder.ComboBoxes = append(ladder.ComboBoxes, cb)
			}
		}

	}

	err = ApplyDocumentUnits(&svg, ladder)
	if err != nil {
		return nil, err
	}

	return ladder, nil
}

func UnmarshalComboBox(cb *ComboBox) error {

	options := ComboOptions{}

	if len(cb.Properties) > 0 {
		err := json.Unmarshal([]byte(cb.Properties), &options)
		if err != nil {
			return err
		}

		cb.Options = options
	}
	return nil

}

func UnmarshalTextPrefill(tp *TextPrefill) error {

	var paragraph Paragraph
	properties := "{\"text\":\"\"}"
	if len(tp.Properties) > 0 {
		properties = tp.Properties
	}
	err := json.Unmarshal([]byte(properties), &paragraph)
	if err != nil {
		return err
	}

	tp.Text = paragraph

	return nil

}

func ApplyDocumentUnits(svg *Csvg__svg, ladder *Ladder) error {

	// iterate through the structure applying the conversion from
	// document units to points

	//note we do NOT apply the modification to ladder.DIM because this has its own
	//units in it and has already been handled.

	units := svg.Cnamedview__sodipodi.AttrInkscapeSpacedocument_dash_units

	sf := float64(1)

	switch units {
	case "mm":
		sf = geo.PPMM
	case "px":
		sf = geo.PPPX
	case "pt":
		sf = 1
	case "in":
		sf = geo.PPIN
	}

	ladder.Anchor.X = sf * ladder.Anchor.X
	ladder.Anchor.Y = sf * ladder.Anchor.Y

	for idx, tf := range ladder.TextFields {
		err := scaleTextFieldUnits(&tf, sf)
		if err != nil {
			return err
		}
		ladder.TextFields[idx] = tf
	}

	for idx, tp := range ladder.TextPrefills {
		err := scaleTextPrefillUnits(&tp, sf)
		if err != nil {
			return err
		}
		ladder.TextPrefills[idx] = tp
	}

	for idx, cb := range ladder.ComboBoxes {
		err := scaleComboBoxUnits(&cb, sf)
		if err != nil {
			return err
		}
		ladder.ComboBoxes[idx] = cb
	}

	return nil
}

func scaleTextFieldUnits(tf *TextField, sf float64) error {
	if tf == nil {
		return errors.New("nil pointer to TextField")
	}

	tf.Rect.Corner.X = sf * tf.Rect.Corner.X
	tf.Rect.Corner.Y = sf * tf.Rect.Corner.Y
	tf.Rect.Dim.Width = sf * tf.Rect.Dim.Width
	tf.Rect.Dim.Height = sf * tf.Rect.Dim.Height

	return nil
}
func scaleComboBoxUnits(cb *ComboBox, sf float64) error {
	if cb == nil {
		return errors.New("nil pointer to ComboBox")
	}

	cb.Rect.Corner.X = sf * cb.Rect.Corner.X
	cb.Rect.Corner.Y = sf * cb.Rect.Corner.Y
	cb.Rect.Dim.Width = sf * cb.Rect.Dim.Width
	cb.Rect.Dim.Height = sf * cb.Rect.Dim.Height

	return nil
}
func scaleTextPrefillUnits(tf *TextPrefill, sf float64) error {
	if tf == nil {
		return errors.New("nil pointer to TextField")
	}

	tf.Rect.Corner.X = sf * tf.Rect.Corner.X
	tf.Rect.Corner.Y = sf * (tf.Rect.Corner.Y)
	tf.Rect.Dim.Width = sf * tf.Rect.Dim.Width
	tf.Rect.Dim.Height = sf * tf.Rect.Dim.Height

	return nil
}

func formRect(rect geo.Rect, dim geo.Dim) []float64 {

	return []float64{rect.Corner.X, dim.Height - rect.Corner.Y, (rect.Corner.X + rect.Dim.Width), dim.Height - (rect.Corner.Y + rect.Dim.Height)}

}

func getTabSequence(r *Crect__svg) int64 {
	var TabSequence = regexp.MustCompile(`(?i:(tab|tab-))([0-9]+)`)
	var SequenceNumber = regexp.MustCompile(`([0-9]+)`)
	//TODO - combine regexp into one
	var n int64
	n, err := strconv.ParseInt(SequenceNumber.FindString(TabSequence.FindString(r.Id)), 10, 64)
	if err != nil {
		return int64(0)
	}
	return n
}

func TranslatePosition(pos, vec geo.Point) geo.Point {

	return geo.Point{X: pos.X + vec.X, Y: pos.Y + vec.Y}

}
func DiffPosition(from, to geo.Point) geo.Point {

	return geo.Point{X: to.X - from.X, Y: to.Y - from.Y}

}
