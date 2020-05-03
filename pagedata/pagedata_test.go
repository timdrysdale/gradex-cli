package pagedata

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/mattetti/filebuffer"
	"github.com/stretchr/testify/assert"
	"github.com/timdrysdale/unipdf/v3/annotator"
	"github.com/timdrysdale/unipdf/v3/creator"
	"github.com/timdrysdale/unipdf/v3/model"
	pdf "github.com/timdrysdale/unipdf/v3/model"
	"github.com/timdrysdale/unipdf/v3/model/optimize"
)

func TestPruneOldRevsions(t *testing.T) {

	testSet := make(map[int][]PageData)

	testSet[0] = []PageData{
		PageData{
			Revision: 0,
		},
		PageData{
			Revision: 2,
		},
		PageData{
			Revision: 1,
		},
	}
	testSet[1] = []PageData{
		PageData{
			Revision: 0,
		},
		PageData{
			Revision: 1,
		},
	}
	testSet[2] = []PageData{
		PageData{
			Revision: 0,
		},
	}
	assert.Equal(t, 6, GetLen(testSet))
	err := PruneOldRevisions(&testSet)
	assert.NoError(t, err)
	assert.Equal(t, GetLen(testSet), 3)
	assert.Equal(t, testSet[0][0].Revision, 2)
	assert.Equal(t, testSet[1][0].Revision, 1)
	assert.Equal(t, testSet[2][0].Revision, 0)
}

func TestSelectPageDataByRevision(t *testing.T) {

	//throw error on empty, obvs
	_, err := SelectPageDataByRevision([]PageData{})
	assert.Error(t, err)

	testSet := []PageData{
		PageData{
			Revision: 0,
		},
	}

	pd, err := SelectPageDataByRevision(testSet)
	assert.NoError(t, err)
	assert.Equal(t, 0, pd.Revision)

	testSet = []PageData{
		PageData{
			Revision: 0,
		},
		PageData{
			Revision: 1,
		},
	}

	pd, err = SelectPageDataByRevision(testSet)
	assert.NoError(t, err)
	assert.Equal(t, 1, pd.Revision)

	testSet = []PageData{
		PageData{
			Revision: 0,
		},
		PageData{
			Revision: 2,
		},
		PageData{
			Revision: 1,
		},
	}
	pd, err = SelectPageDataByRevision(testSet)
	assert.NoError(t, err)
	assert.Equal(t, 2, pd.Revision)

}

func TestSelectProcessByLast(t *testing.T) {

	//throw error on empty, obvs
	_, err := SelectProcessByLast(PageData{})
	assert.Error(t, err)

	now := time.Now().UnixNano()

	testSet := PageData{
		Processing: []ProcessingDetails{
			ProcessingDetails{
				Name:     "first",
				Sequence: 0,
				UnixTime: now - 3600000,
			},
		},
	}

	proc, err := SelectProcessByLast(testSet)
	assert.NoError(t, err)
	assert.Equal(t, "first", proc.Name)

	testSet = PageData{
		Processing: []ProcessingDetails{
			ProcessingDetails{
				Name:     "first",
				Sequence: 0,
				UnixTime: now - 3600000,
			},
			ProcessingDetails{
				Name:     "second",
				Sequence: 1,
				UnixTime: now - 1800000,
			},
			ProcessingDetails{
				Name:     "third",
				Sequence: 2,
				UnixTime: now - 100000,
			},
			ProcessingDetails{
				Name:     "fourth",
				Sequence: 2,
				UnixTime: now,
			},
		},
	}

	proc, err = SelectProcessByLast(testSet)
	assert.NoError(t, err)
	assert.Equal(t, "fourth", proc.Name)

	testSet = PageData{
		Processing: []ProcessingDetails{
			ProcessingDetails{
				Name:     "fifth",
				Sequence: 3,
				UnixTime: now - 3600000,
			},
			ProcessingDetails{
				Name:     "second",
				Sequence: 1,
				UnixTime: now - 1800000,
			},
			ProcessingDetails{
				Name:     "third",
				Sequence: 2,
				UnixTime: now - 100000,
			},
			ProcessingDetails{
				Name:     "fourth",
				Sequence: 2,
				UnixTime: now,
			},
		},
	}

	proc, err = SelectProcessByLast(testSet)
	assert.NoError(t, err)
	assert.Equal(t, "fifth", proc.Name)

}

func TestGetLen(t *testing.T) {

	foo := make(map[int][]PageData)

	bar := make([]PageData, 5)

	dab := make([]PageData, 6)

	foo[10] = bar
	foo[99] = dab

	assert.Equal(t, GetLen(foo), 11)

}

func TestWriteRead(t *testing.T) {

	c := creator.New()
	c.SetPageMargins(0, 0, 0, 0) // we're not printing so use the whole page
	c.SetPageSize(creator.PageSizeA4)
	c.NewPage()
	text1 := "{\"exam\":\"ENGI99887\",\"number\":\"B12345\",\"page\":1}"
	WritePageData(c, text1)
	text2 := "{\"exam\":\"ENGI99887\",\"number\":\"B12345\",\"page\":2}"
	c.NewPage()
	WritePageData(c, text2)
	c.NewPage()
	// write to memory instead of a file
	var buf bytes.Buffer

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

	textp1, err := ReadPageData(page)

	if err != nil {
		t.Error(err)
	}

	page2, err := pdfReader.GetPage(2)
	if err != nil {
		t.Error(err)
	}

	textp2, err := ReadPageData(page2)
	if err != nil {
		t.Error(err)
	}

	assertEqual(t, text1, textp1[0])
	assertEqual(t, text2, textp2[0])

}

func TestWriteReadOptimiser(t *testing.T) {

	c := creator.New()
	c.SetPageMargins(0, 0, 0, 0) // we're not printing so use the whole page
	c.SetPageSize(creator.PageSizeA4)
	c.NewPage()
	text1 := "{\"exam\":\"ENGI99887\",\"number\":\"B12345\",\"page\":1}"
	WritePageData(c, text1)
	text2 := "{\"exam\":\"ENGI99887\",\"number\":\"B12345\",\"page\":2}"
	c.NewPage()
	WritePageData(c, text2)

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

	textp1, err := ReadPageData(page)

	if err != nil {
		t.Error(err)
	}

	page2, err := pdfReader.GetPage(2)
	if err != nil {
		t.Error(err)
	}

	textp2, err := ReadPageData(page2)
	if err != nil {
		t.Error(err)
	}

	assertEqual(t, text1, textp1[0])
	assertEqual(t, text2, textp2[0])

}

func TestWriteReadDouble(t *testing.T) {

	c := creator.New()
	c.SetPageMargins(0, 0, 0, 0) // we're not printing so use the whole page
	c.SetPageSize(creator.PageSizeA4)

	c.NewPage()
	text1a := "{\"exam\":\"ENGI99887\",\"number\":\"B12345\",\"page\":1,\"Batch\":\"a\"}"
	text1b := "{\"exam\":\"ENGI99886\",\"number\":\"B12345\",\"page\":1,\"Batch\":\"xx\"}"
	WritePageData(c, text1a)
	WritePageData(c, text1b)

	c.NewPage()
	text2a := "{\"exam\":\"ENGI99887\",\"number\":\"B12345\",\"page\":2,\"Batch\":\"a\"}"
	text2b := "{\"exam\":\"ENGI99897\",\"number\":\"B12345\",\"page\":2,\"Batch\":\"b\"}"
	WritePageData(c, text2a)
	WritePageData(c, text2b)

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
	pdfReader, err := pdf.NewPdfReader(fbuf)
	if err != nil {
		t.Error(err)
	}
	page, err := pdfReader.GetPage(1)
	if err != nil {
		t.Error(err)
	}

	textp1, err := ReadPageData(page)

	if err != nil {
		t.Error(err)
	}

	page2, err := pdfReader.GetPage(2)
	if err != nil {
		t.Error(err)
	}

	textp2, err := ReadPageData(page2)
	if err != nil {
		t.Error(err)
	}

	if len(textp1) == 2 {

		assert.True(t, itemExists(textp1, text1a))
		assert.True(t, itemExists(textp1, text1b))
	} else {
		t.Error("Wrong number of page data tokens")
	}

	if len(textp2) == 2 {

		assert.True(t, itemExists(textp2, text2a))
		assert.True(t, itemExists(textp2, text2a))

	} else {
		t.Error("Wrong number of page data tokens")
	}

}

func TestWriteReadDoubleLong(t *testing.T) {

	c := creator.New()
	c.SetPageMargins(0, 0, 0, 0) // we're not printing so use the whole page
	c.SetPageSize(creator.PageSizeA4)

	c.NewPage()
	text1a := strings.Repeat("X", 9999)
	text1b := strings.Repeat("Y", 9999)
	WritePageData(c, text1a)
	WritePageData(c, text1b)

	c.NewPage()
	text2a := strings.Repeat("A", 9999)
	text2b := strings.Repeat("B", 9999)
	WritePageData(c, text2a)
	WritePageData(c, text2b)

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

	textp1, err := ReadPageData(page)

	if err != nil {
		t.Error(err)
	}

	page2, err := pdfReader.GetPage(2)
	if err != nil {
		t.Error(err)
	}

	textp2, err := ReadPageData(page2)
	if err != nil {
		t.Error(err)
	}

	if len(textp1) == 2 {

		assert.True(t, itemExists(textp1, text1a))
		assert.True(t, itemExists(textp1, text1b))
	} else {
		t.Error("Wrong number of page data tokens")
	}

	if len(textp2) == 2 {

		assert.True(t, itemExists(textp2, text2a))
		assert.True(t, itemExists(textp2, text2a))

	} else {
		t.Error("Wrong number of page data tokens")
	}

}

func TestWriteReadOtherText(t *testing.T) {

	c := creator.New()
	c.SetPageMargins(0, 0, 0, 0) // we're not printing so use the whole page
	c.SetPageSize(creator.PageSizeA4)
	c.NewPage()
	text1 := "{\"exam\":\"ENGI99887\",\"number\":\"B12345\",\"page\":1}"
	WritePageData(c, text1)
	text2 := "{\"exam\":\"ENGI99887\",\"number\":\"B12345\",\"page\":2}"
	c.NewPage()
	WritePageData(c, text2)

	p := c.NewParagraph("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	p.SetFontSize(12)
	p.SetPos(0, 0)
	c.Draw(p)

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

	textp1, err := ReadPageData(page)

	if err != nil {
		t.Error(err)
	}

	page2, err := pdfReader.GetPage(2)
	if err != nil {
		t.Error(err)
	}

	textp2, err := ReadPageData(page2)
	if err != nil {
		t.Error(err)
	}

	assertEqual(t, text1, textp1[0])
	assertEqual(t, text2, textp2[0])

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

	text1a := strings.Repeat("X", 9999)
	text1b := strings.Repeat("Y", 9999)
	WritePageData(c, text1a)
	WritePageData(c, text1b)

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

	text1a := strings.Repeat("X", 9999)
	text1b := strings.Repeat("Y", 9999)

	inputPath := "./test/adobe-page-data.pdf"
	field, err := getPdfFieldData(inputPath, "viewer")

	if strings.Compare(field, "type viewer name here") == 0 {
		t.Error("Edit the form data in ./test/adode-page-data.txt and then run test again")
	}

	f, err := os.Open(inputPath)
	if err != nil {
		fmt.Println("Can't open pdf")
		os.Exit(1)
	}

	pdfReader, err := pdf.NewPdfReader(f)
	if err != nil {
		fmt.Println("Can't read test pdf")
		os.Exit(1)
	}

	page, err := pdfReader.GetPage(1)
	if err != nil {
		t.Error(err)
	}

	textp1, err := ReadPageData(page)

	if err != nil {
		t.Error(err)
	}

	if len(textp1) == 2 {

		assert.True(t, itemExists(textp1, text1a))
		assert.True(t, itemExists(textp1, text1b))
	} else {
		t.Error("Wrong number of page data tokens")
	}

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

func TestMarshalling(t *testing.T) {

	pd := PageData{
		Exam: ExamDetails{
			CourseCode: "ENGI12123",
			Diet:       "2020-Summer",
			UUID:       "69197384-fd15-42ac-ac16-82dbe4d52dd0",
		},
		Author: AuthorDetails{
			Anonymous: "B12345",
			Identity:  "e4937a51-4a4a-45b8-bb79-1a841f2b0e78",
		},
		Page: PageDetails{
			UUID:   "a94a71f5-b867-45f9-92f6-ddcc8c39bd9c",
			Number: 15,
		},
	}

	c := creator.New()
	c.SetPageMargins(0, 0, 0, 0) // we're not printing so use the whole page
	c.SetPageSize(creator.PageSizeA4)
	c.NewPage()

	MarshalPageData(c, &pd)

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

	pdout, err := UnmarshalPageData(page)

	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(pd, pdout[0]) {
		t.Error("struct doesn't match")
	}

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
func PrettyPrintStruct(layout interface{}) error {

	json, err := json.MarshalIndent(layout, "", "\t")
	if err != nil {
		return err
	}

	fmt.Println(string(json))
	return nil
}
