package ingester

import (
	"fmt"
	"sort"
	"strings"

	"github.com/looplab/fsm"
	"github.com/rs/zerolog"
	"github.com/timdrysdale/gradex-cli/pagedata"
	"github.com/timdrysdale/gradex-cli/parsesvg"
)

// merge processed papers to retain duplicate pages only if they contain "work"

var (
	statusSeen    = "status-seen"
	statusMarked  = "status-marked"
	statusBad     = "status-bad"
	statusSkipped = "status-skipped"
)

type Page struct {
	Path    string
	Message string
}

type PageSummary struct {
	Original   string //unique key (e.g. original path)
	PageNumber int
	OwnPath    string
	Status     string // pageSeen .. pogeSkipped
	WasFor     string // e.g. marker initials
}

type PageCollection struct {
	Seen    []PageSummary
	Marked  []PageSummary
	Bad     []PageSummary
	Skipped []PageSummary
}

func newPageFSM() *fsm.FSM {

	return fsm.NewFSM(
		statusSkipped,
		fsm.Events{
			{Name: statusSeen, Src: []string{statusSkipped}, Dst: statusSeen},
			{Name: statusBad, Src: []string{statusSkipped, statusSeen}, Dst: statusBad},
			{Name: statusMarked, Src: []string{statusSkipped, statusSeen, statusBad}, Dst: statusMarked},
		},
		fsm.Callbacks{},
	)
}

func (g *Ingester) MergeProcessedPapers(exam, stage string) error {

	logger := g.logger.With().Str("process", "merge-processed-papers").Str("stage", stage).Str("exam", exam).Logger()

	stage = strings.ToLower(stage)

	if !validStageForProcessedPapers(stage) {
		logger.Error().Msg("Is not a valid stage")
		return fmt.Errorf("%s is not a valid stage for flatten-processed\n", stage)
	}

	fromDir, err := g.MergeProcessedPapersFromDir(exam, stage)
	if err != nil {
		logger.Error().Msg("Could not get MergeProcessedPapersFromDir")
		return err
	}

	/*
		    toDir, err := g.MergeProcessedPapersToDir(exam, stage)
			if err != nil {
				logger.Error().Msg("Could not get MergeProcessedPapersToDir")
				return err
			}
	*/

	// get all pdf files in the FromDir
	// load pagedata from each file
	// categorising each page
	// add to map
	// when all pages added to map
	// for each document :-
	// decide which pages to merge
	// create merged docs

	paperPaths, err := g.GetFileList(fromDir)
	if err != nil {
		logger.Error().
			Str("source-dir", fromDir).
			Str("error", err.Error()).
			Msg(fmt.Sprintf("Stopping early; couldn't get files because %v\n", err))
		return err
	}

	//paperMap, err := mapProcessedPapers(paperPaths)

	pageSummaries, err := summariseFiles(paperPaths, &logger)

	if err != nil {
		logger.Error().
			Str("error", err.Error()).
			Msg(fmt.Sprintf("Stopping early; couldn't summarise files because %v\n", err))
		return err
	}

	paperMap := createPaperMap(pageSummaries)

	mergePathMap := createMergePathMap(paperMap)

	parsesvg.PrettyPrintStruct(mergePathMap)

	// TODO convert key to output basename (add full path, decoration, find right path?)

	//	for key, Page
	//					err = merge.PDF(mergePaths, ot.OutputPath)

	return nil

}

func summariseFiles(paths []string, logger *zerolog.Logger) ([]PageSummary, error) {

	summaries := []PageSummary{}

	var lastError error
	lastError = nil

	for _, path := range paths {

		pageDataMap, err := pagedata.UnMarshalAllFromFile(path)

		if err != nil {
			lastError = fmt.Errorf("Skipping (%s): error obtaining pagedata", path)
			logger.Error().
				Str("file", path).
				Str("error", err.Error()).
				Msg(fmt.Sprintf("Skipping (%s): error obtaining pagedata\n", path))
			continue
		}

		if pagedata.GetLen(pageDataMap) < 1 {
			lastError = fmt.Errorf("Skipping (%s): no pagedata in file", path)
			logger.Error().
				Str("file", path).
				Msg(fmt.Sprintf("Skipping (%s): no pagedata in file\n", path))
			continue
		}

		for _, pageData := range pageDataMap {

			summary := summarisePage(pageData)

			summaries = append(summaries, summary)

		}
	}

	return summaries, lastError
}

func summarisePage(pageData pagedata.PageData) PageSummary {

	pageFSM := newPageFSM()

	// Original.Number, Original.Of

	for _, item := range pageData.Current.Data {

		if strings.Contains(item.Key, textFieldPrefix) {

			if strings.Contains(item.Key, "page-ok") && item.Value != "" {
				pageFSM.Event(statusSeen)
			}
			if strings.Contains(item.Key, "page-bad") && item.Value != "" {
				pageFSM.Event(statusBad)
			}
			if !strings.Contains(item.Key, "page-bad") && !strings.Contains(item.Key, "page-ok") && item.Value != "" {
				pageFSM.Event(statusMarked)
			}
		}
	}
	summary := PageSummary{
		Original:   getOriginalKey(pageData),
		PageNumber: getPageNumber(pageData),
		OwnPath:    getOwnPath(pageData),
		Status:     pageFSM.Current(),
		WasFor:     getWasFor(pageData),
	}

	return summary
}

func createPaperMap(summaries []PageSummary) map[string]map[int]PageCollection {

	paperMap := make(map[string]map[int]PageCollection)

	for _, summary := range summaries {

		if _, ok := paperMap[summary.Original]; !ok {
			paperMap[summary.Original] = make(map[int]PageCollection)
		}

		if _, ok := paperMap[summary.Original][summary.PageNumber]; !ok {
			var emptyPC PageCollection
			paperMap[summary.Original][summary.PageNumber] = emptyPC
		}

		pc := paperMap[summary.Original][summary.PageNumber]
		switch summary.Status {
		case statusBad:
			pc.Bad = append(pc.Bad, summary)
		case statusMarked:
			pc.Marked = append(pc.Marked, summary)
		case statusSeen:
			pc.Seen = append(pc.Seen, summary)
		case statusSkipped:
			pc.Skipped = append(pc.Skipped, summary)
		}

		paperMap[summary.Original][summary.PageNumber] = pc
	}

	return paperMap

}

func createMergePathMap(paperMap map[string]map[int]PageCollection) map[string][]Page {

	// make decisions about what to keep
	mergePageMap := make(map[string]map[int][]Page) // map by page

	for key, collectionMap := range paperMap {

		if _, ok := mergePageMap[key]; !ok {
			mergePageMap[key] = make(map[int][]Page)
		}

		for pageNumber, pageCollection := range collectionMap {

			pageList := createPageList(pageCollection)

			mergePageMap[key][pageNumber] = pageList

		}

	}

	// order the pages for each file
	mergePathMap := make(map[string][]Page) //sorted list for each file

	for key, pageMap := range mergePageMap {

		pageNumbers := []int{}

		for pageNumber, _ := range pageMap {
			pageNumbers = append(pageNumbers, pageNumber)
		}

		sort.Ints(pageNumbers)

		sortedList := []Page{}

		for _, pageNumber := range pageNumbers { //go by page number

			for _, page := range pageMap[pageNumber] { // collect all copies of this page
				sortedList = append(sortedList, page)
			}

		}

		mergePathMap[key] = sortedList

	}

	return mergePathMap

}

// TODO summarise pageCollection on each page's message
func createPageItem(pageCollection PageCollection, thisPage PageSummary) Page {

	message := "This page " + strings.TrimPrefix(thisPage.Status, "status-") + " by " + thisPage.WasFor + "\nMarked:"

	for _, summary := range pageCollection.Marked {
		message = message + " " + summary.WasFor
	}

	message = message + "\nBad:"
	for _, summary := range pageCollection.Bad {
		message = message + " " + summary.WasFor
	}

	message = message + "\nSeen:"
	for _, summary := range pageCollection.Seen {
		message = message + " " + summary.WasFor
	}

	message = message + "\nSkipped:"
	for _, summary := range pageCollection.Skipped {
		message = message + " " + summary.WasFor
	}

	return Page{
		Path:    thisPage.OwnPath,
		Message: message,
	}

}

func createPageList(pageCollection PageCollection) []Page {

	pageList := []Page{}

	for _, summary := range pageCollection.Marked {
		pageList = append(pageList, createPageItem(pageCollection, summary))
	}

	// the pageList.Message summarises everything else we need to know
	if len(pageList) > 0 {
		return pageList
	}

	// return a single page from any other list of pages
	if len(pageCollection.Bad) > 0 {
		for _, summary := range pageCollection.Bad {
			pageList = append(pageList, createPageItem(pageCollection, summary))
			return pageList
		}
	}

	if len(pageCollection.Seen) > 0 {
		for _, summary := range pageCollection.Seen {
			pageList = append(pageList, createPageItem(pageCollection, summary))
			return pageList
		}
	}

	if len(pageCollection.Skipped) > 0 {
		for _, summary := range pageCollection.Skipped {
			pageList = append(pageList, createPageItem(pageCollection, summary))
			return pageList
		}
	}

	return pageList
}

func getOriginalKey(pageData pagedata.PageData) string {
	return pageData.Current.Original.Path
}

func getPageNumber(pageData pagedata.PageData) int {
	return pageData.Current.Original.Number
}

func getOwnPath(pageData pagedata.PageData) string {
	return pageData.Current.Own.Path
}
func getWasFor(pageData pagedata.PageData) string {
	return (pageData.Previous[len(pageData.Previous)-1]).Process.For
}
