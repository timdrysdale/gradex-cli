package ingester

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog"
	"github.com/timdrysdale/gradex-cli/pagedata"
	"github.com/timdrysdale/gradex-cli/parselearn"
)

// wait for user to press an "do ingest button", then filewalk to get the paths
func (g *Ingester) StageFromIngest() error {

	ingestPath := g.Ingest()

	logger := g.logger.With().Str("process", "stage-from-ingest").Logger()

	logger.Info().Msg("STARTING INGEST")

	// TODO consider listing paths then moving....
	//pdfPaths := []string{}
	//txtPaths := []string{}

LOOP:
	for {
		passAgain := false

		err := filepath.Walk(ingestPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			switch {
			case g.IsArchive(path):
				//passAgain = true
				//g.handleIngestArchive(path, &logger)
				logger.Error().
					Str("file", path).
					Msg("Please extract files from the zip manually")

			case IsTXT(path):
				g.handleTXT(path, &logger)

			case IsPDF(path):
				g.handleIngestPDF(path, &logger)

			case IsCSV(path):
				g.handleIngestCSV(path, &logger)

			}
			return nil
		})

		if err != nil {
			logger.Error().
				Str("error", err.Error()).
				Msg("Filewalk failed")
		}

		if !passAgain {
			break LOOP
		}
	}

	// TODO check raw pdf?

	//TODO some reporting on what is left over? or another tool can do that?
	// and overall file system status tool?
	return nil
}

func (g *Ingester) handleTXT(path string, logger *zerolog.Logger) {

	err := parselearn.CheckFilename(path)

	if err != nil {
		g.logger.Warn().Str("file", path).Msg(fmt.Sprintf("Can't find the file listed in the receipt in the receipt %s", path))
	}

	moved, err := g.MoveIfNewerThanDestinationInDir(path, g.TempTXT(), logger)

	if err != nil {

		logger.Error().Str("file", path).Str("destination", g.TempTXT()).Msg("Could not move TXT file to TempTxt")

	} else {

		if moved {

			logger.Info().Str("file", path).Str("destination", g.TempTXT()).Msg("Moved TXT file to TempTxt")

		} else {

			logger.Info().Str("file", path).Str("destination", g.TempTXT()).Msg("Did not move TXT file to TempTxt (too old)")

		}
	}
}

func (g *Ingester) handleIngestCSV(path string, logger *zerolog.Logger) {

	if strings.ToLower(filepath.Base(path)) == "identity.csv" {

		moved, err := g.MoveIfNewerThanDestinationInDir(path, g.Identity(), logger)

		if err != nil {

			g.logger.Error().Str("file", path).Str("destination", g.Identity()).Msg("Couldn't move new identity csv into position")

		} else {
			if moved {
				g.logger.Info().Str("file", path).Str("destination", g.Identity()).Msg("New identity.csv installed")
			} else {
				g.logger.Info().Str("file", path).Str("destination", g.Identity()).Msg("identity.csv ignored (too old)")
			}

		}
	}
}

// leave file in ingest if not newer - to overwrite current file with an older version
// e.g. to roll back a change, you have to roll forward by modifying the old file,
// saving it to get a new modtime (can change back the mod before ingesting if needed)
// just need the new mod time on the file
func (g *Ingester) handleIngestPDF(path string, logger *zerolog.Logger) {

	ts, err := pagedata.TriageFile(path)

	if err != nil {
		// no page data so either a raw script, file from old gradex tool, or the pagedata has been corrupted
		// put in TempPDF in case it is raw script. If the other cases apply, it will ultimately be rejected
		// and we can have a human sort it from there (TODO pagedata injection tool for these repair jobs!)

		moved, err := g.MoveIfNewerThanDestinationInDir(path, g.TempPDF(), logger)

		if err != nil {

			g.logger.Error().Str("file", path).Str("destination", g.TempPDF()).Msg("Couldn't move raw PDF into TempPDF dir")

		} else {

			if moved {
				g.logger.Info().Str("file", path).Str("destination", g.TempPDF()).Msg("Moved raw PDF into TempPDF Dir")
			} else {
				g.logger.Info().Str("file", path).Str("destination", g.TempPDF()).Msg("Raw PDF NOT moved into TempPDF Dir (too old)")
			}

		}

		return
	}

	if len(ts) > 0 {

		logger.Info().
			Dict("properties", zerolog.Dict().
				Str("Is", ts[1].Is).
				Str("What", ts[1].What).
				Str("For", ts[1].For).
				Str("ToDo", ts[1].ToDo),
			).Msg("Identified a PDF with pagedata, for ingesting")
	}
	t := ts[1]
	switch t.ToDo {

	case "flattening":

		// these aren't usually exported, but we may be repopulating a new ingester or
		// manually correcting something, so we consider our options
		origin := g.GetExamDir(t.What, anonPapers)
		moved, err := g.MoveIfNewerThanDestinationInDir(path, origin, logger)
		if err != nil {
			g.logger.Error().Str("file", path).Str("destination", origin).Msg("Couldn't move flattened PDF into origin dir")
		} else {
			if moved {
				g.logger.Info().Str("file", path).Str("destination", origin).Msg("Moved raw PDF into origin Dir")
			} else {
				g.logger.Info().Str("file", path).Str("destination", origin).Msg("Raw PDF NOT moved into origin Dir (too old)")
			}
		}
		return

		// leave the file in ingest if we don't want it
	case "labelling":
		// these could be marked, or just being returned by DSA if prematurely exported
		origin := g.GetExamDirNamed(t.What, questionSent, t.For)

		preOrigin := g.GetExamDirNamed(t.What, questionReady, t.For)

		if g.IsSameAsSelfInDir(path, origin) {
			// put the file back in Ready (we keep this incoming version _just_in_case_ it had mods
			// despite having original time stamp and size!
			err := os.Rename(path, filepath.Join(preOrigin, filepath.Base(path)))
			if err != nil {
				return
			}

			// delete the version we had "sent" - this could be DSA re-ingesting exports before sending them
			err = os.Remove(filepath.Join(origin, filepath.Base(path)))
			if err != nil {
				return
			}
		} else {
			// it's (probably) been marked at least partly, so see if it is newer
			// than a version we might already have
			destination := g.GetExamDirNamed(t.What, questionBack, t.For)

			moved, err := g.MoveIfNewerThanDestinationInDir(path, destination, logger)

			switch {

			case err == nil && moved:

				g.logger.Info().
					Str("file", path).
					Str("destination", destination).
					Str("stage", "marking").
					Msg("PDF moved to QuestionBack because it looks labelled")

			case err == nil && !moved:

				err := os.Remove(path)

				if err == nil {
					g.logger.Info().
						Str("file", path).
						Str("destination", destination).
						Msg("PDF labelled but we have a newer labelled version, deleted")
				} else {
					g.logger.Error().
						Str("file", path).
						Str("destination", destination).
						Str("error", err.Error()).
						Msg("PDF labelled, but we have a newer labelled version, and ERROR deleting. Sigh. Over to you, human")
				}

			case err != nil:

				g.logger.Error().
					Str("file", path).
					Str("destination", destination).
					Str("error", err.Error()).
					Msg("PDF marked, but ERROR prevented attempted move to marked papers, returning to ingest")

				destination := g.Ingest()

				err := g.MoveToDir(path, destination)

				if err != nil {
					g.logger.Error().
						Str("file", path).
						Str("destination", destination).
						Str("error", err.Error()).
						Msg(fmt.Sprintf("Couldn't put in marked papers or return to ingest. Consider checking %s and moving as needed.", path))
				}

			} //switch

		}
		return
	case "marking":
		// these could be marked, or just being returned by DSA if prematurely exported
		origin := g.GetExamDirNamed(t.What, markerSent, t.For)

		preOrigin := g.GetExamDirNamed(t.What, markerReady, t.For)

		if g.IsSameAsSelfInDir(path, origin) {
			// put the file back in Ready (we keep this incoming version _just_in_case_ it had mods
			// despite having original time stamp and size!
			err := os.Rename(path, filepath.Join(preOrigin, filepath.Base(path)))
			if err != nil {
				return
			}

			// delete the version we had "sent" - this could be DSA re-ingesting exports before sending them
			err = os.Remove(filepath.Join(origin, filepath.Base(path)))
			if err != nil {
				return
			}
		} else {
			// it's (probably) been marked at least partly, so see if it is newer
			// than a version we might already have
			destination := g.GetExamDirNamed(t.What, markerBack, t.For)

			moved, err := g.MoveIfNewerThanDestinationInDir(path, destination, logger)

			switch {

			case err == nil && moved:

				g.logger.Info().
					Str("file", path).
					Str("destination", destination).
					Str("stage", "marking").
					Msg("PDF moved to MarkerBack because it looks marked")

			case err == nil && !moved:

				err := os.Remove(path)

				if err == nil {
					g.logger.Info().
						Str("file", path).
						Str("destination", destination).
						Msg("PDF marked but we have a newer marked version, deleted")
				} else {
					g.logger.Error().
						Str("file", path).
						Str("destination", destination).
						Str("error", err.Error()).
						Msg("PDF marked, but we have a newer marked version, and ERROR deleting. Sigh. Over to you, human")
				}

			case err != nil:

				g.logger.Error().
					Str("file", path).
					Str("destination", destination).
					Str("error", err.Error()).
					Msg("PDF marked, but ERROR prevented attempted move to marked papers, returning to ingest")

				destination := g.Ingest()

				err := g.MoveToDir(path, destination)

				if err != nil {
					g.logger.Error().
						Str("file", path).
						Str("destination", destination).
						Str("error", err.Error()).
						Msg("Could not return PDF to ingest.")
				}

			} //switch

		}
	case "moderating":

		origin := g.GetExamDirNamed(t.What, moderatorSent, t.For)

		preOrigin := g.GetExamDirNamed(t.What, moderatorReady, t.For)

		if g.IsSameAsSelfInDir(path, origin) {
			// put the file back in Ready (we keep this incoming version _just_in_case_ it had mods
			// despite having original time stamp and size!
			err := os.Rename(path, filepath.Join(preOrigin, filepath.Base(path)))
			if err != nil {
				return
			}

			// delete the version we had "sent" - this could be DSA re-ingesting exports before sending them
			err = os.Remove(filepath.Join(origin, filepath.Base(path)))
			if err != nil {
				return
			}
		} else {
			// it's (probably) been marked at least partly, so see if it is newer
			// than a version we might already have
			destination := g.GetExamDirNamed(t.What, moderatorBack, t.For)
			g.MoveIfNewerThanDestinationInDir(path, destination, logger)
			return
		}
	case "checking":

		origin := g.GetExamDirNamed(t.What, checkerSent, t.For)

		preOrigin := g.GetExamDirNamed(t.What, checkerReady, t.For)

		if g.IsSameAsSelfInDir(path, origin) {
			// put the file back in Ready (we keep this incoming version _just_in_case_ it had mods
			// despite having original time stamp and size!
			err := os.Rename(path, filepath.Join(preOrigin, filepath.Base(path)))
			if err != nil {
				return
			}

			// delete the version we had "sent" - this could be DSA re-ingesting exports before sending them
			err = os.Remove(filepath.Join(origin, filepath.Base(path)))
			if err != nil {
				return
			}
		} else {
			// it's (probably) been marked at least partly, so see if it is newer
			// than a version we might already have
			destination := g.GetExamDirNamed(t.What, checkerBack, t.For)
			g.MoveIfNewerThanDestinationInDir(path, destination, logger)
			return
		}
	default:
		// check later to see if it has a learn receipt, etc
		g.MoveIfNewerThanDestinationInDir(path, g.TempPDF(), logger)
		return

	}

	// errors.New("Didn't know how to handle pdf ingest")
	return
}
