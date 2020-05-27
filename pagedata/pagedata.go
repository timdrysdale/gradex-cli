package pagedata

import (
	"encoding/json"
	"fmt"
	"hash/crc32"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/timdrysdale/unipdf/v3/creator"
	"github.com/timdrysdale/unipdf/v3/extractor"
	pdf "github.com/timdrysdale/unipdf/v3/model"
)

func TriageFile(inputPath string) (map[int]Summary, error) {

	s := make(map[int]Summary)

	pds, err := UnMarshalAllFromFile(inputPath)
	if err != nil {
		return s, err
	}

	for n, pd := range pds {
		c := pd.Current
		s[n] = Summary{
			Is:   c.Is,
			What: c.Item.What,
			For:  c.Process.For,
			ToDo: c.Process.ToDo,
		}
	}
	return s, nil

}

func PrettyPrintStruct(layout interface{}) error {

	json, err := json.MarshalIndent(layout, "", "\t")
	if err != nil {
		return err
	}

	fmt.Println(string(json))
	return nil
}

func UnMarshalAllFromFile(inputPath string) (map[int]PageData, error) {

	pdMap := make(map[int]PageData)
	// multiple pagedatas per page is "permitted..." in case of redundancy, compositing, etc.
	// but these are expected to be identical, and we assume (weakly) that damaged pagedata
	// will throw an error at unmarshalling (this requires a key to be corrupted, which is
	// not guaranteed - TODO add CRC or hash check)
	// no disambiguation is provided for different page data on the same page -
	// previous page data should be the array of previous page datas

	textMap, err := extractTextFromPDF(inputPath)
	if err != nil {
		return pdMap, err
	}

	//we get one string per page, which may have multiple pageDatas in it
	for page, text := range textMap {
		var pds []PageData

		rawpds := extractPageDatasFromText(text)

		for _, rawpd := range rawpds {

			pd, err := unMarshalPageData(rawpd)

			if err != nil {
				fmt.Printf("Error extracting page data from Page %d of %s\n", page, inputPath)
				continue
			}
			pds = append(pds, pd)
		}
		if len(pds) > 0 {
			found := false
			for _, pd := range pds {
				if pd.Current.UUID != "" {
					pdMap[page] = pd
					found = true
				}
			}
			if !found {
				fmt.Printf("Error finding non-null page data from Page %d of %s\n", page, inputPath)
			}
		}
	}

	return pdMap, nil

}

func MarshalOneToCreator(c *creator.Creator, pd *PageData) error {

	token, err := json.Marshal(pd)
	if err != nil {
		return err
	}

	writeMarshalledPageDataToCreator(c, string(token))

	return nil

}

func GetLen(input map[int]PageData) int {
	return len(input)
}

//>>>>>>>>>>>>>>>>>>>>>>>> PRIVATE FUNCTIONS >>>>>>>>>>>>>>>>>>>>>>>

/// >>>>>>>>>>>>>>>>>>>>>>>>> TEXT <-> STRUCT >>>>>>>>>>>>>>>>>>>>>>>>

func unMarshalPageData(text string) (PageData, error) {

	var pd PageData

	err := json.Unmarshal([]byte(text), &pd)

	return pd, err

}

func unmarshalPageDatasFromPage(page *pdf.PdfPage) ([]PageData, error) {

	pageDatas := []PageData{}

	tokens, err := readPageDatasFromPage(page)

	if err != nil {
		return pageDatas, err
	}

	var lastError error

	for _, token := range tokens {

		pd, err := unMarshalPageData(token)

		if err != nil {
			lastError = err
			continue
		}

		pageDatas = append(pageDatas, pd)

	}

	return pageDatas, lastError

}

// >>>>>>>>>>>>>>>>> WRITE TO CREATOR >>>>>>>>>>>>>>>>>>>>>>>>>>>>>
func writeMarshalledPageDataToCreator(c *creator.Creator, text string) {

	//drop non-ascii characters to avoid hash issues?
	re := regexp.MustCompile("[[:^ascii:]]")
	text = re.ReplaceAllLiteralString(text, "")

	crc32c := crc32.MakeTable(crc32.Castagnoli)

	hash := fmt.Sprintf("%d", crc32.Checksum([]byte(text), crc32c))

	fulltag := StartTag + text + EndTag + StartHash + hash + EndHash

	check := extractPageDatasFromText(fulltag)

	if text != check[0] {
		fmt.Printf("PageData: Error generating Hash for:n%s\n", text)
	}

	writeTextToCreator(c, fulltag)
}

// We put the text off the page so we are not merged with visible text on the page
// There is a historical challenge in cropping PDF to the visible page
// because, well, it's not that simple to be sure you cropped everything.
// and what do you do with resources that get cut in half.
// so it is robust enough for now to assume that remains an unsolved problem
// and we have tested explicitly that optimisers don't understand
// human visibility limitations
// write one on the page, and one off the page
func writeTextToCreator(c *creator.Creator, text string) {
	p := c.NewParagraph(text)
	p.SetFontSize(0.000001)
	rand.Seed(time.Now().UnixNano())

	//off page
	x := rand.Float64()*0.1 + 99999
	y := rand.Float64()*999 + 99999
	p.SetPos(x, y)
	c.Draw(p)

	//on page
	x = rand.Float64()*0.1 + 1
	y = rand.Float64()*0.1 + 1
	p.SetPos(x, y)
	c.Draw(p)

}

// >>>>>>>>>>>>>>>>>>>> GENERIC READING OPERATION ON TEXT >>>>>>>>>>>>>>>>>>>>>>
// separate for ease of testing
func extractPageDatasFromText(pageText string) []string {

	hashCheckError := false

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

		pageText = pageText[endIndex+EndTagOffset : len(pageText)]

		startIndex = strings.Index(pageText, StartHash)
		if startIndex < 0 {
			break LOOP
		}

		endIndex = strings.Index(pageText, EndHash)
		if endIndex < 0 {
			break LOOP
		}

		actualHash := pageText[startIndex+StartHashOffset : endIndex]

		pageText = pageText[endIndex+EndHashOffset : len(pageText)]

		crc32c := crc32.MakeTable(crc32.Castagnoli)
		checkHash := fmt.Sprintf("%d", crc32.Checksum([]byte(token), crc32c))

		if actualHash != checkHash {
			if hashCheckError == false {
				hashCheckError = true
				fmt.Printf("------------------PageData Hash Check Warning--------------------\n (Got: %s; Want: %s)\n", actualHash, checkHash)
				fmt.Println("Please check the raw data to identify the page, and check it is processed properly:")
				fmt.Printf("%s\n----------------------------------------------------------------\n", token)
			}

		}

		tokens = append(tokens, token)

	}

	return tokens
}

//>>>>>>>>>>>>>>>>>>>>>>>> READ FROM PDF FILE >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>

// outputPdfText produces a map of strings, one per page, indexed starting at 1 for page 1
// mod from https://github.com/unidoc/unipdf-examples/blob/master/text/pdf_extract_text.go
func extractTextFromPDF(path string) (map[int]string, error) {

	textMap := make(map[int]string)

	f, err := os.Open(path)
	if err != nil {
		return textMap, fmt.Errorf("Error opening file %v", err)
	}

	defer f.Close()

	pdfReader, err := pdf.NewPdfReader(f)
	if err != nil {
		return textMap, fmt.Errorf("Error reading PDF from file %v", err)
	}

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return textMap, fmt.Errorf("Error counting pages %v", err)
	}

	for i := 0; i < numPages; i++ {
		pageNum := i + 1

		page, err := pdfReader.GetPage(pageNum)
		if err != nil {
			return textMap, fmt.Errorf("Could not get page %d because %v", pageNum, err)
		}

		ex, err := extractor.New(page)
		if err != nil {
			return textMap, fmt.Errorf("Could not create extractor for page %d because %v", pageNum, err)
		}

		text, err := ex.ExtractText()
		if err != nil {
			return textMap, fmt.Errorf("Could not extract text for page %d because %v", pageNum, err)
		}
		textMap[pageNum] = text
	}

	return textMap, nil
}

// >>>>>>>>>>>>>>>>>>>>>>LESSER-USED FUNCTIONS >>>>>>>>>>>>>>>>>>>>>>>>>>>>>

// >>>>>>>>>>>>>>>>>>>>> WRITE TO PDF FILE >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>

func writeTextToPDF(path string) error {
	// This would push unipdf to the limits of its ability to
	// faithfully duplicate the annotations in a file ...
	// so we call it repair tool and flatten images (see ../repair/)
	return fmt.Errorf("Not implemented - see repair tool")
}

// >>>>>>>>>>>>>>>>>>>>>>>>>> READ FROM CREATOR PAGE >>>>>>>>>>>>>>>>>>>
func readPageDatasFromPage(page *pdf.PdfPage) ([]string, error) {

	text, err := readTextFromPage(page)

	if err != nil {
		return []string{text}, err
	}

	return extractPageDatasFromText(text), nil
}

func readTextFromPage(page *pdf.PdfPage) (string, error) {

	ex, err := extractor.New(page) //this is a pdf function, not pagedata
	if err != nil {
		return "", err
	}

	text, err := ex.ExtractText() // this is a PDF function, not pagedata
	return text, err
}
