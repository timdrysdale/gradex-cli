package ingester

import (
	"errors"
	"fmt"
	"strings"

	"github.com/timdrysdale/gradex-cli/pagedata"
)

func (g *Ingester) SplitForEnter(exam string) error {

	files, err := g.GetFileList(g.GetExamDir(exam, markerProcessed))

	if err != nil {
		return err
	}

	pdfFiles := make(map[string]bool)
	pdByFile := make(map[string]map[int]pagedata.PageData)

	for _, path := range files {

		if g.IsPDF(path) {

			pdfFiles[path] = false

			pageDataMap, _ := pagedata.UnMarshalAllFromFile(path)

			//no page data = do enter!
			if err != nil {

				(pdfFiles)[path] = true

			}

			if pagedata.GetLen(pageDataMap) < 1 {

				(pdfFiles)[path] = true

			}

			pdByFile[path] = pageDataMap

		}

	}
	numErrors := 0
	newCount := 0
	selectByOpticalOnly(&pdfFiles, pdByFile)

	for k, v := range pdfFiles {

		if v {
			copied, err := g.CopyIfNewerThanDestinationInDir(k, g.GetExamDir(exam, enterActive), g.logger)
			if copied {
				g.logger.Info().
					Str("file", k).
					Str("destination", g.GetExamDir(exam, enterActive)).
					Msg("Copied")
				newCount++
			} else {
				g.logger.Info().
					Str("file", k).
					Str("destination", g.GetExamDir(exam, enterActive)).
					Msg("Not copied - not new")
			}

			if err != nil {
				numErrors++
				g.logger.Error().
					Str("file", k).
					Str("error", err.Error()).
					Str("destination", g.GetExamDir(exam, enterActive)).
					Msg("Could not copy to enter-active dir")
			}
		} else {
			copied, err := g.CopyIfNewerThanDestinationInDir(k, g.GetExamDir(exam, enterInactive), g.logger)
			if copied {
				newCount++
				g.logger.Info().
					Str("file", k).
					Str("destination", g.GetExamDir(exam, enterInactive)).
					Msg("Copied")

			} else {
				g.logger.Info().
					Str("file", k).
					Str("destination", g.GetExamDir(exam, enterInactive)).
					Msg("Not copied - not new")
			}

			if err != nil {
				numErrors++
				g.logger.Error().
					Str("file", k).
					Str("error", err.Error()).
					Str("destination", g.GetExamDir(exam, enterInactive)).
					Msg("Could not copy to enter-inactive dir")
			}

		}

	}
	if newCount == 0 {
		g.logger.Info().
			Msg("No new files added to enter task")
	} else {
		g.logger.Info().
			Int("count", newCount).
			Msg(fmt.Sprintf("%d new files added to enter task", newCount))
	}
	if numErrors > 0 {
		return errors.New("Problems moving files into directories - check logfile for details")
	} else {

		return nil
	}

}

func selectByOpticalOnly(fileMap *map[string]bool, pageDataMap map[string]map[int]pagedata.PageData) {

	for path, _ := range *fileMap {

		//see whether there are potical fields without textfields...
		// look for tf-optical. Should this be done at page summary?
	PAGE:
		for _, docMap := range pageDataMap[path] {

			df := docMap.Current.Data

			keyMap := make(map[string]int)

			for _, item := range df {

				if strings.Contains(item.Value, markDetected) && strings.Contains(item.Key, opticalSuffix) && strings.Contains(item.Key, textFieldPrefix) {

					keyMap[strings.TrimSuffix(item.Key, opticalSuffix)] = keyMap[strings.TrimSuffix(item.Key, opticalSuffix)] + 1

				}

				if item.Value == "" && strings.Contains(item.Key, textFieldPrefix) && !strings.Contains(item.Key, opticalSuffix) {

					keyMap[strings.TrimSuffix(item.Key, opticalSuffix)] = keyMap[strings.TrimSuffix(item.Key, opticalSuffix)] + 1

				}
			}

			for _, score := range keyMap {
				if score > 1 {
					(*fileMap)[path] = true
					break PAGE
				}
			}

		}

	}

}
