package ingester

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/timdrysdale/anon"
	"github.com/timdrysdale/gradex-cli/comment"
	"github.com/timdrysdale/gradex-cli/merge"
	"github.com/timdrysdale/gradex-cli/pagedata"
	"github.com/timdrysdale/gradex-cli/parselearn"
	"github.com/timdrysdale/gradex-cli/parsesvg"
	"github.com/timdrysdale/pool"
	pdf "github.com/timdrysdale/unipdf/v3/model"
)

func (g *Ingester) FlattenNewPapers(exam string) error {

	logger := g.logger.With().Str("process", "flatten").Logger()

	//assume someone hits a button to ask us to do this ...

	// load our identity database
	identity, err := anon.New(g.IdentityCSV())
	if err != nil {
		logger.Error().
			Str("file", g.IdentityCSV()).
			Str("course", exam).
			Str("error", err.Error()).
			Msg("Cannot open identity.csv")
		return err
	}

	flattenTasks := []FlattenTask{}

	receipts, err := g.GetFileList(g.AcceptedReceipts(exam))
	if err != nil {
		logger.Error().
			Str("file", g.AcceptedReceipts(exam)).
			Str("course", exam).
			Str("error", err.Error()).
			Msg("Cannot get list of accepted Receipts")
		return errors.New("Can't find this exam - please check spelling or ingest your first papers")
	}

	for _, receipt := range receipts {

		sub, err := parselearn.ParseLearnReceipt(receipt)
		if err != nil {
			logger.Error().
				Str("file", receipt).
				Str("course", exam).
				Str("error", err.Error()).
				Msg(fmt.Sprintf("Cannot parse receipt %s because %v", receipt, err))
			continue
		}

		pdfPath, err := GetPDFPath(sub.Filename, g.AcceptedPapers(exam))

		if err != nil {
			logger.Error().
				Str("file", sub.Filename).
				Str("dir", g.AcceptedPapers(exam)).
				Str("course", exam).
				Str("error", err.Error()).
				Msg(fmt.Sprintf("couldn't get PDF filename for %s because %v", sub.Filename, err))
			continue
		}

		if !g.Redo { //carry on regardless if forcing a redo
			if getDone(pdfPath) { //check for done file - don't process if it exists
				logger.Info().
					Str("file", pdfPath).
					Msg("Skipping flattening - already done")
				continue
			}
		}

		logger.Info().
			Str("file", pdfPath).
			Msg("Flattening new file")

		count, err := CountPages(pdfPath)

		if err != nil {
			logger.Error().
				Str("file", pdfPath).
				Str("course", exam).
				Str("error", err.Error()).
				Msg(fmt.Sprintf("couldn't countPages for %s because %v", pdfPath, err))
			continue
		}

		shortDate, err := GetShortLearnDate(sub)
		if err != nil {
			logger.Error().
				Str("file", sub.OwnPath).
				Str("course", exam).
				Str("error", err.Error()).
				Msg(fmt.Sprintf("couldn't get shortlearndate for %s because %v", receipt, err))
			continue
		}

		anonymousIdentity, err := identity.GetAnonymous(sub.Matriculation)
		if err != nil {
			logger.Error().
				Str("identity", sub.Matriculation).
				Str("course", exam).
				Str("error", err.Error()).
				Msg(fmt.Sprintf("couldn't get anonymous identity for for %s because %v\n", sub.Matriculation, err))
			continue
		}

		// we'll use this same set of procDetails for flattens that we do in this batch
		// that means we can use the uuid to map the processing in graphviz later, for example

		pd := pagedata.PageDetail{}
		pd.Is = pagedata.IsPage
		pd.Follows = "" //first proces
		pd.Revision = 0 //first revision

		pd.Own = pagedata.FileDetail{} //depends on stage, fill in when we flatten

		pd.Original = pagedata.FileDetail{
			Path: pdfPath,
			UUID: safeUUID(),
			Of:   count,
		}

		pd.Process = pagedata.ProcessDetail{
			UUID:     safeUUID(),
			UnixTime: time.Now().UnixNano(),
			Name:     "flatten",
			By:       "gradex-cli",
			For:      "ingester",
			ToDo:     "prepare-for-marking",
		}

		pd.Item = pagedata.ItemDetail{
			What:    sub.Assignment,
			When:    shortDate,
			Who:     anonymousIdentity,
			WhoType: pagedata.IsAnonymous,
		}

		pdataMap := make(map[int]pagedata.PageData)

		for page := 1; page <= count; page++ {

			thisPd := pd

			thisPd.UUID = safeUUID()

			thisPd.Original.Number = page

			thisPageData := pagedata.PageData{
				Current: thisPd,
			}

			pdataMap[page] = thisPageData

		}

		renamedBase := g.GetAnonymousFileName(sub.Assignment, anonymousIdentity)
		outputPath := filepath.Join(g.AnonymousPapers(sub.Assignment), renamedBase)

		flattenTasks = append(flattenTasks, FlattenTask{
			PreparedFor: "ingester",
			ToDo:        "flattening",
			InputPath:   pdfPath,
			OutputPath:  outputPath,
			PageCount:   count,
			PageDataMap: pdataMap})
	}

	// now process the files
	N := len(flattenTasks)

	pcChan := make(chan int, N)

	tasks := []*pool.Task{}

	for i := 0; i < N; i++ {

		inputPath := flattenTasks[i].InputPath
		outputPath := flattenTasks[i].OutputPath
		pdataMap := flattenTasks[i].PageDataMap

		newtask := pool.NewTask(func() error {
			pc, err := g.FlattenOnePDF(inputPath, outputPath, pdataMap, &logger)
			pcChan <- pc
			if err == nil {
				setDone(inputPath, &logger) // so we don't have to do it again
				logger.Info().
					Int("page-count", pc).
					Str("file", inputPath).
					Str("destination", outputPath).
					Msg("Processing finished OK")
			} else {

				logger.Error().
					Int("page-count", pc).
					Str("file", inputPath).
					Str("destination", outputPath).
					Str("error", err.Error()).
					Msg("Processing ERROR")
			}

			return err
		})
		tasks = append(tasks, newtask)
	}

	p := pool.NewPool(tasks, runtime.GOMAXPROCS(-1))

	closed := make(chan struct{})

	p.Run()

	var numErrors int
	for _, task := range p.Tasks {
		if task.Err != nil {
			numErrors++
		}
	}
	if numErrors == 0 {
		logger.Info().
			Int("task-count", len(p.Tasks)).
			Msg(fmt.Sprintf("Processed  %d scripts OK", len(p.Tasks)))
	} else {
		logger.Error().
			Int("task-count", len(p.Tasks)).
			Int("error-count", numErrors).
			Msg(fmt.Sprintf("Processed  %d scripts with %d errors", len(p.Tasks), numErrors))
	}

	close(closed)

	return nil

}

func (g *Ingester) FlattenOnePDF(inputPath, outputPath string, pageDataMap map[int]pagedata.PageData, logger *zerolog.Logger) (int, error) {

	if strings.ToLower(filepath.Ext(inputPath)) != ".pdf" {
		logger.Error().
			Str("file", inputPath).
			Msg(fmt.Sprintf("%s does not appear to be a pdf", inputPath))
		return 0, errors.New(fmt.Sprintf("%s does not appear to be a pdf", inputPath))
	}

	// need page count to find the jpeg files again later
	numPages, err := CountPages(inputPath)

	// render to images
	what := pageDataMap[1].Current.Item.What
	jpegPath := g.AcceptedPaperImages(what) //exam/coursecode

	suffix := filepath.Ext(inputPath)
	basename := strings.TrimSuffix(filepath.Base(inputPath), suffix)
	jpegFileOption := fmt.Sprintf("%s/%s%%04d.jpg", jpegPath, basename)

	f, err := os.Open(inputPath)
	if err != nil {
		logger.Error().
			Str("file", inputPath).
			Str("error", err.Error()).
			Msg(fmt.Sprintf("can't open %s", inputPath))
		return 0, err
	}

	pdfReader, err := pdf.NewPdfReader(f)
	if err != nil {
		logger.Error().
			Str("file", inputPath).
			Str("error", err.Error()).
			Msg(fmt.Sprintf("can't read %s", inputPath))
		return 0, err
	}

	comments, err := comment.GetComments(pdfReader)

	f.Close()

	err = ConvertPDFToJPEGs(inputPath, jpegPath, jpegFileOption)
	if err != nil {
		logger.Error().
			Str("file", inputPath).
			Str("destination", jpegPath).
			Str("naming-format", jpegFileOption).
			Str("error", err.Error()).
			Msg(fmt.Sprintf("can't convert to images %s", inputPath))
		return 0, err
	}

	// convert images to individual pdfs, with form overlay

	pagePath := g.AcceptedPaperPages(what)
	pageFileOption := fmt.Sprintf("%s/%s%%04d.pdf", pagePath, basename)

	mergePaths := []string{}

	// gs starts indexing at 1
	for imgIdx := 1; imgIdx <= numPages; imgIdx = imgIdx + 1 {

		pd := pageDataMap[imgIdx]

		pageNumber := imgIdx - 1

		// construct image name
		previousImagePath := fmt.Sprintf(jpegFileOption, imgIdx)
		pageFilename := fmt.Sprintf(pageFileOption, imgIdx)

		//TODO select Layout to suit landscape or portrait
		svgLayoutPath := g.FlattenLayoutSVG()

		pd.Current.Own = pagedata.FileDetail{
			Path:   pageFilename,
			UUID:   safeUUID(),
			Number: imgIdx,
			Of:     numPages,
		}

		headerPrefills := parsesvg.DocPrefills{}

		headerPrefills[pageNumber] = make(map[string]string)

		headerPrefills[pageNumber]["page-number"] = fmt.Sprintf("%d/%d", imgIdx, numPages)

		headerPrefills[pageNumber]["author"] = pageDataMap[imgIdx].Current.Item.Who

		headerPrefills[pageNumber]["date"] = pageDataMap[imgIdx].Current.Item.When

		headerPrefills[pageNumber]["title"] = pageDataMap[imgIdx].Current.Item.What

		if len(headerPrefills[pageNumber]["title"]) > 12 {
			headerPrefills[pageNumber]["title"] = headerPrefills[pageNumber]["title"][0:13]
		}
		contents := parsesvg.SpreadContents{
			SvgLayoutPath:         svgLayoutPath,
			SpreadName:            "flatten",
			PreviousImagePath:     previousImagePath,
			PageNumber:            pageNumber,
			PdfOutputPath:         pageFilename,
			Comments:              comments,
			PageData:              pd,
			TemplatePathsRelative: true,
			Prefills:              headerPrefills,
		}

		err = parsesvg.RenderSpreadExtra(contents)
		if err != nil {
			logger.Error().
				Str("error", err.Error()).
				Msg(fmt.Sprintf("Error rendering spread for %s", inputPath))
			return 0, err

		}

		mergePaths = append(mergePaths, pageFilename)
	}
	err = merge.PDF(mergePaths, outputPath)
	if err != nil {
		logger.Error().
			Str("error", err.Error()).
			Msg(fmt.Sprintf("Error merging for %s: %s", inputPath, err.Error()))
		return 0, err
	}
	doneFile := doneFilePath(outputPath)
	_, err = os.Stat(doneFile)
	if err == nil {
		err = os.Remove(doneFile)
		if err != nil {
			logger.Error().
				Str("file", outputPath).
				Str("error", err.Error()).
				Msg("Could not delete stale Done File")
		}
	}

	logger.Info().
		Str("file", inputPath).
		Int("page-count", numPages).
		Msg(fmt.Sprintf("processing finished for %s, with %d pages", inputPath, numPages))
	return numPages, nil

}
