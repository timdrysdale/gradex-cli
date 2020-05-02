package ingester

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/timdrysdale/copy"
)

//>>>>>>>>>>>>>>>> GENERAL PATHS >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>

func (g *Ingester) MarkedCombined(exam string) string {
	return filepath.Join(g.Exam(), exam, markedCombined)
}
func (g *Ingester) MarkedMerged(exam string) string {
	return filepath.Join(g.Exam(), exam, markedMerged)
}
func (g *Ingester) MarkedPruned(exam string) string {
	return filepath.Join(g.Exam(), exam, markedPruned)
}
func (g *Ingester) MarkedReady(exam string) string {
	return filepath.Join(g.Exam(), exam, markedReady)
}
func (g *Ingester) ModerateActive(exam string) string {
	return filepath.Join(g.Exam(), exam, moderateActive)
}

func (g *Ingester) ModeratedCombined(exam string) string {
	return filepath.Join(g.Exam(), exam, moderatedCombined)
}
func (g *Ingester) ModeratedMerged(exam string) string {
	return filepath.Join(g.Exam(), exam, moderatedMerged)
}
func (g *Ingester) ModeratedPruned(exam string) string {
	return filepath.Join(g.Exam(), exam, moderatedPruned)
}

func (g *Ingester) ModeratedReady(exam string) string {
	return filepath.Join(g.Exam(), exam, moderatedReady)
}

func (g *Ingester) ModerateInActive(exam string) string {
	return filepath.Join(g.Exam(), exam, moderateInActive)
}

func (g *Ingester) ModeratedInActiveBack(exam string) string {
	return filepath.Join(g.Exam(), exam, moderateInActiveBack)
}

func (g *Ingester) CheckedCombined(exam string) string {
	return filepath.Join(g.Exam(), exam, checkedCombined)
}
func (g *Ingester) CheckedMerged(exam string) string {
	return filepath.Join(g.Exam(), exam, checkedMerged)
}
func (g *Ingester) CheckedPruned(exam string) string {
	return filepath.Join(g.Exam(), exam, checkedPruned)
}
func (g *Ingester) CheckedReady(exam string) string {
	return filepath.Join(g.Exam(), exam, checkedReady)
}

func (g *Ingester) DoneDecoration() string {
	return "d"
}

func (g *Ingester) MarkerABCDecoration(initials string) string {
	return fmt.Sprintf("-ma%s", limit(initials, N))
}

func (g *Ingester) MarkerABCDirName(initials string) string {
	return limit(initials, N)
}

func (g *Ingester) ModeratorABCDecoration(initials string) string {
	return fmt.Sprintf("-mo%s", limit(initials, N))
}

func (g *Ingester) ModeratorABCDirName(initials string) string {
	return limit(initials, N)
}

func (g *Ingester) CheckerABCDecoration(initials string) string {
	return fmt.Sprintf("-c%s", limit(initials, N))
}

func (g *Ingester) CheckerABCDirName(initials string) string {
	return limit(initials, N)
}

func (g *Ingester) MarkerNDecoration(number int) string {
	return fmt.Sprintf("-ma%d", number)
}

func (g *Ingester) MarkerNDirName(number int) string {
	return fmt.Sprintf("marker%d", number)
}

func (g *Ingester) ModeratorNDecoration(number int) string {
	return fmt.Sprintf("-mo%d", number)
}

func (g *Ingester) ModeratorNDirName(number int) string {
	return fmt.Sprintf("moderator%d", number)
}

func (g *Ingester) CheckerNDecoration(number int) string {
	return fmt.Sprintf("-c%d", number)
}

func (g *Ingester) CheckerNDirName(number int) string {
	return fmt.Sprintf("checker%d", number)
}

func (g *Ingester) MarkerReady(exam, marker string) string {
	path := filepath.Join(g.Exam(), exam, markerReady, limit(marker, N))
	g.EnsureDirAll(path)
	return path
}

func (g *Ingester) MarkerSent(exam, marker string) string {
	path := filepath.Join(g.Exam(), exam, markerSent, limit(marker, N))
	g.EnsureDirAll(path)
	return path
}

func (g *Ingester) MarkerBack(exam, marker string) string {
	path := filepath.Join(g.Exam(), exam, markerBack, limit(marker, N))
	g.EnsureDirAll(path)
	return path
}

func (g *Ingester) ModeratorReady(exam, moderator string) string {
	path := filepath.Join(g.Exam(), exam, moderatorReady, limit(moderator, N))
	g.EnsureDirAll(path)
	return path
}

func (g *Ingester) ModeratorSent(exam, moderator string) string {
	path := filepath.Join(g.Exam(), exam, moderatorSent, limit(moderator, N))
	g.EnsureDirAll(path)
	return path
}

func (g *Ingester) ModeratorBack(exam, moderator string) string {
	path := filepath.Join(g.Exam(), exam, moderatorBack, limit(moderator, N))
	g.EnsureDirAll(path)
	return path
}

func (g *Ingester) CheckerReady(exam, checker string) string {
	path := filepath.Join(g.Exam(), exam, checkerReady, limit(checker, N))
	g.EnsureDirAll(path)
	return path
}

func (g *Ingester) CheckerSent(exam, checker string) string {
	path := filepath.Join(g.Exam(), exam, checkerSent, limit(checker, N))
	g.EnsureDirAll(path)
	return path
}

func (g *Ingester) CheckerBack(exam, checker string) string {
	path := filepath.Join(g.Exam(), exam, checkerBack, limit(checker, N))
	g.EnsureDirAll(path)
	return path
}

func (g *Ingester) ReMarkerReady(exam, marker string) string {
	path := filepath.Join(g.Exam(), exam, remarkerReady, limit(marker, N))
	g.EnsureDirAll(path)
	return path
}

func (g *Ingester) ReMarkerSent(exam, marker string) string {
	path := filepath.Join(g.Exam(), exam, remarkerSent, limit(marker, N))
	g.EnsureDirAll(path)
	return path
}

func (g *Ingester) ReMarkerBack(exam, marker string) string {
	path := filepath.Join(g.Exam(), exam, remarkerBack, limit(marker, N))
	g.EnsureDirAll(path)
	return path
}

func (g *Ingester) ReCheckerReady(exam, checker string) string {
	path := filepath.Join(g.Exam(), exam, recheckerReady, limit(checker, N))
	g.EnsureDirAll(path)
	return path
}

func (g *Ingester) ReCheckerSent(exam, checker string) string {
	path := filepath.Join(g.Exam(), exam, recheckerSent, limit(checker, N))
	EnsureDirAll(path)
	return path
}

func (g *Ingester) ReCheckerBack(exam, checker string) string {
	path := filepath.Join(g.Exam(), exam, recheckerBack, limit(checker, N))
	g.EnsureDirAll(path)
	return path
}

func (g *Ingester) FlattenLayoutSVG() string {
	return filepath.Join(g.IngestTemplate(), "layout-flatten-312pt.svg")
}

func (g *Ingester) OverlayLayoutSVG() string {
	return filepath.Join(g.OverlayTemplate(), "layout.svg")
}

func (g *Ingester) AcceptedPapers(exam string) string {
	return filepath.Join(g.Exam(), exam, acceptedPapers)
}

func (g *Ingester) AcceptedReceipts(exam string) string {
	return filepath.Join(g.Exam(), exam, acceptedReceipts)
}

//TODO in flatten, swap these paths for the general named ones below
func (g *Ingester) AcceptedPaperImages(exam string) string {
	return filepath.Join(g.Exam(), exam, tempImages)
}
func (g *Ingester) AcceptedPaperPages(exam string) string {
	return filepath.Join(g.Exam(), exam, tempPages)
}
func (g *Ingester) PaperImages(exam string) string {
	return filepath.Join(g.Exam(), exam, tempImages)
}
func (g *Ingester) PaperPages(exam string) string {
	return filepath.Join(g.Exam(), exam, tempPages)
}

func (g *Ingester) AnonymousPapers(exam string) string {
	return filepath.Join(g.Exam(), exam, anonPapers)
}

func (g *Ingester) Identity() string {
	return filepath.Join(g.Etc(), "identity")
}

func (g *Ingester) IdentityCSV() string {
	return filepath.Join(g.Identity(), "identity.csv")
}

func (g *Ingester) Ingest() string {
	return filepath.Join(g.Root(), "ingest")
}

func (g *Ingester) IngestTemplate() string {
	return filepath.Join(g.IngestConf(), "template")
}

func (g *Ingester) OverlayTemplate() string {
	return filepath.Join(g.OverlayConf(), "template")

}
func (g *Ingester) TempPDF() string {
	return filepath.Join(g.Root(), "temp-pdf")
}

func (g *Ingester) TempTXT() string {
	return filepath.Join(g.Root(), "temp-txt")
}

func (g *Ingester) Export() string {
	return filepath.Join(g.Root(), "export")
}

func (g *Ingester) Etc() string {
	return filepath.Join(g.Root(), "etc")
}

func (g *Ingester) Var() string {
	return filepath.Join(g.Root(), "var")
}

func (g *Ingester) Usr() string {
	return filepath.Join(g.Root(), "usr")
}

func (g *Ingester) Exam() string {
	return filepath.Join(g.Usr(), "exam")
}

func (g *Ingester) IngestConf() string {
	return filepath.Join(g.Etc(), "ingest")
}

func (g *Ingester) OverlayConf() string {
	return filepath.Join(g.Etc(), "overlay")
}

func (g *Ingester) ExtractConf() string {
	return filepath.Join(g.Etc(), "extract")
}

func (g *Ingester) SetupConf() string {
	return filepath.Join(g.Etc(), "setup")
}

func (g *Ingester) SetTesting() { //need this when testing other tools
	isTesting = true
}

func (g *Ingester) Root() string {
	return g.root
}

func (g *Ingester) GetExamPath(name string) string {
	return filepath.Join(g.Exam(), name)
}

func (g *Ingester) GetExamStagePath(name, stage string) string {
	return filepath.Join(g.Exam(), name, stage)
}

func (g *Ingester) SetupGradexPaths() error {

	paths := []string{
		g.Root(),
		g.Ingest(),
		g.Identity(),
		g.Export(),
		g.Var(),
		g.Usr(),
		g.Exam(),
		g.TempPDF(),
		g.TempTXT(),
		g.Etc(),
		g.IngestConf(),
		g.OverlayConf(),
		g.IngestTemplate(),
		g.OverlayTemplate(),
		g.ExtractConf(),
		g.SetupConf(),
	}

	for _, path := range paths {

		err := g.EnsureDirAll(path)
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *Ingester) SetupExamPaths(exam string) error {
	// don't use EnsureDirAll so it flags if we are not otherwise setup
	err := g.EnsureDir(g.GetExamPath(exam))
	if err != nil {
		return err
	}

	for _, stage := range ExamStage {
		err := g.EnsureDir(g.GetExamStagePath(exam, stage))
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *Ingester) IsSameAsSelfInDir(source, destinationDir string) bool {

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

func (g *Ingester) IsPDF(path string) bool {
	suffix := strings.ToLower(filepath.Ext(path))
	return strings.Compare(suffix, ".pdf") == 0
}

func (g *Ingester) IsTxt(path string) bool {
	suffix := strings.ToLower(filepath.Ext(path))
	return strings.Compare(suffix, ".txt") == 0
}

func (g *Ingester) IsZip(path string) bool {
	suffix := strings.ToLower(filepath.Ext(path))
	return strings.Compare(suffix, ".zip") == 0
}

func (g *Ingester) IsCsv(path string) bool {
	suffix := strings.ToLower(filepath.Ext(path))
	return strings.Compare(suffix, ".csv") == 0
}

func (g *Ingester) IsArchive(path string) bool {
	archiveExt := []string{".zip", ".tar", ".rar", ".gz", ".br", ".gzip", ".sz", ".zstd", ".lz4", ".xz"}
	return g.ItemExists(archiveExt, filepath.Ext(path))
}

func (g *Ingester) Copy(source, destination string) error {
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

func (g *Ingester) BareFile(name string) string {
	return strings.TrimSuffix(filepath.Base(name), filepath.Ext(name))
}

func (g *Ingester) EnsureDir(dirName string) error {

	err := os.Mkdir(dirName, 0755) //probably umasked with 22 not 02

	os.Chmod(dirName, 0755)

	if err == nil || os.IsExist(err) {
		return nil
	} else {
		return err
	}
}

func (g *Ingester) EnsureDirAll(dirName string) error {

	err := os.MkdirAll(dirName, 0755) //probably umasked with 22 not 02

	os.Chmod(dirName, 0755)

	if err == nil || os.IsExist(err) {
		return nil
	} else {
		return err
	}
}

func (g *Ingester) GetFileListThisDir(dir string) ([]string, error) {

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

func (g *Ingester) GetFileList(dir string) ([]string, error) {

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

func (g *Ingester) CopyIsComplete(source, dest []string) bool {

	sourceBase := g.BaseList(source)
	destBase := g.BaseList(dest)

	for _, item := range sourceBase {

		if !ItemExists(destBase, item) {
			return false
		}
	}

	return true

}

func (g *Ingester) BaseList(paths []string) []string {

	bases := []string{}

	for _, path := range paths {
		bases = append(bases, filepath.Base(path))
	}

	return bases
}

// Mod from array to slice,
// from https://www.golangprograms.com/golang-check-if-array-element-exists.html
func (g *Ingester) ItemExists(sliceType interface{}, item interface{}) bool {
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

func (g *Ingester) GetAnonymousFileName(course, anonymousIdentity string) string {

	return course + "-" + anonymousIdentity + ".pdf"
}
