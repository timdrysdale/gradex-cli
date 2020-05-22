package parsesvg

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/timdrysdale/gradex-cli/comment"
	"github.com/timdrysdale/gradex-cli/geo"
	"github.com/timdrysdale/gradex-cli/pagedata"
	"github.com/timdrysdale/unipdf/v3/annotator"
	"github.com/timdrysdale/unipdf/v3/creator"
	"github.com/timdrysdale/unipdf/v3/model"
	"github.com/timdrysdale/unipdf/v3/model/optimize"
)

func RenderSpread(svgLayoutPath string, spreadName string, previousImagePath string, pageNumber int, pdfOutputPath string) error {

	contents := SpreadContents{
		SvgLayoutPath:     svgLayoutPath,
		SpreadName:        spreadName,
		PreviousImagePath: previousImagePath,
		PageNumber:        pageNumber,
		PdfOutputPath:     pdfOutputPath,
	}
	return RenderSpreadExtra(contents)

}

func RenderSpreadExtra(contents SpreadContents) error {

	svgLayoutPath := contents.SvgLayoutPath
	spreadName := contents.SpreadName
	previousImagePath := contents.PreviousImagePath
	prefillImagePaths := contents.PrefillImagePaths
	comments := contents.Comments
	pageNumber := contents.PageNumber
	pdfOutputPath := contents.PdfOutputPath

	svgBytes, err := ioutil.ReadFile(svgLayoutPath)

	if err != nil {
		return errors.New(fmt.Sprintf("Error opening layout file %s: %v\n", svgLayoutPath, err))
	}

	layout, err := DefineLayoutFromSVG(svgBytes)
	if err != nil {
		return errors.New(fmt.Sprintf("Error obtaining layout from svg %s\n", svgLayoutPath))
	}

	spread := Spread{}

	spread.Name = spreadName

	foundPage := false
	for k, v := range layout.PageDims {
		if strings.Contains(k, spread.Name) {
			spread.Dim = v
			foundPage = true
		}
	}

	if !foundPage {
		return errors.New(fmt.Sprintf("No page size info for spread %s\n", spread.Name))
	}

	// find svg & img elements for this name
	var svgFilenames, imgFilenames []string

	for k, _ := range layout.Filenames {
		if strings.Contains(k, spread.Name) {

			// assume jpg- or no prefix is image; svg- is ladder (image plus acroforms)
			if strings.HasPrefix(k, geo.SVGElement) {
				svgFilenames = append(svgFilenames, k) //we'll get the contents later
			} else {

				imgFilenames = append(imgFilenames, k)
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
		imgfilename := fmt.Sprintf("%s.jpg", layout.Filenames[svgname]) //TODO check again library is jpg-only?

		if contents.TemplatePathsRelative {
			svgfilename = filepath.Join(filepath.Dir(svgLayoutPath), svgfilename)
			imgfilename = filepath.Join(filepath.Dir(svgLayoutPath), imgfilename)
		}

		svgBytes, err := ioutil.ReadFile(svgfilename)
		if err != nil {
			return errors.New(fmt.Sprintf("Entity %s: error opening svg file %s", svgname, svgfilename))
		}

		ladder, err := DefineLadderFromSVG(svgBytes)
		if err != nil {
			return errors.New(fmt.Sprintf("Ladder %s: Error defining ladder from svg because %v", svgname, err))
		}

		if ladder == nil {
			continue //throw error?
		}
		spread.Ladders = append(spread.Ladders, *ladder)

		// append chrome image to the images list
		image := ImageInsert{
			Filename: imgfilename,
			Corner:   corner,
			Dim:      ladder.Dim,
		}

		spread.Images = append(spread.Images, image) //add chrome to list of images to include

		//append TextFields to the Textfield list
		for _, tf := range ladder.TextFields {

			//shift the text field and add it to the list
			//let engine take care of mangling name to suit page
			tf.Rect.Corner = TranslatePosition(corner, tf.Rect.Corner)
			spread.TextFields = append(spread.TextFields, tf)
		}
		//append TextPrefills to the TextPrefill list
		for _, tp := range ladder.TextPrefills {

			//shift the text field and add it to the list
			//let engine take care of mangling name to suit page
			tp.Rect.Corner = TranslatePosition(corner, tp.Rect.Corner)
			spread.TextPrefills = append(spread.TextPrefills, tp)
		}
		//append ComboBoxes to the ComboBox list
		for _, cb := range ladder.ComboBoxes {

			//shift the text field and add it to the list
			//let engine take care of mangling name to suit page
			cb.Rect.Corner = TranslatePosition(corner, cb.Rect.Corner)
			spread.ComboBoxes = append(spread.ComboBoxes, cb)
		}

	}

	// get all the static images that decorate this page, but not the special script "previous-image"

	//fmt.Println(prefillImagePaths)
	//fmt.Println(imgFilenames)
	for _, imgname := range imgFilenames {

		if _, ok := layout.ImageDims[imgname]; !ok {
			return errors.New(fmt.Sprintf("No size for image %s (must be provided in layout - check you have a correctly named box on the images layer in Inkscape)\n", imgname))
		}

		imgfilename := imgname //in case not specified, e.g. previous image

		if filename, ok := layout.Filenames[imgname]; ok {
			imgfilename = fmt.Sprintf("%s.jpg", filename)
		}

		if contents.TemplatePathsRelative {
			imgfilename = filepath.Join(filepath.Dir(svgLayoutPath), imgfilename)
		}

		// overwrite filename with dynamically supplied one, if supplied
		if filename, ok := prefillImagePaths[imgname]; ok {

			imgfilename = fmt.Sprintf("%s.jpg", filename)
		}

		if contents.PrefillImagePathsRelative {
			imgfilename = filepath.Join(filepath.Dir(svgLayoutPath), imgfilename)
		}

		corner := layout.Anchor

		if thisAnchor, ok := layout.Anchors[imgname]; ok {
			corner = thisAnchor
		}

		// append chrome image to the images list
		image := ImageInsert{
			Filename: imgfilename,
			Corner:   corner,
			Dim:      layout.ImageDims[imgname],
		}

		spread.Images = append(spread.Images, image) //add chrome to list of images to include
	}

	// Obtain the special "previous-image" which is flattened/rendered to image version of this page at the last step

	previousImageAnchorName := fmt.Sprintf("img-previous-%s", spread.Name)

	previousImageDimName := fmt.Sprintf("previous-%s", spread.Name)

	corner := layout.Anchors[previousImageAnchorName] //DiffPosition(layout.Anchors[previousImageAnchorName], layout.Anchor)

	previousImage := ImageInsert{
		Filename: previousImagePath,
		Corner:   corner,
		Dim:      layout.ImageDims[previousImageDimName],
	}

	// We do NOT add the previousImage to spread.Images because we treat it differently

	// We do things in a funny order here so that we can load the previous-image
	// and set the dynamic page size if needed

	c := creator.New()
	c.SetPageMargins(0, 0, 0, 0) // we're not printing so use the whole page
	var page *model.PdfPage
	if strings.Compare(previousImage.Filename, "") != 0 {

		img, err := c.NewImageFromFile(previousImage.Filename)

		if err != nil {
			return errors.New(fmt.Sprintf("Error opening spread %s previous-image file %s: %v", spread.Name, previousImage.Filename, err))
		}

		// Now we do the scaling to fit the page - see timdrysdale/pagescale for a demo
		if spread.Dim.DynamicWidth {
			img.ScaleToHeight(spread.Dim.Height)
			spread.ExtraWidth = img.Width() //we'll increase the page size by the image size
		} else {
			imgScaledWidth := img.Width() * previousImage.Dim.Height / img.Height()

			if imgScaledWidth > previousImage.Dim.Width {
				// oops, we're too big, so scale using width instead
				img.ScaleToWidth(previousImage.Dim.Width)
			} else {
				img.ScaleToHeight(previousImage.Dim.Height)
			}

		}
		img.SetPos(previousImage.Corner.X, previousImage.Corner.Y)
		// we use GetWidth() so value includes fixed width plus extra width
		c.SetPageSize(creator.PageSize{spread.GetWidth(), spread.Dim.Height})

		page = c.NewPage()

		c.Draw(img) //draw previous image
	} else {
		c.SetPageSize(creator.PageSize{spread.GetWidth(), spread.Dim.Height})

		page = c.NewPage()

	}

	// pagedata used to go in here

	for _, v := range spread.Images {
		img, err := c.NewImageFromFile(v.Filename)

		if err != nil {
			return errors.New(fmt.Sprintf("Error opening image file %s: %s", v.Filename, err))
		}
		// all these images are static so we set dims directly
		// user needs to spot if they did their artwork to the wrong spec
		// or maybe they want it that way - we'd never know...
		// TODO consider logging a warning here for GUI etc
		img.SetWidth(v.Dim.Width)
		img.SetHeight(v.Dim.Height)
		if spread.Dim.DynamicWidth {
			img.SetPos(v.Corner.X+spread.ExtraWidth, v.Corner.Y)
		} else {
			img.SetPos(v.Corner.X, v.Corner.Y) //TODO check this has correct sense for non-zero offsets
		}
		c.Draw(img)
	}

	newComments := []comment.Comment{}

	// expect our calling function to have pre-loaded any old comments
	// into the pagedata.Comment array, so we know to print our new
	// comments above them.
	numOldComments := len(contents.PageData.Current.Comments)
	numNewComments := float64(len(comments.GetByPage(pageNumber)))
	numTotalComments := numOldComments + numNewComments

	// Draw in our flattened comments
	rowHeight := 12.0
	x := 0.3 * rowHeight
	y := c.Height() - ((0.3 + numTotalComments) * rowHeight)
	y = y + numOldComments*rowHeight

	// figure out who edited last, and hence made any new comments
	numOldPageDatas := len(contents.PageData.Previous)
	lastEditor := "" //just show a number if editor not named
	if numOldPageDatas > 0 {
		previousPageData := contents.PageData.Previous[numOldPageDatas-1]
		lastEditor = "-" + limit(previousPageData.Process.For, 3)
	}

	for i, cmt := range comments.GetByPage(pageNumber) {
		cmt.Label = fmt.Sprintf("%d%s", lastEditor, i+numOldComments)
		comment.DrawComment(c, cmt, x, y)
		y = y + rowHeight
		newComments = append(newComments, cmt)
	}

	// add these comments and labels to the page data
	contents.PageData.Current.Comments = append(content.PageData.Current.Comments, newComments)

	if !reflect.DeepEqual(contents.PageData, pagedata.PageData{}) {
		pagedata.MarshalOneToCreator(c, &contents.PageData)
	}

	for _, tp := range spread.TextPrefills {
		//update prefill contents from info given
		if val, ok := contents.Prefills[pageNumber][tp.ID]; ok {
			tp.Text.Text = val
		}
		// update our prefill text
		p := c.NewParagraph(tp.Text.Text)

		p.SetFontSize(tp.Text.TextSize)

		p.SetPos(tp.Rect.Corner.X, tp.Rect.Corner.Y)

		c.Draw(p)

	}

	/*******************************************************************************
	  Note that multipage acroforms are a wriggly issue!
	  This code is intended for single-page demos - check gradex-overlay for the
	  multipage method
	  ******************************************************************************/
	form := model.NewPdfAcroForm()

	for _, tf := range spread.TextFields {

		tfopt := annotator.TextFieldOptions{Value: tf.Prefill} //TODO - MaxLen?!
		// TODO consider allowing a more templated mangling of the ID number
		// For multi-student entries (although, OTH, there will be per-page ID data etc embedded too
		// which may be more useful in this regard, rather than overloading the textfield id)
		name := fmt.Sprintf("page-%03d-%s", pageNumber+1, tf.ID) //match physical page number

		if spread.Dim.DynamicWidth {
			tf.Rect.Corner.X = tf.Rect.Corner.X + spread.ExtraWidth
		}

		textf, err := annotator.NewTextField(page, name, formRect(tf.Rect, layout.Dim), tfopt)
		if err != nil {
			panic(err)
		}
		*form.Fields = append(*form.Fields, textf.PdfField)
		page.AddAnnotation(textf.Annotations[0].PdfAnnotation)
	}

	for _, cb := range spread.ComboBoxes { //also a forms object
		//update combobox options (possible values) from info given
		if val, ok := contents.ComboBoxes[pageNumber][cb.ID]; ok {
			cb.Options.Options = val.Options
		}
		name := fmt.Sprintf("page-%03d-%s", pageNumber+1, cb.ID) //match physical page number
		opt := annotator.ComboboxFieldOptions{Choices: cb.Options.Options}
		comboboxf, err := annotator.NewComboboxField(page, name, formRect(cb.Rect, layout.Dim), opt)
		if err != nil {
			panic(err)
		}

		*form.Fields = append(*form.Fields, comboboxf.PdfField)
		page.AddAnnotation(comboboxf.Annotations[0].PdfAnnotation)
	}

	err = c.SetForms(form)
	if err != nil {
		return errors.New(fmt.Sprintf("Error: %v\n", err))
	}

	c.SetOptimizer(optimize.New(optimize.Options{
		CombineDuplicateDirectObjects:   true,
		CombineIdenticalIndirectObjects: true,
		CombineDuplicateStreams:         true,
		CompressStreams:                 true,
		UseObjectStreams:                true,
		ImageQuality:                    90,
		ImageUpperPPI:                   150,
	}))

	c.WriteToFile(pdfOutputPath)
	return nil
}

func limit(initials string, N int) string {
	if len(initials) < 3 {
		N = len(initials)
	}
	return strings.ToUpper(initials[0:N])
}
