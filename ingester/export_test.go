package ingester

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/timdrysdale/chmsg"
)

func createFile(t *testing.T, path string) {
	emptyFile, err := os.Create(path)
	assert.NoError(t, err)
	emptyFile.Close()
}

func mustExist(t *testing.T, path string) {

	_, err := os.Stat(path)

	if os.IsNotExist(err) {
		t.Error(fmt.Sprintf("Missing %s", path))
	}
}

func mustNotExist(t *testing.T, path string) {

	_, err := os.Stat(path)

	if os.IsExist(err) {
		t.Error(fmt.Sprintf("Found (unwanted) %s", path))
	}
}

func (g *Ingester) MoveExportedFilesToIngest(exam, stage, actor string) error {
	// this is for testing, so we don't try to salvage any further file moves
	// after we get an error (keep function simple + fail "noisy" in test!)
	_, _, exportDir, err := g.GetExportDirs(exam, stage, actor)

	if err != nil {
		return err
	}
	files, err := g.GetFileList(exportDir)

	if err != nil {
		return err
	}

	for _, file := range files {
		err := UpdateModTime(file) //else it won't looked actioned to stage
		if err != nil {
			return err
		}
		err = g.MoveToDir(file, g.Ingest())
		if err != nil {
			return err
		}
	}
	return nil
}

//https://socketloop.com/tutorials/golang-change-a-file-last-modified-date-and-time
func UpdateModTime(filename string) error {

	_, err := os.Stat(filename)

	if err != nil {
		return err
	}

	// get current timestamp

	currenttime := time.Now().Local()

	// change both atime and mtime to currenttime

	err = os.Chtimes(filename, currenttime, currenttime)

	if err != nil {
		return err
	}
	return nil
}

func TestExport(t *testing.T) {

	logFile := "./ingester-testing.log"

	f, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer f.Close()

	logger := zerolog.New(f).With().Timestamp().Logger()

	mch := make(chan chmsg.MessageInfo)

	g, err := New("./tmp-delete-me", mch, &logger)

	assert.NoError(t, err)

	assert.Equal(t, "./tmp-delete-me", g.Root())

	//>>>>>>>>>>>>>>>>>>>>>>>>> SETUP >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>

	g.EnsureDirectoryStructure()

	exam := "export-test"
	actor := "tester"

	g.SetupExamDirs(exam)

	fn := "test.pdf"

	//LABELLING
	fileReady := filepath.Join(g.GetExamDirNamed(exam, questionReady, actor), fn)
	fileSent := filepath.Join(g.GetExamDirNamed(exam, questionSent, actor), fn)
	fileExport := filepath.Join(g.GetExportDir(exam, labelling, actor), fn)

	createFile(t, fileReady)

	g.ExportFiles(exam, labelling, actor)

	mustNotExist(t, fileReady)
	mustExist(t, fileSent)
	mustExist(t, fileExport)

	// MARKING
	fileReady = filepath.Join(g.GetExamDirNamed(exam, markerReady, actor), fn)
	fileSent = filepath.Join(g.GetExamDirNamed(exam, markerSent, actor), fn)
	fileExport = filepath.Join(g.GetExportDir(exam, marking, actor), fn)

	createFile(t, fileReady)

	g.ExportFiles(exam, marking, actor)

	mustNotExist(t, fileReady)
	mustExist(t, fileSent)
	mustExist(t, fileExport)

	// MODERATING
	fileReady = filepath.Join(g.GetExamDirNamed(exam, moderatorReady, actor), fn)
	fileSent = filepath.Join(g.GetExamDirNamed(exam, moderatorSent, actor), fn)
	fileExport = filepath.Join(g.GetExportDir(exam, moderating, actor), fn)

	createFile(t, fileReady)

	g.ExportFiles(exam, moderating, actor)

	mustNotExist(t, fileReady)
	mustExist(t, fileSent)
	mustExist(t, fileExport)

	// CHECKING
	fileReady = filepath.Join(g.GetExamDirNamed(exam, checkerReady, actor), fn)
	fileSent = filepath.Join(g.GetExamDirNamed(exam, checkerSent, actor), fn)
	fileExport = filepath.Join(g.GetExportDir(exam, checking, actor), fn)

	createFile(t, fileReady)

	g.ExportFiles(exam, checking, actor)

	mustNotExist(t, fileReady)
	mustExist(t, fileSent)
	mustExist(t, fileExport)

	// REMARKING
	fileReady = filepath.Join(g.GetExamDirNamed(exam, reMarkerReady, actor), fn)
	fileSent = filepath.Join(g.GetExamDirNamed(exam, reMarkerSent, actor), fn)
	fileExport = filepath.Join(g.GetExportDir(exam, remarking, actor), fn)

	createFile(t, fileReady)

	g.ExportFiles(exam, remarking, actor)

	mustNotExist(t, fileReady)
	mustExist(t, fileSent)
	mustExist(t, fileExport)

	// RECHECKING
	fileReady = filepath.Join(g.GetExamDirNamed(exam, reCheckerReady, actor), fn)
	fileSent = filepath.Join(g.GetExamDirNamed(exam, reCheckerSent, actor), fn)
	fileExport = filepath.Join(g.GetExportDir(exam, rechecking, actor), fn)

	createFile(t, fileReady)

	g.ExportFiles(exam, rechecking, actor)

	mustNotExist(t, fileReady)
	mustExist(t, fileSent)
	mustExist(t, fileExport)

}
