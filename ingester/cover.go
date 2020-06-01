package ingester

import (
	"errors"
	"fmt"
	"time"

	"github.com/rs/zerolog"
	"github.com/timdrysdale/chmsg"
	"github.com/timdrysdale/gradex-cli/pagedata"
)

//not implemented
//func getQNumberfromDataKey(key string) int {
//	re := regexp.MustCompile("/tf-q([0-9]*)(.*)")
//	return 0
//}

func (g *Ingester) CoverPage(cp OverlayCommand, logger *zerolog.Logger) error {
	return errors.New("Not Implemented")
}

// Add a cover page summarising the marking done so far
func (g *Ingester) AddCheckCoverBar(exam string, checker string) error {
	logger := g.logger.With().Str("process", "add-check-cover-bar").Logger()
	mc := chmsg.MessagerConf{
		ExamName:     exam,
		FunctionName: "overlay",
		TaskName:     "add-check-cover-bar",
	}

	cm := chmsg.New(mc, g.msgCh, g.timeout)

	procDetail := pagedata.ProcessDetail{
		UUID:     safeUUID(),
		UnixTime: time.Now().UnixNano(),
		Name:     "check-cover",
		By:       "gradex-cli",
		ToDo:     "checking",
		For:      checker,
	}

	cp := OverlayCommand{
		FromPath:       g.GetExamDir(exam, enterProcessed),
		ToPath:         g.GetExamDirNamed(exam, checkerCover, checker),
		ExamName:       exam,
		TemplatePath:   g.OverlayLayoutSVG(),
		SpreadName:     "addition",
		ProcessDetail:  procDetail,
		Msg:            cm,
		PathDecoration: g.GetNamedTaskDecoration(checking, checker),
	}

	err := g.CoverPage(cp, &logger)
	if err == nil {
		cm.Send(fmt.Sprintf("Finished check-cover UUID=%s\n", procDetail.UUID))
		logger.Info().
			Str("UUID", procDetail.UUID).
			Str("checker", checker).
			Str("exam", exam).
			Msg("Finished add-check-cover")
	} else {
		logger.Error().
			Str("UUID", procDetail.UUID).
			Str("checker", checker).
			Str("exam", exam).
			Str("error", err.Error()).
			Msg("Error add-check-cover")
	}

	procDetail = pagedata.ProcessDetail{
		UUID:     safeUUID(),
		UnixTime: time.Now().UnixNano(),
		Name:     "check-bar",
		By:       "gradex-cli",
		ToDo:     "checking",
		For:      checker,
	}

	oc := OverlayCommand{
		CoverPath:      g.GetExamDirNamed(exam, checkerCover, checker),
		FromPath:       g.GetExamDir(exam, enterProcessed),
		ToPath:         g.GetExamDirNamed(exam, checkerReady, checker),
		ExamName:       exam,
		TemplatePath:   g.OverlayLayoutSVG(),
		SpreadName:     "check",
		ProcessDetail:  procDetail,
		Msg:            cm,
		PathDecoration: g.GetNamedTaskDecoration(checking, checker),
	}

	err = g.OverlayPapers(oc, &logger)

	if err == nil {
		cm.Send(fmt.Sprintf("Finished Processing add-check-cover-bar UUID=%s\n", procDetail.UUID))
		logger.Info().
			Str("UUID", procDetail.UUID).
			Str("checker", checker).
			Str("exam", exam).
			Msg("Finished add-check-bar")
	} else {
		logger.Error().
			Str("UUID", procDetail.UUID).
			Str("checker", checker).
			Str("exam", exam).
			Str("error", err.Error()).
			Msg("Error add-check-bar")
	}
	return err

}
