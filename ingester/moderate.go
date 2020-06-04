package ingester

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/timdrysdale/gradex-cli/util"
)

func (g *Ingester) SplitForModeration(exam string, minFiles int, minPercent float64) error {

	dir := g.GetExamDir(exam, markerProcessed)

	files, err := g.GetFileList(dir)

	if err != nil {
		return err
	}

	pdfFiles := make(map[string]bool)

	inputCount := 0

	for _, file := range files {

		if g.IsPDF(file) {
			pdfFiles[file] = false
			inputCount++
		}
	}

	fmt.Printf("We think we have %d files to split for moderating\n", inputCount)

	reqdPercent := requiredPercent(len(pdfFiles), minFiles, minPercent)

	predictedActiveCount := float64(inputCount) * reqdPercent / 100

	fmt.Printf("We think we want %f percent of the files, which equates to  %f files, to be active; minpercent was %f\n", reqdPercent, predictedActiveCount, minPercent)

	selectByPercent(&pdfFiles, reqdPercent)

	util.PrettyPrintStruct(pdfFiles)

	activeCount := 0
	inactiveCount := 0

	for _, isActive := range pdfFiles {
		if isActive {
			activeCount++
		} else {
			inactiveCount++
		}
	}

	predictedInactiveCount := inputCount - int(math.Round(predictedActiveCount))

	fmt.Printf("We think we have %d active files (wanted %d), and %d inactive files (should be %d)\n", activeCount, int(predictedActiveCount), inactiveCount, predictedInactiveCount)
	numErrors := 0
	newCount := 0
	for k, v := range pdfFiles {

		if v {
			copied, err := g.CopyIfNewerThanDestinationInDir(k, g.GetExamDir(exam, moderatorActive), g.logger)
			if copied {
				g.logger.Info().
					Str("file", k).
					Str("destination", g.GetExamDir(exam, moderatorActive)).
					Msg("Copied")
				newCount++
			} else {
				g.logger.Info().
					Str("file", k).
					Str("destination", g.GetExamDir(exam, moderatorActive)).
					Msg("Not copied - not new")
			}

			if err != nil {
				numErrors++
				g.logger.Error().
					Str("file", k).
					Str("error", err.Error()).
					Str("destination", g.GetExamDir(exam, moderatorActive)).
					Msg("Could not move to moderate-active dir")
			}
		} else {
			copied, err := g.CopyIfNewerThanDestinationInDir(k, g.GetExamDir(exam, moderatorInactive), g.logger)
			if copied {
				newCount++
				g.logger.Info().
					Str("file", k).
					Str("destination", g.GetExamDir(exam, moderatorInactive)).
					Msg("Copied")

			} else {
				g.logger.Info().
					Str("file", k).
					Str("destination", g.GetExamDir(exam, moderatorInactive)).
					Msg("Not copied - not new")
			}
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

	if newCount == 0 {
		g.logger.Info().
			Msg("No new files added to moderation task")
	} else {
		g.logger.Info().
			Int("count", newCount).
			Msg(fmt.Sprintf("%d new files added to moderation task", newCount))
	}

	if numErrors > 0 {
		return errors.New("Problems moving files into directories - check logfile for details")
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

	numRequired := int(math.Ceil(float64(len(*fileMap))*percent/100.0 - 1e-9))

	// is the percentage high enough that selecting one less than the total is not enough?

	nMinusOnePercentage := float64(len(*fileMap)-1) / float64(len(*fileMap))

	nMinusOneIsNotEnough := nMinusOnePercentage < (percent / 100)

	if nMinusOneIsNotEnough {
		// we need to select everything to meet the criteria
		// so just do that, save time!
		for k, _ := range *fileMap {

			(*fileMap)[k] = true
		}
		return

	}

	// perform the smallest number of selections possible
	// Note that inverting the percentage is all we need to do
	// We do NOT need to invert the logic values afterwards
	// that will just mean we select fewer than intended
	if percent > 50 {
		percent = 100 - percent
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

}
