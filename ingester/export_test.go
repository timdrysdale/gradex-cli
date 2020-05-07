package ingester

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

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

	g.SetupExamPaths(exam)

	fn := "test.pdf"

	//LABELLING
	fileReady := filepath.Join(g.QuestionReady(exam, actor), fn)
	fileSent := filepath.Join(g.QuestionSent(exam, actor), fn)
	fileExport := filepath.Join(g.ExportLabelling(exam, actor), fn)

	createFile(t, fileReady)

	g.ExportForLabelling(exam, actor)

	mustNotExist(t, fileReady)
	mustExist(t, fileSent)
	mustExist(t, fileExport)

	// MARKING
	fileReady = filepath.Join(g.MarkerReady(exam, actor), fn)
	fileSent = filepath.Join(g.MarkerSent(exam, actor), fn)
	fileExport = filepath.Join(g.ExportMarking(exam, actor), fn)

	createFile(t, fileReady)

	g.ExportForMarking(exam, actor)

	mustNotExist(t, fileReady)
	mustExist(t, fileSent)
	mustExist(t, fileExport)

	// MODERATING
	fileReady = filepath.Join(g.ModeratorReady(exam, actor), fn)
	fileSent = filepath.Join(g.ModeratorSent(exam, actor), fn)
	fileExport = filepath.Join(g.ExportModerating(exam, actor), fn)

	createFile(t, fileReady)

	g.ExportForModerating(exam, actor)

	mustNotExist(t, fileReady)
	mustExist(t, fileSent)
	mustExist(t, fileExport)

	// CHECKING
	fileReady = filepath.Join(g.CheckerReady(exam, actor), fn)
	fileSent = filepath.Join(g.CheckerSent(exam, actor), fn)
	fileExport = filepath.Join(g.ExportChecking(exam, actor), fn)

	createFile(t, fileReady)

	g.ExportForChecking(exam, actor)

	mustNotExist(t, fileReady)
	mustExist(t, fileSent)
	mustExist(t, fileExport)

	// REMARKING
	fileReady = filepath.Join(g.ReMarkerReady(exam, actor), fn)
	fileSent = filepath.Join(g.ReMarkerSent(exam, actor), fn)
	fileExport = filepath.Join(g.ExportReMarking(exam, actor), fn)

	createFile(t, fileReady)

	g.ExportForReMarking(exam, actor)

	mustNotExist(t, fileReady)
	mustExist(t, fileSent)
	mustExist(t, fileExport)

	// RECHECKING
	fileReady = filepath.Join(g.ReCheckerReady(exam, actor), fn)
	fileSent = filepath.Join(g.ReCheckerSent(exam, actor), fn)
	fileExport = filepath.Join(g.ExportReChecking(exam, actor), fn)

	createFile(t, fileReady)

	g.ExportForReChecking(exam, actor)

	mustNotExist(t, fileReady)
	mustExist(t, fileSent)
	mustExist(t, fileExport)

}
