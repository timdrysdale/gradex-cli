package ingester

import (
	"fmt"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func (g *Ingester) ExportForMarking(exam, marker string, logger *zerolog.Logger) {

	source := g.MarkerReady(exam, marker)

	logger.Info().
		Str("process", "export").
		Str("source", source).
		Msg("Exporting Marking")

	files, err := GetFileList(source)

	if err != nil {

		log.Error().
			Str("process", "export").
			Str("course", exam).
			Str("marker", marker).
			Str("error", err.Error()).
			Msg("could not get file list for exporting")
	}

	numErrors := 0

	destination := g.ExportMarking(exam, marker)

	for _, file := range files {

		if !g.IsPDF(file) {
			continue
		}

		destination = g.ExportMarking(exam, marker)

		err := g.CopyToDir(file, destination)
		if err == nil {

			destination = g.MarkerSent(exam, marker)

			err = g.MoveToDir(file, destination)

			if err != nil {
				logger.Error().
					Str("process", "export").
					Str("course", exam).
					Str("marker", marker).
					Str("file", file).
					Str("destination", destination).
					Str("error", err.Error()).
					Msg("could not copy file to MarkerSent")

			}

		} else {
			numErrors++
			logger.Error().
				Str("process", "export").
				Str("course", exam).
				Str("marker", marker).
				Str("file", file).
				Str("destination", destination).
				Str("error", err.Error()).
				Msg("could not copy file to export it")

		}

	}
	if numErrors == 0 {
		logger.Info().
			Str("process", "export").
			Str("course", exam).
			Str("marker", marker).
			Int("count", len(files)).
			Str("destination", destination).
			Msg(fmt.Sprintf("Exported %d files to %s", len(files), destination))
	}

}

func (g *Ingester) ExportForModerating(exam, moderator string, logger *zerolog.Logger) {

	files, err := GetFileList(g.ModeratorReady(exam, moderator))

	if err != nil {

		log.Error().
			Str("process", "export").
			Str("course", exam).
			Str("moderator", moderator).
			Str("error", err.Error()).
			Msg("could not get file list for exporting")
	}

	numErrors := 0

	destination := g.ExportModerating(exam, moderator)

	for _, file := range files {

		if !g.IsPDF(file) {
			continue
		}

		destination = g.ExportModerating(exam, moderator)

		err := g.CopyToDir(file, destination)
		if err == nil {

			destination = g.ModeratorSent(exam, moderator)
			err = g.MoveToDir(file, destination)
			if err != nil {
				logger.Error().
					Str("process", "export").
					Str("course", exam).
					Str("moderator", moderator).
					Str("file", file).
					Str("destination", destination).
					Str("error", err.Error()).
					Msg("could not copy file to ModeratorSent")

			}

		} else {
			numErrors++
			logger.Error().
				Str("process", "export").
				Str("course", exam).
				Str("moderator", moderator).
				Str("file", file).
				Str("destination", destination).
				Str("error", err.Error()).
				Msg("could not copy file to export it")

		}
	}
	if numErrors == 0 {
		logger.Info().
			Str("process", "export").
			Str("course", exam).
			Str("moderator", moderator).
			Int("count", len(files)).
			Str("destination", destination).
			Msg(fmt.Sprintf("Exported %d files to %s", len(files), destination))
	}

}

func (g *Ingester) ExportForChecking(exam, checker string, logger *zerolog.Logger) {

	files, err := GetFileList(g.CheckerReady(exam, checker))
	if err != nil {

		log.Error().
			Str("process", "export").
			Str("course", exam).
			Str("checker", checker).
			Str("error", err.Error()).
			Msg("could not get file list for exporting")
	}

	numErrors := 0

	destination := g.ExportChecking(exam, checker)

	for _, file := range files {

		if !g.IsPDF(file) {
			continue
		}
		destination = g.ExportChecking(exam, checker)
		err := g.CopyToDir(file, destination)
		if err == nil {

			destination = g.CheckerSent(exam, checker)
			err = g.MoveToDir(file, destination)
			if err != nil {
				logger.Error().
					Str("process", "export").
					Str("course", exam).
					Str("checker", checker).
					Str("file", file).
					Str("destination", destination).
					Str("error", err.Error()).
					Msg("could not copy file to CheckerSent")

			}

		} else {
			numErrors++
			logger.Error().
				Str("process", "export").
				Str("course", exam).
				Str("checker", checker).
				Str("file", file).
				Str("destination", destination).
				Str("error", err.Error()).
				Msg("could not copy file to export it")

		}
	}
	if numErrors == 0 {
		logger.Info().
			Str("process", "export").
			Str("course", exam).
			Str("checker", checker).
			Int("count", len(files)).
			Str("destination", destination).
			Msg(fmt.Sprintf("Exported %d files to %s", len(files), destination))
	}

}

func (g *Ingester) ExportForReMarking(exam, marker string, logger *zerolog.Logger) {

	files, err := GetFileList(g.ReMarkerReady(exam, marker))
	if err != nil {

		log.Error().
			Str("process", "export").
			Str("course", exam).
			Str("remarker", marker).
			Str("error", err.Error()).
			Msg("could not get file list for exporting")
	}
	numErrors := 0

	destination := g.ExportReMarking(exam, marker)

	for _, file := range files {

		if !g.IsPDF(file) {
			continue
		}
		destination = g.ExportReMarking(exam, marker)
		err := g.CopyToDir(file, destination)
		if err == nil {
			destination = g.ReMarkerSent(exam, marker)

			err = g.MoveToDir(file, destination)
			if err != nil {
				logger.Error().
					Str("process", "export").
					Str("course", exam).
					Str("marker", marker).
					Str("file", file).
					Str("destination", destination).
					Str("error", err.Error()).
					Msg("could not copy file to ReMarkerSent")

			}

		} else {
			numErrors++
			logger.Error().
				Str("course", exam).
				Str("marker", marker).
				Str("file", file).
				Str("destination", destination).
				Str("error", err.Error()).
				Msg("could not copy file to export it")

		}
	}
	if numErrors == 0 {
		logger.Info().
			Str("course", exam).
			Str("marker", marker).
			Int("count", len(files)).
			Str("destination", destination).
			Msg(fmt.Sprintf("Exported %d files to %s", len(files), destination))
	}

}

func (g *Ingester) ExportForReChecking(exam, checker string, logger *zerolog.Logger) {

	files, err := GetFileList(g.ReCheckerReady(exam, checker))
	if err != nil {

		log.Error().
			Str("process", "export").
			Str("course", exam).
			Str("rechecker", checker).
			Str("error", err.Error()).
			Msg("could not get file list for exporting")
	}
	numErrors := 0

	destination := g.ExportReChecking(exam, checker)

	for _, file := range files {

		if !g.IsPDF(file) {
			continue
		}
		destination = g.ExportReChecking(exam, checker)

		err := g.CopyToDir(file, destination)
		if err == nil {

			destination = g.CheckerSent(exam, checker)
			err = g.MoveToDir(file, destination)
			if err != nil {
				logger.Error().
					Str("process", "export").
					Str("course", exam).
					Str("checker", checker).
					Str("file", file).
					Str("destination", destination).
					Str("error", err.Error()).
					Msg("could not copy file to CheckerSent")

			}

		} else {
			numErrors++
			logger.Error().
				Str("process", "export").
				Str("course", exam).
				Str("checker", checker).
				Str("file", file).
				Str("destination", destination).
				Str("error", err.Error()).
				Msg("could not copy file to export it")

		}
	}

	if numErrors == 0 {
		logger.Info().
			Str("course", exam).
			Str("checker", checker).
			Int("count", len(files)).
			Str("destination", destination).
			Msg(fmt.Sprintf("Exported %d files to %s", len(files), destination))
	}

}
