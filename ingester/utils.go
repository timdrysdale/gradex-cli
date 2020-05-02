package ingester

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/timdrysdale/copy"
	"github.com/timdrysdale/gradexpath"
	"github.com/timdrysdale/parselearn"
	pdf "github.com/timdrysdale/unipdf/v3/model"
)

func (g *Ingester) CleanFromIngest() error {

	files, err := g.GetFileList(g.Ingest())
	if err != nil {
		return err
	}
	errorCache := error(nil)

	for _, file := range files {
		err = os.Remove(file)
		if err != nil {
			//count errors?
			errorCache = err
			g.logger.Error().Str("file", file).Msg("Could not delete from Ingest")
		}
	}
	return errorCache
}

// for validate's receipt map
func fileKey(path string) string {
	return strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
}

// This must remain idempotent so we can call it every startup
func (g *Ingester) EnsureDirectoryStructure() error {
	return g.SetupGradexPaths()
}

//need to be case insensitive
func IsPDF(path string) bool {
	suffix := strings.ToLower(filepath.Ext(path))
	return strings.Compare(suffix, ".pdf") == 0
}

func IsTXT(path string) bool {
	suffix := strings.ToLower(filepath.Ext(path))
	return strings.Compare(suffix, ".txt") == 0
}

func IsZIP(path string) bool {
	suffix := strings.ToLower(filepath.Ext(path))
	return strings.Compare(suffix, ".zip") == 0
}

func IsCSV(path string) bool {
	suffix := strings.ToLower(filepath.Ext(path))
	return strings.Compare(suffix, ".csv") == 0
}

func IsArchive(path string) bool {
	archiveExt := []string{".zip", ".tar", ".rar", ".gz", ".br", ".gzip", ".sz", ".zstd", ".lz4", ".xz"}
	return ItemExists(archiveExt, filepath.Ext(path))
}

// when we read the Learn receipt, we might get a suffix for a word doc etc
// so find the pdf file in the target directory with the same base prefix name
// but possibly variable capitalisation of the suffix (handmade file!)
func GetPDFPath(filename, directory string) (string, error) {

	// if the original receipt says the submission was not pdf
	// we need to find a handmade PDF with possibly non-lower case suffix
	// so search for matching basename
	if !IsPDF(filename) {

		possibleFiles, err := gradexpath.GetFileList(directory)
		if err != nil {
			return "", err
		}

	LOOP:
		for _, file := range possibleFiles {
			want := gradexpath.BareFile(filename)
			got := gradexpath.BareFile(file)
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

func GetShortLearnDate(sub parselearn.Submission) (string, error) {

	if sub == (parselearn.Submission{}) {
		return "", errors.New("Empty submission")
	}
	newDate := sub.DateSubmitted
	//Example: "Tuesday, 23 April 2020 10-43-23 o'clock BST"

	tokens := strings.Split(newDate, " ")

	if len(tokens) == 7 {
		day := tokens[1]
		month := tokens[2]
		year := tokens[3]
		if len(month) >= 3 {
			month = month[0:3]
		}
		s := []string{day, month, year}
		newDate = strings.Join(s, "-")
	}

	// no change if we don't understand the format?
	return newDate, nil
}

func checkMatriculation(m string) (bool, error) {
	expectedLength := 8
	actualLength := len(m)
	if actualLength != expectedLength {
		return false, errors.New(fmt.Sprintf("Wrong length got %d not %d", actualLength, expectedLength))
	}
	if strings.HasPrefix(strings.ToLower(m), "s") {
		return false, errors.New(fmt.Sprintf("Does not start with s"))
	}
	return true, nil
}

func checkExamNumber(m string) (bool, error) {
	expectedLength := 7
	actualLength := len(m)
	if actualLength != expectedLength {
		return false, errors.New(fmt.Sprintf("Wrong length got %d not %d", actualLength, expectedLength))
	}

	if strings.HasPrefix(strings.ToLower(m), "b") {
		return false, errors.New(fmt.Sprintf("Does not start with b"))

	}
	return true, nil
}

func Copy(source, destination string) error {
	// last param is buffer size ...
	info, err := os.Stat(source)
	if err != nil {
		return err
	}
	if info.Size() > 1024*1024 {
		return copy.Copy(source, destination, 32*1024)
	} else {
		return copy.Copy(source, destination, 1024*1024)
	}
}

func BareFile(name string) string {
	return strings.TrimSuffix(filepath.Base(name), filepath.Ext(name))
}

func EnsureDir(dirName string) error {

	err := os.Mkdir(dirName, 0755) //probably umasked with 22 not 02

	os.Chmod(dirName, 0755)

	if err == nil || os.IsExist(err) {
		return nil
	} else {
		return err
	}
}

func EnsureDirAll(dirName string) error {

	err := os.MkdirAll(dirName, 0755) //probably umasked with 22 not 02

	os.Chmod(dirName, 0755)

	if err == nil || os.IsExist(err) {
		return nil
	} else {
		return err
	}
}

func GetFileListThisDir(dir string) ([]string, error) {

	paths := []string{}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return filepath.SkipDir
		}

		paths = append(paths, path)

		return nil
	})

	return paths, err

}

func GetFileList(dir string) ([]string, error) {

	paths := []string{}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			paths = append(paths, path)
		}

		return nil
	})

	return paths, err

}

func CopyIsComplete(source, dest []string) bool {

	sourceBase := BaseList(source)
	destBase := BaseList(dest)

	for _, item := range sourceBase {

		if !ItemExists(destBase, item) {
			return false
		}
	}

	return true

}

func BaseList(paths []string) []string {

	bases := []string{}

	for _, path := range paths {
		bases = append(bases, filepath.Base(path))
	}

	return bases
}

// Mod from array to slice,
// from https://www.golangprograms.com/golang-check-if-array-element-exists.html
func ItemExists(sliceType interface{}, item interface{}) bool {
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

func GetAnonymousFileName(course, anonymousIdentity string) string {

	return course + "-" + anonymousIdentity + ".pdf"
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

// create or time from https://golangbyexample.com/touch-file-golang/
func setDone(path string, logger *zerolog.Logger) {

	doneFile := doneFilePath(path)

	_, err := os.Stat(doneFile)
	if os.IsNotExist(err) {
		file, err := os.Create(doneFile)
		if err != nil {
			logger.Error().
				Str("file", doneFile).
				Str("error", err.Error()).
				Msg("Could not write done file")
		}
		defer file.Close()
	} else {
		currentTime := time.Now().Local()
		err = os.Chtimes(doneFile, currentTime, currentTime)
		if err != nil {
			logger.Error().
				Str("file", doneFile).
				Str("error", err.Error()).
				Msg("Could not update time of done file")
		}
	}
}

func getDone(path string) bool {

	donefile := doneFilePath(path)

	_, err := os.Stat(donefile)

	if err == nil {
		return true //done file exists, is done
	} else {
		return false
	}

}

func doneFilePath(path string) string {

	base := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))

	return filepath.Join(filepath.Dir(path), "."+base+".done")
}

// if the source file is not newer, it's not an error
// we just won't move it - anything left we deal with later
// also, we delete "<file>.done" indicator
func (g *Ingester) MoveIfNewerThanDestination(source, destination string, logger *zerolog.Logger) (bool, error) {

	//check both exist
	sourceInfo, err := os.Stat(source)

	if err != nil {
		return false, err
	}

	destinationInfo, err := os.Stat(destination)

	// source newer by definition if destination does not exist
	if os.IsNotExist(err) {
		return true, os.Rename(source, destination)

	}
	if err != nil {
		//TODO  what does this case represent? an error at destination?
		return false, err
	}

	if sourceInfo.ModTime().After(destinationInfo.ModTime()) {

		err = os.Rename(source, destination)

		if err == nil {
			doneFile := doneFilePath(destination)
			_, err := os.Stat(doneFile)

			if err == nil {

				err = os.Remove(doneFile)

				if err != nil {
					logger.Error().
						Str("file", doneFile).
						Str("error", err.Error()).
						Msg("Error removing done file")
					return true, err
				}

			} // no done file to remove
			return true, nil
		} else {
			return false, err
		}
	} else {
		return false, nil //too old
	}
}

//returns true if moved
func (g *Ingester) MoveIfNewerThanDestinationInDir(source, destinationDir string, logger *zerolog.Logger) (bool, error) {

	destination := filepath.Join(destinationDir, filepath.Base(source))

	moved, err := g.MoveIfNewerThanDestination(source, destination, logger)

	return moved, err
}

func (g *Ingester) MoveToDir(source, destinationDir string) error {

	destination := filepath.Join(destinationDir, filepath.Base(source))

	return os.Rename(source, destination)
}
