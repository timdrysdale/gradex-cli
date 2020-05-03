package pagedata

import (
	"encoding/json"
	"errors"
	"math/rand"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/timdrysdale/unipdf/v3/creator"
	"github.com/timdrysdale/unipdf/v3/extractor"
	pdf "github.com/timdrysdale/unipdf/v3/model"
)

// UnixTime is not all that portable as it is ...
// so only using it as a tiebreaker
// when likely that the process was repeated on the
// same machine
// For cases where marking is repeated after checking
// The overlay SHOULD just use a higher sequence number,
// sequence number is a step in sequence, not another
// name for a particular process
const (
	Raw             = iota
	ReadyToMark     = iota
	Marked          = iota
	ReadyToModerate = iota
	Moderated       = iota
	ReadyToCheck    = iota
	Checked         = iota
)

type PdfSummary struct {
	CourseCode  string
	PreparedFor string
	ToDo        string
}

func TriagePdf(inputPath string) (PdfSummary, error) {

	pdfs := PdfSummary{}

	pdm, err := GetPageDataFromFile(inputPath)
	if err != nil {
		return pdfs, err
	}

	err = PruneOldRevisions(&pdm)
	if err != nil {
		return pdfs, err
	}

OUTER:
	for _, v := range pdm {
		for _, pd := range v {
			pdfs.PreparedFor = pd.PreparedFor
			pdfs.ToDo = pd.ToDo
			pdfs.CourseCode = pd.Exam.CourseCode
			break OUTER
		}
	}
	return pdfs, nil
}

func PruneOldRevisions(pdmap *map[int][]PageData) error {
	for k, v := range *pdmap {
		pd, err := SelectPageDataByRevision(v)
		if err != nil {
			return err
		}
		(*pdmap)[k] = []PageData{pd}
	}
	return nil
}

func SelectPageDataByRevision(pds []PageData) (PageData, error) {
	if len(pds) < 1 {
		return PageData{}, errors.New("empty")
	}
	if len(pds) == 1 {
		return pds[0], nil
	}

	sort.SliceStable(pds, func(i, j int) bool {
		return pds[i].Revision > pds[j].Revision

	})

	return pds[0], nil
}
func SelectQuestionByLast(pd PageData) (QuestionDetails, error) {
	if len(pd.Questions) < 1 {
		return QuestionDetails{}, errors.New("empty")
	}
	if len(pd.Questions) == 1 {
		return pd.Questions[0], nil
	}

	Q := pd.Questions
	sort.SliceStable(Q, func(i, j int) bool {
		if Q[i].Sequence == Q[j].Sequence {
			return Q[i].UnixTime > Q[j].UnixTime
		} else {
			return Q[i].Sequence > Q[j].Sequence
		}
	})

	return Q[0], nil
}

func SelectProcessByLast(pd PageData) (ProcessingDetails, error) {
	if len(pd.Processing) < 1 {
		return ProcessingDetails{}, errors.New("empty")
	}
	if len(pd.Processing) == 1 {
		return pd.Processing[0], nil
	}

	Process := pd.Processing
	sort.SliceStable(Process, func(i, j int) bool {
		if Process[i].Sequence == Process[j].Sequence {
			return Process[i].UnixTime > Process[j].UnixTime
		} else {
			return Process[i].Sequence > Process[j].Sequence
		}
	})

	return Process[0], nil

}

func GetLen(input map[int][]PageData) int {
	items := 0
	for _, v := range input {
		for _ = range v {
			items++
		}
	}
	return items
}

func GetPageDataFromFile(inputPath string) (map[int][]PageData, error) {

	docData := make(map[int][]PageData)

	texts, err := OutputPdfText(inputPath)
	if err != nil {
		return docData, err
	}

	//one text per page
	for i, text := range texts {
		var pds []PageData

		strs := ExtractPageData(text)

		for _, str := range strs {
			var pd PageData

			if err := json.Unmarshal([]byte(str), &pd); err != nil {
				continue
			}

			pds = append(pds, pd)
		}

		docData[i] = pds
	}

	return docData, nil

}

func UnmarshalPageData(page *pdf.PdfPage) ([]PageData, error) {

	pageDatas := []PageData{}

	tokens, err := ReadPageData(page)

	if err != nil {
		return pageDatas, err
	}

	var lastError error

	for _, token := range tokens {

		var pd PageData

		if err := json.Unmarshal([]byte(token), &pd); err != nil {
			lastError = err
			continue
		}

		pageDatas = append(pageDatas, pd)

	}

	return pageDatas, lastError

}

func MarshalPageData(c *creator.Creator, pd *PageData) error {

	token, err := json.Marshal(pd)
	if err != nil {
		return err
	}

	WritePageData(c, string(token))

	return nil

}

func ReadPageData(page *pdf.PdfPage) ([]string, error) {

	text, err := ReadPageString(page)

	if err != nil {
		return []string{text}, err
	}

	return ExtractPageData(text), nil

}

func ExtractPageData(pageText string) []string {

	var tokens []string

LOOP:
	for {

		startIndex := strings.Index(pageText, StartTag)
		if startIndex < 0 {
			break LOOP
		}

		endIndex := strings.Index(pageText, EndTag)
		if endIndex < 0 {
			break LOOP
		}

		token := pageText[startIndex+StartTagOffset : endIndex]

		tokens = append(tokens, token)

		pageText = pageText[endIndex+EndTagOffset : len(pageText)]

	}

	return tokens
}

func ReadPageString(page *pdf.PdfPage) (string, error) {

	ex, err := extractor.New(page)
	if err != nil {
		return "", err
	}

	text, err := ex.ExtractText()
	return text, err
}

func WritePageData(c *creator.Creator, text string) {
	WritePageString(c, StartTag+text+EndTag)
}

func WritePageString(c *creator.Creator, text string) {
	p := c.NewParagraph(text)
	p.SetFontSize(0.000001)
	rand.Seed(time.Now().UnixNano())
	x := rand.Float64()*0.1 + 99999 //0.3
	y := rand.Float64()*999 + 99999 //0.3
	p.SetPos(x, y)
	c.Draw(p)
}

// this function is for use in a co-operative
// environment - you can slip one past the gaolie
// in the custom fields in Questions/Processing/Custom
func StripAuthorIdentity(pd PageData) PageData {

	safe := PageData{}

	safe.Exam = pd.Exam
	safe.Author = AuthorDetails{Anonymous: pd.Author.Anonymous}
	safe.Page = pd.Page
	safe.Contact = pd.Contact
	safe.Submission = SubmissionDetails{} //nothing!
	safe.Questions = pd.Questions
	safe.Processing = pd.Processing
	safe.Custom = pd.Custom

	return safe
}

// outputPdfText produces array of strings, one string per page
// mod from https://github.com/unidoc/unipdf-examples/blob/master/text/pdf_extract_text.go
func OutputPdfText(inputPath string) ([]string, error) {

	texts := []string{}

	f, err := os.Open(inputPath)
	if err != nil {
		return texts, err
	}

	defer f.Close()

	pdfReader, err := pdf.NewPdfReader(f)
	if err != nil {
		return texts, err
	}

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return texts, err
	}

	for i := 0; i < numPages; i++ {
		pageNum := i + 1

		page, err := pdfReader.GetPage(pageNum)
		if err != nil {
			return texts, err
		}

		ex, err := extractor.New(page)
		if err != nil {
			return texts, err
		}

		text, err := ex.ExtractText()
		if err != nil {
			return texts, err
		}

		texts = append(texts, text)
	}

	return texts, nil
}
