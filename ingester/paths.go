package ingester

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/timdrysdale/copy"
)

//>>>>>>>>>>>>>> MERGE PROCESSED PAPERS >>>>>>>>>>>>>>>>>>>>>>

func (g *Ingester) MergeProcessedPapersFromDir(exam, stage string) (string, error) {
	return g.FlattenProcessedPapersToDir(exam, stage)
}

func (g *Ingester) MergeProcessedPapersToDir(exam, stage string) (string, error) {

	var dir string

	switch stage {

	case "marked":
		dir = markerProcessed

	case "remarked":
		dir = reMarkerProcessed

	case "moderated":
		dir = moderatorProcessed

	case "remoderated":
		dir = reModeratorProcessed

	case "entered":
		dir = enterProcessed

	case "reentered":
		dir = reEnterProcessed

	case "checked":
		dir = checkerProcessed

	case "rechecked":
		dir = reCheckerProcessed
	default:
		return "", fmt.Errorf("unknown stage %s", stage)
	}

	path := g.GetExamDir(exam, dir)
	g.EnsureDirAll(path)
	return path, nil
}

//>>>>>>>>>>>>>> FLATTEN PROCESSED PAPERS >>>>>>>>>>>>>>>>>>>>>>

func (g *Ingester) FlattenProcessedPapersFromDir(exam, stage string) (string, error) {

	var dir string

	switch stage {

	case "marked":
		dir = markerBack

	case "remarked":
		dir = reMarkerBack

	case "moderated":
		dir = moderatorBack

	case "remoderated":
		dir = reModeratorBack

	case "entered":
		dir = enterBack

	case "reentered":
		dir = reEnterBack

	case "checked":
		dir = checkerBack

	case "rechecked":
		dir = reCheckerBack

	default:
		return "", fmt.Errorf("unknown stage %s", stage)
	}

	path := g.GetExamDir(exam, dir)
	g.EnsureDirAll(path)
	return path, nil
}

func (g *Ingester) FlattenProcessedPapersToDir(exam, stage string) (string, error) {

	var dir string

	switch stage {

	case "marked":
		dir = markerFlattened

	case "remarked":
		dir = reMarkerFlattened

	case "moderated":
		dir = moderatorFlattened

	case "remoderated":
		dir = reModeratorFlattened

	case "entered":
		dir = enterFlattened

	case "reentered":
		dir = reEnterFlattened

	case "checked":
		dir = checkerFlattened

	case "rechecked":
		dir = reCheckerFlattened

	default:
		return "", fmt.Errorf("unknown stage %s", stage)
	}

	path := g.GetExamDir(exam, dir)
	g.EnsureDirAll(path)
	return path, nil
}

func (g *Ingester) GetExamRoot(exam string) string {
	path := filepath.Join(g.Exam(), exam)
	g.EnsureDirAll(path)
	return path
}

func (g *Ingester) GetExamDir(exam, dir string) string {
	path := filepath.Join(g.Exam(), exam, dir)
	g.EnsureDirAll(path)
	return path
}

// note that inactive moderator back would use this function as
// destination := GetExamDir(exam, moderatorBack, inactive)
func (g *Ingester) GetExamDirSub(exam, dir, sub string) string {
	path := filepath.Join(g.Exam(), exam, dir, sub)
	g.EnsureDirAll(path)
	return path
}

func (g *Ingester) GetExportDir(exam, dir, name string) string {
	path := filepath.Join(g.Export(), exam+"-"+dir+"-"+g.GetShortActorName(name))
	g.EnsureDirAll(path)
	return path
}

/* Not sure if this is actually used? replace with doneDecoration
func (g *Ingester) DoneDecoration() string {
	return "d"
}
*/

func (g *Ingester) GetShortActorName(name string) string {
	return GetShortActorName(name)
}

//some external libraries use this, like pagedata
func GetShortActorName(name string) string {
	return limitToUpper(name, 3)
}

func (g *Ingester) GetNamedTaskDecoration(task, name string) string {
	return "-" + limitToLower(task, 2) + GetShortActorName(name)
}

func (g *Ingester) GetExamDirNamed(exam, dir, name string) string {
	path := filepath.Join(g.Exam(), exam, dir, GetShortActorName(name))
	g.EnsureDirAll(path)
	return path
}

// note these rely on info contained in the instantiated ingester
// because they can be set on the command line
func (g *Ingester) FlattenLayoutSVG() string {
	return filepath.Join(g.IngestTemplate(), g.ingestTemplatePath)
}

func (g *Ingester) OverlayLayoutSVG() string {
	return filepath.Join(g.OverlayTemplate(), g.overlayTemplatePath)
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

func (g *Ingester) SetupGradexDirs() error {

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

func (g *Ingester) MigrateVersionDirStruct(exam string, test bool) error {

	err := g.MigrateTempImages(exam, test)

	var lastError error

	if err != nil {
		lastError = fmt.Errorf("Could not migrate temp Images because %s", err)
	}

	err = g.RemoveEmptySubDirs(exam, test)

	if err != nil {
		return fmt.Errorf("Could not remove old sub-directories because %s", err)
	}

	if !test {
		err = g.SetupExamDirs(filepath.Base(exam))

		if err != nil {
			return fmt.Errorf("Could not set up new sub-directories because %s", err)
		}
	}
	return lastError

}

// on update from 0.4 -0.5 temp images moved
func (g *Ingester) MigrateTempImages(exam string, test bool) error {

	from := filepath.Join(exam, "03-temporary-images")

	_, err := os.Stat(from)

	if os.IsNotExist(err) {
		// already done
		return nil
	}

	to := filepath.Join(exam, tempImages)
	var lastError error

	files, err := g.GetFileList(from)

	if err != nil {
		return fmt.Errorf("Could not get list of temporary image files to migrate because %s", err)
	}

	if test {
		fmt.Printf("TEST MODE (not moving)\n Migrating %05d FILES\n FROM: %s\n   TO: %s\n", len(files), from, to)

	} else {

		g.EnsureDir(to)

		for _, file := range files {

			err := g.MoveToDir(file, to)
			if err != nil {
				fmt.Printf("Could not migrate temporary image file %s to %s because %s\n", file, to, err)
				lastError = fmt.Errorf("Could not migrate temporary image file %s to %s because %s", file, to, err)
			}
		}
	}

	return lastError

}

func GetSubDirList(dir string) ([]string, error) {

	paths := []string{}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && path != dir {

			paths = append(paths, path)

			return filepath.SkipDir
		}

		return nil
	})

	return paths, err

}

func (g *Ingester) RemoveEmptySubDirs(dir string, test bool) error {

	emptySubDirs, err := g.GetEmptySubDirs(dir)

	if err != nil {
		return err
	}

	var lastError error

	newStageMap := make(map[string]bool)

	for _, newstage := range ExamStage {
		newStageMap[filepath.Join(dir, newstage)] = true
	}

	for _, emptyDir := range emptySubDirs {

		if emptyDir == dir {
			//don't erase parent
			continue
		}

		if _, ok := newStageMap[emptyDir]; ok {
			// we'll use this dir in new version
			continue
		}

		if test {
			fmt.Printf("TEST MODE (not deleting): %s\n", emptyDir)
		} else {
			thisError := os.Remove(emptyDir)

			if thisError != nil {
				lastError = thisError
			}
		}
	}

	return lastError

}

//https://rosettacode.org/wiki/Empty_directory#Go
func IsEmptyDir(name string) (bool, error) {
	entries, err := ioutil.ReadDir(name)
	if err != nil {
		return false, err
	}
	return len(entries) == 0, nil
}

func (g *Ingester) GetEmptySubDirs(dir string) ([]string, error) {

	paths := []string{}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {

			empty, err := IsEmptyDir(path)
			if err == nil && empty {
				paths = append(paths, path)
			}
		}

		return nil
	})

	return paths, err
}

func (g *Ingester) SetupExamDirs(exam string) error {
	// don't use EnsureDirAll so it flags if we are not otherwise setup
	err := g.EnsureDir(g.GetExamRoot(exam))
	if err != nil {
		return err
	}

	err = g.RemoveEmptySubDirs(g.GetExamRoot(exam), false) //not a test

	if err != nil {
		fmt.Printf("Error cleaning unused directories because %s\n", err.Error())
	}

	for _, stage := range ExamStage {
		err := g.EnsureDirAll(g.GetExamDir(exam, stage))
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
