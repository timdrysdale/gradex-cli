package pagedata

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/mattetti/filebuffer"
	"github.com/stretchr/testify/assert"
	"github.com/timdrysdale/unipdf/v3/annotator"
	"github.com/timdrysdale/unipdf/v3/creator"
	"github.com/timdrysdale/unipdf/v3/model"
	pdf "github.com/timdrysdale/unipdf/v3/model"
	"github.com/timdrysdale/unipdf/v3/model/optimize"
)

func TestGetLen(t *testing.T) {

	foo := make(map[int]PageData)

	var bar, dab, dib, dob PageData

	foo[10] = bar
	foo[99] = dab
	foo[12] = dib
	foo[35] = dob

	assert.Equal(t, GetLen(foo), 4)

}

/*
	p := c.NewParagraph("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	p.SetFontSize(12)
	p.SetPos(0.1, 0.1)
	c.Draw(p)

*/
func TestWriteRead(t *testing.T) {

	c := creator.New()
	c.SetPageMargins(0, 0, 0, 0) // we're not printing so use the whole page
	c.SetPageSize(creator.PageSizeA4)
	c.NewPage()
	text1 := "{\"exam\":\"ENGI99887\",\"number\":\"B12345\",\"page\":1}"
	writeTextToCreator(c, text1)
	text2 := "{\"exam\":\"ENGI99887\",\"number\":\"B12345\",\"page\":2}"
	c.NewPage()
	writeTextToCreator(c, text2)
	c.NewPage()
	// write to memory instead of a file
	var buf bytes.Buffer

	c.SetOptimizer(optimize.New(optimize.Options{
		CombineDuplicateDirectObjects:   true,
		CombineIdenticalIndirectObjects: true,
		CombineDuplicateStreams:         true,
		CompressStreams:                 true,
		UseObjectStreams:                true,
		ImageQuality:                    90,
		ImageUpperPPI:                   150,
	}))

	err := c.Write(&buf)
	if err != nil {
		t.Error(err)
	}

	// convert buffer to readseeker
	var bufslice []byte
	fbuf := filebuffer.New(bufslice)
	fbuf.Write(buf.Bytes())

	// read in from memory
	pdfReader, err := model.NewPdfReader(fbuf)
	if err != nil {
		t.Error(err)
	}

	page, err := pdfReader.GetPage(1)
	if err != nil {
		t.Error(err)
	}

	textp1, err := readTextFromPage(page)

	if err != nil {
		t.Error(err)
	}

	page2, err := pdfReader.GetPage(2)
	if err != nil {
		t.Error(err)
	}

	textp2, err := readTextFromPage(page2)
	if err != nil {
		t.Error(err)
	}

	assertEqual(t, text1+"\n"+text1, textp1) // we write two copies:
	assertEqual(t, text2+"\n"+text2, textp2) // one on, one off, page

}

func TestMarshalling(t *testing.T) {

	pd := PageData{
		Current: PageDetail{
			Is:   IsPage,
			UUID: "69197384-fd15-42ac-ac16-82dbe4d52dd0",
		},
		Previous: []PageDetail{
			PageDetail{},
			PageDetail{},
		},
	}

	c := creator.New()
	c.SetPageMargins(0, 0, 0, 0) // we're not printing so use the whole page
	c.SetPageSize(creator.PageSizeA4)
	c.NewPage()

	MarshalOneToCreator(c, &pd)

	// write to memory instead of a file
	var buf bytes.Buffer

	c.SetOptimizer(optimize.New(optimize.Options{
		CombineDuplicateDirectObjects:   true,
		CombineIdenticalIndirectObjects: true,
		CombineDuplicateStreams:         true,
		CompressStreams:                 true,
		UseObjectStreams:                true,
		ImageQuality:                    90,
		ImageUpperPPI:                   150,
	}))

	err := c.Write(&buf)
	if err != nil {
		t.Error(err)
	}

	// convert buffer to readseeker
	var bufslice []byte
	fbuf := filebuffer.New(bufslice)
	fbuf.Write(buf.Bytes())

	// read in from memory
	pdfReader, err := model.NewPdfReader(fbuf)
	if err != nil {
		t.Error(err)
	}
	page, err := pdfReader.GetPage(1)
	if err != nil {
		t.Error(err)
	}

	pdout, err := unmarshalPageDatasFromPage(page)

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, pd, pdout[0])

}

func TestWriteOutputForAdobe(t *testing.T) {

	outputPath := "./test/adobe-page-data.pdf"
	_, err := os.Stat(outputPath)
	if !os.IsNotExist(err) {
		return
	}
	fmt.Printf("Creating %s  please edit in your editor and run tests again\n", outputPath)
	c := creator.New()
	c.SetPageMargins(0, 0, 0, 0) // we're not printing so use the whole page
	c.SetPageSize(creator.PageSizeA4)

	c.NewPage()

	p := c.NewParagraph("Test page for github.com/timdrysdale/pdfpagedata")
	p.SetFontSize(12)
	p.SetPos(200, 10)
	c.Draw(p)

	pd := PageData{
		Current: PageDetail{
			Is:   IsPage,
			UUID: "69197384-fd15-42ac-ac16-82dbe4d52dd0",
		},
		Previous: []PageDetail{
			PageDetail{
				Is:   IsPage,
				UUID: "69197384-fd15-42ac-ac16-82dbe4d52dd0",
				Data: []Field{
					Field{
						Key:   "key",
						Value: "value",
					},
				},
			},
		},
	}

	MarshalOneToCreator(c, &pd)

	// write to memory
	var buf bytes.Buffer

	err = c.Write(&buf)
	if err != nil {
		t.Error(err)
	}

	// convert buffer to readseeker
	var bufslice []byte
	fbuf := filebuffer.New(bufslice)
	fbuf.Write(buf.Bytes())

	// read in from memory
	pdfReader, err := model.NewPdfReader(fbuf)
	if err != nil {
		t.Error(err)
	}

	pdfWriter := pdf.NewPdfWriter()

	page, err := pdfReader.GetPage(1)
	if err != nil {
		t.Error(err)
	}

	form := pdf.NewPdfAcroForm()

	tfopt := annotator.TextFieldOptions{Value: "type viewer name here"}
	name := fmt.Sprintf("viewer")

	textf, err := annotator.NewTextField(page, name, []float64{100, 200, 150, 250}, tfopt)
	if err != nil {
		t.Error(err)
	}

	*form.Fields = append(*form.Fields, textf.PdfField)
	page.AddAnnotation(textf.Annotations[0].PdfAnnotation)

	err = pdfWriter.SetForms(form)
	if err != nil {
		t.Error(err)
	}

	err = pdfWriter.AddPage(page)
	if err != nil {
		t.Error(err)
	}

	of, err := os.Create(outputPath)
	if err != nil {
		t.Error(err)
	}

	defer of.Close()

	pdfWriter.SetOptimizer(optimize.New(optimize.Options{
		CombineDuplicateDirectObjects:   true,
		CombineIdenticalIndirectObjects: true,
		CombineDuplicateStreams:         true,
		CompressStreams:                 true,
		UseObjectStreams:                true,
		ImageQuality:                    90,
		ImageUpperPPI:                   150,
	}))

	pdfWriter.Write(of)

}

// see if data off the page survives a form edit with Adobe....
func TestAdobe(t *testing.T) {

	pd := PageData{
		Current: PageDetail{
			Is:   IsPage,
			UUID: "69197384-fd15-42ac-ac16-82dbe4d52dd0",
		},
		Previous: []PageDetail{
			PageDetail{
				Is:   IsPage,
				UUID: "69197384-fd15-42ac-ac16-82dbe4d52dd0",
				Data: []Field{
					Field{
						Key:   "key",
						Value: "value",
					},
				},
			},
		},
	}

	inputPath := "./test/adobe-page-data.pdf"
	field, err := getPdfFieldData(inputPath, "viewer")

	assert.NoError(t, err)

	if strings.Compare(field, "type viewer name here") == 0 {
		t.Error("Edit the form data in ./test/adode-page-data.txt and then run test again")
	}

	pdMap, err := UnMarshalAllFromFile(inputPath)

	assert.Equal(t, 1, len(pdMap))
	assert.Equal(t, pd, pdMap[1])

}

func getPdfFieldData(inputPath, targetFieldName string) (string, error) {
	f, err := os.Open(inputPath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	pdfReader, err := pdf.NewPdfReader(f)
	if err != nil {
		return "", err
	}

	acroForm := pdfReader.AcroForm
	if acroForm == nil {
		return "", nil
	}

	match := false
	fields := acroForm.AllFields()
	for _, field := range fields {
		fullname, err := field.FullName()
		if err != nil {
			return "", err
		}
		if fullname == targetFieldName {
			match = true

			if field.V != nil {
				return field.V.String(), nil
			} else {
				return "", nil
			}
		}
	}

	if !match {
		return "", errors.New("field not found")
	}
	return "", nil
}

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Fatalf("%s != %s", a, b)
	}
}

// Mod from array to slice,
// from https://www.golangprograms.com/golang-check-if-array-element-exists.html
func itemExists(sliceType interface{}, item interface{}) bool {
	slice := reflect.ValueOf(sliceType)

	if slice.Kind() != reflect.Slice {
		panic("Invalid data-type")
	}

	for i := 0; i < slice.Len(); i++ {
		if slice.Index(i).Interface() == item {
			return true
		}
	}

	return false
}
