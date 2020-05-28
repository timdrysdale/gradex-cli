package ingester

import (
	"fmt"
	"path/filepath"
	"sort"

	"github.com/timdrysdale/gradex-cli/extract"
	"github.com/timdrysdale/gradex-cli/merge"
	"github.com/timdrysdale/gradex-cli/pagedata"
	"github.com/timdrysdale/gradex-cli/parsesvg"
)

// get all files in QuestionBack and split according to labelling into
// subfolders in QuestionSplit
// then apply a batching tool to prepare marking sets
// we only re-sort Q if there is no done file

func (g *Ingester) SortQuestions(exam string) error {

	files, err := GetFileList(g.GetExamDir(exam, questionBack))

	if err != nil {

		g.logger.Error().
			Msg("Error getting question back files")
	}

	//question, who, pageNum, childfile
	qfm := make(map[string]map[string]map[int]string)

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

		for k, v := range tfm { //for each page of textfields

			if pd, ok := pdm[k]; ok { //get the page data if exists, or skip

				thisPath := pd.Current.Own.Path
				thisWho := pd.Current.Item.Who
				thisNumber := pd.Current.Own.Number

				Q1 := v["question-first-number"]
				Q1skip := v["question-first-skip"]
				Q2 := v["question-second-number"]
				Q2skip := v["question-second-skip"]

				if Q1skip != "SKIP" {
					if Q1 != "" {
						if _, ok := qfm[Q1]; !ok { //check Q map exists
							qfm[Q1] = make(map[string]map[int]string)
						}
						if _, ok := qfm[Q1][thisWho]; !ok { //check identity exists
							qfm[Q1][thisWho] = make(map[int]string)
						}
						qfm[Q1][thisWho][thisNumber] = thisPath //add this path
					}
				}
				if Q2skip != "SKIP" {
					if Q2 != "" {
						if _, ok := qfm[Q2]; !ok { //check Q map exists
							qfm[Q2] = make(map[string]map[int]string)
						}
						if _, ok := qfm[Q2][thisWho]; !ok { //check identity exists
							qfm[Q2][thisWho] = make(map[int]string)
						}
						qfm[Q2][thisWho][thisNumber] = thisPath //add this path
					}

				}
				if v["page-bad"] != "" {
					pagebad = append(pagebad, file)
				}

			}

		} //for tfm
	} //for file

	parsesvg.PrettyPrintStruct(qfm)

	batchSize := 20 //how many complete sets of answers to that question, per batch

	for question, fileMap := range qfm {
		//fmt.Println("QUESTION %s\n", question)
		var pagePaths []string

		batchCount := 0
		Qcount := 0

		for _, pageset := range fileMap {

			//fmt.Printf("WHO %s\n", who)

			pageNums := []int{}

			//parsesvg.PrettyPrintStruct(pageset)
			for N, _ := range pageset {
				//fmt.Printf("K:%v\n", K)
				//fmt.Printf("V:%v\n", V)
				pageNums = append(pageNums, N) // identify all the pagenumbers for this Q
			}
			//fmt.Println(pageNums)
			sort.Ints(pageNums) // sort page numbers into order
			//fmt.Println(pageNums)

			for _, N := range pageNums {
				pagePaths = append(pagePaths, pageset[N])
				fmt.Printf("%d: %s\n", N, pageset[N])
			}

			if Qcount >= batchSize { //if we've amassed enough answers to make a batch, save it

				filename := fmt.Sprintf("%s-%02d.pdf", question, batchCount)

				destination := filepath.Join(g.GetExamDir(exam, questionSplit), filename)

				merge.PDF(pagePaths, destination)

				Qcount = 0
				pagePaths = []string{}
				batchCount++

			} else {
				Qcount++
			}

		} // files for this question all done

		if Qcount > 0 { // catch the last partial batch of files, if any

			filename := fmt.Sprintf("%s-%02d.pdf", question, batchCount)

			destination := filepath.Join(g.GetExamDir(exam, questionSplit), filename)

			merge.PDF(pagePaths, destination)
		}

	} // for each question

	for _, file := range pagebad {
		err = g.CopyToDir(file, g.GetExamDir(exam, pageBad))
		if err != nil {
			g.logger.Error().
				Str("file", file).
				Str("error", err.Error()).
				Msg("Error copying file to pageBad directory")
		}
	}

	return nil
}

func (g *Ingester) SortCheck(exam string) error {

	files, err := GetFileList(g.GetExamDir(exam, anonPapers))

	usedMap := make(map[string]int)

	if err != nil {

		g.logger.Error().
			Msg("Error getting question back files")
	}

	//WHAT WE STARTED WITH
	fileCount := 0
	pageCount := 0
	for _, file := range files {
		if g.IsPDF(file) {
			fileCount++
			pdm, err := pagedata.UnMarshalAllFromFile(file) //(map[int]PageData, error)
			if err != nil {
				g.logger.Error().
					Str("file", file).
					Msg("Error getting pagedata")
				continue
			}

			for _, pd := range pdm {
				pageCount++
				thisPath := pd.Current.Own.Path
				if _, ok := usedMap[thisPath]; ok {
					fmt.Printf("Found duplicate input page %s\n", thisPath)
					g.logger.Error().
						Str("process", "sort-check").
						Str("file", thisPath).
						Msg("Duplicate page in input")
				} else {
					usedMap[thisPath] = 0
				}

			}
		}
	}

	fmt.Printf("Inputs: Found %d scripts with %d pages\n", fileCount, pageCount)

	files, err = GetFileList(g.GetExamDir(exam, questionSplit))
	ofileCount := 0
	opageCount := 0

	for _, file := range files {
		if g.IsPDF(file) {
			ofileCount++
			pdm, err := pagedata.UnMarshalAllFromFile(file) //(map[int]PageData, error)
			if err != nil {
				g.logger.Error().
					Str("file", file).
					Msg("Error getting pagedata")
				continue
			}

			for _, pd := range pdm {
				opageCount++
				thisPath := pd.Current.Own.Path
				if _, ok := usedMap[thisPath]; !ok {
					fmt.Printf("Found unexpected page in split-by-question %s\n", thisPath)
					g.logger.Error().
						Str("process", "sort-check").
						Str("file", thisPath).
						Msg("Unexpected page in split-by-question")
				} else {
					usedMap[thisPath] += 1
				}
			}
		}
	}

	fmt.Printf("Qsplit: Found %d batches with %d pages\n", ofileCount, opageCount)
	duplicateCount := 0
	totalUsageCount := 0

	for pagePath, useCount := range usedMap {
		totalUsageCount = totalUsageCount + useCount

		if useCount < 1 {
			fmt.Printf("Unused page: %s\n", pagePath)
			g.logger.Warn().
				Str("process", "sort-check").
				Str("exam", exam).
				Str("file", pagePath).
				Msg("Page missing from split-by-question")
		}
		if useCount > 1 {
			fmt.Printf("Duplicate page: %s\n", pagePath)
			duplicateCount = duplicateCount + (useCount - 1)
			g.logger.Warn().
				Str("process", "sort-check").
				Str("exam", exam).
				Str("file", pagePath).
				Msg("Page used more than once in split-by-question")
		}

	}

	pagesMissingCount := pageCount - (totalUsageCount - duplicateCount)

	if pagesMissingCount != 0 {

		fmt.Printf("WARNING: There are %d missing pages\n", pagesMissingCount)
		g.logger.Warn().
			Str("process", "sort-check").
			Str("exam", exam).
			Int("missing", pagesMissingCount).
			Msg(fmt.Sprintf("There were %d Pages missing from split questions", pagesMissingCount))
	} else {
		fmt.Printf("INFO: NO missing pages\n")
		g.logger.Info().
			Str("process", "sort-check").
			Str("exam", exam).
			Msg("All pages accounted for")
	}

	return nil

}
