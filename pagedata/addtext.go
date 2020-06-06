package pagedata

import (
	"os"

	"github.com/timdrysdale/unipdf/v3/creator"
	pdf "github.com/timdrysdale/unipdf/v3/model"
)

// TODO check what this can handle:
//Text: Yes
//Image: Yes
//TextField: NO

//modified from https://github.com/unidoc/unipdf-examples/blob/master/text/pdf_insert_text.go
func AddPageDataToPDF(inputPath string, outputPath string, pdMap map[int]PageData) error {
	// Read the input pdf file.
	f, err := os.Open(inputPath)
	if err != nil {
		return err
	}

	pdfReader, err := pdf.NewPdfReader(f)
	if err != nil {
		return err
	}

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return err
	}

	c := creator.New()

	// Load the pages.
	for i := 0; i < numPages; i++ {
		page, err := pdfReader.GetPage(i + 1)
		if err != nil {
			return err
		}

		err = c.AddPage(page)
		if err != nil {
			return err
		}

		if pdfReader.AcroForm != nil {
			err = c.SetForms(pdfReader.AcroForm)
			if err != nil {
				return err
			}

		}

		if pd, ok := pdMap[i+1]; ok { //book pages for index

			err = MarshalOneToCreator(c, &pd)
			if err != nil {
				return err
			}
		}

	}

	f.Close() // just in case we want to overwrite

	return c.WriteToFile(outputPath)
}
