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

//"<gradex-pagedata>{"current":{"is":"page","own":{"path":"/usr/local/gradex/usr/exam/MECE11016 Advanced Composite Materials 5 - Exam Dropbox/04-temporary-pages/MECE11016 Advanced Composite Materials 5 - Exam Dropbox_s1518636_attempt_2020-04-29-20-42-31_MECE11016-B0871560015.pdf","UUID":"490aec82-df2b-4776-bf62-1bbb63bab63b","number":15,"of":15},"original":{"path":"/usr/local/gradex/usr/exam/MECE11016 Advanced Composite Materials 5 - Exam Dropbox/03-accepted-papers/MECE11016 Advanced Composite Materials 5 - Exam Dropbox_s1518636_attempt_2020-04-29-20-42-31_MECE11016-B087156.pdf","UUID":"c85d0369-6cb0-4a91-bdf8-c059a7dd8c37","number":15,"of":15},"current":{"path":"","UUID":"","number":0,"of":0},"item":{"what":"MECE11016 Advanced Composite Materials 5 - Exam Dropbox","when":"29-Apr-2020","who":"B087156","UUID":"","whoType":"anonymous"},"process":{"name":"mark-bar","UUID":"f8a76765-d6d2-4ac9-a729-813d2b3900fc","unixTime":1588682003585183986,"for":"DR","toDo":"marking","by":"gradex-cli","data":null},"UUID":"4c4b2944-1aa4-4e7a-921a-150b87aaf6a9","follows":"7ca27a4f-2098-4a1a-aadf-7ca0e2fa9a10","revision":0,"data":null},"previous":[{"is":"page","own":{"path":"/usr/local/gradex/usr/exam/MECE11016 Advanced Composite Materials 5 - Exam Dropbox/04-temporary-pages/MECE11016 Advanced Composite Materials 5 - Exam Dropbox_s1518636_attempt_2020-04-29-20-42-31_MECE11016-B0871560015.pdf","UUID":"490aec82-df2b-4776-bf62-1bbb63bab63b","number":15,"of":15},"original":{"path":"/usr/local/gradex/usr/exam/MECE11016 Advanced Composite Materials 5 - Exam Dropbox/03-accepted-papers/MECE11016 Advanced Composite Materials 5 - Exam Dropbox_s1518636_attempt_2020-04-29-20-42-31_MECE11016-B087156.pdf","UUID":"c85d0369-6cb0-4a91-bdf8-c059a7dd8c37","number":15,"of":15},"current":{"path":"","UUID":"","number":0,"of":0},"item":{"what":"MECE11016 Advanced Composite Materials 5 - Exam Dropbox","when":"29-Apr-2020","who":"B087156","UUID":"","whoType":"anonymous"},"process":{"name":"flatten","UUID":"b86b96eb-a192-4d63-9e8d-e451d016835e","unixTime":1588680034255432373,"for":"ingester","toDo":"prepare-for-marking","by":"gradex-cli","data":null},"UUID":"7ca27a4f-2098-4a1a-aadf-7ca0e2fa9a10","follows":"","revision":0,"data":null}]}</gradex-pagedata><hash>1640479646</hash>

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
