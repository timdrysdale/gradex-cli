package ingester

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/timdrysdale/gradex-cli/parselearn"
)

func limit(initials string, N int) string {
	return limitToUpper(initials, N)
}

func limitToUpper(initials string, N int) string {
	if len(initials) < 3 {
		N = len(initials)
	}
	return strings.ToUpper(initials[0:N])
}
func limitToLower(initials string, N int) string {
	if len(initials) < 3 {
		N = len(initials)
	}
	return strings.ToLower(initials[0:N])
}

func shortenBaseFileName(baseFileName string) string {

	getShortLearnName := regexp.MustCompile("\\_([a-zA-Z0-9]*\\_attempt_[0-9-]*)\\_")
	//Before: PGEEnnnn A Super Long Exam Name - Exam Dropbox_s0000000_attempt_2020-05-01-02-00-00_PGEEnnnn-B000000.pdf
	//After _s0000000_attempt_2020-05-01-02-00-00_
	// without picking up any more non-digit therefore safe info due to trailing _ in original filename
	shortLearnNameMatches := getShortLearnName.FindStringSubmatch(baseFileName)

	shortName := baseFileName

	if len(shortLearnNameMatches) > 1 {
		shortName = shortLearnNameMatches[1]
	}

	return shortName

}

//looks for the word between - and .pdf (case insensitive) at the end of the string
func GetAnonymousFromPath(path string) string {
	var anonymous string

	re := regexp.MustCompile("-(B[0-9]*)[-.]")
	matches := re.FindStringSubmatch(filepath.Base(path))
	if len(matches) > 1 {
		anonymous = matches[1]
	}
	return anonymous

}

//looks for the word between - and .pdf (case insensitive) at the end of the string
func GetAnonymousFromPathBasic(path string) string {
	var anonymous string

	re := regexp.MustCompile("-(\\w*).[pP][dD][fF]$")
	matches := re.FindStringSubmatch(path)
	if len(matches) > 1 {
		anonymous = matches[1]
	}
	return anonymous

}

func shortenAssignment(name string) string {

	tokens := strings.Split(name, " ")

	newName := tokens[0]

	if len(newName) > 12 {
		newName = newName[0:12]
	}

	return newName
}

func safeUUID() string {
	UUIDBytes, err := uuid.NewRandom()
	uuidStr := UUIDBytes.String()
	if err != nil {
		uuidStr = fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return uuidStr
}

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
	return g.SetupGradexDirs()
}

func CountPDFInDir(dir string) (int, error) {

	count := 0

	files, err := GetFileList(dir)

	for _, file := range files {
		if IsPDF(file) {
			count++
		}
	}

	return count, err

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
