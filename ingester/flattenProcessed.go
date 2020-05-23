package ingester

import (
	"errors"
	"fmt"
	"strings"
)

// This file is to be like add bars ....

// initial sanity check on stage that has been specified
func (g *Ingester) validStageForFlattenProcessPapers(stage string) bool {

	switch strings.ToLower(stage) {

	case "marked", "remarked", "moderated", "checked", "rechecked":
		return true
	default:
		return false
	}
}

func (g *Ingester) FlattenProcessedPapers(exam, stage string) error {

	logger := g.logger.With().Str("process", "flatten-processed-papers").Str("stage", stage).Str("exam", exam).Logger()

	stage = strings.ToLower(stage)

	if !g.validStageForFlattenProcessPapers(stage) {
		logger.Error().Msg("Is not a valid stage")
		return fmt.Errorf("%s is not a valid stage for flatten-processed\n", stage)
	}

	fromDir, err := g.FlattenProcessedPapersFromDir(exam, stage)
	if err != nil {
		logger.Error().Msg("Could not get FlattenProcessedPapersFromDir")
		return err
	}

	toDir, err := g.FlattenProcessedPapersToDir(exam, stage)
	if err != nil {
		logger.Error().Msg("Could not get FlattenProcessedPapersToDir")
		return err
	}

	fmt.Printf("%s->%s\n", fromDir, toDir)

	return errors.New("not implemented yet")

}
