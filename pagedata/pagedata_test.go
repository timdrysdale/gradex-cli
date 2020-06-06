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
	"github.com/timdrysdale/gradex-cli/image"
	"github.com/timdrysdale/gradex-cli/merge"
	"github.com/timdrysdale/unipdf/v3/annotator"
	"github.com/timdrysdale/unipdf/v3/creator"
	"github.com/timdrysdale/unipdf/v3/model"
	pdf "github.com/timdrysdale/unipdf/v3/model"
	"github.com/timdrysdale/unipdf/v3/model/optimize"
)

var oldPageDataText = `<gradex-pagedata>{"current":{"is":"page","own":{"path":"/home/tim/gradex/usr/exam/Demo/04-temporary-pages/Demo-B123456-maAAAd0004-me1.pdf","UUID":"b1e6e022-e350-456a-8d30-e868ee3664ec","number":4,"of":4},"original":{"path":"/home/tim/gradex/usr/exam/Demo/04-temporary-pages/Demo-B123456-maAAAd0004.pdf","UUID":"ea404ce4-bb29-49b6-9eff-6fe1aa6b0e98","number":4,"of":4},"current":{"path":"","UUID":"","number":0,"of":0},"item":{"what":"Demo","when":"21-Apr-2020","who":"B123456","UUID":"","whoType":"anonymous"},"process":{"name":"merge-marked","UUID":"1e94533f-84f3-4647-a3a5-61715daa17b7","unixTime":1591245518790914813,"for":"ingester","toDo":"further-processing","by":"gradex-cli","data":null},"UUID":"d35328df-6717-4ed7-a812-4dbe0e7be905","follows":"92858188-485e-48c8-93fd-9b0a52909dae","revision":0,"data":[{"k":"tf-subtotal-05","v":""},{"k":"tf-q2-section","v":""},{"k":"tf-subtotal-09","v":""},{"k":"tf-page-ok","v":""},{"k":"tf-q1-number","v":""},{"k":"tf-subtotal-12","v":""},{"k":"tf-subtotal-07","v":""},{"k":"tf-subtotal-11","v":""},{"k":"tf-subtotal-04","v":""},{"k":"tf-subtotal-01","v":""},{"k":"tf-subtotal-00","v":""},{"k":"tf-subtotal-06","v":""},{"k":"tf-subtotal-10","v":""},{"k":"tf-subtotal-16","v":""},{"k":"tf-subtotal-03","v":""},{"k":"tf-subtotal-14","v":""},{"k":"tf-subtotal-08","v":""},{"k":"tf-subtotal-15","v":""},{"k":"tf-q1-mark","v":""},{"k":"tf-subtotal-02","v":""},{"k":"tf-page-bad","v":""},{"k":"tf-subtotal-13","v":""},{"k":"tf-q2-mark","v":""},{"k":"tf-q1-section","v":""},{"k":"tf-q2-number","v":""},{"k":"tf-subtotal-07-optical","v":""},{"k":"tf-subtotal-11-optical","v":""},{"k":"tf-subtotal-04-optical","v":""},{"k":"tf-subtotal-01-optical","v":""},{"k":"tf-subtotal-00-optical","v":""},{"k":"tf-subtotal-06-optical","v":""},{"k":"tf-subtotal-10-optical","v":""},{"k":"tf-subtotal-16-optical","v":""},{"k":"tf-subtotal-03-optical","v":""},{"k":"tf-subtotal-14-optical","v":""},{"k":"tf-subtotal-08-optical","v":""},{"k":"tf-subtotal-15-optical","v":""},{"k":"tf-q1-mark-optical","v":""},{"k":"tf-subtotal-02-optical","v":""},{"k":"tf-page-bad-optical","v":""},{"k":"tf-subtotal-13-optical","v":""},{"k":"tf-q2-mark-optical","v":""},{"k":"tf-q1-section-optical","v":""},{"k":"tf-q2-number-optical","v":""},{"k":"tf-subtotal-05-optical","v":""},{"k":"tf-q2-section-optical","v":""},{"k":"tf-subtotal-09-optical","v":""},{"k":"tf-page-ok-optical","v":""},{"k":"tf-q1-number-optical","v":""},{"k":"tf-subtotal-12-optical","v":""},{"k":"merge-message","v":"This page skipped by AAA. Marked:[] Bad:[] Seen:[] Skipped:[ AAA]"}],"comments":null},"previous":[{"is":"page","own":{"path":"/home/tim/gradex/usr/exam/Demo/04-temporary-pages/demo-b.pdf0004.pdf","UUID":"c8cf190d-907b-41dc-bf19-69ea12273e83","number":4,"of":4},"original":{"path":"/home/tim/gradex/usr/exam/Demo/03-accepted-papers/demo-b.pdf.pdf","UUID":"56758f14-5a90-4838-baee-7391eee60592","number":4,"of":4},"current":{"path":"","UUID":"","number":0,"of":0},"item":{"what":"Demo","when":"21-Apr-2020","who":"B123456","UUID":"","whoType":"anonymous"},"process":{"name":"flatten","UUID":"d550ef43-51e9-4ed6-8a09-3b76d2d999d5","unixTime":1590980699215789245,"for":"ingester","toDo":"prepare-for-marking","by":"gradex-cli","data":null},"UUID":"83177c32-6284-4613-bfc8-552764ddcd1f","follows":"","revision":0,"data":null,"comments":null},{"is":"page","own":{"path":"/home/tim/gradex/usr/exam/Demo/04-temporary-pages/Demo-B1234560004.pdf","UUID":"609805ce-3183-4031-8c25-594fecfe09c2","number":4,"of":4},"original":{"path":"/home/tim/gradex/usr/exam/Demo/04-temporary-pages/demo-b.pdf0004.pdf","UUID":"c8cf190d-907b-41dc-bf19-69ea12273e83","number":4,"of":4},"current":{"path":"","UUID":"","number":0,"of":0},"item":{"what":"Demo","when":"21-Apr-2020","who":"B123456","UUID":"","whoType":"anonymous"},"process":{"name":"mark-bar","UUID":"d9b73d86-4405-4278-b6c5-d2446c135f94","unixTime":1590980714980176108,"for":"AAA","toDo":"marking","by":"gradex-cli","data":null},"UUID":"16f0b480-2821-4d39-988b-2dd8532ceaf4","follows":"83177c32-6284-4613-bfc8-552764ddcd1f","revision":0,"data":null,"comments":null},{"is":"page","own":{"path":"/home/tim/gradex/usr/exam/Demo/04-temporary-pages/Demo-B123456-maAAAd0004.pdf","UUID":"ea404ce4-bb29-49b6-9eff-6fe1aa6b0e98","number":4,"of":4},"original":{"path":"/home/tim/gradex/usr/exam/Demo/04-temporary-pages/Demo-B1234560004.pdf","UUID":"609805ce-3183-4031-8c25-594fecfe09c2","number":4,"of":4},"current":{"path":"","UUID":"","number":0,"of":0},"item":{"what":"Demo","when":"21-Apr-2020","who":"B123456","UUID":"","whoType":"anonymous"},"process":{"name":"flatten-processed-papers","UUID":"1d04d8d8-590c-47b7-9d40-f0abe3bd5421","unixTime":1590980843335767840,"for":"ingester","toDo":"further-processing","by":"gradex-cli","data":null},"UUID":"92858188-485e-48c8-93fd-9b0a52909dae","follows":"16f0b480-2821-4d39-988b-2dd8532ceaf4","revision":0,"data":[{"k":"tf-subtotal-05","v":""},{"k":"tf-q2-section","v":""},{"k":"tf-subtotal-09","v":""},{"k":"tf-page-ok","v":""},{"k":"tf-q1-number","v":""},{"k":"tf-subtotal-12","v":""},{"k":"tf-subtotal-07","v":""},{"k":"tf-subtotal-11","v":""},{"k":"tf-subtotal-04","v":""},{"k":"tf-subtotal-01","v":""},{"k":"tf-subtotal-00","v":""},{"k":"tf-subtotal-06","v":""},{"k":"tf-subtotal-10","v":""},{"k":"tf-subtotal-16","v":""},{"k":"tf-subtotal-03","v":""},{"k":"tf-subtotal-14","v":""},{"k":"tf-subtotal-08","v":""},{"k":"tf-subtotal-15","v":""},{"k":"tf-q1-mark","v":""},{"k":"tf-subtotal-02","v":""},{"k":"tf-page-bad","v":""},{"k":"tf-subtotal-13","v":""},{"k":"tf-q2-mark","v":""},{"k":"tf-q1-section","v":""},{"k":"tf-q2-number","v":""},{"k":"tf-subtotal-07-optical","v":""},{"k":"tf-subtotal-11-optical","v":""},{"k":"tf-subtotal-04-optical","v":""},{"k":"tf-subtotal-01-optical","v":""},{"k":"tf-subtotal-00-optical","v":""},{"k":"tf-subtotal-06-optical","v":""},{"k":"tf-subtotal-10-optical","v":""},{"k":"tf-subtotal-16-optical","v":""},{"k":"tf-subtotal-03-optical","v":""},{"k":"tf-subtotal-14-optical","v":""},{"k":"tf-subtotal-08-optical","v":""},{"k":"tf-subtotal-15-optical","v":""},{"k":"tf-q1-mark-optical","v":""},{"k":"tf-subtotal-02-optical","v":""},{"k":"tf-page-bad-optical","v":""},{"k":"tf-subtotal-13-optical","v":""},{"k":"tf-q2-mark-optical","v":""},{"k":"tf-q1-section-optical","v":""},{"k":"tf-q2-number-optical","v":""},{"k":"tf-subtotal-05-optical","v":""},{"k":"tf-q2-section-optical","v":""},{"k":"tf-subtotal-09-optical","v":""},{"k":"tf-page-ok-optical","v":""},{"k":"tf-q1-number-optical","v":""},{"k":"tf-subtotal-12-optical","v":""}],"comments":null}]}</gradex-pagedata><hash>3053422766</hash>`

func TestBackwardsCompatible(t *testing.T) {

	// marshall a string that is missing the Revision field
	rawpds := extractPageDatasFromText(oldPageDataText)
	count := 0
	for _, rawpd := range rawpds {

		pd, err := unMarshalPageData(rawpd)
		assert.NoError(t, err)
		assert.Equal(t, 0, pd.Revision)
		count++

	}
	assert.Equal(t, 1, count)
}

func TestSelectHighestRevision(t *testing.T) {

	// make pagedata with different revisions,
	// Check we get the highest revision when we read the file

	pdrev1 := PageData{
		Revision: 1,
		Current: PageDetail{
			Is: IsPage,
			Item: ItemDetail{
				What: "Short name",
			},
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

	outputPath := "./test/revision-test.pdf"

	_, err := os.Stat(outputPath)

	if os.IsNotExist(err) {

		fmt.Printf("Creating %s  please edit in your editor and run tests again\n", outputPath)
		c := creator.New()
		c.SetPageMargins(0, 0, 0, 0) // we're not printing so use the whole page
		c.SetPageSize(creator.PageSizeA4)

		c.NewPage()

		p := c.NewParagraph("Revision Test page for github.com/timdrysdale/pdfpagedata")
		p.SetFontSize(12)
		p.SetPos(200, 10)
		c.Draw(p)

		pd := PageData{
			Current: PageDetail{
				Is: IsPage,
				Item: ItemDetail{
					What: "Long name",
				},
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

		MarshalOneToCreator(c, &pdrev1)

		// write to memory
		var buf bytes.Buffer

		err = c.Write(&buf)
		if err != nil {
			t.Error(err)
		}

		of, err := os.Create(outputPath)
		if err != nil {
			t.Error(err)
		}

		defer of.Close()

		c.SetOptimizer(optimize.New(optimize.Options{
			CombineDuplicateDirectObjects:   true,
			CombineIdenticalIndirectObjects: true,
			CombineDuplicateStreams:         true,
			CompressStreams:                 true,
			UseObjectStreams:                true,
			ImageQuality:                    90,
			ImageUpperPPI:                   150,
		}))

		c.Write(of)
	}

	pdMap, err := UnMarshalAllFromFile(outputPath)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(pdMap))
	assert.Equal(t, pdrev1, pdMap[1])

}

func TestUpdatePageData(t *testing.T) {

	// start with existing file with image, sidebar textfields, and pagedata
	// load, add pagedata, save
	// read pagedata and check updated
	// do visual difference check
	// do textfield check

	textFieldValue := "It's a test field!"
	originalPDF := "./test/test-update.pdf"
	imagePath := "./img/test.jpg"
	updatedPDF := "./test/test-update-rev1.pdf"

	c := creator.New()
	c.SetPageMargins(0, 0, 0, 0) // we're not printing so use the whole page
	c.SetPageSize(creator.PageSizeA4)

	page := c.NewPage()

	img, err := c.NewImageFromFile(imagePath)
	assert.NoError(t, err)

	img.SetPos(150, 50)
	img.SetWidth(300)
	img.SetHeight(300)
	c.Draw(img)

	p := c.NewParagraph("Test update page for github.com/timdrysdale/gradex-cli/pdfpagedata")
	p.SetFontSize(12)
	p.SetPos(100, 10)
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

	form := pdf.NewPdfAcroForm()

	tfopt := annotator.TextFieldOptions{Value: textFieldValue}
	name := fmt.Sprintf("viewer")

	textf, err := annotator.NewTextField(page, name, []float64{100, 200, 150, 250}, tfopt)
	assert.NoError(t, err)

	*form.Fields = append(*form.Fields, textf.PdfField)
	page.AddAnnotation(textf.Annotations[0].PdfAnnotation)

	err = c.SetForms(form)
	assert.NoError(t, err)

	of, err := os.Create(originalPDF)
	assert.NoError(t, err)

	defer of.Close()

	c.SetOptimizer(optimize.New(optimize.Options{
		CombineDuplicateDirectObjects:   true,
		CombineIdenticalIndirectObjects: true,
		CombineDuplicateStreams:         true,
		CompressStreams:                 true,
		UseObjectStreams:                true,
		ImageQuality:                    90,
		ImageUpperPPI:                   150,
	}))

	c.Write(of)

	pdMap := make(map[int]PageData)

	pdMap[1] = PageData{
		Revision: 1,
		Current: PageDetail{
			Is: IsPage,
			Item: ItemDetail{
				What: "Short name",
			},
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

	err = AddPageDataToPDF(originalPDF, updatedPDF, pdMap)

	result, err := image.VisuallyIdenticalMultiPagePDFByConvert(originalPDF, updatedPDF)
	assert.NoError(t, err)
	assert.Equal(t, true, result)

	pdMapUpdated, err := UnMarshalAllFromFile(updatedPDF)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(pdMapUpdated))
	assert.Equal(t, pdMap[1], pdMapUpdated[1])

	field, err := getPdfFieldData(updatedPDF, "viewer")

	assert.NoError(t, err)
	assert.Equal(t, field, textFieldValue)

}

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

func TestSurviveAdvancedMerge(t *testing.T) {

	// write two files with pagedata
	// merge
	// check pagedata ok
	outputPath1 := "./test/merge-test-page-1.pdf"
	outputPath2 := "./test/merge-test-page-2.pdf"
	mergePath := "./test/merge-test-merged.pdf"
	// write page 1
	c := creator.New()
	c.SetPageMargins(0, 0, 0, 0) // we're not printing so use the whole page
	c.SetPageSize(creator.PageSizeA4)

	c.NewPage()

	p := c.NewParagraph("Test page for github.com/timdrysdale/pdfpagedata PAGE 1")
	p.SetFontSize(12)
	p.SetPos(200, 10)
	c.Draw(p)

	pd1 := PageData{
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

	MarshalOneToCreator(c, &pd1)

	of, err := os.Create(outputPath1)
	if err != nil {
		t.Error(err)
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

	c.Write(of)
	of.Close()

	//write page 2
	c = creator.New()
	c.SetPageMargins(0, 0, 0, 0) // we're not printing so use the whole page
	c.SetPageSize(creator.PageSizeA4)

	c.NewPage()

	p = c.NewParagraph("Test page for github.com/timdrysdale/pdfpagedata PAGE 1")
	p.SetFontSize(12)
	p.SetPos(200, 10)
	c.Draw(p)

	pd2 := PageData{
		Current: PageDetail{
			Is:   IsPage,
			UUID: "062510b9-bb95-40b3-bce8-455406f730df",
		},
		Previous: []PageDetail{
			PageDetail{
				Is:   IsPage,
				UUID: "37f631b2-5707-40a9-a2cb-bd08dd830b9e",
				Data: []Field{
					Field{
						Key:   "foo",
						Value: "bar",
					},
				},
			},
		},
	}

	MarshalOneToCreator(c, &pd2)

	of, err = os.Create(outputPath2)
	if err != nil {
		t.Error(err)
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

	c.Write(of)
	of.Close()

	// merge
	inputPaths := []string{outputPath1, outputPath2}
	err = merge.PDF(inputPaths, mergePath)
	assert.NoError(t, err)

	// check

	pdMap, err := UnMarshalAllFromFile(mergePath)

	assert.Equal(t, 2, len(pdMap))
	assert.Equal(t, pd1, pdMap[1])
	assert.Equal(t, pd2, pdMap[2])

}

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
