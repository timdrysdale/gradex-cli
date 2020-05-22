package ingester

import (
	"os"
	"path/filepath"
	"strings"

	pdf "github.com/timdrysdale/unipdf/v3/model"
)

// when we read the Learn receipt, we might get a suffix for a word doc etc
// so find the pdf file in the target directory with the same base prefix name
// but possibly variable capitalisation of the suffix (handmade file!)
func GetPDFPath(filename, directory string) (string, error) {

	// if the original receipt says the submission was not pdf
	// we need to find a handmade PDF with possibly non-lower case suffix
	// so search for matching basename
	if !IsPDF(filename) {

		possibleFiles, err := GetFileList(directory)
		if err != nil {
			return "", err
		}

	LOOP:
		for _, file := range possibleFiles {
			want := BareFile(filename)
			got := BareFile(file)
			equal := strings.Compare(want, got) == 0
			if equal {
				filename = file
				break LOOP
			}
		}

	} else { //assume the file is there
		filename = filepath.Join(directory, filename)
	}
	return filename, nil
}

func CountPages(inputPath string) (int, error) {

	numPages := 0

	f, err := os.Open(inputPath)
	if err != nil {
		return numPages, err
	}

	pdfReader, err := pdf.NewPdfReader(f)
	if err != nil {
		return numPages, err
	}

	defer f.Close()

	isEncrypted, err := pdfReader.IsEncrypted()
	if err != nil {
		return numPages, err
	}

	if isEncrypted {
		_, err = pdfReader.Decrypt([]byte(""))
		if err != nil {
			return numPages, err
		}
	}

	numPages, err = pdfReader.GetNumPages()
	if err != nil {
		return numPages, err
	}

	return numPages, nil

}
