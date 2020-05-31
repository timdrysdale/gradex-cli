package ingester

import (
	"fmt"
	"strings"
)

var (
	exportStages = []string{
		"labelling",
		"marking",
		"remarking",
		"moderating",
		"remoderating",
		"entering",
		"reentering",
		"checking",
		"rechecking",
	}
)

func ValidStageForExport(stage string) bool {

	switch strings.ToLower(stage) {

	case labelling, marking, remarking, moderating, remoderating, entering, reentering, checking, rechecking:
		return true
	default:
		return false
	}
}

func (g *Ingester) GetExportDirs(exam, stage, actor string) (string, string, string, error) {

	var ready, sent string

	switch stage {
	case labelling:
		ready = questionReady
		sent = questionSent
	case marking:
		ready = markerReady
		sent = markerSent

	case remarking:
		ready = reMarkerReady
		sent = reMarkerSent

	case moderating:
		ready = moderatorReady
		sent = moderatorSent

	case remoderating:
		ready = reModeratorReady
		sent = reModeratorSent

	case entering:
		ready = enterReady
		sent = enterSent

	case reentering:
		ready = reEnterReady
		sent = reEnterSent

	case checking:
		ready = checkerReady
		sent = checkerSent

	case rechecking:
		ready = reCheckerReady
		sent = reCheckerSent

	default:
		return "", "", "", fmt.Errorf("unknown stage %s.\n Try: [%s]", stage, strings.Join(exportStages, ","))
	}

	readyDir := g.GetExamDirNamed(exam, ready, actor)
	sentDir := g.GetExamDirNamed(exam, sent, actor)
	exportDir := g.GetExportDir(exam, stage, actor)

	g.EnsureDirAll(readyDir)
	g.EnsureDirAll(sentDir)
	g.EnsureDirAll(exportDir)

	return readyDir, sentDir, exportDir, nil
}

func (g *Ingester) ExportFiles(exam, stage, actor string) error {

	logger := g.logger.With().
		Str("exam", exam).
		Str("stage", stage).
		Str("actor", actor).Logger()

	logger.Info().Msg("Exporting")

	readyDir, sentDir, exportDir, err := g.GetExportDirs(exam, stage, actor)

	if err != nil {
		logger.Error().
			Str("error", err.Error()).
			Msg("Could not get export directories")
		return err
	}

	logger = logger.With().
		Str("ready", readyDir).
		Str("sent", sentDir).
		Str("export", exportDir).
		Logger()

	files, err := GetFileList(readyDir)

	if err != nil {

		logger.Error().
			Str("error", err.Error()).
			Msg("could not get file list for exporting")
	}

	numErrors := 0

	for _, file := range files {

		if !g.IsPDF(file) {
			continue
		}

		_, err := g.CopyIfNewerThanDestinationInDir(file, exportDir, &logger)

		if err == nil {

			err = g.MoveToDir(file, sentDir)

			if err != nil {
				g.logger.Error().
					Str("error", err.Error()).
					Str("file", file).
					Msg("Could not move file to sent directory")
			}

		} else {
			numErrors++
			g.logger.Error().
				Str("file", file).
				Str("error", err.Error()).
				Msg("could not copy file to export directory")

		}

	}
	if numErrors == 0 {
		g.logger.Info().
			Int("count", len(files)).
			Msg(fmt.Sprintf("Exported %d files", len(files)))

		return nil
	}

	return fmt.Errorf("%d errors in exporting - see logfile for details", numErrors)
}
