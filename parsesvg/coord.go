package parsesvg

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/timdrysdale/gradex-cli/geo"
	"github.com/timdrysdale/gradex-cli/optical"
)

// zero,zero is upper right for this (swap in X) because the page width is dynamic,
// coordinate with respect left hand side of page varies from page to page depending on page si
// this is a cut-down version of render i.e. same logic
// textfield coordinates look obfuscated in the unipdf extraction tool
// so we use our layout to find them, rather than read them from the file
func GetTextFieldSpread(svgLayoutPath, spreadName string) (Spread, error) {

	spread := Spread{}

	svgBytes, err := ioutil.ReadFile(svgLayoutPath)

	if err != nil {
		return spread, errors.New(fmt.Sprintf("Error opening layout file %s: %v\n", svgLayoutPath, err))
	}

	layout, err := DefineLayoutFromSVG(svgBytes)
	if err != nil {
		return spread, errors.New(fmt.Sprintf("Error obtaining layout from svg %s\n", svgLayoutPath))
	}

	spread.Name = spreadName

	foundPage := false
	for k, v := range layout.PageDims {
		if strings.Contains(k, spread.Name) {
			spread.Dim = v
			foundPage = true
		}
	}

	if !foundPage {
		return spread, errors.New(fmt.Sprintf("No page size info for spread %s\n", spread.Name))
	}

	// find svg & img elements for this name
	var svgFilenames []string

	for k, _ := range layout.Filenames {
		if strings.Contains(k, spread.Name) {

			// assume jpg- or no prefix is image; svg- is ladder (image plus acroforms)
			if strings.HasPrefix(k, geo.SVGElement) {
				svgFilenames = append(svgFilenames, k) //we'll get the contents later
			}
		}
	}

	// get all the textfields (and put image of associated chrome into images list)
	// note that if page dynamic, textfields are ALL dynamically shifting wrt to dynamic page edge,
	// no matter what side of the previous image edge they are. This means we only need one set of dims
	// the layout engine will just add the amount of the previous image's size in the dynamic dimension
	// We need to add the anchor position to the textfield positions (which are relative to that anchor)

	for _, svgname := range svgFilenames {

		corner := geo.Point{X: 0, Y: 0} //default to layout anchor if not in the list - keeps layout drawing cleaner

		if thisAnchor, ok := layout.Anchors[svgname]; ok {
			corner = thisAnchor
		}

		svgfilename := fmt.Sprintf("%s.svg", layout.Filenames[svgname])

		// assume relative paths (absolute paths are not practial for production anyway)
		svgfilename = filepath.Join(filepath.Dir(svgLayoutPath), svgfilename)

		svgBytes, err := ioutil.ReadFile(svgfilename)
		if err != nil {
			return spread, errors.New(fmt.Sprintf("Entity %s: error opening svg file %s", svgname, svgfilename))
		}

		ladder, err := DefineLadderFromSVG(svgBytes)
		if err != nil {
			return spread, errors.New(fmt.Sprintf("Ladder %s: Error defining ladder from svg because %v", svgname, err))
		}

		if ladder == nil {
			continue //throw error?
		}
		spread.Ladders = append(spread.Ladders, *ladder)

		//append TextFields to the Textfield list
		for _, tf := range ladder.TextFields {

			//shift the text field and add it to the list
			//let engine take care of mangling name to suit page
			tf.Rect.Corner = TranslatePosition(corner, tf.Rect.Corner)
			spread.TextFields = append(spread.TextFields, tf)
		}
	}
	return spread, nil
}

func SwapTextFieldXCoords(spread *Spread) error {

	if spread == nil {
		return errors.New("nil pointer to spread")
	}

	width := (*spread).Dim.Width

	for idx, tf := range (*spread).TextFields {

		tf.Rect.Corner.X = width - tf.Rect.Corner.X

		(*spread).TextFields[idx] = tf

	}

	return nil

}

func GetTextFieldsByTopRight(svgLayoutPath, spreadName string) ([]TextField, error) {

	tf := []TextField{}

	spread, err := GetTextFieldSpread(svgLayoutPath, spreadName)

	if err != nil {
		return tf, err
	}

	err = SwapTextFieldXCoords(&spread)

	if err != nil {
		return tf, err
	}

	return spread.TextFields, nil

}

func SwitchTextFieldOrigin(spread *Spread, width, height float64) error {

	if spread == nil {
		return errors.New("nil pointer to spread")
	}

	for idx, tf := range (*spread).TextFields {

		tf.Rect.Corner.X = width - tf.Rect.Corner.X

		(*spread).TextFields[idx] = tf
	}

	return nil

}

func ScaleTextFieldGeometry(spread *Spread, scaleFactor float64) error {

	if spread == nil {
		return errors.New("nil pointer to spread")
	}

	for idx, tf := range (*spread).TextFields {
		err := scaleTextFieldUnits(&tf, scaleFactor)
		if err != nil {
			return err
		}
		(*spread).TextFields[idx] = tf
	}

	return nil

}

func GetImageBoxesForTextFields(svgLayoutPath, spreadName string, widthPx, heightPx int, vanilla bool) ([]optical.Box, error) {

	boxes := []optical.Box{}

	spread, err := GetTextFieldSpread(svgLayoutPath, spreadName)

	if err != nil {
		return boxes, err
	}

	// coords from top-right corner (known point)
	err = SwapTextFieldXCoords(&spread)

	if err != nil {
		return boxes, err
	}

	// convert to units of pixels
	// only width is dynamic, so heightPx and spread.Dim.Width give the scaling

	scaleFactor := float64(heightPx) / spread.Dim.Height

	err = ScaleTextFieldGeometry(&spread, scaleFactor)

	// convert origin to bottomLeft

	err = SwitchTextFieldOrigin(&spread, float64(widthPx), float64(heightPx))

	// now we have the coordinates in origin lower left, units of pixels
	// all that is left is to transfer these into the optical.Box struct

	for _, tf := range spread.TextFields {

		box := optical.Box{
			Vanilla: vanilla,
			ID:      tf.ID,
			Bounds:  geo.ConvertToImageRectangle(tf.Rect),
		}

		boxes = append(boxes, box)

	}

	return boxes, nil

}
