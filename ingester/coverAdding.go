package ingester

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/gocarina/gocsv"
	"github.com/timdrysdale/gradex-cli/pagedata"
	"vbom.ml/util/sortorder"
)

type Q struct {
	Number  string
	Section string
	Mark    string
}

type coverQ struct {
	Line int    //line number on cover, starting at top/zero
	Q    string //question
	OK   bool   //it is ok
	Fix  bool   //needs fixing
	Rule string //how to fix it
	Mark string //mark value
}

//([0-9]*\\.?[0-9]*)
//q, newValue, is := isMarkSubRule(cq.Rule)
func isMarkSubRule(cq coverQ) (string, string, bool) {

	if cq.Rule == "-" { // substitute for no attempt
		return cq.Q, "-", true
	}

	// first look for a double rule, and extract the mark sub bit
	// then look for a single rule, and extract the mark sub bit

	// double rule: contains <question>-<mark>[/<outof>] (we throw away the <outof>)
	// applies to the new question included in the rule
	reMQsub := regexp.MustCompile("([a-zA-Z0-9]*)-([0-9]*\\.?[0-9]*)[\\/]?[0-9]*")

	// single rule: contains only an int or floating point number, no letter, and we ditch the <outof>
	// applies to the original question number on this line
	// unless we demand a digit at the start of the string, this returns
	// false positive on Qsub rule being found
	reMSub := regexp.MustCompile("^([0-9][0-9]*\\.?[0-9]*)[\\/]?[0-9]*")

	tokens := reMQsub.FindStringSubmatch(cq.Rule)

	if len(tokens) == 3 {
		//fmt.Printf("MQ:%s: %d,%v\n", cq.Rule, len(tokens), tokens)
		num, err := getNumStr(tokens[2])
		if err != nil {
			return "", "", false
		}
		return cq.Q, num, true
	}

	if len(tokens) == 0 { // not a double rule

		tokens = reMSub.FindStringSubmatch(cq.Rule)
		//fmt.Printf("M:%s: %d,%v\n", cq.Rule, len(tokens), tokens)

		if len(tokens) == 2 {
			num, err := getNumStr(tokens[0])
			if err != nil {
				return "", "", false
			}
			return cq.Q, num, true
		}

	}

	return "", "", false

}

//oldQ, newQ, is := isQSubRule(cq.Rule)
func isQSubRule(cq coverQ) (string, string, bool) {
	if cq.Rule == "-" {
		return "", "", false //this is only a mark subst rule
	}
	// first look for a double rule, and extract the Q sub bit
	// then look for a single rule, and extract the Q sub bit

	// double rule: contains <question>-<mark>[/<outof>] (we throw away the <outof>)
	// applies to the new question included in the rule
	reMQsub := regexp.MustCompile("([a-zA-Z0-9]*)-([0-9]*\\.?[0-9]*)[\\/]?[0-9]*")

	// single rule: contains only the new question label
	// applies to the original question number on this line
	reQSub := regexp.MustCompile("[a-zA-Z][0-9][0-9]*")

	tokens := reMQsub.FindStringSubmatch(cq.Rule)

	// we don't use the second sub token tokens[2], but detecting three parts helps
	// identify we have a double rule

	if len(tokens) == 3 {
		return cq.Q, tokens[1], true
	}

	if len(tokens) == 0 { //not a double rule

		tokens = reQSub.FindStringSubmatch(cq.Rule)

		if len(tokens) == 1 {
			return cq.Q, tokens[0], true
		}

	}

	return "", "", false

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

// for getting ok/dofix/fixrule info from textfields on cover page
func isCoverInfo(key string) (string, int, bool) {

	if strings.Contains(key, "optical") {
		return "", -1, false
	}

	re := regexp.MustCompile("tf-mark-(\\w*)-([0-9]*)")

	tokens := re.FindStringSubmatch(key)

	if len(tokens) == 3 {
		// type, line,

		line, err := strconv.ParseInt(strings.TrimSpace(tokens[2]), 10, 64)

		if err != nil {
			return "", -1, false
		}

		return tokens[1], int(line), true

	} else {

		return "", -1, false
	}
}

// for getting values from textprefill
func isCoverQ(key string) (string, bool) {

	if strings.Contains(key, "optical") {
		return "", false
	}

	re := regexp.MustCompile("pf-q-(\\w*)")

	tokens := re.FindStringSubmatch(key)

	if len(tokens) == 2 {
		// questionNumber, type,
		return tokens[1], true
	} else {

		return "", false
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

func getNumStr(mark string) (string, error) {

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
		return tokens[1], nil
	} else {
		return "", fmt.Errorf("Mark %s not recognised as a number or fraction", mark)
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

//>>>>>>>>>>>>>>>>>>>>>>>>>> COVER QMAP >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>

// get info from textfields, we'll figure out Question later on from prefills....
// pass in the pagedata from the cover page only
func getCoverQMap(pd pagedata.PageDetail) map[int]coverQ {

	coverQMap := make(map[int]coverQ)

	for _, item := range pd.Data {

		// piece together the elements in a Q struct
		// one by one as we find the textfields
		// in the pagedata

		if item.Value != "" {
			what, line, is := isCoverInfo(item.Key)

			if is {

				if _, ok := coverQMap[line]; !ok {
					coverQMap[line] = coverQ{
						Line: line,
					}
				}

				//get struct, update, and put back later
				cq := coverQMap[line]

				switch what {
				case "ok":
					cq.OK = !(item.Value == "") // true if non-empty
				case "fix":
					cq.Fix = !(item.Value == "")
				case "new":
					cq.Rule = strings.ToUpper(item.Value)
				}

				coverQMap[line] = cq

			}

		}
	}

	// we've no adding to do in this map
	// because these lines _should_ be unique
	// we will handle any adding/changes using the rules
	return coverQMap

}

func finalMarksMap(pd pagedata.PageData, questions []string) (map[string]string, map[string]bool) {

	// we KNOW (we hope!) that the marks are on the front page, in a
	// combination of pagedata and textfields (textfields only if
	// there are corrections).
	// Silly old self forgot to link explicitly the rows in the prefills
	// the explicitly labelled text fields, BUT, that was probably because
	// deep down we knew we could just do the same natural order key soring
	// and end up where we need to be .... whilst keeping the page data
	// all clean and tidy *cough*

	// regexp to get Question label
	// map the questions and values
	// order the keys to make a list
	// get the text field values
	// act on any that are marked fix
	// do regexp on fix.
	// this is a strict-match scenario.
	// IFF there is - with at least one word char either side, THEN it is a mark sub then Q sub
	// if there is a - with info missing on one side then we throw an error, probably - too
	// hard to be sure we know that the checker knew how we'd interpret it!
	// we can always check with them, and reprocess it with it corrected.
	// better that than getting a mark wrong through miscommunication
	// even if this will be way slower for the process.
	// If there are two hyphens, then we assume the required mark is a hyphen.
	// If there is just a floating or integer number, we assume it is a new mark
	// if it is not a number, we assume it is a label.
	// if we have a question file, and it is not in the map, throw an error?
	// maybe punt the error onto the console and carry on, just in case someone
	// went for a tea break....

	// we want the text fields from the current pagedata
	coverQMap := getCoverQMap(pd.Current) //map of textfield info by line number

	// find the pagedata.PageDetail with the front cover question labels and marks
	var mergeEnteredPd pagedata.PageDetail
	for _, previousPd := range pd.Previous {

		if previousPd.Process.Name == "merge-entered" {
			mergeEnteredPd = previousPd

		}
	}

	// get the labels and marks into a map by question label
	prefillMap := make(map[string]string)
	for _, prefill := range mergeEnteredPd.Data {
		q, is := isCoverQ(prefill.Key)
		if is {
			prefillMap[q] = prefill.Value
		}
	}

	// recreate the order of the questions on the cover
	// using the same sort algorithm

	var pkeys []string
	for k := range prefillMap {
		pkeys = append(pkeys, k)
	}

	sort.Sort(sortorder.Natural(pkeys))

	// Add the question label and marks to our coverQMap

	// note we will have prefills in strict contiguous sequence
	// but textfield lines may have been skipped....i.e. OK = Fix = false
	for n, key := range pkeys {
		if _, ok := coverQMap[n]; !ok {
			coverQMap[n] = coverQ{
				Line: n,
			}
		}
		cq := coverQMap[n]
		cq.Q = key
		cq.Mark = prefillMap[key]
		coverQMap[n] = cq
	}

	// Find any substitution rules

	MarkSubMap := make(map[string]string)

	QSubMap := make(map[string]string)

	for _, cq := range coverQMap {
		if cq.Fix {
			Q, newValue, is := isMarkSubRule(cq)
			if is {
				MarkSubMap[Q] = newValue
			}
			oldQ, newQ, is := isQSubRule(cq)
			if is {
				QSubMap[oldQ] = newQ
			}
		}
	}

	//fmt.Println("Mark sub")
	//util.PrettyPrintStruct(MarkSubMap)
	//fmt.Println("Q sub")
	//util.PrettyPrintStruct(QSubMap)

	QMap := make(map[string]string)

	for _, cq := range coverQMap {

		QMap[cq.Q] = cq.Mark
	}

	// the order matters here, because we add subparts
	QMap = applyMarkSubMap(QMap, MarkSubMap)
	QMap = applyQSubMap(QMap, QSubMap)

	// put in "-" for any standard Q we are missing in this map
	for _, k := range questions {
		k = strings.TrimSpace(strings.ToUpper(k))
		if _, ok := QMap[k]; !ok {
			QMap[k] = "-"
		}
	}

	skipMap := make(map[string]bool)

	for _, cq := range coverQMap {
		skipMap[cq.Q] = (!cq.OK) && (!cq.Fix)
	}

	return QMap, skipMap
}
