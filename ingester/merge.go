package ingester

import (
	"fmt"
	"strings"

	"github.com/looplab/fsm"
	"github.com/rs/zerolog"
	"github.com/timdrysdale/gradex-cli/pagedata"
)

// merge processed papers to retain duplicate pages only if they contain "work"

var (
	statusSeen    = "status-seen"
	statusMarked  = "status-marked"
	statusBad     = "status-bad"
	statusSkipped = "status-skipped"
)

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

	fmt.Println(pageSummaries)

	return nil

}

//func mapProcessedPapers(paths string) (map[string]map[int]PageCollection, error) {
//		err = addPagesToPaperMap(&papers, pageDataMap)
//}

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

//func addPagesToPaperMap(paperMap *map[string]map[int]PageCollection, pageDataMap map[int]pagedata.PageData) error {
//
//	if paperMap == nil {
//		return errors.New("Nil pointer to paperMap")
//	}
//
//	for pageNumber, pageData := range pageDataMap {
//
//		//assume pages in doc could have come from anywhere, so get original for each page
//
//		key := getOriginalKey(pageData)
//
//		if _, ok := paperMap[key]; !ok { // check if first entry for this doc
//			paperMap[key] = make(map[int]PageCollection)
//		}
//
//		pageMap := paperMap[key]
//
//		pageMap[pageNumber]
//	}
//
//	return nil
//
//}

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
