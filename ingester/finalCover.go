package ingester

import (
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	"github.com/fvbommel/sortorder"
	"github.com/rs/zerolog"
	"github.com/timdrysdale/gradex-cli/pagedata"
	"github.com/timdrysdale/gradex-cli/parsesvg"
	"github.com/timdrysdale/pool"
)

func (g *Ingester) FinalCoverPage(cp CoverPageCommand, logger *zerolog.Logger) error {
	//find pages in processed fir
	// for each page, mangle the name to get the coverpage name
	// sum up all the questions
	// make sure only questions in questions conf are included in the cover page

	files, err := g.GetFileList(cp.FromPath)

	if err != nil {
		logger.Error().
			Str("dir", cp.FromPath).
			Str("error", err.Error()).
			Msg("Error getting files")

		return err
	}

	cpTasks := []CoverPageTask{}

	for _, path := range files {
		if !IsPDF(path) {
			continue
		}

		cpTasks = append(cpTasks, CoverPageTask{
			Path:    path,
			Command: cp,
		})

	}

	//>>>>>>>>>>>>>>>>>>>>> TASKS READY

	N := len(cpTasks)

	tasks := []*pool.Task{}

	for i := 0; i < N; i++ {

		cpt := cpTasks[i]

		newtask := pool.NewTask(func() error {
			err := DoOneFinalCoverPage(cpt, logger)

			if err == nil {
				setDone(cpt.Path, logger)
				logger.Debug().Str("file", cpt.Path).Msg("set done file at source")
				logger.Info().
					Str("file", cpt.Path).
					Msg(fmt.Sprintf("Finished processing %s", cpt.Path))

				return nil
			} else {
				logger.Error().
					Str("file", cpt.Path).
					Str("error", err.Error()).
					Msg(fmt.Sprintf("Error processing %s", cpt.Path))
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
			Msg(fmt.Sprintf("Processing finished with final coverpage tasks returning <%d> errors from <%d> scripts\n", numErrors, N))
	} else {
		logger.Info().
			Int("error-count", numErrors).
			Int("script-count", N).
			Msg(fmt.Sprintf("Processing finished <%d> final scripts without any errors\n", N))
	}

	return nil

}

func DoOneFinalCoverPage(ct CoverPageTask, logger *zerolog.Logger) error {

	pdMap := make(map[int]pagedata.PageData)

	path := ct.Path
	cp := ct.Command

	pdMap, err := pagedata.UnMarshalAllFromFile(path)

	if err != nil {
		logger.Error().
			Str("file", path).
			Str("error", err.Error()).
			Msg("Error getting pagedata")
	}

	if _, ok := pdMap[1]; !ok {
		return errors.New("No pagedata at page one; needed for cover")
	}
	QMap, skipMap := finalMarksMap(pdMap[1], cp.Questions)

	for q, skip := range skipMap {
		if skip {
			fmt.Printf("%s: %s not addition checked\n", filepath.Base(path), q)
		}
	}

	pageNumber := 0 //starts at zero

	Prefills := parsesvg.DocPrefills{}

	Prefills[pageNumber] = make(map[string]string)

	Prefills[pageNumber]["page-number"] = "final"

	var thisPageData pagedata.PageData
	// get first page data
	for _, pdm := range pdMap {
		thisPageData = pdm
		break
	}

	fields := []pagedata.Field{}

	for k, v := range QMap {
		fields = append(fields, pagedata.Field{
			Key:   "pf-q-" + k,
			Value: v,
		})
	}

	thisPageData.Current.Data = fields

	Prefills[pageNumber]["author"] = thisPageData.Current.Item.Who

	Prefills[pageNumber]["date"] = thisPageData.Current.Item.When

	Prefills[pageNumber]["title"] = shortenAssignment(thisPageData.Current.Item.What)

	Prefills[pageNumber]["for"] = thisPageData.Current.Process.For

	var qkeys []string
	for k := range QMap {
		qkeys = append(qkeys, k)
	}

	sort.Sort(sortorder.Natural(qkeys))

	for idx, qk := range qkeys {
		question := fmt.Sprintf("question-%02d", idx)
		mark := fmt.Sprintf("mark-awarded-%02d", idx)
		Prefills[pageNumber][question] = qk
		Prefills[pageNumber][mark] = QMap[qk]

	}

	pageFilename := filepath.Join(cp.ToPath, strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))+"-cover.pdf")

	contents := parsesvg.SpreadContents{
		SvgLayoutPath:         cp.TemplatePath,
		SpreadName:            cp.SpreadName,
		PageNumber:            pageNumber,
		PdfOutputPath:         pageFilename,
		PageData:              thisPageData,
		TemplatePathsRelative: true,
		Prefills:              Prefills,
	}

	err = parsesvg.RenderSpreadExtra(contents)
	if err != nil {
		msg := fmt.Sprintf("Error rendering spread for cover page for (%s) because %v\n", path, err)
		logger.Error().
			Str("file", path).
			Str("error", err.Error()).
			Msg(msg)
		fmt.Println(msg)
	}
	return err
}
