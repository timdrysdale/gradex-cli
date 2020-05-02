package ingester

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/timdrysdale/chmsg"
	"github.com/timdrysdale/gradexpath"
	"github.com/timdrysdale/pdfpagedata"
)

func CollectFilesFrom(path string) error {
	files, err := gradexpath.GetFileList(path)
	if err != nil {
		return err
	}
	for _, file := range files {
		destination := filepath.Join("./example-output", filepath.Base(file))
		err := gradexpath.Copy(file, destination)
		if err != nil {
			fmt.Printf("ERROR COPYING FILES %v %s %s\n", err, file, destination)
		}
	}
	return err //only tracking last error for this out of convenience
}

func TestAddBars(t *testing.T) {

	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	verbose := false
	collectOutputs := true

	mch := make(chan chmsg.MessageInfo)

	closed := make(chan struct{})
	defer close(closed)
	go func() {
		for {
			select {
			case <-closed:
				break
			case msg := <-mch:
				if verbose {
					fmt.Printf("MC:%s\n", msg.Message)
				}
			}

		}
	}()

	logFile := "./ingester-testing.log"

	f, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer f.Close()

	logger := zerolog.New(f).With().Timestamp().Logger()

	//logger := zerolog.Nop()

	g, err := New("./tmp-delete-me", mch, &logger)
	assert.NoError(t, err)

	assert.Equal(t, "./tmp-delete-me", g.Root())

	//>>>>>>>>>>>>>>>>>>>>>>>>> SETUP >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
	// don't use GetRoot() here
	// JUST in case we kill a whole working installation

	if collectOutputs {
		err := os.RemoveAll("./example-output")
		assert.NoError(t, err)
		err = g.EnsureDir("./example-output")
		assert.NoError(t, err)
	}

	os.RemoveAll("./tmp-delete-me")

	g.EnsureDirectoryStructure()

	testfiles, err := g.GetFileList("./test")

	assert.NoError(t, err)

	for _, file := range testfiles {
		destination := filepath.Join(g.Ingest(), filepath.Base(file))
		err := Copy(file, destination)
		assert.NoError(t, err)

	}

	templateFiles, err := g.GetFileList("./test-fs/etc/ingest/template")
	assert.NoError(t, err)

	for _, file := range templateFiles {
		destination := filepath.Join(g.IngestTemplate(), filepath.Base(file))
		err := Copy(file, destination)
		assert.NoError(t, err)
	}

	ingestfiles, err := g.GetFileList(g.Ingest())
	assert.NoError(t, err)

	assert.True(t, CopyIsComplete(testfiles, ingestfiles))

	//>>>>>>>>>>>>>>>>>>>>>>>>> INGEST >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>

	g.StageFromIngest()

	expectedRejects, err := g.GetFileList("./expected/rejects")
	assert.NoError(t, err)

	actualRejects, err := g.GetFileList(g.Ingest())
	assert.NoError(t, err)

	assert.True(t, len(expectedRejects) == len(actualRejects))
	assert.True(t, CopyIsComplete(expectedRejects, actualRejects))

	expectedTxt, err := g.GetFileList("./expected/temp-txt-after-stage")
	assert.NoError(t, err)

	actualTxt, err := g.GetFileList(g.TempTXT())
	assert.NoError(t, err)

	assert.True(t, len(expectedTxt) == len(actualTxt))
	assert.True(t, CopyIsComplete(expectedTxt, actualTxt))

	expectedPdf, err := g.GetFileList("./expected/temp-pdf-after-stage")
	assert.NoError(t, err)

	actualPdf, err := g.GetFileList(g.TempPDF())
	assert.NoError(t, err)

	assert.Equal(t, len(expectedPdf), len(actualPdf))
	assert.True(t, CopyIsComplete(expectedPdf, actualPdf))

	//>>>>>>>>>>>>>>>>>>>>>>>>> VALIDATE >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
	expectedPdf, err = g.GetFileList("./expected/temp-pdf-after-validation")
	assert.NoError(t, err)
	assert.NoError(t, g.ValidateNewPapers())

	exam := "Practice Exam Drop Box"

	actualPdf, err = g.GetFileList(g.AcceptedPapers(exam))
	assert.NoError(t, err)
	assert.Equal(t, len(expectedPdf), len(actualPdf))
	assert.True(t, CopyIsComplete(expectedPdf, actualPdf))

	expectedTxt, err = g.GetFileList("./expected/temp-txt-after-validation")
	assert.NoError(t, err)

	actualTxt, err = g.GetFileList(g.AcceptedReceipts(exam))
	assert.NoError(t, err)
	assert.Equal(t, len(expectedTxt), len(actualTxt))
	assert.True(t, CopyIsComplete(expectedTxt, actualTxt))

	tempPdf, err := g.GetFileList(g.TempPDF())
	assert.NoError(t, err)
	assert.Equal(t, 0, len(tempPdf))

	tempTxt, err := g.GetFileList(g.TempTXT())
	assert.NoError(t, err)
	assert.Equal(t, 0, len(tempTxt))

	expectedRejects, err = g.GetFileList("./expected/rejects-after-validation")
	actualRejects, err = g.GetFileList(g.Ingest())
	assert.NoError(t, err)
	assert.Equal(t, len(expectedRejects), len(actualRejects))
	assert.True(t, CopyIsComplete(expectedRejects, actualRejects))

	//>>>>>>>>>>>>>>>>>>>>>>>>> SETUP FOR FLATTEN/RENAME  >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
	// Now we test Flatten
	//copy in the identity database
	src := "./test-fs/etc/identity/identity.csv"
	dest := g.IdentityCSV()
	err = Copy(src, dest)
	assert.NoError(t, err)
	_, err = os.Stat(dest)

	//>>>>>>>>>>>>>>>>>>>>>>>>> FLATTEN/RENAME  >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
	err = g.FlattenNewPapers("Practice Exam Drop Box")
	assert.NoError(t, err)

	// check files exist

	expectedAnonymousPdf := []string{
		"Practice Exam Drop Box-B999995.pdf",
		"Practice Exam Drop Box-B999997.pdf",
		"Practice Exam Drop Box-B999998.pdf",
		"Practice Exam Drop Box-B999999.pdf",
	}

	anonymousPdf, err := g.GetFileList(g.AnonymousPapers(exam))
	assert.NoError(t, err)

	assert.Equal(t, len(anonymousPdf), len(expectedAnonymousPdf))

	assert.True(t, CopyIsComplete(expectedAnonymousPdf, anonymousPdf))

	// check data extraction

	pds, err := pdfpagedata.GetPageDataFromFile(anonymousPdf[0])
	assert.NoError(t, err)
	pd := pds[0]
	assert.Equal(t, pd[0].Exam.CourseCode, "Practice Exam Drop Box")

	CollectFilesFrom(g.AnonymousPapers(exam))
	assert.NoError(t, err)
	//>>>>>>>>>>>>>>>>>>>>>>>>> SETUP FOR OVERLAY (via ADDBARS) >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>

	templateFiles, err = g.GetFileList("./test-fs/etc/overlay/template")
	assert.NoError(t, err)

	for _, file := range templateFiles {
		destination := filepath.Join(g.OverlayTemplate(), filepath.Base(file))
		err := Copy(file, destination)
		assert.NoError(t, err)
	}

	//>>>>>>>>>>>>>>>>>>>>>>>>> ADD MARKBAR  >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>

	marker := "tddrysdale"
	err = g.AddMarkBar(exam, marker)
	assert.NoError(t, err)

	expectedMarker1Pdf := []string{
		"Practice Exam Drop Box-B999995-maTDD.pdf",
		"Practice Exam Drop Box-B999997-maTDD.pdf",
		"Practice Exam Drop Box-B999998-maTDD.pdf",
		"Practice Exam Drop Box-B999999-maTDD.pdf",
	}

	CollectFilesFrom(g.MarkerReady(exam, marker))
	assert.NoError(t, err)

	readyPdf, err := g.GetFileList(g.MarkerReady(exam, marker))

	assert.NoError(t, err)

	assert.Equal(t, len(expectedMarker1Pdf), len(readyPdf))

	assert.True(t, CopyIsComplete(expectedMarker1Pdf, readyPdf))

	pds, err = pdfpagedata.GetPageDataFromFile(readyPdf[0])
	assert.NoError(t, err)
	pd = pds[0]
	assert.Equal(t, pd[0].Questions[0].Name, "marking")

	for _, file := range readyPdf[0:2] {
		destination := filepath.Join(g.ModerateActive(exam), filepath.Base(file))
		err := Copy(file, destination)
		assert.NoError(t, err)
	}
	for _, file := range readyPdf[2:4] {
		destination := filepath.Join(g.ModerateInActive(exam), filepath.Base(file))
		err := Copy(file, destination)
		assert.NoError(t, err)
	}

	//>>>>>>>>>>>>>>>>>>>>>>>>> ADD ACTIVE MODERATE BAR  >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
	moderator := "ABC"
	err = g.AddModerateActiveBar(exam, moderator)
	assert.NoError(t, err)

	expectedActive := []string{ //note the d is missing for convenience here
		"Practice Exam Drop Box-B999995-maTDD-moABC.pdf",
		"Practice Exam Drop Box-B999997-maTDD-moABC.pdf",
	}

	activePdf, err := g.GetFileList(g.ModeratorReady(exam, moderator))
	assert.NoError(t, err)

	assert.Equal(t, len(expectedActive), len(activePdf))

	assert.True(t, CopyIsComplete(expectedActive, activePdf))

	CollectFilesFrom(g.ModeratorReady(exam, moderator))
	assert.NoError(t, err)
	//>>>>>>>>>>>>>>>>>>>>>>>>> ADD INACTIVE MODERATE BAR  >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
	err = g.AddModerateInActiveBar(exam)
	assert.NoError(t, err)

	expectedInActive := []string{ //note the d is missing for convenience here
		"Practice Exam Drop Box-B999998-maTDD-moX.pdf",
		"Practice Exam Drop Box-B999999-maTDD-moX.pdf",
	}

	inActivePdf, err := g.GetFileList(g.ModeratedInActiveBack(exam))
	assert.NoError(t, err)

	CollectFilesFrom(g.ModeratedInActiveBack(exam))
	assert.NoError(t, err)

	assert.Equal(t, len(expectedInActive), len(inActivePdf))

	assert.True(t, CopyIsComplete(expectedInActive, inActivePdf))

	// copy files to common area (as if have processed them - not checked here)

	for _, file := range activePdf {
		destination := filepath.Join(g.ModeratedReady(exam), filepath.Base(file))
		err := Copy(file, destination)
		assert.NoError(t, err)
	}

	for _, file := range inActivePdf {
		destination := filepath.Join(g.ModeratedReady(exam), filepath.Base(file))
		err := Copy(file, destination)
		assert.NoError(t, err)
	}

	expectedModeratedReadyPdf := []string{ //note the d is missing for convenience here
		"Practice Exam Drop Box-B999995-maTDD-moABC.pdf",
		"Practice Exam Drop Box-B999997-maTDD-moABC.pdf",
		"Practice Exam Drop Box-B999998-maTDD-moX.pdf",
		"Practice Exam Drop Box-B999999-maTDD-moX.pdf",
	}

	moderatedReadyPdf, err := g.GetFileList(g.ModeratedReady(exam))
	assert.NoError(t, err)

	assert.Equal(t, len(expectedModeratedReadyPdf), len(moderatedReadyPdf))

	assert.True(t, CopyIsComplete(expectedModeratedReadyPdf, moderatedReadyPdf))

	//>>>>>>>>>>>>>>>>>>>>>>>>> ADD CHECK BAR  >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>

	checker := "LD"

	err = g.AddCheckBar(exam, checker)
	assert.NoError(t, err)
	expectedChecked := []string{ //note the d is missing for convenience here
		"Practice Exam Drop Box-B999995-maTDD-moABC-cLD.pdf",
		"Practice Exam Drop Box-B999997-maTDD-moABC-cLD.pdf",
		"Practice Exam Drop Box-B999998-maTDD-moX-cLD.pdf",
		"Practice Exam Drop Box-B999999-maTDD-moX-cLD.pdf",
	}

	checkedPdf, err := g.GetFileList(g.CheckerReady(exam, checker))
	assert.NoError(t, err)

	assert.Equal(t, len(expectedChecked), len(checkedPdf))

	assert.True(t, CopyIsComplete(expectedChecked, checkedPdf))
	CollectFilesFrom(g.CheckerReady(exam, checker))
	assert.NoError(t, err)
}
