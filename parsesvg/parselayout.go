package parsesvg

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/timdrysdale/gradex-cli/geo"
)

func DefineLayoutFromSVG(input []byte) (*Layout, error) {

	var svg Csvg__svg
	layout := &Layout{}

	err := xml.Unmarshal(input, &svg)

	if err != nil {
		return nil, err
	}

	// get title
	if svg.Cmetadata__svg.CRDF__rdf != nil {
		if svg.Cmetadata__svg.CRDF__rdf.CWork__cc != nil {
			if svg.Cmetadata__svg.CRDF__rdf.CWork__cc.Ctitle__dc != nil {
				layout.ID = svg.Cmetadata__svg.CRDF__rdf.CWork__cc.Ctitle__dc.String
			}
		}
	}

	layout.Anchor = geo.Point{X: 0, Y: 0}

	layoutDim, err := getLadderDim(&svg)
	if err != nil {
		return nil, err
	}

	layout.Dim = layoutDim

	// look for reference & header/ladder anchor positions
	// these also contain the base filename in the description
	for _, g := range svg.Cg__svg {
		// get transform applied to layer, if any
		if g.AttrInkscapeSpacelabel == geo.AnchorsLayer {
			gdx, gdy := getTranslate(g.Transform)

			layout.Anchors = make(map[string]geo.Point)
			layout.Filenames = make(map[string]string)

			for _, r := range g.Cpath__svg {
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

				if r.Title != nil {
					if r.Title.String == geo.AnchorReference {

						layout.Anchor = geo.Point{X: newX, Y: newY}
					} else {

						layout.Anchors[r.Title.String] = geo.Point{X: newX, Y: newY}

						if r.Desc != nil {
							layout.Filenames[r.Title.String] = r.Desc.String
						}
					}
				} else {
					log.Errorf("Anchor at (%f,%f) has no title, so ignoring\n", newX, newY)
				}
			}
		}
	}

	// look for pageDims
	layout.PageDims = make(map[string]geo.Dim)
	for _, g := range svg.Cg__svg {

		if g.AttrInkscapeSpacelabel == geo.PagesLayer {
			for _, r := range g.Crect__svg {
				w, err := strconv.ParseFloat(r.Width, 64)
				if err != nil {
					return nil, err
				}
				h, err := strconv.ParseFloat(r.Height, 64)
				if err != nil {
					return nil, err
				}

				if r.Title != nil { //avoid seg fault, obvs

					fullname := r.Title.String
					name := ""
					isDynamic := false

					switch {
					case strings.HasPrefix(fullname, "page-dynamic-"):
						name = strings.TrimPrefix(fullname, "page-dynamic-")
						isDynamic = true
					case strings.HasPrefix(fullname, "page-static-"):
						name = strings.TrimPrefix(fullname, "page-static-")
					default:
						// unadorned pages are considered static
						// because this is the least surprising behaviour
						name = strings.TrimPrefix(fullname, "page-")
					}

					if name != "" { //reject anonymous pages
						layout.PageDims[name] = geo.Dim{Width: w, Height: h, DynamicWidth: isDynamic}
					}

				} else {
					log.Errorf("Page at with size (%f,%f) has no title, so ignoring\n", w, h)
				}
			}
		}
	}
	// look for previousImageDims
	layout.ImageDims = make(map[string]geo.Dim)
	for _, g := range svg.Cg__svg {
		if g.AttrInkscapeSpacelabel == geo.ImagesLayer {
			for _, r := range g.Crect__svg {
				w, err := strconv.ParseFloat(r.Width, 64)
				if err != nil {
					return nil, err
				}
				h, err := strconv.ParseFloat(r.Height, 64)
				if err != nil {
					return nil, err
				}

				if r.Title != nil { //avoid seg fault, obvs

					fullname := r.Title.String
					name := ""
					isDynamic := false

					switch {
					case strings.HasPrefix(fullname, "image-dynamic-"):
						name = strings.TrimPrefix(fullname, "image-dynamic-")
						name = strings.TrimPrefix(name, "width-")  //we may want this later, so leave in API
						name = strings.TrimPrefix(name, "height-") //getting info from box size for now
						isDynamic = true
					case strings.HasPrefix(fullname, "image-static-"):
						name = strings.TrimPrefix(fullname, "image-static-")
					default:
						// we're just trying to strip off prefixes,
						// not prevent underadorned names from working
						name = strings.TrimPrefix(fullname, "image-")
					}

					if name != "" { //reject anonymous images - can't place them
						layout.ImageDims[name] = geo.Dim{Width: w, Height: h, DynamicWidth: isDynamic}
					}

				} else {
					log.Errorf("Page at with size (%f,%f) has no title, so ignoring\n", w, h)
				}
			}
		}
	}

	err = ApplyDocumentUnitsScaleLayout(&svg, layout)
	if err != nil {
		return nil, err
	}

	return layout, nil
}

func ApplyDocumentUnitsScaleLayout(svg *Csvg__svg, layout *Layout) error {

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

	layout.Anchor.X = sf * layout.Anchor.X
	layout.Anchor.Y = sf * layout.Anchor.Y

	for k, v := range layout.Anchors {
		v.X = sf * v.X
		v.Y = sf * v.Y
		layout.Anchors[k] = v

	}
	for k, v := range layout.PageDims {
		v.Width = sf * v.Width
		v.Height = sf * v.Height
		layout.PageDims[k] = v

	}

	for k, v := range layout.ImageDims {
		v.Width = sf * v.Width
		v.Height = sf * v.Height
		layout.ImageDims[k] = v
	}

	return nil
}

func PrettyPrintLayout(layout *Layout) error {

	json, err := json.MarshalIndent(layout, "", "\t")
	if err != nil {
		return err
	}

	fmt.Println(string(json))
	return nil
}

func PrintLayout(layout *Layout) error {

	json, err := json.Marshal(layout)
	if err != nil {
		return err
	}

	fmt.Println(string(json))
	return nil
}

func PrettyPrintStruct(layout interface{}) error {

	json, err := json.MarshalIndent(layout, "", "\t")
	if err != nil {
		return err
	}

	fmt.Println(string(json))
	return nil
}

type Action struct {
	Verb   string
	When   time.Time
	Who    string
	Params map[string]string
}

type MetaData struct {
	Exam      string
	Candidate string
	Diet      string
	Actions   []Action
}
