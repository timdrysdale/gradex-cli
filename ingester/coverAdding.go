package ingester

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/gocarina/gocsv"
	"github.com/timdrysdale/gradex-cli/pagedata"
)

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

	//TODO select option based on whether enter active or inactive ....
	activeOption0 := "merge-entered"
	activeOption1 := "flatten-entered"
	inactiveOption0 := "merge-marked"
	inactiveOption1 := "flatten-marked"

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

	// chose the most appropriate process (merge-entered is more recent than merge-marked.)
	// enter-active-bar won't be present if the paper was marked with the keyboard
	pageDetails := []pagedata.PageDetail{}

	for _, pm := range processMap {

		// check all processes in this page, and look for enter-active-bar

		active := false

		for key, _ := range pm {

			if key == "enter-active-bar" {
				active = true
			}

		}

		switch active {

		case true:

			if pd, ok := pm[activeOption0]; ok {

				pageDetails = append(pageDetails, pd)

				//logger.Info().
				//	Str("file", path).
				//	Int("page", pageNumber).
				//	Msg("Using enter-active-bar for add-cover question data")

			} else if pd, ok := pm[activeOption1]; ok {

				pageDetails = append(pageDetails, pd)
				//logger.Info().
				//	Str("file", path).
				//	Int("page", pageNumber).
				//	Msg("Using merge-marked for add-cover question data")

			}
		case false:

			if pd, ok := pm[inactiveOption0]; ok {

				pageDetails = append(pageDetails, pd)

				//logger.Info().
				//	Str("file", path).
				//	Int("page", pageNumber).
				//	Msg("Using enter-active-bar for add-cover question data")

			} else if pd, ok := pm[inactiveOption1]; ok {

				pageDetails = append(pageDetails, pd)
				//logger.Info().
				//	Str("file", path).
				//	Int("page", pageNumber).
				//	Msg("Using merge-marked for add-cover question data")

			}

		}

	}
	return pageDetails
}

type QuestionSub struct {
	OldQ string `csv:"oldQ"`
	NewQ string `csv:"newQ"`
}

type MarkSub struct {
	Q    string `csv:"Q"`
	Mark string `csv:"mark"`
}

func applyQSubMap(Qmap map[string]string, QSubMap map[string]string) map[string]string {

	for Q, mark := range Qmap {

		if newQ, ok := QSubMap[Q]; ok {
			if currentMark, ok := Qmap[newQ]; ok {

				// asssume denominators have gone!

				var currentVal float64

				currentVal, err := getNum(currentMark)

				if err != nil {
					//fmt.Printf("Error applying QSubMap, %s is not a float64\n", currentMark)
					currentVal = 0
				}

				var markVal float64

				markVal, err = getNum(mark)
				if err != nil {
					//fmt.Printf("Error applying QSubMap, %s is not a float64\n", mark)
					markVal = 0
				}

				if currentMark == "-" && mark == "-" {
					Qmap[newQ] = "-"
				} else {

					Qmap[newQ] = fmt.Sprintf("%g", currentVal+markVal) //add to an existing mark if present
				}

			} else {
				Qmap[newQ] = mark
			}
			delete(Qmap, Q)
		}

	}

	return Qmap

}
func applyMarkSubMap(Qmap map[string]string, MarkSubMap map[string]string) map[string]string {

	for Q, _ := range Qmap {

		if newMark, ok := MarkSubMap[Q]; ok {
			Qmap[Q] = newMark
		}

	}

	return Qmap

}

func getQSubMap(configPath string) (map[string]string, error) {

	subs := []QuestionSub{}

	subMap := make(map[string]string)

	subPath := filepath.Join(configPath, "question-substitutions.csv")

	subFile, err := os.Open(subPath)
	if err != nil {
		return subMap, err
	}

	err = gocsv.UnmarshalFile(subFile, &subs)

	for _, line := range subs {
		subMap[line.OldQ] = line.NewQ
	}

	return subMap, err
}

func getMarkSubMap(configPath, who string) (map[string]string, error) {

	subs := []MarkSub{}
	subMap := make(map[string]string)

	subPath := filepath.Join(configPath, "mark-substitutions-"+who+".csv")

	subFile, err := os.Open(subPath)
	if err != nil {
		return subMap, err
	}

	err = gocsv.UnmarshalFile(subFile, &subs)

	for _, line := range subs {
		subMap[line.Q] = line.Mark
	}

	return subMap, err
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
