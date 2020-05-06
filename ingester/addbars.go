package ingester

import (
	"fmt"
	"time"

	"github.com/timdrysdale/chmsg"
	"github.com/timdrysdale/gradex-cli/pagedata"
)

func (g *Ingester) AddMarkBar(exam string, marker string) error {

	logger := g.logger.With().Str("process", "add-mark-bar").Logger()

	mc := chmsg.MessagerConf{
		ExamName:     exam,
		FunctionName: "overlay",
		TaskName:     "add-mark-bar",
	}

	cm := chmsg.New(mc, g.msgCh, g.timeout)

	procDetail := pagedata.ProcessDetail{
		UUID:     safeUUID(),
		UnixTime: time.Now().UnixNano(),
		Name:     "mark-bar",
		By:       "gradex-cli",
		ToDo:     "marking",
		For:      marker,
	}

	oc := OverlayCommand{
		FromPath:       g.AnonymousPapers(exam),
		ToPath:         g.MarkerReady(exam, marker),
		ExamName:       exam,
		TemplatePath:   g.OverlayLayoutSVG(),
		SpreadName:     "mark",
		ProcessDetail:  procDetail,
		Msg:            cm,
		PathDecoration: g.MarkerABCDecoration(marker),
	}

	err := g.OverlayPapers(oc, &logger)

	if err == nil {
		cm.Send(fmt.Sprintf("Finished Processing markbar UUID=%s\n", procDetail.UUID))
		logger.Info().
			Str("UUID", procDetail.UUID).
			Str("marker", marker).
			Str("exam", exam).
			Msg("Finished add-mark-bar")
	} else {
		logger.Error().
			Str("UUID", procDetail.UUID).
			Str("marker", marker).
			Str("exam", exam).
			Str("error", err.Error()).
			Msg("Error add-mark-bar")
	}
	return err
}

func (g *Ingester) AddModerateActiveBar(exam string, moderator string) error {

	logger := g.logger.With().Str("process", "add-moderate-active-bar").Logger()

	mc := chmsg.MessagerConf{
		ExamName:     exam,
		FunctionName: "overlay",
		TaskName:     "add-moderate-active-bar",
	}

	cm := chmsg.New(mc, g.msgCh, g.timeout)

	procDetail := pagedata.ProcessDetail{
		UUID:     safeUUID(),
		UnixTime: time.Now().UnixNano(),
		Name:     "moderate-active-bar",
		By:       "gradex-cli",
		ToDo:     "moderating",
		For:      moderator,
	}

	oc := OverlayCommand{
		FromPath:       g.ModerateActive(exam),
		ToPath:         g.ModeratorReady(exam, moderator),
		ExamName:       exam,
		TemplatePath:   g.OverlayLayoutSVG(),
		SpreadName:     "moderate-active",
		ProcessDetail:  procDetail,
		Msg:            cm,
		PathDecoration: g.ModeratorABCDecoration(moderator),
	}

	err := g.OverlayPapers(oc, &logger)

	if err == nil {
		cm.Send(fmt.Sprintf("Finished Processing add-moderate-active-bar UUID=%s\n", procDetail.UUID))
		logger.Info().
			Str("UUID", procDetail.UUID).
			Str("moderator", moderator).
			Str("exam", exam).
			Msg("Finished add-moderate-active-bar")
	} else {
		logger.Error().
			Str("UUID", procDetail.UUID).
			Str("moderator", moderator).
			Str("exam", exam).
			Str("error", err.Error()).
			Msg("Error add-moderate-active-bar")
	}

	return err
}

func (g *Ingester) AddModerateInActiveBar(exam string) error {

	logger := g.logger.With().Str("process", "add-moderate-inactive-bar").Logger()

	mc := chmsg.MessagerConf{
		ExamName:     exam,
		FunctionName: "overlay",
		TaskName:     "add-moderate-inactive-bar",
	}

	cm := chmsg.New(mc, g.msgCh, g.timeout)

	procDetail := pagedata.ProcessDetail{
		UUID:     safeUUID(),
		UnixTime: time.Now().UnixNano(),
		Name:     "moderate-active-bar",
		By:       "gradex-cli",
		ToDo:     "moderating",
		For:      "X",
	}

	oc := OverlayCommand{
		FromPath:       g.ModerateInActive(exam),
		ToPath:         g.ModeratedInActiveBack(exam),
		ExamName:       exam,
		TemplatePath:   g.OverlayLayoutSVG(),
		SpreadName:     "moderate-inactive",
		ProcessDetail:  procDetail,
		Msg:            cm,
		PathDecoration: g.ModeratorABCDecoration("X"),
	}
	err := g.OverlayPapers(oc, &logger)

	if err == nil {
		cm.Send(fmt.Sprintf("Finished Processing add-moderate-inactive-bar UUID=%s\n", procDetail.UUID))
		logger.Info().
			Str("UUID", procDetail.UUID).
			Str("exam", exam).
			Msg("Finished add-moderate-inactive-bar")
	} else {
		logger.Error().
			Str("UUID", procDetail.UUID).
			Str("exam", exam).
			Str("error", err.Error()).
			Msg("Error add-moderate-inactive-bar")
	}

	return err
}

func (g *Ingester) AddCheckBar(exam string, checker string) error {
	logger := g.logger.With().Str("process", "add-check-bar").Logger()
	mc := chmsg.MessagerConf{
		ExamName:     exam,
		FunctionName: "overlay",
		TaskName:     "add-check-bar",
	}

	cm := chmsg.New(mc, g.msgCh, g.timeout)

	procDetail := pagedata.ProcessDetail{
		UUID:     safeUUID(),
		UnixTime: time.Now().UnixNano(),
		Name:     "check-bar",
		By:       "gradex-cli",
		ToDo:     "checking",
		For:      checker,
	}

	oc := OverlayCommand{
		FromPath:       g.ModeratedReady(exam),
		ToPath:         g.CheckerReady(exam, checker),
		ExamName:       exam,
		TemplatePath:   g.OverlayLayoutSVG(),
		SpreadName:     "check",
		ProcessDetail:  procDetail,
		Msg:            cm,
		PathDecoration: g.CheckerABCDecoration(checker),
	}

	err := g.OverlayPapers(oc, &logger)

	if err == nil {
		cm.Send(fmt.Sprintf("Finished Processing add-check-bar UUID=%s\n", procDetail.UUID))
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
			Msg("Error add-checker-bar")
	}
	return err

}