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

	pdf "github.com/timdrysdale/unipdf/v3/model"
)

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

func PrettyPrintStruct(layout interface{}) error {

	json, err := json.MarshalIndent(layout, "", "\t")
	if err != nil {
		return err
	}

	fmt.Println(string(json))
	return nil
}
