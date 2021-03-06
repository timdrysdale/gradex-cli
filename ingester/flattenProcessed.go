package ingester

import (
	"fmt"
	"strings"
	"time"

	"github.com/timdrysdale/chmsg"
	"github.com/timdrysdale/gradex-cli/pagedata"
)

// This file is to be like add bars ....

// initial sanity check on stage that has been specified
// also used by merge "half" of the process (see merge.go)
func ValidStageForProcessedPapers(stage string) bool {

	switch strings.ToLower(stage) {

	case "marked", "remarked", "moderated", "remoderated", "entered", "reentered", "checked", "rechecked":
		return true
	default:
		return false
	}
}

func getSpreadForBoxes(stage string) string {

	switch stage {

	case "marked":
		return "mark"
	case "remarked":
		return "remark"
	case "moderated":
		return "moderate-active" //we don't get boxes for inactive - NOTE we can know this because no textfields either!
	case "entered":
		return "enter-active"
	case "checked":
		return "check"
	case "rechecked":
		return "recheck"
	default:
		return ""
	}
}

func (g *Ingester) FlattenProcessedPapers(exam, stage string) error {

	logger := g.logger.With().Str("process", "flatten-processed-papers").Str("stage", stage).Str("exam", exam).Logger()

	stage = strings.ToLower(stage)

	if !ValidStageForProcessedPapers(stage) {
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

	taskName := fmt.Sprintf("flatten-%s", stage)

	mc := chmsg.MessagerConf{
		ExamName:     exam,
		FunctionName: "overlay",
		TaskName:     taskName,
	}

	cm := chmsg.New(mc, g.msgCh, g.timeout)

	procDetail := pagedata.ProcessDetail{
		UUID:     safeUUID(),
		UnixTime: time.Now().UnixNano(),
		Name:     taskName,
		By:       "gradex-cli",
		ToDo:     "further-processing",
		For:      "ingester",
	}

	oc := OverlayCommand{
		FromPath:             fromDir,
		ToPath:               toDir,
		AncestorPath:         g.GetExamDir(exam, anonPapers),
		ExamName:             exam,
		TemplatePath:         g.OverlayLayoutSVG(),
		SpreadName:           "flatten-processed",
		ProcessDetail:        procDetail,
		Msg:                  cm,
		OpticalBoxSpread:     getSpreadForBoxes(stage),
		ReadOpticalBoxes:     true,
		OmitPreviousComments: true, //avoid QBOX line in report checked from previous stage's comments
	}

	err = g.OverlayPapers(oc, &logger)

	if err == nil {
		cm.Send(fmt.Sprintf("Finished Processing %s UUID=%s\n", taskName, procDetail.UUID))
		logger.Info().
			Str("UUID", procDetail.UUID).
			Msg(fmt.Sprintf("Finished %s", taskName))
	} else {
		logger.Error().
			Str("UUID", procDetail.UUID).
			Str("error", err.Error()).
			Msg(fmt.Sprintf("Error %s", taskName))
	}

	return err

}
