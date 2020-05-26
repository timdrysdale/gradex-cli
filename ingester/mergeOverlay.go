package ingester

import (
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/rs/zerolog"
	"github.com/timdrysdale/gradex-cli/merge"
	"github.com/timdrysdale/gradex-cli/pagedata"
	"github.com/timdrysdale/gradex-cli/parsesvg"
	"github.com/timdrysdale/pool"
)

type MergeFile struct {
	OutputPath string // should be basename with .pdf ext
	InputPages []Page // should be absolute paths to pdfs (we'll get pagedata, then work out the image filename)
}

type MergeCommand struct {
	MergeFiles    []MergeFile
	ToDir         string
	Template      string
	SpreadName    string
	ProcessDetail pagedata.ProcessDetail
}

type MergeTask struct {
	MergeFile     MergeFile
	ToDir         string
	ProcessDetail pagedata.ProcessDetail
	SpreadName    string
	Template      string
}

// we pass pointer to logger that has a processing stage string pre-prended to it
// so we can tell what stage overlay is being used at
func (g *Ingester) MergeOverlayPapers(mc MergeCommand, logger *zerolog.Logger) error {

	// create an array of tasks

	mergeTasks := []MergeTask{}

	for _, mergeFile := range mc.MergeFiles {

		// no clear way to use the done files here, because what is "done" could change depending on
		// small changes to the marking status of other files
		// for simplicity, and least surprise, accept that this process will always run anew
		// on ALL files when called - so responsibility for any efficiency savings lies
		// with the calling function to prune entries from the MergeMap that are already done

		mergeTasks = append(mergeTasks, MergeTask{
			MergeFile:     mergeFile,
			ToDir:         mc.ToDir,
			SpreadName:    mc.SpreadName,
			Template:      mc.Template,
			ProcessDetail: mc.ProcessDetail,
		})
		logger.Info().
			Str("output", mergeFile.OutputPath).
			Int("page-count", len(mergeFile.InputPages)).
			Msg("Preparing to do merge overlay")

	}

	// now process the files

	N := len(mergeTasks)

	tasks := []*pool.Task{}

	for i := 0; i < N; i++ {

		mt := mergeTasks[i]

		newtask := pool.NewTask(func() error {
			pc, err := g.MergeOverlayOnePDF(mt, logger)
			if err == nil {
				logger.Info().
					Str("file", mt.MergeFile.OutputPath).
					Int("page-count", pc).
					Msg("Done Merge Overlay")
				return nil
			} else {
				logger.Error().
					Str("file", mt.MergeFile.OutputPath).
					Int("page-count", pc).
					Str("error", err.Error()).
					Msg("Error with Merge Overlay")
				return err
			}
		})
		tasks = append(tasks, newtask)
	}

	p := pool.NewPool(tasks, runtime.GOMAXPROCS(-1))
	logger.Info().
		Int("speed-up", runtime.GOMAXPROCS(-1)).
		Msg(fmt.Sprintf("Using parallel processing to get x%d speed-up\n", runtime.GOMAXPROCS(-1)))

	p.Run()

	var numErrors int
	for _, task := range p.Tasks {
		if task.Err != nil {
			logger.Error().
				Str("error", task.Err.Error()).
				Msg(fmt.Sprintf("Processing problem %v", task.Err))
			numErrors++
		}
	}

	// report how we did
	if numErrors > 0 {
		logger.Error().
			Int("error-count", numErrors).
			Int("script-count", N).
			Msg(fmt.Sprintf("Processing finished with merge-overlay tasks returning <%d> errors from <%d> scripts\n", numErrors, N))

	} else {
		logger.Info().
			Int("error-count", numErrors).
			Int("script-count", N).
			Msg(fmt.Sprintf("Processing finished <%d> scripts without any errors\n", N))

	}
	return nil
}

//---------------- MergeOverlayOnePDF -----------------------------------------------------
//
//
//-----------------------------------------------------------------------------------------
func (g *Ingester) MergeOverlayOnePDF(mt MergeTask, logger *zerolog.Logger) (int, error) {

	mergePaths := []string{}

	for idx, page := range mt.MergeFile.InputPages {

		// get pagedata for this page, which will be on page 1

		inPath := page.Path

		pageDataMap, err := pagedata.UnMarshalAllFromFile(inPath)

		if err != nil {
			logger.Error().
				Str("file", inPath).
				Str("error", err.Error()).
				Msg(fmt.Sprintf("%s: error obtaining pagedata\n", inPath))
			return 0, err

		}

		if pagedata.GetLen(pageDataMap) < 1 {
			logger.Error().
				Str("file", inPath).
				Msg(fmt.Sprintf("%s: no pagedata in file\n", inPath))
			return 0, err
		}

		if _, ok := pageDataMap[1]; !ok {
			logger.Error().
				Str("file", inPath).
				Msg(fmt.Sprintf("%s: no pagedata for page 1\n", inPath))
			return 0, err
		}

		thisPageData := pageDataMap[1]

		//-------------- Update PageData --------------

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

		newThisPageDataCurrent.Process = mt.ProcessDetail

		thisPageData.Current = newThisPageDataCurrent

		//---------------------------------------WORK OUT COURSE CODE--------------------------
		// find coursecode in these pages - assume belong SAME assignment. Not true of future bu-page mode

		courseCode := thisPageData.Current.Item.What

		if courseCode == "" {
			logger.Error().
				Str("file", inPath).
				Msg(fmt.Sprintf("Can't figure out the course code for file (%s) - not present in PageData?\n", inPath))
			return 0, errors.New("Couldn't find a course code")
		}

		// flatten the page - seems wasteful, but we can't get the flattened comments unless we do that
		// the previous image for this flattened page does not contain the comments, because they are added as text
		//---------------------------------------------------------------------------------------------------

		jpegDir := g.PaperImages(courseCode)
		basename := strings.TrimSuffix(filepath.Base(inPath), filepath.Ext(inPath))
		jpegPathOption := fmt.Sprintf("%s/%s-me%%d.jpg", jpegDir, basename)

		//should just be one page, but keep index just in case (but no zero padding)
		//so that we can see the problem if we look in the temp-images dir

		/*no new comments expected! so leave this commented out for now
		f, err := os.Open(inPath)
		if err != nil {
			logger.Error().
				Str("file", inPath).
				Str("error", err.Error()).
				Msg(fmt.Sprintf("Can't open file (%s) because: %v\n", inPath, err))
			ot.Msg.Send(fmt.Sprintf("Can't open file (%s) because: %v\n", inPath, err))
			return 0, err
		}

		pdfReader, err := pdf.NewPdfReader(f)
		if err != nil {
			logger.Error().
				Str("file", inPath).
				Str("error", err.Error()).
				Msg(fmt.Sprintf("Can't read from file (%s) because: %v\n", inPath, err))
			ot.Msg.Send(fmt.Sprintf("Can't read from file (%s) because: %v\n", inPath, err))
			return 0, err
		}

		comments, err := comment.GetComments(pdfReader)

		f.Close()
		*/

		err = ConvertPDFToJPEGs(inPath, jpegDir, jpegPathOption)
		if err != nil {
			logger.Error().
				Str("file", inPath).
				Str("error", err.Error()).
				Msg(fmt.Sprintf("Can't flatten file (%s) to images because: %v\n", inPath, err))
			return 0, err
		}

		// convert images to individual pdfs, with form overlay

		pageDir := g.PaperPages(courseCode)
		pagePathOption := fmt.Sprintf("%s/%s-me%%d.pdf", pageDir, basename)

		pagePath := fmt.Sprintf(pagePathOption, 1) //this is the merge page index -me1, we get the 000n index from the Own.Path
		jpegPath := fmt.Sprintf(jpegPathOption, 1) //first page starts at 1
		//----------------------------------------------------------------------------------------------------

		prefills := parsesvg.DocPrefills{}

		prefills[idx] = make(map[string]string)

		prefills[idx]["message"] = page.Message

		prefills[idx]["page-number"] = fmt.Sprintf("%d/%d", idx+1, len(mt.MergeFile.InputPages)) //add one we so we get a display pagenumber that starts at one

		prefills[idx]["author"] = thisPageData.Current.Item.Who

		prefills[idx]["date"] = thisPageData.Current.Item.When

		prefills[idx]["title"] = thisPageData.Current.Item.What

		prefills[idx]["for"] = thisPageData.Current.Process.For

		contents := parsesvg.SpreadContents{
			SvgLayoutPath:         mt.Template,
			SpreadName:            mt.SpreadName,
			PreviousImagePath:     jpegPath,
			PageNumber:            idx,
			PdfOutputPath:         pagePath,
			PageData:              thisPageData, //no pageNumber index needed
			TemplatePathsRelative: true,
			Prefills:              prefills,
		}

		err = parsesvg.RenderSpreadExtra(contents)
		if err != nil {
			logger.Error().
				Str("file", pagePath).
				Int("page-number", idx+1).
				Str("error", err.Error()).
				Msg(fmt.Sprintf("Error rendering spread for page <%d> of (%s) because %v\n", idx+1, mt.MergeFile.OutputPath, err))

			return 0, err

		}

		mergePaths = append(mergePaths, pagePath)

	} //for

	outputPath := filepath.Join(mt.ToDir, mt.MergeFile.OutputPath)

	err := merge.PDF(mergePaths, outputPath)
	if err != nil {
		logger.Error().
			Str("file", outputPath).
			Str("error", err.Error()).
			Msg(fmt.Sprintf("Error merging processed pages for (%s) because %v\n", mt.MergeFile.OutputPath, err))
		return 0, err
	}

	logger.Info().
		Str("file", outputPath).
		Int("page-count", len(mergePaths)).
		Str("spread-name", mt.SpreadName).
		Msg(fmt.Sprintf("Finished rendering merge-overlay for (%s) which had <%d> pages\n", mt.MergeFile.OutputPath, len(mergePaths)))

	return len(mergePaths), nil

} //func
