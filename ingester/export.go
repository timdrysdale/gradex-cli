package ingester

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

func (g *Ingester) ExportForLabelling(exam, labeller string) {

	source := g.GetExamDirNamed(exam, questionReady, labeller)

	g.logger.Info().Msg("Exporting")

	files, err := GetFileList(source)

	if err != nil {

		g.logger.Error().
			Str("error", err.Error()).
			Msg("could not get file list for exporting")
	}

	g.logger.Info().Msg("Exporting")

	numErrors := 0

	destination := g.GetExportDir(exam, questionReady, labeller)

	for _, file := range files {

		if !g.IsPDF(file) {
			continue
		}

		destination := g.GetExportDir(exam, questionReady, labeller)

		err := g.CopyToDir(file, destination)
		if err == nil {

			destination = g.GetExamDirNamed(exam, questionSent, labeller)

			err = g.MoveToDir(file, destination)

			if err != nil {
				g.logger.Error().
					Str("file", file).
					Str("destination", destination).
					Str("error", err.Error()).
					Msg("could not copy file to LabellerSent")

			}

		} else {
			numErrors++
			g.logger.Error().
				Str("file", file).
				Str("destination", destination).
				Str("error", err.Error()).
				Msg("could not copy file to export it")

		}

	}
	if numErrors == 0 {
		g.logger.Info().
			Int("count", len(files)).
			Str("destination", destination).
			Msg(fmt.Sprintf("Exported %d files to %s", len(files), destination))
	}

}

func (g *Ingester) ExportForMarking(exam, marker string) {

	source := g.GetExamDirNamed(exam, markerReady, marker)

	files, err := GetFileList(source)

	if err != nil {

		g.logger.Error().
			Str("error", err.Error()).
			Msg("could not get file list for exporting")
	}

	g.logger.Info().Msg("Exporting")

	numErrors := 0

	destination := g.GetExportDir(exam, markerReady, marker)

	for _, file := range files {

		if !g.IsPDF(file) {
			continue
		}

		destination := g.GetExportDir(exam, markerReady, marker)

		err := g.CopyToDir(file, destination)
		if err == nil {

			destination = g.GetExamDirNamed(exam, markerSent, marker)

			err = g.MoveToDir(file, destination)

			if err != nil {
				g.logger.Error().
					Str("file", file).
					Str("destination", destination).
					Str("error", err.Error()).
					Msg("could not copy file to MarkerSent")

			}

		} else {
			numErrors++
			g.logger.Error().
				Str("file", file).
				Str("destination", destination).
				Str("error", err.Error()).
				Msg("could not copy file to export it")

		}

	}
	if numErrors == 0 {
		g.logger.Info().
			Int("count", len(files)).
			Str("destination", destination).
			Msg(fmt.Sprintf("Exported %d files to %s", len(files), destination))
	}

}

func (g *Ingester) ExportForModerating(exam, moderator string) {

	files, err := GetFileList(g.GetExamDirNamed(exam, moderatorReady, moderator))

	if err != nil {

		g.logger.Error().
			Str("error", err.Error()).
			Msg("could not get file list for exporting")
	}
	g.logger.Info().Msg("Exporting")
	numErrors := 0

	destination := g.GetExportDir(exam, moderatorReady, moderator)

	for _, file := range files {

		if !g.IsPDF(file) {
			continue
		}
		destination := g.GetExportDir(exam, moderatorReady, moderator)

		err := g.CopyToDir(file, destination)
		if err == nil {

			destination = g.GetExamDirNamed(exam, moderatorSent, moderator)
			err = g.MoveToDir(file, destination)
			if err != nil {
				g.logger.Error().
					Str("file", file).
					Str("destination", destination).
					Str("error", err.Error()).
					Msg("could not copy file to ModeratorSent")

			}

		} else {
			numErrors++
			g.logger.Error().
				Str("file", file).
				Str("destination", destination).
				Str("error", err.Error()).
				Msg("could not copy file to export it")

		}
	}
	if numErrors == 0 {
		g.logger.Info().
			Int("count", len(files)).
			Str("destination", destination).
			Msg(fmt.Sprintf("Exported %d files to %s", len(files), destination))
	}

}

func (g *Ingester) ExportForReModerating(exam, moderator string) {

	files, err := GetFileList(g.GetExamDirNamed(exam, reModeratorReady, moderator))

	if err != nil {

		g.logger.Error().
			Str("error", err.Error()).
			Msg("could not get file list for exporting")
	}
	g.logger.Info().Msg("Exporting")
	numErrors := 0

	destination := g.GetExportDir(exam, reModeratorReady, moderator)

	for _, file := range files {

		if !g.IsPDF(file) {
			continue
		}
		destination := g.GetExportDir(exam, reModeratorReady, moderator)

		err := g.CopyToDir(file, destination)
		if err == nil {

			destination = g.GetExamDirNamed(exam, reModeratorSent, moderator)
			err = g.MoveToDir(file, destination)
			if err != nil {
				g.logger.Error().
					Str("file", file).
					Str("destination", destination).
					Str("error", err.Error()).
					Msg("could not copy file to reModeratorSent")

			}

		} else {
			numErrors++
			g.logger.Error().
				Str("file", file).
				Str("destination", destination).
				Str("error", err.Error()).
				Msg("could not copy file to export it")

		}
	}
	if numErrors == 0 {
		g.logger.Info().
			Int("count", len(files)).
			Str("destination", destination).
			Msg(fmt.Sprintf("Exported %d files to %s", len(files), destination))
	}

}

func (g *Ingester) ExportForChecking(exam, checker string) {

	files, err := GetFileList(g.GetExamDirNamed(exam, checkerReady, checker))
	if err != nil {

		g.logger.Error().
			Str("error", err.Error()).
			Msg("could not get file list for exporting")
	}

	g.logger.Info().Msg("Exporting")

	numErrors := 0

	destination := g.GetExportDir(exam, checkerReady, checker)

	for _, file := range files {

		if !g.IsPDF(file) {
			continue
		}
		destination := g.GetExportDir(exam, checkerReady, checker)

		err := g.CopyToDir(file, destination)
		if err == nil {

			destination = g.GetExamDirNamed(exam, checkerSent, checker)
			err = g.MoveToDir(file, destination)
			if err != nil {
				g.logger.Error().
					Str("file", file).
					Str("destination", destination).
					Str("error", err.Error()).
					Msg("could not copy file to CheckerSent")

			}

		} else {
			numErrors++
			g.logger.Error().
				Str("file", file).
				Str("destination", destination).
				Str("error", err.Error()).
				Msg("could not copy file to export it")

		}
	}
	if numErrors == 0 {
		g.logger.Info().
			Int("count", len(files)).
			Str("destination", destination).
			Msg(fmt.Sprintf("Exported %d files to %s", len(files), destination))
	}

}

func (g *Ingester) ExportForReMarking(exam, marker string) {

	files, err := GetFileList(g.GetExamDirNamed(exam, reMarkerReady, marker))
	if err != nil {

		log.Error().
			Str("error", err.Error()).
			Msg("could not get file list for exporting")
	}
	g.logger.Info().Msg("Exporting")
	numErrors := 0

	destination := g.GetExportDir(exam, reMarkerReady, marker)

	for _, file := range files {

		if !g.IsPDF(file) {
			continue
		}
		destination := g.GetExportDir(exam, reMarkerReady, marker)

		err := g.CopyToDir(file, destination)
		if err == nil {

			destination = g.GetExamDirNamed(exam, reMarkerSent, marker)

			err = g.MoveToDir(file, destination)
			if err != nil {
				g.logger.Error().
					Str("file", file).
					Str("destination", destination).
					Str("error", err.Error()).
					Msg("could not copy file to ReMarkerSent")

			}

		} else {
			numErrors++
			g.logger.Error().
				Str("file", file).
				Str("destination", destination).
				Str("error", err.Error()).
				Msg("could not copy file to export it")

		}
	}
	if numErrors == 0 {
		g.logger.Info().
			Int("count", len(files)).
			Str("destination", destination).
			Msg(fmt.Sprintf("Exported %d files to %s", len(files), destination))
	}

}

func (g *Ingester) ExportForReChecking(exam, checker string) {

	files, err := GetFileList(g.GetExamDirNamed(exam, reCheckerReady, checker))
	if err != nil {

		log.Error().
			Str("error", err.Error()).
			Msg("could not get file list for exporting")
	}
	g.logger.Info().Msg("Exporting")
	numErrors := 0

	destination := g.GetExportDir(exam, reCheckerReady, checker)

	for _, file := range files {

		if !g.IsPDF(file) {
			continue
		}
		destination := g.GetExportDir(exam, reCheckerReady, checker)

		err := g.CopyToDir(file, destination)
		if err == nil {

			destination = g.GetExamDirNamed(exam, reCheckerSent, checker)

			err = g.MoveToDir(file, destination)
			if err != nil {
				g.logger.Error().
					Str("file", file).
					Str("destination", destination).
					Str("error", err.Error()).
					Msg("could not copy file to CheckerSent")

			}

		} else {
			numErrors++
			g.logger.Error().
				Str("file", file).
				Str("destination", destination).
				Str("error", err.Error()).
				Msg("could not copy file to export it")

		}
	}

	if numErrors == 0 {
		g.logger.Info().
			Int("count", len(files)).
			Str("destination", destination).
			Msg(fmt.Sprintf("Exported %d files to %s", len(files), destination))
	}

}
