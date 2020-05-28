package ingester

import (
	"errors"
	"math"
	"math/rand"
	"time"
)

func (g *Ingester) SplitForModeration(exam string, minFiles int, minPercent float64) error {

	files, err := g.GetFileList(g.GetExamDir(exam, markerReady))

	if err != nil {
		return err
	}

	pdfFiles := make(map[string]bool)

	for _, file := range files {

		if g.IsPDF(file) {
			pdfFiles[file] = false
		}

	}

	reqdPercent := requiredPercent(len(pdfFiles), minFiles, minPercent)

	selectByPercent(&pdfFiles, reqdPercent)

	numErrors := 0

	for k, v := range pdfFiles {

		if v {
			err = g.MoveToDir(k, g.GetExamDir(exam, moderatorActive))
			if err != nil {
				numErrors++
				g.logger.Error().
					Str("file", k).
					Str("error", err.Error()).
					Str("destination", g.GetExamDir(exam, moderatorActive)).
					Msg("Could not move to moderate-active dir")
			}
		} else {
			err = g.MoveToDir(k, g.GetExamDir(exam, moderatorInactive))
			if err != nil {
				numErrors++
				g.logger.Error().
					Str("file", k).
					Str("error", err.Error()).
					Str("destination", g.GetExamDir(exam, moderatorInactive)).
					Msg("Could not move to moderate-inactive dir")
			}

		}

	}

	if numErrors > 0 {
		return errors.New("Problems moving files into directories")
	} else {

		return nil
	}

}

func requiredPercent(numFiles, minFiles int, minPercent float64) float64 {

	// we want whatever is greater, minFiles/numFiles, or minPercent
	// returned to us as a percentage

	byMinFilesPercent := float64(minFiles) / float64(numFiles) * 100

	if byMinFilesPercent > minPercent {
		minPercent = byMinFilesPercent
	}

	return minPercent

}

func selectByPercent(fileMap *map[string]bool, percent float64) {

	numSelected := int(0)
	s := rand.NewSource(time.Now().Unix())

	r := rand.New(s) // initialize local pseudorandom generator

	numRequired := int(math.Ceil(float64(len(*fileMap)) * percent / 100.0))
	inverted := false

	if float64(len(*fileMap)-1) < percent {
		// we need to select everything to meet the criteria
		// so just do that, save time!
		for k, _ := range *fileMap {

			(*fileMap)[k] = true
		}
		return

	}

	// perform the smallest number of selections possible
	if percent > 50 {
		percent = 100 - percent
		inverted = true
	}

LOOP:
	for {

		for k, _ := range *fileMap {

			if r.Float64()*100 < percent {
				if (*fileMap)[k] == false {
					(*fileMap)[k] = true
					numSelected++
				}
			}
			if numSelected >= numRequired {
				break LOOP
			}
		}

	}

	if inverted {
		for k, v := range *fileMap {

			(*fileMap)[k] = !v
		}
	}

}
