package ingester

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/timdrysdale/gradex-cli/csv"
	"github.com/timdrysdale/gradex-cli/pagedata"
)

type Mark struct {
	Q string
	V string
}

type MarkCompare struct {
	Who        string
	What       string
	When       string
	Final      []Mark
	Draft      []Mark
	DraftTotal float64
	FinalTotal float64
	Comment    string
}

func (g *Ingester) FinalReport(exam string) error {

	s := csv.New()

	s.SetFixedHeader([]string{"what", "when", "who", "comment", "total", "total-draft"})

	qfile := filepath.Join(g.GetExamDir(exam, config), "questions.csv")

	reqdQ, err := GetRequiredQuestions(qfile)
	if err != nil {
		fmt.Println(err)
	}

	combinedReqdQ := []string{}

	for _, q := range reqdQ {

		combinedReqdQ = append(combinedReqdQ, q+"-draft")
		combinedReqdQ = append(combinedReqdQ, q)
	}

	s.SetRequiredHeader(combinedReqdQ)

	markMap := make(map[string]MarkCompare)

	files, err := g.GetFileList(g.GetExamDir(exam, finalCover))
	if err != nil {
		return err
	}

	fmt.Println(g.GetExamDir(exam, finalCover))

	for _, file := range files {

		fmt.Println(file)

		if !IsPDF(file) {
			continue
		}

		marks, item, err := GetMarksFromCoverPage(file)

		if err != nil {
			return fmt.Errorf("ERROR getting marks from %s", file)
		}

		who := item.Who

		markMap[who] = MarkCompare{
			Who:   who,
			When:  item.When,
			What:  item.What,
			Final: marks,
		}

	}

	files, err = g.GetFileList(g.GetExamDir(exam, checkerCover))
	if err != nil {
		return err
	}
	fmt.Println(g.GetExamDir(exam, finalCover))

	for _, file := range files {

		fmt.Println(file)

		if !IsPDF(file) {
			continue
		}

		marks, item, err := GetMarksFromCoverPage(file)

		if err != nil {
			return fmt.Errorf("ERROR getting marks from %s", file)
		}

		who := item.Who

		if _, ok := markMap[who]; !ok {
			markMap[who] = MarkCompare{
				Who:   who,
				When:  item.When,
				What:  item.What,
				Draft: marks,
			}
		} else {
			mm := markMap[who]
			mm.Draft = marks
			markMap[who] = mm
		}
	}

	for _, m := range markMap {

		line := s.Add()

		line.Add("what", m.What)
		line.Add("who", m.Who)
		line.Add("when", m.When)

		thisTotal := 0.0

		for _, mark := range m.Draft {
			thisVal := 0.0
			if !(mark.V == "-" || mark.V == "") {
				thisVal, err = strconv.ParseFloat(mark.V, 64)
				if err != nil {
					return err
				}
			}

			thisTotal = thisTotal + thisVal
			line.Add(mark.Q+"-draft", mark.V)

		}

		thisTotalStr := fmt.Sprintf("%g", thisTotal)
		line.Add("total-draft", thisTotalStr)
		m.DraftTotal = thisTotal

		thisTotal = 0.0

		for _, mark := range m.Final {
			thisVal := 0.0
			if !(mark.V == "-" || mark.V == "") {
				thisVal, err = strconv.ParseFloat(mark.V, 64)
				if err != nil {
					return err
				}
			}
			thisTotal = thisTotal + thisVal
			line.Add(mark.Q, mark.V)

		}

		thisTotalStr = fmt.Sprintf("%g", thisTotal)
		line.Add("total", thisTotalStr)

		m.FinalTotal = thisTotal

		totalChanged := m.DraftTotal != m.FinalTotal

		if totalChanged {
			line.Add("comment", "TOTAL-CHANGED")
		} else {

			splitChanged := false

			for _, mark := range m.Final {

				for _, oldMark := range m.Draft {
					if oldMark.Q == mark.Q {
						if oldMark.V != mark.V {
							splitChanged = true
						}
					}
				}
			}

			if splitChanged {
				line.Add("comment", "SPLIT-CHANGED")
			} else {

				line.Add("comment", "")
			}
		}

	}

	reportBase := fmt.Sprintf("FinalMarks-%s-%d.csv", shortenAssignment(exam), time.Now().Unix())
	reportPath := filepath.Join(g.GetExamDir(exam, reports), reportBase)

	f, err := os.OpenFile(reportPath, os.O_RDWR|os.O_CREATE, os.ModePerm)

	defer f.Close()

	_, err = s.WriteCSV(f)

	return err

}

func (g *Ingester) CheckReport(exam string) error {

	s := csv.New()

	s.SetFixedHeader([]string{"what", "who", "when"})

	qfile := filepath.Join(g.GetExamDir(exam, config), "questions.csv")

	reqdQ, err := GetRequiredQuestions(qfile)
	if err != nil {
		fmt.Println(err)
	}
	s.SetRequiredHeader(reqdQ)

	files, err := g.GetFileList(g.GetExamDir(exam, checkerCover))

	for _, file := range files {

		if !IsPDF(file) {
			continue
		}

		marks, item, err := GetMarksFromCoverPage(file)

		if err != nil {
			return fmt.Errorf("ERROR getting marks from %s", file)
		}

		line := s.Add()

		line.Add("what", item.What)
		line.Add("who", item.Who)
		line.Add("when", item.When)

		for _, mark := range marks {
			line.Add(mark.Q, mark.V)
		}

	}

	reportBase := fmt.Sprintf("ProvisionalMarks-%s-%d.csv", shortenAssignment(exam), time.Now().Unix())
	reportPath := filepath.Join(g.GetExamDir(exam, reports), reportBase)

	f, err := os.OpenFile(reportPath, os.O_RDWR|os.O_CREATE, os.ModePerm)

	defer f.Close()

	_, err = s.WriteCSV(f)

	return err

}

func GetRequiredQuestions(qfile string) ([]string, error) {

	questions := []string{}
	qbytes, err := ioutil.ReadFile(qfile)
	if err != nil {
		return questions, err
	}
	questions = strings.Split(string(qbytes), ",")

	for i, q := range questions {

		questions[i] = strings.TrimSpace(strings.ToUpper(q))

	}

	return questions, nil
}

func GetMarksFromCoverPage(path string) ([]Mark, pagedata.ItemDetail, error) {

	marks := []Mark{}
	item := pagedata.ItemDetail{}

	pdMap, err := pagedata.UnMarshalAllFromFile(path)

	if err != nil {
		return marks, item, err
	}

	for _, pd := range pdMap {

		item = pd.Current.Item

		for _, field := range pd.Current.Data {

			if question, err := getQ(field.Key); err == nil {

				marks = append(marks, Mark{
					Q: question,
					V: field.Value})

			}

		}

	}

	return marks, item, nil

}

func getQ(q string) (string, error) {

	re := regexp.MustCompile("pf-q-(\\w*)")

	tokens := re.FindStringSubmatch(strings.TrimSpace(q))

	if len(tokens) == 2 {
		return tokens[1], nil
	} else {
		return "", fmt.Errorf("Key  %s not recognised as a cover page question", q)
	}
}
