package ingester

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/timdrysdale/chmsg"
	"github.com/timdrysdale/gradex-cli/pagedata"
	"github.com/timdrysdale/gradex-cli/parsesvg"
	"vbom.ml/util/sortorder"
)

//not implemented
//func getQNumberfromDataKey(key string) int {
//	re := regexp.MustCompile("/tf-q([0-9]*)(.*)")
//	return 0
//}

type Q struct {
	Number  string
	Section string
	Mark    string
}

func isQ(key string) (string, string, bool) {

	if strings.Contains(key, "optical") {
		return "", "", false
	}

	re := regexp.MustCompile("tf-q([0-9]*)-(\\w*)")

	tokens := re.FindStringSubmatch(key)

	if len(tokens) == 3 {
		// questionNumber, type,
		return tokens[1], tokens[2], true
	} else {

		return "", "", false
	}
}

func getQNum(num string) (string, error) {

	re := regexp.MustCompile("^(\\w?[0-9]*)")

	tokens := re.FindStringSubmatch(strings.TrimSpace(num))

	if len(tokens) == 2 {
		return tokens[1], nil
	} else {
		return "", fmt.Errorf("Number %s not recognised as a Qnumber", num)
	}

}

func getNum(mark string) (float64, error) {

	re := regexp.MustCompile("^([0-9]*\\.?[0-9]*)")

	/*
				    ^(\d*\.?\d*)\/?\\?-?\d*

				     This regexp finds the numerator in fractions marked
					with forward slash, backward slash or hypen
					the denominator is ignored
						^ beginning of line
						\d* any number of digits
						\.? optional dot
						\d* any number of digits

		             e.g. .5/23 -> .5 (half a mark!)

	*/

	tokens := re.FindStringSubmatch(strings.TrimSpace(mark))

	if len(tokens) == 2 {
		return strconv.ParseFloat(tokens[1], 64)
	} else {
		return 0, fmt.Errorf("Mark %s not recognised as a number or fraction", mark)
	}
}

func selectPageDetailsWithMarks(pdMap map[int]pagedata.PageData) []pagedata.PageDetail {

	option1 := "enter-active-bar"
	option2 := "merge-marked"
	option3 := "flatten-marked"

	// get the custom data fields for each process
	// page -> process -> PageDetail
	processMap := make(map[int]map[string]pagedata.PageDetail)

	for pageNumber, singlePagePD := range pdMap {

		processMap[pageNumber] = make(map[string]pagedata.PageDetail)

		processMap[pageNumber][singlePagePD.Current.Process.Name] = singlePagePD.Current

		for _, pd := range singlePagePD.Previous {
			processMap[pageNumber][pd.Process.Name] = pd
		}

	}

	// chose the most appropriate process (enter-active-bar is more recent than merge-marked.)
	// enter-active-bar won'y be present if the paper was marked with the keyboard
	pageDetails := []pagedata.PageDetail{}

	for _, pm := range processMap {

		if pd, ok := pm[option1]; ok {

			pageDetails = append(pageDetails, pd)

			//logger.Info().
			//	Str("file", path).
			//	Int("page", pageNumber).
			//	Msg("Using enter-active-bar for add-cover question data")

		} else if pd, ok := pm[option2]; ok {

			pageDetails = append(pageDetails, pd)
			//logger.Info().
			//	Str("file", path).
			//	Int("page", pageNumber).
			//	Msg("Using merge-marked for add-cover question data")

		} else if pd, ok := pm[option3]; ok {

			pageDetails = append(pageDetails, pd)

		} else {

			//logger.Error().
			//	Str("file", path).
			//	Int("page", pageNumber).
			//	Msg("Error no recognised source of marks - skipping page marks")
			//fmt.Printf("WARN: cover-page for %s: page %d: no recognised source of marks; skipping\n", path, pageNumber)

		}

	}
	return pageDetails
}

func getQMap(pageDetails []pagedata.PageDetail) map[string]string {

	// make a separate interim map for each PageDetail, to avoid collision between keys
	// (the keys repeat on each page)
	// note textfield-question-number is not the exam question number
	// it just comes from the textfield name on the sheet
	// page -> textfield-question-number -> Q struct

	pageQmap := make(map[int]map[string]Q)

	for page, detail := range pageDetails {

		pqm := make(map[string]Q) //qnumber is string format

		for _, item := range detail.Data {

			// piece together the elements in a Q struct
			// one by one as we find the textfields
			// in the pagedata

			if item.Value != "" {
				n, what, is := isQ(item.Key)

				if is {

					if _, ok := pqm[n]; !ok {
						pqm[n] = Q{}
					}

					qn := pqm[n] //get struct, update, and put back

					Val := strings.ToUpper(item.Value)

					switch what {
					case "mark":
						qn.Mark = item.Value
					case "section":
						qn.Section = Val
					case "number":
						num, err := getQNum(Val)
						if err == nil {
							qn.Number = num
						} else {
							qn.Number = Val
						}
					}

					pqm[n] = qn

				}

			}
		}
		pageQmap[page] = pqm
	}

	// we might have more than one mark per question
	// esp if we have split markers
	// so get all the marks values into arrays,
	// one array per question number
	partsQmap := make(map[string][]string)

	for _, qm := range pageQmap {
		for _, q := range qm {

			key := q.Section + q.Number

			if parts, ok := partsQmap[key]; ok {
				partsQmap[key] = append(parts, q.Mark)
			} else {
				partsQmap[key] = []string{q.Mark}
			}
		}
	}

	// add up all the values of those marks
	// interpreting fractions etc
	// convert back to float
	finalQmap := make(map[string]string)

	for q, parts := range partsQmap {

		var val float64

		for _, part := range parts {
			partVal, err := getNum(part)
			if err == nil {
				val = val + partVal
			}
		}

		finalQmap[q] = fmt.Sprintf("%g", val)
	}

	return finalQmap

}

func (g *Ingester) CoverPage(cp CoverPageCommand, logger *zerolog.Logger) error {
	//find pages in processed fir
	// for each page, mangle the name to get the coverpage name
	// sum up all the questions
	// make sure only questions in questions conf are included in the cover page

	EmptyPartsQmap := make(map[string][]string)

	for _, q := range cp.Questions {
		EmptyPartsQmap[q] = []string{"-"} //get all the strings, then sort it out.... "-" for not attempted
	}

	files, err := g.GetFileList(cp.FromPath)

	if err != nil {
		logger.Error().
			Str("dir", cp.FromPath).
			Str("error", err.Error()).
			Msg("Error getting files")

		return err
	}

	for _, path := range files {
		if !IsPDF(path) {
			continue
		}

		pdMap := make(map[int]pagedata.PageData)

		pdMap, err := pagedata.UnMarshalAllFromFile(path)

		if err != nil {
			logger.Error().
				Str("file", path).
				Str("error", err.Error()).
				Msg("Error getting pagedata")
		}

		pageDetails := selectPageDetailsWithMarks(pdMap)

		Qmap := getQMap(pageDetails)

		for _, k := range cp.Questions {
			k = strings.TrimSpace(strings.ToUpper(k))
			if _, ok := Qmap[k]; !ok {
				Qmap[k] = "-"
			}
		}

		pageNumber := 0 //starts at zero

		Prefills := parsesvg.DocPrefills{}

		Prefills[pageNumber] = make(map[string]string)

		Prefills[pageNumber]["page-number"] = "add"

		var thisPageData pagedata.PageData
		// get first page data
		for _, pdm := range pdMap {
			thisPageData = pdm
			break
		}

		fields := []pagedata.Field{}

		for k, v := range Qmap {
			fields = append(fields, pagedata.Field{
				Key:   "pf-q-" + k,
				Value: v,
			})
		}

		thisPageData.Current.Data = fields

		Prefills[pageNumber]["author"] = thisPageData.Current.Item.Who

		Prefills[pageNumber]["date"] = thisPageData.Current.Item.When

		Prefills[pageNumber]["title"] = shortenAssignment(thisPageData.Current.Item.What)

		Prefills[pageNumber]["for"] = thisPageData.Current.Process.For

		var qkeys []string
		for k := range Qmap {
			qkeys = append(qkeys, k)
		}

		sort.Sort(sortorder.Natural(qkeys))

		for idx, qk := range qkeys {
			question := fmt.Sprintf("question-%02d", idx)
			mark := fmt.Sprintf("mark-awarded-%02d", idx)
			Prefills[pageNumber][question] = qk
			Prefills[pageNumber][mark] = Qmap[qk]

		}

		//idx := 0
		//for k, v := range Qmap {
		//	question := fmt.Sprintf("question-%02d", idx)
		//	mark := fmt.Sprintf("mark-awarded-%02d", idx)
		//	Prefills[pageNumber][question] = k
		//	Prefills[pageNumber][mark] = v
		//	idx++
		//}

		pageFilename := filepath.Join(cp.ToPath, strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))+"-cover.pdf")

		contents := parsesvg.SpreadContents{
			SvgLayoutPath:         cp.TemplatePath,
			SpreadName:            cp.SpreadName,
			PageNumber:            pageNumber,
			PdfOutputPath:         pageFilename,
			PageData:              thisPageData,
			TemplatePathsRelative: true,
			Prefills:              Prefills,
		}

		err = parsesvg.RenderSpreadExtra(contents)
		if err != nil {
			msg := fmt.Sprintf("Error rendering spread for cover page for (%s) because %v\n", path, err)
			logger.Error().
				Str("file", path).
				Str("error", err.Error()).
				Msg(msg)
			fmt.Println(msg)
		}
	}

	//sorting!! use order in the csv file, in a list, to get the keys out
	return nil

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
	qfile := filepath.Join(g.GetExamDir(exam, config), "questions.csv")
	qbytes, err := ioutil.ReadFile(qfile)
	if err != nil {
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
	}

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

	err = g.CoverPage(cp, &logger)
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
