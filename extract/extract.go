/*
 * Get form field data for a specific field from a PDF file.
 *
 * Run as: go run pdf_form_get_field_data <input.pdf> [full field name]
 * If no field specified will output values for all fields.
 */

package extract

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/timdrysdale/gradex-cli/geo"
	"github.com/timdrysdale/unipdf/v3/core"
	pdf "github.com/timdrysdale/unipdf/v3/model"
)

type TextField struct {
	Name    string
	Key     string
	PageNum int
	Value   string
	Rect    []float64
	PageDim geo.Dim
}

func ExtractTextFieldsFromPDF(path string) (map[int]map[string]string, error) {

	fieldsByPage := make(map[int]map[string]string)

	fields, err := mapPdfFieldData(path)
	if err != nil {
		return fieldsByPage, err
	}

	for pagekey, val := range fields {
		p, key := getPageAndKey(pagekey)
		if p > -1 {
			if _, ok := fieldsByPage[p]; !ok { //init map for this page if not present
				fieldsByPage[p] = make(map[string]string)
			}
			fieldsByPage[p][key] = val
		}
	}

	return fieldsByPage, nil
}

// put together the textfields with the right pagesize
func ExtractTextFieldsStructFromPDF(path string) (map[int]map[string]TextField, error) {

	fieldsByPage := make(map[int]map[string]TextField)

	fields, pageSizeList, err := mapTextFields(path)
	if err != nil {
		return fieldsByPage, err
	}

	for pagekey, val := range fields {
		p, key := getPageAndKey(pagekey)
		if p > -1 {
			if _, ok := fieldsByPage[p]; !ok { //init map for this page if not present
				fieldsByPage[p] = make(map[string]TextField)
			}
			val.Key = key
			val.PageNum = p
			pageIndex := p - 1
			if pageIndex < len(pageSizeList) {
				val.PageDim = pageSizeList[pageIndex]
			}
			fieldsByPage[p][key] = val
		}
	}

	return fieldsByPage, nil
}

func boolVal(str string) bool {
	return strings.Compare(str, "") != 0
}

func getPageAndKey(pagekey string) (int, string) {

	// we're looking for the prefix code page-nnn-
	// which may be prepended by prefix docn.

	r := regexp.MustCompile("(?:page-)(\\d{3})-(.*)")

	tokens := r.FindStringSubmatch(pagekey)

	if len(tokens) >= 3 {
		pageStr := tokens[1]
		pageNum, err := strconv.ParseInt(pageStr, 10, 64)
		key := tokens[2]
		if err == nil {
			return int(pageNum), key
		}
	}

	return -1, ""

}

func printPdfFieldData(inputPath, targetFieldName string) error {
	f, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	fmt.Printf("Input file: %s\n", inputPath)

	defer f.Close()

	pdfReader, err := pdf.NewPdfReader(f)
	if err != nil {
		return err
	}

	acroForm := pdfReader.AcroForm
	if acroForm == nil {
		fmt.Printf(" No formdata present\n")
		return nil
	}

	match := false
	fields := acroForm.AllFields()
	for _, field := range fields {
		fullname, err := field.FullName()
		if err != nil {
			return err
		}
		if fullname == targetFieldName || targetFieldName == "" {
			match = true
			if field.V != nil {
				fmt.Printf("Field '%s': '%v' (%T)\n", fullname, field.V, field.V)
			} else {
				fmt.Printf("Field '%s': not filled\n", fullname)
			}
		}
	}

	if !match {
		return errors.New("field not found")
	}
	return nil
}

func mapPdfFieldData(inputPath string) (map[string]string, error) {

	textfields := make(map[string]string)

	f, err := os.Open(inputPath)
	if err != nil {
		return textfields, errors.New(fmt.Sprintf("Problem opening file %s", inputPath))
	}
	defer f.Close()

	pdfReader, err := pdf.NewPdfReader(f)
	if err != nil {
		return textfields, errors.New(fmt.Sprintf("Problem creating reader %s", inputPath))
	}

	acroForm := pdfReader.AcroForm
	if acroForm == nil {
		return textfields, nil
	}

	fields := acroForm.AllFields()
	for _, field := range fields {
		fullname, err := field.FullName()
		if err != nil {
			continue
		}

		val := ""

		if field.V != nil {
			val = field.V.String()
		}

		textfields[fullname] = val

	}

	return textfields, nil
}

func mapTextFields(inputPath string) (map[string]TextField, []geo.Dim, error) {

	textfields := make(map[string]TextField)
	pageSizeList := []geo.Dim{}

	f, err := os.Open(inputPath)
	if err != nil {
		return textfields, pageSizeList, errors.New(fmt.Sprintf("Problem opening file %s", inputPath))
	}
	defer f.Close()

	pdfReader, err := pdf.NewPdfReader(f)
	if err != nil {
		return textfields, pageSizeList, errors.New(fmt.Sprintf("Problem creating reader %s", inputPath))
	}

	for _, page := range pdfReader.PageList {

		rect, err := page.GetMediaBox()

		if err == nil {

			pageSizeList = append(pageSizeList, geo.Dim{
				Width:  rect.Width(),
				Height: rect.Height(),
			})
		}
	}

	acroForm := pdfReader.AcroForm
	if acroForm == nil {
		return textfields, pageSizeList, nil
	}

	fields := acroForm.AllFields()
	for _, field := range fields {
		fullname, err := field.FullName()
		if err != nil {
			continue
		}

		val := ""

		if field.V != nil {
			val = field.V.String()
		}

		annots := field.Annotations

		rect := []float64{}
		// this has always been length one in docs produced so far, so
		// this length check / indexing is just to avoid nul pointer
		// any additional annotations will be ignored, and anyway
		// what would we do with a field two with rects? mmmm
		if len(annots) > 0 {
			rect, err = annots[0].Rect.(*core.PdfObjectArray).ToFloat64Array()
		}

		textfields[fullname] = TextField{
			Name:  fullname,
			Value: val,
			Rect:  rect,
		}

	}

	return textfields, pageSizeList, nil
}

func PrettyPrintStruct(layout interface{}) error {

	json, err := json.MarshalIndent(layout, "", "\t")
	if err != nil {
		return err
	}

	fmt.Println(string(json))
	return nil
}
