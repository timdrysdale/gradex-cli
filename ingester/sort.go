package ingester

import (
	"path/filepath"

	"github.com/timdrysdale/gradex-cli/extract"
	"github.com/timdrysdale/gradex-cli/merge"
	"github.com/timdrysdale/gradex-cli/pagedata"
)

// get all files in QuestionBack and split according to labelling into
// subfolders in QuestionSplit
// then apply a batching tool to prepare marking sets
// we only re-sort Q if there is no done file

func (g *Ingester) SortQuestions(exam string) error {

	files, err := GetFileList(g.QuestionBack(exam, ""))

	if err != nil {

		g.logger.Error().
			Msg("Error getting question back files")
	}

	qfm := make(map[string][]string)

	pagebad := []string{}

	for _, file := range files {

		pdm, err := pagedata.UnMarshalAllFromFile(file) //(map[int]PageData, error)
		if err != nil {
			g.logger.Error().
				Str("file", file).
				Msg("Error getting pagedata")
			continue
		}
		tfm, err := extract.ExtractTextFieldsFromPDF(file) //(map[int]map[string]string, error)
		if err != nil {
			g.logger.Error().
				Str("file", file).
				Msg("Error getting textfields")
			continue
		}

		//if Q to be marked, get Current.Own.Path...

		for k, v := range tfm {

			if pd, ok := pdm[k]; ok {

				thisFile := pd.Current.Own.Path
				Q1 := v["question-first-number"]
				Q1skip := v["question-first-skip"]
				Q2 := v["question-second-number"]
				Q2skip := v["question-second-skip"]

				if Q1skip != "SKIP" {
					if Q1 != "" {
						qfm[Q1] = append(qfm[Q1], thisFile)
					}
				}
				if Q2skip != "SKIP" {
					if Q2 != "" {
						qfm[Q2] = append(qfm[Q2], thisFile)
					}
				}
				if v["page-bad"] != "" {
					pagebad = append(pagebad, file)
				}

			}

		}

		for question, files := range qfm {

			destination := filepath.Join(g.QuestionSplit(exam, ""), question+".pdf")
			merge.PDF(files, destination)

			//for _, file := range files {
			//err = g.CopyToDir(file, g.QuestionSplit(exam, question))
			//if err != nil {
			//	g.logger.Error().
			//		Str("file", file).
			//		Str("question", question).
			//		Str("error", err.Error()).
			//		Msg("Error copying file to question directory")
			//}
			//}

		}

		for _, file := range pagebad {
			err = g.CopyToDir(file, g.PageBad(exam))
			if err != nil {
				g.logger.Error().
					Str("file", file).
					Str("error", err.Error()).
					Msg("Error copying file to pageBad directory")
			}
		}

	}
	return nil
}
