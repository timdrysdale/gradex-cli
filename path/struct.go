package gradexpath

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/timdrysdale/copy"
)

type GradexPath struct {
	root string
}

type GradexExam struct {
	exam string
	gp   *GradexPath
}

func New(path string) (*GradexPath, error) {

	gp := &GradexPath{}

	_, err := os.Stat(path)
	if err == nil {
		gp.root = path
		err = gp.setupGradexPaths()
	}
	return gp, err
}

func (gp *GradexPath) NewExam(exam string) (*GradexExam, error) {
	ge := &GradexExam{}
	ge.exam = exam
	ge.gp = gp
	err := ge.gp.SetupExamPaths(exam)
	return ge, err

}

func limit(initials string, N int) string {
	if len(initials) < 3 {
		N = len(initials)
	}
	return strings.ToUpper(initials[0:N])
}

//>>>>>>>>>>>>>>>> GENERAL PATHS >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>

func (gp *GradexPath) MarkedCombined(exam string) string {
	return filepath.Join(gp.Exam(), exam, markedCombined)
}
func (gp *GradexPath) MarkedMerged(exam string) string {
	return filepath.Join(gp.Exam(), exam, markedMerged)
}
func (gp *GradexPath) MarkedPruned(exam string) string {
	return filepath.Join(gp.Exam(), exam, markedPruned)
}
func (gp *GradexPath) MarkedReady(exam string) string {
	return filepath.Join(gp.Exam(), exam, markedReady)
}
func (gp *GradexPath) ModerateActive(exam string) string {
	return filepath.Join(gp.Exam(), exam, moderateActive)
}

func (gp *GradexPath) ModeratedCombined(exam string) string {
	return filepath.Join(gp.Exam(), exam, moderatedCombined)
}
func (gp *GradexPath) ModeratedMerged(exam string) string {
	return filepath.Join(gp.Exam(), exam, moderatedMerged)
}
func (gp *GradexPath) ModeratedPruned(exam string) string {
	return filepath.Join(gp.Exam(), exam, moderatedPruned)
}

func (gp *GradexPath) ModeratedReady(exam string) string {
	return filepath.Join(gp.Exam(), exam, moderatedReady)
}

func (gp *GradexPath) ModerateInActive(exam string) string {
	return filepath.Join(gp.Exam(), exam, moderateInActive)
}

func (gp *GradexPath) ModeratedInActiveBack(exam string) string {
	return filepath.Join(gp.Exam(), exam, moderateInActiveBack)
}

func (gp *GradexPath) CheckedCombined(exam string) string {
	return filepath.Join(gp.Exam(), exam, checkedCombined)
}
func (gp *GradexPath) CheckedMerged(exam string) string {
	return filepath.Join(gp.Exam(), exam, checkedMerged)
}
func (gp *GradexPath) CheckedPruned(exam string) string {
	return filepath.Join(gp.Exam(), exam, checkedPruned)
}
func (gp *GradexPath) CheckedReady(exam string) string {
	return filepath.Join(gp.Exam(), exam, checkedReady)
}

func (gp *GradexPath) DoneDecoration() string {
	return "d"
}

func (gp *GradexPath) MarkerABCDecoration(initials string) string {
	return fmt.Sprintf("-ma%s", limit(initials, N))
}

func (gp *GradexPath) MarkerABCDirName(initials string) string {
	return limit(initials, N)
}

func (gp *GradexPath) ModeratorABCDecoration(initials string) string {
	return fmt.Sprintf("-mo%s", limit(initials, N))
}

func (gp *GradexPath) ModeratorABCDirName(initials string) string {
	return limit(initials, N)
}

func (gp *GradexPath) CheckerABCDecoration(initials string) string {
	return fmt.Sprintf("-c%s", limit(initials, N))
}

func (gp *GradexPath) CheckerABCDirName(initials string) string {
	return limit(initials, N)
}

func (gp *GradexPath) MarkerNDecoration(number int) string {
	return fmt.Sprintf("-ma%d", number)
}

func (gp *GradexPath) MarkerNDirName(number int) string {
	return fmt.Sprintf("marker%d", number)
}

func (gp *GradexPath) ModeratorNDecoration(number int) string {
	return fmt.Sprintf("-mo%d", number)
}

func (gp *GradexPath) ModeratorNDirName(number int) string {
	return fmt.Sprintf("moderator%d", number)
}

func (gp *GradexPath) CheckerNDecoration(number int) string {
	return fmt.Sprintf("-c%d", number)
}

func (gp *GradexPath) CheckerNDirName(number int) string {
	return fmt.Sprintf("checker%d", number)
}

func (gp *GradexPath) MarkerReady(exam, marker string) string {
	path := filepath.Join(gp.Exam(), exam, markerReady, limit(marker, N))
	gp.EnsureDirAll(path)
	return path
}

func (gp *GradexPath) MarkerSent(exam, marker string) string {
	path := filepath.Join(gp.Exam(), exam, markerSent, limit(marker, N))
	gp.EnsureDirAll(path)
	return path
}

func (gp *GradexPath) MarkerBack(exam, marker string) string {
	path := filepath.Join(gp.Exam(), exam, markerBack, limit(marker, N))
	gp.EnsureDirAll(path)
	return path
}

func (gp *GradexPath) ModeratorReady(exam, moderator string) string {
	path := filepath.Join(gp.Exam(), exam, moderatorReady, limit(moderator, N))
	gp.EnsureDirAll(path)
	return path
}

func (gp *GradexPath) ModeratorSent(exam, moderator string) string {
	path := filepath.Join(gp.Exam(), exam, moderatorSent, limit(moderator, N))
	gp.EnsureDirAll(path)
	return path
}

func (gp *GradexPath) ModeratorBack(exam, moderator string) string {
	path := filepath.Join(gp.Exam(), exam, moderatorBack, limit(moderator, N))
	gp.EnsureDirAll(path)
	return path
}

func (gp *GradexPath) CheckerReady(exam, checker string) string {
	path := filepath.Join(gp.Exam(), exam, checkerReady, limit(checker, N))
	gp.EnsureDirAll(path)
	return path
}

func (gp *GradexPath) CheckerSent(exam, checker string) string {
	path := filepath.Join(gp.Exam(), exam, checkerSent, limit(checker, N))
	gp.EnsureDirAll(path)
	return path
}

func (gp *GradexPath) CheckerBack(exam, checker string) string {
	path := filepath.Join(gp.Exam(), exam, checkerBack, limit(checker, N))
	gp.EnsureDirAll(path)
	return path
}

func (gp *GradexPath) ReMarkerReady(exam, marker string) string {
	path := filepath.Join(gp.Exam(), exam, remarkerReady, limit(marker, N))
	gp.EnsureDirAll(path)
	return path
}

func (gp *GradexPath) ReMarkerSent(exam, marker string) string {
	path := filepath.Join(gp.Exam(), exam, remarkerSent, limit(marker, N))
	gp.EnsureDirAll(path)
	return path
}

func (gp *GradexPath) ReMarkerBack(exam, marker string) string {
	path := filepath.Join(gp.Exam(), exam, remarkerBack, limit(marker, N))
	gp.EnsureDirAll(path)
	return path
}

func (gp *GradexPath) ReCheckerReady(exam, checker string) string {
	path := filepath.Join(gp.Exam(), exam, recheckerReady, limit(checker, N))
	gp.EnsureDirAll(path)
	return path
}

func (gp *GradexPath) ReCheckerSent(exam, checker string) string {
	path := filepath.Join(gp.Exam(), exam, recheckerSent, limit(checker, N))
	EnsureDirAll(path)
	return path
}

func (gp *GradexPath) ReCheckerBack(exam, checker string) string {
	path := filepath.Join(gp.Exam(), exam, recheckerBack, limit(checker, N))
	gp.EnsureDirAll(path)
	return path
}

func (gp *GradexPath) FlattenLayoutSVG() string {
	return filepath.Join(gp.IngestTemplate(), "layout-flatten-312pt.svg")
}

func (gp *GradexPath) OverlayLayoutSVG() string {
	return filepath.Join(gp.OverlayTemplate(), "layout.svg")
}

func (gp *GradexPath) AcceptedPapers(exam string) string {
	return filepath.Join(gp.Exam(), exam, acceptedPapers)
}

func (gp *GradexPath) AcceptedReceipts(exam string) string {
	return filepath.Join(gp.Exam(), exam, acceptedReceipts)
}

//TODO in flatten, swap these paths for the general named ones below
func (gp *GradexPath) AcceptedPaperImages(exam string) string {
	return filepath.Join(gp.Exam(), exam, tempImages)
}
func (gp *GradexPath) AcceptedPaperPages(exam string) string {
	return filepath.Join(gp.Exam(), exam, tempPages)
}
func (gp *GradexPath) PaperImages(exam string) string {
	return filepath.Join(gp.Exam(), exam, tempImages)
}
func (gp *GradexPath) PaperPages(exam string) string {
	return filepath.Join(gp.Exam(), exam, tempPages)
}

func (gp *GradexPath) AnonymousPapers(exam string) string {
	return filepath.Join(gp.Exam(), exam, anonPapers)
}

func (gp *GradexPath) Identity() string {
	return filepath.Join(gp.Etc(), "identity")
}

func (gp *GradexPath) IdentityCSV() string {
	return filepath.Join(gp.Identity(), "identity.csv")
}

func (gp *GradexPath) Ingest() string {
	return filepath.Join(gp.Root(), "ingest")
}

func (gp *GradexPath) IngestTemplate() string {
	return filepath.Join(gp.IngestConf(), "template")
}

func (gp *GradexPath) OverlayTemplate() string {
	return filepath.Join(gp.OverlayConf(), "template")

}
func (gp *GradexPath) TempPdf() string {
	return filepath.Join(gp.Root(), "temp-pdf")
}

func (gp *GradexPath) TempTxt() string {
	return filepath.Join(gp.Root(), "temp-txt")
}

func (gp *GradexPath) Export() string {
	return filepath.Join(gp.Root(), "export")
}

func (gp *GradexPath) Etc() string {
	return filepath.Join(gp.Root(), "etc")
}

func (gp *GradexPath) Var() string {
	return filepath.Join(gp.Root(), "var")
}

func (gp *GradexPath) Usr() string {
	return filepath.Join(gp.Root(), "usr")
}

func (gp *GradexPath) Exam() string {
	return filepath.Join(gp.Usr(), "exam")
}

func (gp *GradexPath) IngestConf() string {
	return filepath.Join(gp.Etc(), "ingest")
}

func (gp *GradexPath) OverlayConf() string {
	return filepath.Join(gp.Etc(), "overlay")
}

func (gp *GradexPath) ExtractConf() string {
	return filepath.Join(gp.Etc(), "extract")
}

func (gp *GradexPath) SetupConf() string {
	return filepath.Join(gp.Etc(), "setup")
}

func (gp *GradexPath) SetTesting() { //need this when testing other tools
	isTesting = true
}

func (gp *GradexPath) Root() string {
	if isTesting {
		return testroot
	}
	return root
}

func (gp *GradexPath) GetExamPath(name string) string {
	return filepath.Join(gp.Exam(), name)
}

func (gp *GradexPath) GetExamStagePath(name, stage string) string {
	return filepath.Join(gp.Exam(), name, stage)
}

func (gp *GradexPath) setupGradexPaths() error {

	paths := []string{
		gp.Root(),
		gp.Ingest(),
		gp.Identity(),
		gp.Export(),
		gp.Var(),
		gp.Usr(),
		gp.Exam(),
		gp.TempPdf(),
		gp.TempTxt(),
		gp.Etc(),
		gp.IngestConf(),
		gp.OverlayConf(),
		gp.IngestTemplate(),
		gp.OverlayTemplate(),
		gp.ExtractConf(),
		gp.SetupConf(),
	}

	for _, path := range paths {

		err := gp.EnsureDirAll(path)
		if err != nil {
			return err
		}
	}
	return nil
}

func (gp *GradexPath) SetupExamPaths(exam string) error {
	// don't use EnsureDirAll so it flags if we are not otherwise setup
	err := gp.EnsureDir(gp.GetExamPath(exam))
	if err != nil {
		return err
	}

	for _, stage := range ExamStage {
		err := gp.EnsureDir(gp.GetExamStagePath(exam, stage))
		if err != nil {
			return err
		}
	}
	return nil
}

// if the source file is not newer, it's not an error
// we just won't move it - anything left we deal with later
func (gp *GradexPath) MoveIfNewerThanDestination(source, destination string) error {

	//check both exist
	sourceInfo, err := os.Stat(source)

	if err != nil {
		return err
	}

	destinationInfo, err := os.Stat(destination)

	// source newer by definition if destination does not exist
	if os.IsNotExist(err) {
		err = os.Rename(source, destination)
		return err
	}
	if err != nil {
		return err
	}
	if sourceInfo.ModTime().After(destinationInfo.ModTime()) {
		err = os.Rename(source, destination)
		return err
	}

	return nil

}

func (gp *GradexPath) IsSameAsSelfInDir(source, destinationDir string) bool {

	//check both exist
	sourceInfo, err := os.Stat(source)

	if err != nil {
		return false
	}

	destination := filepath.Join(destinationDir, filepath.Base(source))

	destinationInfo, err := os.Stat(destination)

	if err != nil {
		return false
	}

	timeEqual := sourceInfo.ModTime().Equal(destinationInfo.ModTime())
	sizeEqual := sourceInfo.Size() == destinationInfo.Size()
	return timeEqual && sizeEqual

}

func (gp *GradexPath) MoveIfNewerThanDestinationInDir(source, destinationDir string) error {

	//check both exist
	sourceInfo, err := os.Stat(source)

	if err != nil {
		return err
	}

	destination := filepath.Join(destinationDir, filepath.Base(source))

	destinationInfo, err := os.Stat(destination)

	// source newer by definition if destination does not exist
	if os.IsNotExist(err) {
		err = os.Rename(source, destination)
		return err
	}
	if err != nil {
		return err
	}
	if sourceInfo.ModTime().After(destinationInfo.ModTime()) {
		err = os.Rename(source, destination)
		return err
	}

	return nil

}

func (gp *GradexPath) IsPdf(path string) bool {
	suffix := strings.ToLower(filepath.Ext(path))
	return strings.Compare(suffix, ".pdf") == 0
}

func (gp *GradexPath) IsTxt(path string) bool {
	suffix := strings.ToLower(filepath.Ext(path))
	return strings.Compare(suffix, ".txt") == 0
}

func (gp *GradexPath) IsZip(path string) bool {
	suffix := strings.ToLower(filepath.Ext(path))
	return strings.Compare(suffix, ".zip") == 0
}

func (gp *GradexPath) IsCsv(path string) bool {
	suffix := strings.ToLower(filepath.Ext(path))
	return strings.Compare(suffix, ".csv") == 0
}

func (gp *GradexPath) IsArchive(path string) bool {
	archiveExt := []string{".zip", ".tar", ".rar", ".gz", ".br", ".gzip", ".sz", ".zstd", ".lz4", ".xz"}
	return gp.ItemExists(archiveExt, filepath.Ext(path))
}

func (gp *GradexPath) Copy(source, destination string) error {
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

func (gp *GradexPath) BareFile(name string) string {
	return strings.TrimSuffix(filepath.Base(name), filepath.Ext(name))
}

func (gp *GradexPath) EnsureDir(dirName string) error {

	err := os.Mkdir(dirName, 0755) //probably umasked with 22 not 02

	os.Chmod(dirName, 0755)

	if err == nil || os.IsExist(err) {
		return nil
	} else {
		return err
	}
}

func (gp *GradexPath) EnsureDirAll(dirName string) error {

	err := os.MkdirAll(dirName, 0755) //probably umasked with 22 not 02

	os.Chmod(dirName, 0755)

	if err == nil || os.IsExist(err) {
		return nil
	} else {
		return err
	}
}

func (gp *GradexPath) GetFileListThisDir(dir string) ([]string, error) {

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

func (gp *GradexPath) GetFileList(dir string) ([]string, error) {

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

func (gp *GradexPath) CopyIsComplete(source, dest []string) bool {

	sourceBase := gp.BaseList(source)
	destBase := gp.BaseList(dest)

	for _, item := range sourceBase {

		if !ItemExists(destBase, item) {
			return false
		}
	}

	return true

}

func (gp *GradexPath) BaseList(paths []string) []string {

	bases := []string{}

	for _, path := range paths {
		bases = append(bases, filepath.Base(path))
	}

	return bases
}

// Mod from array to slice,
// from https://www.golangprograms.com/golang-check-if-array-element-exists.html
func (gp *GradexPath) ItemExists(sliceType interface{}, item interface{}) bool {
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

func (gp *GradexPath) GetAnonymousFileName(course, anonymousIdentity string) string {

	return course + "-" + anonymousIdentity + ".pdf"
}
