package ingester

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/timdrysdale/gradex-cli/parselearn"
)

func (g *Ingester) ValidateNewPapers() error {

	logger := g.logger.With().Str("process", "validate-new-papers").Logger()

	// wait for user to press an "do import new scripts button", then check the temp-txt and temp-pdf dirs
	possibleReceipts, err := g.GetFileList(g.TempTXT())
	if err != nil {
		logger.Error().
			Str("source", g.TempTXT()).
			Msg("Could not get list of possible receipts")
		return err
	}

	// Map receipts, keeping only the latest revision for any given filename, ignoring dir and ext
	// so as to capture files in different dirs e.g. patch dirs, and with renamed filetypes
	receiptMap := make(map[string]parselearn.Submission)

	for _, receipt := range possibleReceipts {

		sub, err := parselearn.ParseLearnReceipt(receipt)

		if err != nil {
			logger.Error().
				Str("file", receipt).
				Msg("Did not parse as a Learn receipt")
			continue // assume there may be others uses for txt, and that clean up will happen at end of the ingest
		}

		if existingSub, ok := receiptMap[fileKey(sub.Filename)]; ok {
			if sub.Revision > existingSub.Revision {
				receiptMap[fileKey(sub.Filename)] = sub
			}
		} else {
			receiptMap[fileKey(sub.Filename)] = sub
		}

	}

	// >>>>>>>>>>>>> drop IGNORE receipts >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
	parselearn.HandleIgnoreReceipts(&receiptMap)

	// >>>>>>>>>>>>>>> drop multiple file submissions>>>>>>>>>>>>>>>>>>>>
	// look for, and reject, any multiple file submissions
	// these need flattening before merging so automatic merging
	// is a TODO - automatic flatten and merge multple  pdf submission
	// these need explicit patching to distinguish between us just taking
	// the first file before it is merged, with taking a merged file named
	// the same as the first file - with a patch receipt and a manual merge
	// we can name it what we like
	for k, v := range receiptMap {
		if v.NumberOfFiles > 1 {
			list, err := parselearn.GetFilePaths(v.OwnPath)
			if err != nil {
				logger.Error().
					Str("receipt", v.OwnPath).
					Msg("Could not check for multiple files in subsmission")
				continue
			}
			if len(list) > 1 {
				logger.Error().
					Str("receipt", v.OwnPath).
					Str("files", strings.Join(list, ";")).
					Msg("Rejecting because need to merge the submission into one file")
				delete(receiptMap, k)
			}
		}
	}

	for _, sub := range receiptMap {

		// assume we want to process this exam at some point - so set up the structure now
		// if it does not exist already
		_, err = os.Stat(g.GetExamPath(sub.Assignment))
		if os.IsNotExist(err) {
			err = g.SetupExamPaths(sub.Assignment)
			if err != nil {
				g.logger.Error().
					Str("course", sub.Assignment).
					Msg("Could not ensure directory structure was set up. Yikes, disk full? Bailing out!")
				return err // If we can't set up a new exam, we may as well bail out
			}
		}

		pdfFilename, err := GetPDFPath(sub.Filename, g.TempPDF())
		if err != nil {
			logger.Error().
				Str("file", sub.Filename).
				Str("location", g.TempPDF()).
				Msg("Error figuring out PDF filename, skipping this submission")
			continue
		}

		// file we want to get from the temp-pdf dir
		currentPath := filepath.Join(g.TempPDF(), filepath.Base(pdfFilename))
		destinationDir := g.AcceptedPapers(sub.Assignment)

		baseFileName := filepath.Base(pdfFilename)
		ShortLearnName := regexp.MustCompile("(\\_.*\\_{1})")
		//Before: PGEEnnnn A Super Long Exam Name - Exam Dropbox_s0000000_attempt_2020-05-01-02-00-00_PGEEnnnn-B000000.pdf
		//After _s0000000_attempt_2020-05-01-02-00-00_
		ShortLearnName := LearnOnly.FindString(baseFileName)
		destination := filepath.Join(destinationDir, ShortLearnName+filepath.Ext(pdfFilename))

		logger.Info().
			Str("before", baseFileName).
			Str("after", ShortLearnLearn).
			Msg("Using LEARN-specific name shortener")

		_, err = os.Stat(currentPath)

		if !os.IsNotExist(err) { //PDF file exists, move it to accepted papers

			//moved, err := g.MoveIfNewerThanDestinationInDir(currentPath, destination, &logger)
			moved, err := g.MoveIfNewerThanDestination(currentPath, destination, &logger)

			switch {

			case err == nil && moved:

				g.logger.Info().
					Str("file", currentPath).
					Str("course", sub.Assignment).
					Str("destination", destination).
					Msg("PDF validated and moved to accepted papers")

				destination := g.AcceptedReceipts(sub.Assignment)

				// this is not move-if-newer because it should match the pdf?
				err = g.MoveToDir(sub.OwnPath,
					destination)

				if err == nil {
					g.logger.Info().
						Str("file", sub.OwnPath).
						Str("course", sub.Assignment).
						Str("destination", destination).
						Msg("Moved receipt to Accepted Receipts")

				} else {
					g.logger.Error().
						Str("file", sub.OwnPath).
						Str("course", sub.Assignment).
						Str("destination", destination).
						Str("error", err.Error()).
						Msg("Could not put receipt in Accepted Receipts")
				}

			case err == nil && !moved:

				err := os.Remove(currentPath)

				if err == nil {
					g.logger.Info().
						Str("file", currentPath).
						Str("course", sub.Assignment).
						Str("destination", destination).
						Msg("PDF validated but TOO OLD; deleted")
				} else {
					g.logger.Error().
						Str("file", currentPath).
						Str("course", sub.Assignment).
						Str("destination", destination).
						Str("error", err.Error()).
						Msg("PDF validated but TOO OLD, and error deleting. Sigh. Over to you, human")
				}

			case err != nil:

				g.logger.Error().
					Str("file", currentPath).
					Str("course", sub.Assignment).
					Str("destination", destination).
					Str("error", err.Error()).
					Msg("PDF validated but ERROR prevented attempted move to accepted papers, returning to ingest")

				destination := g.Ingest()

				err := g.MoveToDir(currentPath, destination)

				if err != nil {
					g.logger.Error().
						Str("file", currentPath).
						Str("course", sub.Assignment).
						Str("destination", destination).
						Str("error", err.Error()).
						Msg("Could not return PDF to ingest.")
				}

				err = g.MoveToDir(sub.OwnPath, destination)

				if err != nil {
					g.logger.Error().
						Str("file", sub.OwnPath).
						Str("course", sub.Assignment).
						Str("destination", destination).
						Str("error", err.Error()).
						Msg("Could not return receipt to ingest")
				}
			} //switch

		} else {
			logger.Error().
				Str("file", currentPath).
				Str("course", sub.Assignment).
				Str("destination", destination).
				Str("error", err.Error()).
				Msg("Could not find this PDF to put in accepted papers.")
		}

	}

	// reject back to ingest anything we didn't take further
	rejectPDF, err := g.GetFileList(g.TempPDF())

	for _, reject := range rejectPDF {

		err := g.MoveToDir(reject, g.Ingest())

		if err == nil {

			g.logger.Info().
				Str("file", reject).
				Str("destination", g.Ingest()).
				Msg("PDF rejected")
		} else {
			g.logger.Error().
				Str("file", reject).
				Str("destination", g.Ingest()).
				Str("error", err.Error()).
				Msg("PDF rejectd, but ERROR returning to ingest")
		}

	}

	rejectTXT, err := g.GetFileList(g.TempTXT())

	for _, reject := range rejectTXT {
		err := g.MoveToDir(reject, g.Ingest())

		if err == nil {

			g.logger.Info().
				Str("file", reject).
				Str("destination", g.Ingest()).
				Msg("TXT rejected")
		} else {
			g.logger.Error().
				Str("file", reject).
				Str("destination", g.Ingest()).
				Str("error", err.Error()).
				Msg("TXT rejectd, but ERROR returning to ingest")
		}

	}

	return nil
}
