package ingester

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/rs/zerolog"
	"github.com/timdrysdale/gradex-cli/comment"
	"github.com/timdrysdale/gradex-cli/extract"
	"github.com/timdrysdale/gradex-cli/merge"
	"github.com/timdrysdale/gradex-cli/pagedata"
	"github.com/timdrysdale/gradex-cli/parsesvg"
	"github.com/timdrysdale/pool"
	pdf "github.com/timdrysdale/unipdf/v3/model"
)

// This function places content - but needs a wrapper to generate that content
// appropriately for each step
// It will update the values of some fields in PageData with the appropriate
// dynamically generated data
// Read page data
// flatten file
// start new doc
// import old image
// add overlay form
// update page data with this process' info
// write pagedata
// write doc
// tell someone about it!

func (g *Ingester) OutputPath(dir, inPath, decoration string) string {

	ext := filepath.Ext(inPath)
	base := strings.TrimSuffix(filepath.Base(inPath), ext)

	return filepath.Join(dir, base+decoration+ext)

}

// we pass pointer to logger that has a processing stage string pre-prended to it
// so we can tell what stage overlay is being used at
func (g *Ingester) OverlayPapers(oc OverlayCommand, logger *zerolog.Logger) error {

	// assume someone hits a button to ask us to do this ...
	// we'll operate on the directory that is associated with
	// this task, for this exam

	overlayTasks := []OverlayTask{}

	inPaths, err := g.GetFileList(oc.FromPath)
	if err != nil {
		oc.Msg.Send(fmt.Sprintf("Stopping early; couldn't get files because %v\n", err))
		logger.Error().
			Str("source-dir", oc.FromPath).
			Str("error", err.Error()).
			Msg(fmt.Sprintf("Stopping early; couldn't get files because %v\n", err))
		return err
	}

	for _, inPath := range inPaths {

		if !g.IsPDF(inPath) { //ignore the done files
			continue
		}

		if getDone(inPath) {
			logger.Info().
				Str("file", inPath).
				Msg("Skipping because already done")
			continue
		}

		count, err := CountPages(inPath)

		if err != nil {
			oc.Msg.Send(fmt.Sprintf("Skipping (%s): error counting pages because %v\n", inPath, err))
			logger.Error().
				Str("file", inPath).
				Str("error", err.Error()).
				Msg(fmt.Sprintf("Skipping (%s): error counting pages because %v\n", inPath, err))
			continue
		}

		pageDataMap, err := pagedata.UnMarshalAllFromFile(inPath)

		if err != nil {
			logger.Error().
				Str("file", inPath).
				Str("error", err.Error()).
				Msg(fmt.Sprintf("Skipping (%s): error obtaining pagedata\n", inPath))
			oc.Msg.Send(fmt.Sprintf("Skipping (%s): error obtaining pagedata\n", inPath))
			continue
		}

		if pagedata.GetLen(pageDataMap) < 1 {
			oc.Msg.Send(fmt.Sprintf("Skipping (%s): no pagedata in file\n", inPath))
			continue
		}

		// get all textfields from the file

		fieldsMapByPage, err := extract.ExtractTextFieldsFromPDF(inPath)

		// TODO Capture optical check box values here - Fairly expensive process,
		// so trigger if there are no marks in the whole doc?
		// anyone marking half and half we will just end up flagging due to "surprise" missing data
		logger.Warn().Msg("NOT CAPTURING OPTICAL CHECK BOX DATA AT THIS TIME!!")

		if err == nil {

			for page, fields := range fieldsMapByPage {
				//https: //stackoverflow.com/questions/17438253/accessing-struct-fields-inside-a-map-value-without-copying
				pd := pageDataMap[page]

				var data []pagedata.Field

				data = pd.Current.Data

				for key, value := range fields {

					data =
						append(data,
							pagedata.Field{
								Key:   key,
								Value: value,
							})
				}
				pd.Current.Data = data
				pageDataMap[page] = pd
			}

		} else {

			logger.Error().
				Str("file", inPath).
				Str("error", err.Error()).
				Msg("Expected text fields but couldn't get them")

		}

		// this is a file-level task, so we we will sort per-page updates
		// to pageData at the child step
		overlayTasks = append(overlayTasks, OverlayTask{
			InputPath: inPath,
			PageCount: count,
			// These now in SharedPD
			//PreparedFor:   oc.PreparedFor,
			//ToDo:          oc.ToDo,
			//NewProcessing: oc.ProcessingDetails, //do dynamic update when processing
			//NewQuestion:   oc.QuestionDetails,   //do dynamic update when processing
			ProcessDetail:  oc.ProcessDetail,
			OldPageDataMap: pageDataMap,
			OutputPath:     g.OutputPath(oc.ToPath, inPath, oc.PathDecoration),
			SpreadName:     oc.SpreadName,
			Template:       oc.TemplatePath,
			Msg:            oc.Msg,
		})
		logger.Info().
			Str("file", inPath).
			Int("page-count", count).
			Int("page-data-count", pagedata.GetLen(pageDataMap)).
			Msg(fmt.Sprintf("Preparing to process: file (%s) has <%d> pages and [%d] pageDatas\n", inPath, count, pagedata.GetLen(pageDataMap)))
		oc.Msg.Send(fmt.Sprintf("Preparing to process: file (%s) has <%d> pages and [%d] pageDatas\n", inPath, count, pagedata.GetLen(pageDataMap)))

	} // for loop through all files

	// now process the files
	N := len(overlayTasks)

	//pcChan := make(chan int, N)

	tasks := []*pool.Task{}

	for i := 0; i < N; i++ {

		ot := overlayTasks[i]

		newtask := pool.NewTask(func() error {
			pc, err := g.OverlayOnePDF(ot, logger)
			if err == nil {
				setDone(ot.InputPath, logger)
				logger.Info().
					Str("file", ot.InputPath).
					Str("destination", ot.OutputPath).
					Int("page-count", pc).
					Msg(fmt.Sprintf("Finished processing (%s) into (%s)which had <%d> pages", ot.InputPath, ot.OutputPath, pc))

				oc.Msg.Send(fmt.Sprintf("Finished processing (%s) into (%s)which had <%d> pages", ot.InputPath, ot.OutputPath, pc))
				oc.Msg.Send(fmt.Sprintf("pages(%d)", pc))
				return nil
			} else {
				logger.Error().
					Str("file", ot.InputPath).
					Str("destination", ot.OutputPath).
					Int("page-count", pc).
					Str("error", err.Error()).
					Msg(fmt.Sprintf("Error processing (%s) into (%s)which had <%d> pages", ot.InputPath, ot.OutputPath, pc))
				return err
			}
		})
		tasks = append(tasks, newtask)
	}

	p := pool.NewPool(tasks, runtime.GOMAXPROCS(-1))
	logger.Info().
		Int("speed-up", runtime.GOMAXPROCS(-1)).
		Msg(fmt.Sprintf("Using parallel processing to get x%d speed-up\n", runtime.GOMAXPROCS(-1)))

	oc.Msg.Send(fmt.Sprintf("Using parallel processing to get x%d speed-up\n", runtime.GOMAXPROCS(-1)))

	closed := make(chan struct{})

	p.Run()

	var numErrors int
	for _, task := range p.Tasks {
		if task.Err != nil {
			logger.Error().
				Str("error", task.Err.Error()).
				Msg(fmt.Sprintf("Processing problem %v", task.Err))
			oc.Msg.Send(fmt.Sprintf("Processing error: %v", task.Err))
			numErrors++
		}
	}
	close(closed)

	// report how we did
	if numErrors > 0 {
		logger.Error().
			Int("error-count", numErrors).
			Int("script-count", N).
			Msg(fmt.Sprintf("Processing finished with overlay tasks returning <%d> errors from <%d> scripts\n", numErrors, N))
		oc.Msg.Send(fmt.Sprintf("Processing finished with overlay tasks returning <%d> errors from <%d> scripts\n", numErrors, N))
	} else {
		logger.Info().
			Int("error-count", numErrors).
			Int("script-count", N).
			Msg(fmt.Sprintf("Processing finished <%d> scripts without any errors\n", N))
		oc.Msg.Send(fmt.Sprintf("Processing finished, completed <%d> scripts\n", N))
	}
	return nil
}

// do one file, dynamically assembling the data we need make the latest pagedata
// from what we get in the OverlayTask struct
// return the number of pages
func (g *Ingester) OverlayOnePDF(ot OverlayTask, logger *zerolog.Logger) (int, error) {

	// need page count to find the jpeg files again later
	numPages, err := CountPages(ot.InputPath)

	// find coursecode in these pages - assume belong SAME assignment. Not true of future bu-page mode
	var courseCode string
OUTER:
	for _, v := range ot.OldPageDataMap {
		if v.Current.Item.What != "" {
			courseCode = v.Current.Item.What
			break OUTER
		}
	}

	if courseCode == "" {
		logger.Error().
			Str("file", ot.InputPath).
			Msg(fmt.Sprintf("Can't figure out the course code for file (%s) - not present in PageData?\n", ot.InputPath))
		ot.Msg.Send(fmt.Sprintf("Can't figure out the course code for file (%s) - not present in PageData?\n", ot.InputPath))
		return 0, errors.New("Couldn't find a course code")
	}

	// render to images

	jpegPath := g.PaperImages(courseCode) //ot.PageDataMap[0].Exam.CourseCode)

	suffix := filepath.Ext(ot.InputPath)
	basename := strings.TrimSuffix(filepath.Base(ot.InputPath), suffix)
	jpegFileOption := fmt.Sprintf("%s/%s%%04d.jpg", jpegPath, basename)

	f, err := os.Open(ot.InputPath)
	if err != nil {
		logger.Error().
			Str("file", ot.InputPath).
			Str("error", err.Error()).
			Msg(fmt.Sprintf("Can't open file (%s) because: %v\n", ot.InputPath, err))
		ot.Msg.Send(fmt.Sprintf("Can't open file (%s) because: %v\n", ot.InputPath, err))
		return 0, err
	}

	pdfReader, err := pdf.NewPdfReader(f)
	if err != nil {
		logger.Error().
			Str("file", ot.InputPath).
			Str("error", err.Error()).
			Msg(fmt.Sprintf("Can't read from file (%s) because: %v\n", ot.InputPath, err))
		ot.Msg.Send(fmt.Sprintf("Can't read from file (%s) because: %v\n", ot.InputPath, err))
		return 0, err
	}

	comments, err := comment.GetComments(pdfReader)

	f.Close()

	err = ConvertPDFToJPEGs(ot.InputPath, jpegPath, jpegFileOption)
	if err != nil {
		logger.Error().
			Str("file", ot.InputPath).
			Str("error", err.Error()).
			Msg(fmt.Sprintf("Can't flatten file (%s) to images because: %v\n", ot.InputPath, err))
		ot.Msg.Send(fmt.Sprintf("Can't flatten file (%s) to images because: %v\n", ot.InputPath, err))
		return 0, err
	}

	// convert images to individual pdfs, with form overlay

	pagePath := g.PaperPages(courseCode)
	pageFileOption := fmt.Sprintf("%s/%s%%04d.pdf", pagePath, basename)

	mergePaths := []string{}

	// gs starts indexing at 1
	for imgIdx := 1; imgIdx <= numPages; imgIdx = imgIdx + 1 {

		// construct image name
		previousImagePath := fmt.Sprintf(jpegFileOption, imgIdx)
		pageFilename := fmt.Sprintf(pageFileOption, imgIdx)

		pageNumber := imgIdx - 1 //pageNumber starts at zero

		//if len(ot.PageDataMap[imgIdx]) < 1 {
		//	logger.Error().
		//		Str("file", ot.InputPath).
		//		Int("page-number", imgIdx).
		//		Msg(fmt.Sprintf("Info: no existing page data for file (%s) on page <%d>\n", ot.InputPath, imgIdx))
		//	ot.Msg.Send(fmt.Sprintf("Info: no existing page data for file (%s) on page <%d>\n", ot.InputPath, imgIdx))
		//}

		// FOR THIS PAGE INDIVIDUALLY
		// we take the "current" pagedata loaded with the extracted field data,
		// and tack it on the end of the list of previous PageDatas

		thisPageData, ok := ot.OldPageDataMap[imgIdx] //on book number, start at 1
		if !ok {
			logger.Error().
				Str("file", ot.InputPath).
				Msg("No pagedata in file")
		}

		oldThisPageDataCurrent := thisPageData.Current

		previousPageData := thisPageData.Previous

		previousPageData = append(previousPageData, oldThisPageDataCurrent)

		thisPageData.Previous = previousPageData

		// Now we CONSTRUCT the NEW current PageData
		// copy over what we had, first:

		newThisPageDataCurrent := oldThisPageDataCurrent

		// now add in things we can only know now
		// like page number, UUID etc.

		newThisPageDataCurrent.UUID = safeUUID()

		newThisPageDataCurrent.Follows = oldThisPageDataCurrent.UUID

		newThisPageDataCurrent.Process = ot.ProcessDetail

		// TODO - do we need to update Own? Host?

		thisPageData.Current = newThisPageDataCurrent

		headerPrefills := parsesvg.DocPrefills{}

		headerPrefills[pageNumber] = make(map[string]string)

		headerPrefills[pageNumber]["page-number"] = fmt.Sprintf("%d/%d", pageNumber+1, numPages) //add one we so we get a display pagenumber that starts at one

		headerPrefills[pageNumber]["author"] = thisPageData.Current.Item.Who

		headerPrefills[pageNumber]["date"] = thisPageData.Current.Item.When

		headerPrefills[pageNumber]["title"] = thisPageData.Current.Item.What

		contents := parsesvg.SpreadContents{
			SvgLayoutPath:         ot.Template,
			SpreadName:            ot.SpreadName,
			PreviousImagePath:     previousImagePath,
			PageNumber:            pageNumber,
			PdfOutputPath:         pageFilename,
			Comments:              comments,
			PageData:              thisPageData,
			TemplatePathsRelative: true,
			Prefills:              headerPrefills,
		}

		err = parsesvg.RenderSpreadExtra(contents)
		if err != nil {
			logger.Error().
				Str("file", ot.InputPath).
				Int("page-number", imgIdx).
				Str("error", err.Error()).
				Msg(fmt.Sprintf("Error rendering spread for page <%d> of (%s) because %v\n", imgIdx, ot.InputPath, err))

			ot.Msg.Send(fmt.Sprintf("Error rendering spread for page <%d> of (%s) because %v\n", imgIdx, ot.InputPath, err))
			return 0, err

		}

		mergePaths = append(mergePaths, pageFilename)
	}
	err = merge.PDF(mergePaths, ot.OutputPath)
	if err != nil {
		logger.Error().
			Str("file", ot.InputPath).
			Str("error", err.Error()).
			Msg(fmt.Sprintf("Error merging processed pages for (%s) because %v\n", ot.InputPath, err))

		ot.Msg.Send(fmt.Sprintf("Error merging processed pages for (%s) because %v\n", ot.InputPath, err))
		return 0, err
	}

	doneFile := doneFilePath(ot.OutputPath)
	_, err = os.Stat(doneFile)
	if err == nil {
		err = os.Remove(doneFile)
		if err != nil {
			logger.Error().
				Str("file", ot.OutputPath).
				Str("error", err.Error()).
				Msg("Could not delete stale Done File")
		}
	}

	logger.Info().
		Str("file", ot.InputPath).
		Int("page-count", ot.PageCount).
		Str("spread-name", ot.SpreadName).
		Msg(fmt.Sprintf("Finished rendering [%s] overlay for (%s) which had <%d> pages\n", ot.SpreadName, ot.InputPath, ot.PageCount))

	ot.Msg.Send(fmt.Sprintf("Finished rendering [%s] overlay for (%s) which had <%d> pages\n", ot.SpreadName, ot.InputPath, ot.PageCount))
	return numPages, nil

}
