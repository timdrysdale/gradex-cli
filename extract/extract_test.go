package extract

import (
	"os"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/timdrysdale/gradex-cli/merge"
	"github.com/timdrysdale/unipdf/v3/annotator"
	"github.com/timdrysdale/unipdf/v3/creator"
	pdf "github.com/timdrysdale/unipdf/v3/model"
	"github.com/timdrysdale/unipdf/v3/model/optimize"
)

func TestWhatPageIsThis(t *testing.T) {

	r := regexp.MustCompile("(?:page-)(\\d{3})-(.*)")

	assert.Equal(t, "001", r.FindStringSubmatch("doc2.page-001-question")[1])
	assert.Equal(t, "question", r.FindStringSubmatch("doc2.page-001-question")[2])

	assert.Equal(t, "002", r.FindStringSubmatch("doc001.page-002-banana-apple")[1])
	assert.Equal(t, "banana-apple", r.FindStringSubmatch("doc001.page-002-banana-apple")[2])

	assert.Equal(t, "002", r.FindStringSubmatch("doc001.page-002-question-003")[1])
	assert.Equal(t, "question-003", r.FindStringSubmatch("doc001.page-002-question-003")[2])
}

func writePage(path, key, value, message string) error {

	c := creator.New()

	c.SetPageMargins(0, 0, 0, 0) // we're not printing so use the whole page

	c.SetPageSize(creator.PageSizeA4)

	page := c.NewPage()

	p := c.NewParagraph(message)

	p.SetFontSize(12)

	p.SetPos(200, 10)

	c.Draw(p)

	form := pdf.NewPdfAcroForm()

	tfopt := annotator.TextFieldOptions{Value: value}

	textf, err := annotator.NewTextField(page, key, []float64{100, 200, 150, 250}, tfopt)
	if err != nil {
		return err
	}

	*form.Fields = append(*form.Fields, textf.PdfField)
	page.AddAnnotation(textf.Annotations[0].PdfAnnotation)

	err = c.SetForms(form)
	if err != nil {
		return err
	}

	c.SetOptimizer(optimize.New(optimize.Options{
		CombineDuplicateDirectObjects:   true,
		CombineIdenticalIndirectObjects: true,
		CombineDuplicateStreams:         true,
		CompressStreams:                 true,
		UseObjectStreams:                true,
		ImageQuality:                    90,
		ImageUpperPPI:                   150,
	}))

	err = c.WriteToFile(path)

	return err
}

func TestExtractTextFieldFromFile(t *testing.T) {

	path := "./test/with-fields.pdf"
	path1 := "./test/with-field-p1.pdf"
	path2 := "./test/with-field-p2.pdf"
	path3 := "./test/with-field-p3.pdf"

	textp1 := "What's your favourite colour?"
	textp2 := "What's your favourite food?"
	textp3 := "What's your favourite number?"

	message1 := "TEST PAGE ONE"
	message2 := "TEST PAGE TWO"
	message3 := "TEST PAGE THREE"

	// page-nnn- is an overlay-specifc-feature
	// which ensures we can track the page number in these fields

	writePage(path1, "page-001-question", textp1, message1)
	writePage(path2, "page-002-question", textp2, message2)
	writePage(path3, "page-003-question", textp3, message3)

	err := merge.PDF([]string{path1, path2, path3}, path)
	assert.NoError(t, err)

	_, err = os.Stat(path)

	assert.NoError(t, err)

	fieldsMap, err := ExtractTextFieldsFromPDF(path)

	assert.NoError(t, err)

	expectedMap := make(map[int]map[string]string)

	expectedMap[1] = make(map[string]string)
	expectedMap[2] = make(map[string]string)
	expectedMap[3] = make(map[string]string)

	expectedMap[1]["question"] = textp1
	expectedMap[2]["question"] = textp2
	expectedMap[3]["question"] = textp3

	assert.Equal(t, expectedMap, fieldsMap)

	os.Remove(path)
	os.Remove(path1)
	os.Remove(path2)
	os.Remove(path3)
}
