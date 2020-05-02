package ingester

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/timdrysdale/anon"
	"github.com/timdrysdale/parselearn"
	"github.com/timdrysdale/parsesvg"
	"github.com/timdrysdale/pdfcomment"
	"github.com/timdrysdale/pdfpagedata"
	"github.com/timdrysdale/pool"
	pdf "github.com/timdrysdale/unipdf/v3/model"
)

func (g *Ingester) FlattenNewPapers(exam string) error {

	logger := g.logger.With().Str("process", "flatten").Logger()

	//assume someone hits a button to ask us to do this ...

	// we'll use this same set of procDetails for flattens that we do in this batch
	// that means we can use the uuid to map the processing in graphviz later, for example
	var UUIDBytes uuid.UUID

	UUIDBytes, err := uuid.NewRandom()
	uuid := UUIDBytes.String()
	if err != nil {
		uuid = fmt.Sprintf("%d", time.Now().UnixNano())
	}
	procDetails := pdfpagedata.ProcessingDetails{
		UUID:     uuid,
		Previous: "none",
		UnixTime: time.Now().UnixNano(),
		Name:     "flatten",
		By:       pdfpagedata.ContactDetails{Name: "ingester"},
		Sequence: 0,
	}

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
		return err
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

		if getDone(pdfPath) { //check for done file - don't process if it exists
			logger.Info().
				Str("file", pdfPath).
				Msg("Skipping flattening - already done")
			continue
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

		pagedata := pdfpagedata.PageData{
			ToDo:        "flattening",
			PreparedFor: "ingester",

			Exam: pdfpagedata.ExamDetails{
				CourseCode: sub.Assignment,
				Date:       shortDate,
			},
			Author: pdfpagedata.AuthorDetails{
				Anonymous: anonymousIdentity,
			},
			Processing: []pdfpagedata.ProcessingDetails{procDetails},
		}

		renamedBase := g.GetAnonymousFileName(sub.Assignment, anonymousIdentity)
		outputPath := filepath.Join(g.AnonymousPapers(sub.Assignment), renamedBase)

		flattenTasks = append(flattenTasks, FlattenTask{
			PreparedFor: "ingester",
			ToDo:        "flattening",
			InputPath:   pdfPath,
			OutputPath:  outputPath,
			PageCount:   count,
			Data:        pagedata})
	}

	// now process the files
	N := len(flattenTasks)

	pcChan := make(chan int, N)

	tasks := []*pool.Task{}

	for i := 0; i < N; i++ {

		inputPath := flattenTasks[i].InputPath
		outputPath := flattenTasks[i].OutputPath
		pd := flattenTasks[i].Data

		newtask := pool.NewTask(func() error {
			pc, err := g.FlattenOnePDF(inputPath, outputPath, pd, &logger)
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

	//	h := thist.NewHist(nil, "Page count", "fixed", 10, false)
	//
	//	go func() {
	//	LOOP:
	//		for {
	//			select {
	//			case pc := <-pcChan:
	//				h.Update(float64(pc))
	//				fmt.Println(h.Draw())
	//			case <-closed:
	//				break LOOP
	//			}
	//		}
	//	}()
	//
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

func (g *Ingester) FlattenOnePDF(inputPath, outputPath string, pageData pdfpagedata.PageData, logger *zerolog.Logger) (int, error) {

	if strings.ToLower(filepath.Ext(inputPath)) != ".pdf" {
		logger.Error().
			Str("file", inputPath).
			Msg(fmt.Sprintf("%s does not appear to be a pdf", inputPath))
		return 0, errors.New(fmt.Sprintf("%s does not appear to be a pdf", inputPath))
	}

	// need page count to find the jpeg files again later
	numPages, err := CountPages(inputPath)

	// render to images
	jpegPath := g.AcceptedPaperImages(pageData.Exam.CourseCode)

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

	comments, err := pdfcomment.GetComments(pdfReader)

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

	pagePath := g.AcceptedPaperPages(pageData.Exam.CourseCode)
	pageFileOption := fmt.Sprintf("%s/%s%%04d.pdf", pagePath, basename)

	mergePaths := []string{}

	pageData.Page.Of = numPages

	// gs starts indexing at 1
	for imgIdx := 1; imgIdx <= numPages; imgIdx = imgIdx + 1 {

		// construct image name
		previousImagePath := fmt.Sprintf(jpegFileOption, imgIdx)
		pageFilename := fmt.Sprintf(pageFileOption, imgIdx)

		//TODO select Layout to suit landscape or portrait
		svgLayoutPath := g.FlattenLayoutSVG()

		pageNumber := imgIdx - 1

		pageData.Page.Number = pageNumber + 1
		pageData.Page.Filename = filepath.Base(pageFilename)

		var pageUUIDBytes uuid.UUID

		pageUUIDBytes, err = uuid.NewRandom()

		pageUUID := pageUUIDBytes.String()

		if err != nil {
			pageUUID = fmt.Sprintf("%d", time.Now().UnixNano())
		}

		pageData.Page.UUID = pageUUID

		headerPrefills := parsesvg.DocPrefills{}

		headerPrefills[pageNumber] = make(map[string]string)

		headerPrefills[pageNumber]["page-number"] = fmt.Sprintf("%d/%d", pageNumber+1, numPages)

		headerPrefills[pageNumber]["author"] = pageData.Author.Anonymous

		headerPrefills[pageNumber]["date"] = pageData.Exam.Date

		headerPrefills[pageNumber]["title"] = pageData.Exam.CourseCode
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
			PageData:              pageData,
			TemplatePathsRelative: true,
			Prefills:              headerPrefills,
		}

		err := parsesvg.RenderSpreadExtra(contents)
		if err != nil {
			logger.Error().
				Str("error", err.Error()).
				Msg(fmt.Sprintf("Error rendering spread for %s", inputPath))
			return 0, err

		}

		mergePaths = append(mergePaths, pageFilename)
	}
	err = MergePDF(mergePaths, outputPath)
	if err != nil {
		logger.Error().
			Str("error", err.Error()).
			Msg(fmt.Sprintf("Error merging for %s: %s", inputPath, err.Error()))
		return 0, err
	}

	logger.Info().
		Str("file", inputPath).
		Int("page-count", numPages).
		Msg(fmt.Sprintf("processing finished for %s, with %d pages", inputPath, numPages))
	return numPages, nil

}
