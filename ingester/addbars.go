package ingester

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	"github.com/timdrysdale/chmsg"
	"github.com/timdrysdale/gradex-cli/pagedata"
)

func (g *Ingester) AddLabelBar(exam, labeller string) error {

	logger := g.logger.With().Str("process", "add-label-bar").Logger()

	mc := chmsg.MessagerConf{
		ExamName:     exam,
		FunctionName: "overlay",
		TaskName:     "add-label-bar",
	}

	cm := chmsg.New(mc, g.msgCh, g.timeout)

	procDetail := pagedata.ProcessDetail{
		UUID:     safeUUID(),
		UnixTime: time.Now().UnixNano(),
		Name:     "label-bar",
		By:       "gradex-cli",
		ToDo:     "labelling",
		For:      labeller,
	}

	oc := OverlayCommand{
		FromPath:       g.GetExamDir(exam, anonPapers),
		ToPath:         g.GetExamDirNamed(exam, questionReady, labeller),
		ExamName:       exam,
		TemplatePath:   g.OverlayLayoutSVG(),
		SpreadName:     "label",
		ProcessDetail:  procDetail,
		Msg:            cm,
		PathDecoration: g.GetNamedTaskDecoration(labelling, labeller),
	}

	err := g.OverlayPapers(oc, &logger)

	if err == nil {
		cm.Send(fmt.Sprintf("Finished Processing markbar UUID=%s\n", procDetail.UUID))
		logger.Info().
			Str("UUID", procDetail.UUID).
			Str("labeller", labeller).
			Str("exam", exam).
			Msg("Finished add-label-bar")
	} else {
		logger.Error().
			Str("UUID", procDetail.UUID).
			Str("labeller", labeller).
			Str("exam", exam).
			Str("error", err.Error()).
			Msg("Error add-label-bar")
	}
	return err
}

func (g *Ingester) AddMarkBarByQ(exam string, marker string) error {

	logger := g.logger.With().Str("process", "add-mark-bar").Logger()

	mc := chmsg.MessagerConf{
		ExamName:     exam,
		FunctionName: "overlay",
		TaskName:     "add-mark-bar-byQ",
	}

	cm := chmsg.New(mc, g.msgCh, g.timeout)

	procDetail := pagedata.ProcessDetail{
		UUID:     safeUUID(),
		UnixTime: time.Now().UnixNano(),
		Name:     "mark-bar-byQ",
		By:       "gradex-cli",
		ToDo:     "marking",
		For:      marker,
	}

	oc := OverlayCommand{
		FromPath:       g.GetExamDir(exam, questionSplit),
		ToPath:         g.GetExamDirNamed(exam, markerReady, marker),
		ExamName:       exam,
		TemplatePath:   g.OverlayLayoutSVG(),
		SpreadName:     "mark",
		ProcessDetail:  procDetail,
		Msg:            cm,
		PathDecoration: g.GetNamedTaskDecoration(marking, marker),
	}

	err := g.OverlayPapers(oc, &logger)

	if err == nil {
		cm.Send(fmt.Sprintf("Finished Processing markbar UUID=%s\n", procDetail.UUID))
		logger.Info().
			Str("UUID", procDetail.UUID).
			Str("marker", marker).
			Str("exam", exam).
			Msg("Finished add-mark-bar-byQ")
	} else {
		logger.Error().
			Str("UUID", procDetail.UUID).
			Str("marker", marker).
			Str("exam", exam).
			Str("error", err.Error()).
			Msg("Error add-mark-bar-byQ")
	}
	return err
}

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
		FromPath:       g.GetExamDir(exam, anonPapers),
		ToPath:         g.GetExamDirNamed(exam, markerReady, marker),
		ExamName:       exam,
		TemplatePath:   g.OverlayLayoutSVG(),
		SpreadName:     "mark",
		ProcessDetail:  procDetail,
		Msg:            cm,
		PathDecoration: g.GetNamedTaskDecoration(marking, marker),
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
		FromPath:       g.GetExamDir(exam, moderatorActive),
		ToPath:         g.GetExamDirNamed(exam, moderatorReady, moderator),
		ExamName:       exam,
		TemplatePath:   g.OverlayLayoutSVG(),
		SpreadName:     "moderate-active",
		ProcessDetail:  procDetail,
		Msg:            cm,
		PathDecoration: g.GetNamedTaskDecoration(moderating, moderator),
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
		Name:     "moderate-inactive-bar",
		By:       "gradex-cli",
		ToDo:     "moderating",
		For:      "X",
	}

	oc := OverlayCommand{
		FromPath:       g.GetExamDir(exam, moderatorInactive),
		ToPath:         g.GetExamDirSub(exam, moderatorBack, inactive),
		ExamName:       exam,
		TemplatePath:   g.OverlayLayoutSVG(),
		SpreadName:     "moderate-inactive",
		ProcessDetail:  procDetail,
		Msg:            cm,
		PathDecoration: g.GetNamedTaskDecoration(moderating, "X"),
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

func (g *Ingester) AddEnterActiveBar(exam string, enterer string) error {
	logger := g.logger.With().Str("process", "add-enter-active-bar").Logger()
	mc := chmsg.MessagerConf{
		ExamName:     exam,
		FunctionName: "overlay",
		TaskName:     "add-enter-bar",
	}

	cm := chmsg.New(mc, g.msgCh, g.timeout)

	procDetail := pagedata.ProcessDetail{
		UUID:     safeUUID(),
		UnixTime: time.Now().UnixNano(),
		Name:     "enter-active-bar",
		By:       "gradex-cli",
		ToDo:     "entering",
		For:      enterer,
	}

	oc := OverlayCommand{
		FromPath:       g.GetExamDir(exam, enterActive),
		ToPath:         g.GetExamDirNamed(exam, enterReady, enterer),
		ExamName:       exam,
		TemplatePath:   g.OverlayLayoutSVG(),
		SpreadName:     "enter-active",
		ProcessDetail:  procDetail,
		Msg:            cm,
		PathDecoration: g.GetNamedTaskDecoration(entering, enterer),
	}

	err := g.OverlayPapers(oc, &logger)

	if err == nil {
		cm.Send(fmt.Sprintf("Finished Processing add-enter-active-bar UUID=%s\n", procDetail.UUID))
		logger.Info().
			Str("UUID", procDetail.UUID).
			Str("enterer", enterer).
			Str("exam", exam).
			Msg("Finished add-enter-active-bar")
	} else {
		logger.Error().
			Str("UUID", procDetail.UUID).
			Str("enterer", enterer).
			Str("exam", exam).
			Str("error", err.Error()).
			Msg("Error add-enter-active-bar")
	}
	return err

}

func (g *Ingester) AddEnterInactiveBar(exam string) error {
	logger := g.logger.With().Str("process", "add-check-bar").Logger()
	mc := chmsg.MessagerConf{
		ExamName:     exam,
		FunctionName: "overlay",
		TaskName:     "add-enter-inactive-bar",
	}

	cm := chmsg.New(mc, g.msgCh, g.timeout)
	enterer := "X"
	procDetail := pagedata.ProcessDetail{
		UUID:     safeUUID(),
		UnixTime: time.Now().UnixNano(),
		Name:     "enter-inactive-bar",
		By:       "gradex-cli",
		ToDo:     "entering",
		For:      enterer,
	}

	oc := OverlayCommand{
		FromPath:       g.GetExamDir(exam, enterInactive),
		ToPath:         g.GetExamDirSub(exam, enterBack, inactive),
		ExamName:       exam,
		TemplatePath:   g.OverlayLayoutSVG(),
		SpreadName:     "enter-inactive", //same design
		ProcessDetail:  procDetail,
		Msg:            cm,
		PathDecoration: g.GetNamedTaskDecoration(entering, enterer),
	}

	err := g.OverlayPapers(oc, &logger)

	if err == nil {
		cm.Send(fmt.Sprintf("Finished Processing add-enter-inactive-bar UUID=%s\n", procDetail.UUID))
		logger.Info().
			Str("UUID", procDetail.UUID).
			Str("enterer", enterer).
			Str("exam", exam).
			Msg("Finished add-enter-inactive-bar")
	} else {
		logger.Error().
			Str("UUID", procDetail.UUID).
			Str("enterer", enterer).
			Str("exam", exam).
			Str("error", err.Error()).
			Msg("Error add-enter-inactive-bar")
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
		FromPath:       g.GetExamDir(exam, enterProcessed),
		ToPath:         g.GetExamDirNamed(exam, checkerReady, checker),
		ExamName:       exam,
		TemplatePath:   g.OverlayLayoutSVG(),
		SpreadName:     "check",
		ProcessDetail:  procDetail,
		Msg:            cm,
		PathDecoration: g.GetNamedTaskDecoration(checking, checker),
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

	questions := []string{}
	if !g.SkipQuestionFile {
		qfile := filepath.Join(g.GetExamDir(exam, config), "questions.csv")
		qbytes, err := ioutil.ReadFile(qfile)
		if err == nil {
			questions = strings.Split(string(qbytes), ",")
			logger.Info().
				Str("UUID", procDetail.UUID).
				Str("checker", checker).
				Str("exam", exam).
				Str("file", qfile).
				Str("questions", string(qbytes)).
				Msg("Got questions for cover page")
		} else {
			logger.Info().
				Str("UUID", procDetail.UUID).
				Str("checker", checker).
				Str("exam", exam).
				Str("file", qfile).
				Str("questions", string(qbytes)).
				Msg("Error opening questions file for cover page")
			return fmt.Errorf("Error opening questions file %s for cover page", qfile)
		}
	}

	fmt.Printf("Questions: %s\n", strings.Join(questions, ","))

	cp := CoverPageCommand{
		Questions:      questions,
		FromPath:       g.GetExamDir(exam, enterProcessed),
		ToPath:         g.GetExamDir(exam, checkerCover),
		ExamName:       exam,
		TemplatePath:   g.OverlayLayoutSVG(),
		SpreadName:     "addition",
		ProcessDetail:  procDetail,
		PathDecoration: "-cover",
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
		CoverPath:      g.GetExamDir(exam, checkerCover),
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
