package ingester

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/timdrysdale/chmsg"
	"github.com/timdrysdale/pdfpagedata"
)

func (g *Ingester) AddMarkBar(exam string, marker string) error {

	logger := g.logger.With().Str("process", "add-mark-bar").Logger()

	mc := chmsg.MessagerConf{
		ExamName:     exam,
		FunctionName: "overlay",
		TaskName:     "add-mark-bar",
	}

	cm := chmsg.New(mc, g.msgCh, g.timeout)

	var UUIDBytes uuid.UUID

	UUIDBytes, err := uuid.NewRandom()
	uuidStr := UUIDBytes.String()
	if err != nil {
		uuidStr = fmt.Sprintf("%d", time.Now().UnixNano())
	}

	procDetails := pdfpagedata.ProcessingDetails{
		UUID:     uuidStr,
		Previous: "", //dynamic
		UnixTime: time.Now().UnixNano(),
		Name:     "mark-bar",
		By:       pdfpagedata.ContactDetails{Name: "ingester"},
		Sequence: 0, //dynamic
	}

	UUIDBytes, err = uuid.NewRandom()
	uuidStr = UUIDBytes.String()
	if err != nil {
		uuidStr = fmt.Sprintf("%d", time.Now().UnixNano())
	}

	markDetails := pdfpagedata.QuestionDetails{
		UUID:     uuidStr,
		Name:     "marking",
		UnixTime: time.Now().UnixNano(),
		Marking: []pdfpagedata.MarkingAction{
			pdfpagedata.MarkingAction{
				Actor: marker,
			},
		},
	}

	oc := OverlayCommand{
		PreparedFor:       marker,
		ToDo:              "marking",
		FromPath:          g.AnonymousPapers(exam),
		ToPath:            g.MarkerReady(exam, marker),
		ExamName:          exam,
		TemplatePath:      g.OverlayLayoutSVG(),
		SpreadName:        "mark",
		ProcessingDetails: procDetails,
		QuestionDetails:   markDetails,
		Msg:               cm,
		PathDecoration:    g.MarkerABCDecoration(marker),
	}

	err = g.OverlayPapers(oc, &logger)

	if err == nil {
		cm.Send(fmt.Sprintf("Finished Processing markbar UUID=%s\n", uuidStr))
		logger.Info().
			Str("UUID", uuidStr).
			Str("marker", marker).
			Str("exam", exam).
			Msg("Finished add-mark-bar")
	} else {
		logger.Error().
			Str("UUID", uuidStr).
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

	var UUIDBytes uuid.UUID

	UUIDBytes, err := uuid.NewRandom()
	uuidStr := UUIDBytes.String()
	if err != nil {
		uuidStr = fmt.Sprintf("%d", time.Now().UnixNano())
	}

	procDetails := pdfpagedata.ProcessingDetails{
		UUID:     uuidStr,
		Previous: "", //dynamic
		UnixTime: time.Now().UnixNano(),
		Name:     "moderate-active-bar",
		By:       pdfpagedata.ContactDetails{Name: "ingester"},
		Sequence: 0, //dynamic
	}

	UUIDBytes, err = uuid.NewRandom()
	uuidStr = UUIDBytes.String()
	if err != nil {
		uuidStr = fmt.Sprintf("%d", time.Now().UnixNano())
	}

	markDetails := pdfpagedata.QuestionDetails{
		UUID:     uuidStr,
		Name:     "moderating",
		UnixTime: time.Now().UnixNano(),
		Moderating: []pdfpagedata.MarkingAction{
			pdfpagedata.MarkingAction{
				Actor: moderator,
			},
		},
	}

	oc := OverlayCommand{
		PreparedFor:       moderator,
		ToDo:              "moderating",
		FromPath:          g.ModerateActive(exam),
		ToPath:            g.ModeratorReady(exam, moderator),
		ExamName:          exam,
		TemplatePath:      g.OverlayLayoutSVG(),
		SpreadName:        "moderate-active",
		ProcessingDetails: procDetails,
		QuestionDetails:   markDetails,
		Msg:               cm,
		PathDecoration:    g.ModeratorABCDecoration(moderator),
	}

	err = g.OverlayPapers(oc, &logger)

	if err == nil {
		cm.Send(fmt.Sprintf("Finished Processing add-moderate-active-bar UUID=%s\n", uuidStr))
		logger.Info().
			Str("UUID", uuidStr).
			Str("moderator", moderator).
			Str("exam", exam).
			Msg("Finished add-moderate-active-bar")
	} else {
		logger.Error().
			Str("UUID", uuidStr).
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

	var UUIDBytes uuid.UUID

	UUIDBytes, err := uuid.NewRandom()
	uuidStr := UUIDBytes.String()
	if err != nil {
		uuidStr = fmt.Sprintf("%d", time.Now().UnixNano())
	}

	procDetails := pdfpagedata.ProcessingDetails{
		UUID:     uuidStr,
		Previous: "", //dynamic
		UnixTime: time.Now().UnixNano(),
		Name:     "moderate-inactive-bar",
		By:       pdfpagedata.ContactDetails{Name: "ingester"},
		Sequence: 0, //dynamic
	}

	UUIDBytes, err = uuid.NewRandom()
	uuidStr = UUIDBytes.String()
	if err != nil {
		uuidStr = fmt.Sprintf("%d", time.Now().UnixNano())
	}

	markDetails := pdfpagedata.QuestionDetails{
		UUID:     uuidStr,
		Name:     "moderating",
		UnixTime: time.Now().UnixNano(),
		Moderating: []pdfpagedata.MarkingAction{
			pdfpagedata.MarkingAction{
				Actor: "none",
			},
		},
	}

	oc := OverlayCommand{
		PreparedFor:       "",
		ToDo:              "moderating",
		FromPath:          g.ModerateInActive(exam),
		ToPath:            g.ModeratedInActiveBack(exam),
		ExamName:          exam,
		TemplatePath:      g.OverlayLayoutSVG(),
		SpreadName:        "moderate-inactive",
		ProcessingDetails: procDetails,
		QuestionDetails:   markDetails,
		Msg:               cm,
		PathDecoration:    g.ModeratorABCDecoration("X"),
	}

	err = g.OverlayPapers(oc, &logger)

	if err == nil {
		cm.Send(fmt.Sprintf("Finished Processing add-moderate-inactive-bar UUID=%s\n", uuidStr))
		logger.Info().
			Str("UUID", uuidStr).
			Str("exam", exam).
			Msg("Finished add-moderate-inactive-bar")
	} else {
		logger.Error().
			Str("UUID", uuidStr).
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

	var UUIDBytes uuid.UUID

	UUIDBytes, err := uuid.NewRandom()
	uuidStr := UUIDBytes.String()
	if err != nil {
		uuidStr = fmt.Sprintf("%d", time.Now().UnixNano())
	}

	procDetails := pdfpagedata.ProcessingDetails{
		UUID:     uuidStr,
		Previous: "", //dynamic
		UnixTime: time.Now().UnixNano(),
		Name:     "check-bar",
		By:       pdfpagedata.ContactDetails{Name: "ingester"},
		Sequence: 0, //dynamic
	}

	UUIDBytes, err = uuid.NewRandom()
	uuidStr = UUIDBytes.String()
	if err != nil {
		uuidStr = fmt.Sprintf("%d", time.Now().UnixNano())
	}

	markDetails := pdfpagedata.QuestionDetails{
		UUID:     uuidStr,
		Name:     "checking",
		UnixTime: time.Now().UnixNano(),
		Checking: []pdfpagedata.MarkingAction{
			pdfpagedata.MarkingAction{
				Actor: checker,
			},
		},
	}

	oc := OverlayCommand{
		PreparedFor:       checker,
		ToDo:              "checking",
		FromPath:          g.ModeratedReady(exam),
		ToPath:            g.CheckerReady(exam, checker),
		ExamName:          exam,
		TemplatePath:      g.OverlayLayoutSVG(),
		SpreadName:        "check",
		ProcessingDetails: procDetails,
		QuestionDetails:   markDetails,
		Msg:               cm,
		PathDecoration:    g.CheckerABCDecoration(checker),
	}

	err = g.OverlayPapers(oc, &logger)

	if err == nil {
		cm.Send(fmt.Sprintf("Finished Processing add-check-bar UUID=%s\n", uuidStr))
		logger.Info().
			Str("UUID", uuidStr).
			Str("checker", checker).
			Str("exam", exam).
			Msg("Finished add-check-bar")
	} else {
		logger.Error().
			Str("UUID", uuidStr).
			Str("checker", checker).
			Str("exam", exam).
			Str("error", err.Error()).
			Msg("Error add-checker-bar")
	}
	return err

}
