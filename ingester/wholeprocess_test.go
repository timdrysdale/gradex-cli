package ingester

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/timdrysdale/chmsg"
	"github.com/timdrysdale/gradex-cli/pagedata"
	"github.com/timdrysdale/gradexpath"
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

	expectedTestFiles := []string{
		"bar.jpg",
		"foo.doc",
		"Practice Exam Drop Box_s00000000_attempt_2020-04-23-09-51-20_some thing.pdf",
		"Practice Exam Drop Box_s00000000_attempt_2020-04-23-09-51-20.txt",
		"Practice Exam Drop Box_s00000001_attempt_2020-04-22-08-25-32_a paper.pdf",
		"Practice Exam Drop Box_s00000001_attempt_2020-04-22-08-25-32.txt",
		"Practice Exam Drop Box_s00000002_attempt_2020-04-22-10-43-23_my exam.doc",
		"Practice Exam Drop Box_s00000002_attempt_2020-04-22-10-43-23_my exam.pdf",
		"Practice Exam Drop Box_s00000002_attempt_2020-04-22-10-43-23.txt",
		"Practice Exam Drop Box_s00000003_attempt_2020_one (copy).pdf",
		"Practice Exam Drop Box_s00000003_attempt_2020_one.pdf",
		"Practice Exam Drop Box_s00000003_attempt_2020_one.txt",
		"Practice Exam Drop Box_s00000005_attempt_2020-04-22-11-58-24_Practice Online Exam - Copy (copy).jpg",
		"Practice Exam Drop Box_s00000005_attempt_2020-04-22-11-58-24_Practice Online Exam - Copy.jpg",
		"Practice Exam Drop Box_s00000005_attempt_2020-04-22-11-58-24_Practice Online Exam.jpg",
		"Practice Exam Drop Box_s00000005_attempt_2020-04-22-11-58-24_Practice Online Exam.pdf",
		"Practice Exam Drop Box_s00000005_attempt_2020-04-22-11-58-24 rev.txt",
		"Practice Exam Drop Box_s00000005_attempt_2020-04-22-11-58-24.txt",
	}

	// if you get extra files in the ingest, it can disrupt the tests in this section
	assert.Equal(t, len(expectedTestFiles), len(testfiles))

	assert.True(t, CopyIsComplete(expectedTestFiles, testfiles))

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

	exam := "Practice"

	actualPdf, err = g.GetFileList(g.GetExamDir(exam, acceptedPapers))
	assert.NoError(t, err)
	assert.Equal(t, len(expectedPdf), len(actualPdf))
	assert.True(t, CopyIsComplete(expectedPdf, actualPdf))

	expectedTxt, err = g.GetFileList("./expected/temp-txt-after-validation")
	assert.NoError(t, err)

	actualTxt, err = g.GetFileList(g.GetExamDir(exam, acceptedReceipts))
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
	err = g.FlattenNewPapers(exam)
	assert.NoError(t, err)

	// check files exist

	expectedAnonymousPdf := []string{
		"Practice-B999995.pdf",
		"Practice-B999997.pdf",
		"Practice-B999998.pdf",
		"Practice-B999999.pdf",
	}

	anonymousPdf, err := g.GetFileList(g.GetExamDir(exam, anonPapers))
	assert.NoError(t, err)

	assert.Equal(t, len(expectedAnonymousPdf), len(anonymousPdf))

	assert.True(t, CopyIsComplete(expectedAnonymousPdf, anonymousPdf))

	// check data extraction

	pds, err := pagedata.UnMarshalAllFromFile(anonymousPdf[0])
	assert.NoError(t, err)
	pd := pds[1] //book number 1 for page 1
	assert.Equal(t, exam, pd.Current.Item.What)

	CollectFilesFrom(g.GetExamDir(exam, anonPapers))
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
		"Practice-B999995-maTDD.pdf",
		"Practice-B999997-maTDD.pdf",
		"Practice-B999998-maTDD.pdf",
		"Practice-B999999-maTDD.pdf",
	}

	CollectFilesFrom(g.GetExamDirNamed(exam, markerReady, marker))
	assert.NoError(t, err)

	readyPdf, err := g.GetFileList(g.GetExamDirNamed(exam, markerReady, marker))

	assert.NoError(t, err)

	assert.Equal(t, len(expectedMarker1Pdf), len(readyPdf))

	assert.True(t, CopyIsComplete(expectedMarker1Pdf, readyPdf))

	pds, err = pagedata.UnMarshalAllFromFile(readyPdf[0])
	assert.NoError(t, err)
	pd = pds[1] //book number, 1 for page 1
	assert.Equal(t, pd.Current.Process.ToDo, "marking")

	for _, file := range readyPdf[0:2] {
		destination := filepath.Join(g.GetExamDir(exam, moderatorActive), filepath.Base(file))
		err := Copy(file, destination)
		assert.NoError(t, err)
	}
	for _, file := range readyPdf[2:4] {
		destination := filepath.Join(g.GetExamDir(exam, moderatorInactive), filepath.Base(file))
		err := Copy(file, destination)
		assert.NoError(t, err)
	}

	//>>>>>>>>>>>>>>>>>>>>>>>>> ADD ACTIVE MODERATE BAR  >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
	moderator := "ABC"
	err = g.AddModerateActiveBar(exam, moderator)
	assert.NoError(t, err)

	expectedActive := []string{ //note the d is missing for convenience here
		"Practice-B999995-maTDD-moABC.pdf",
		"Practice-B999997-maTDD-moABC.pdf",
	}

	activePdf, err := g.GetFileList(g.GetExamDirNamed(exam, moderatorReady, moderator))
	assert.NoError(t, err)

	assert.Equal(t, len(expectedActive), len(activePdf))

	assert.True(t, CopyIsComplete(expectedActive, activePdf))

	CollectFilesFrom(g.GetExamDirNamed(exam, moderatorReady, moderator))
	assert.NoError(t, err)
	//>>>>>>>>>>>>>>>>>>>>>>>>> ADD INACTIVE MODERATE BAR  >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
	err = g.AddModerateInActiveBar(exam)
	assert.NoError(t, err)

	expectedInActive := []string{ //note the d is missing for convenience here
		"Practice-B999998-maTDD-moX.pdf",
		"Practice-B999999-maTDD-moX.pdf",
	}

	inActivePdf, err := g.GetFileList(g.GetExamDirSub(exam, moderatorBack, inactive))
	assert.NoError(t, err)

	CollectFilesFrom(g.GetExamDirSub(exam, moderatorBack, inactive))
	assert.NoError(t, err)

	assert.Equal(t, len(expectedInActive), len(inActivePdf))

	assert.True(t, CopyIsComplete(expectedInActive, inActivePdf))

	// copy files to common area (as if have processed them - not checked here)

	for _, file := range activePdf {
		err := g.CopyToDir(file, g.GetExamDir(exam, moderatorProcessed))
		assert.NoError(t, err)
	}

	for _, file := range inActivePdf {
		err := g.CopyToDir(file, g.GetExamDir(exam, moderatorProcessed))
		assert.NoError(t, err)
	}

	expectedModeratedReadyPdf := []string{ //note the d is missing for convenience here
		"Practice-B999995-maTDD-moABC.pdf",
		"Practice-B999997-maTDD-moABC.pdf",
		"Practice-B999998-maTDD-moX.pdf",
		"Practice-B999999-maTDD-moX.pdf",
	}

	moderatedReadyPdf, err := g.GetFileList(g.GetExamDir(exam, moderatorProcessed))
	assert.NoError(t, err)

	assert.Equal(t, len(expectedModeratedReadyPdf), len(moderatedReadyPdf))

	assert.True(t, CopyIsComplete(expectedModeratedReadyPdf, moderatedReadyPdf))

	///>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>

	for _, path := range moderatedReadyPdf {
		err := g.CopyToDir(path, g.GetExamDir(exam, enterProcessed))
		assert.NoError(t, err)
	}

	//>>>>>>>>>>>>>>>>>>>>>>>>> ADD CHECK BAR  >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>

	checker := "LD"

	err = g.AddCheckBar(exam, checker)
	assert.NoError(t, err)
	expectedChecked := []string{ //note the d is missing for convenience here
		"Practice-B999995-maTDD-moABC-chLD.pdf",
		"Practice-B999997-maTDD-moABC-chLD.pdf",
		"Practice-B999998-maTDD-moX-chLD.pdf",
		"Practice-B999999-maTDD-moX-chLD.pdf",
	}

	checkedPdf, err := g.GetFileList(g.GetExamDirNamed(exam, checkerReady, checker))
	assert.NoError(t, err)

	assert.Equal(t, len(expectedChecked), len(checkedPdf))

	assert.True(t, CopyIsComplete(expectedChecked, checkedPdf))
	CollectFilesFrom(g.GetExamDirNamed(exam, checkerReady, checker))
	assert.NoError(t, err)

	// Now do visual checks

	actualPdfs := []string{
		"./tmp-delete-me/usr/exam/Practice/05-anonymous-papers/Practice-B999999.pdf",
		"./tmp-delete-me/usr/exam/Practice/20-marker-ready/TDD/Practice-B999999-maTDD.pdf",
		"./tmp-delete-me/usr/exam/Practice/30-moderator-ready/ABC/Practice-B999995-maTDD-moABC.pdf",
		"./tmp-delete-me/usr/exam/Practice/32-moderator-back/inactive/Practice-B999999-maTDD-moX.pdf",
		"./tmp-delete-me/usr/exam/Practice/50-checker-ready/LD/Practice-B999995-maTDD-moABC-chLD.pdf",
		"./tmp-delete-me/usr/exam/Practice/50-checker-ready/LD/Practice-B999999-maTDD-moX-chLD.pdf",
	}

	expectedPdfs := []string{
		"./expected/visual/Practice-B999999.pdf",
		"./expected/visual/Practice-B999999-maTDD.pdf",
		"./expected/visual/Practice-B999995-maTDD-moABC.pdf",
		"./expected/visual/Practice-B999999-maTDD-moX.pdf",
		"./expected/visual/Practice-B999995-maTDD-moABC-chLD.pdf",
		"./expected/visual/Practice-B999999-maTDD-moX-chLD.pdf", //deliberate mistake
	}

	for i := 0; i < len(actualPdfs); i++ {
		_, err := os.Stat(actualPdfs[i])
		assert.NoError(t, err)
		_, err = os.Stat(expectedPdfs[i])
		assert.NoError(t, err)
		result, err := visuallyIdenticalMultiPagePDF(actualPdfs[i], expectedPdfs[i])
		assert.NoError(t, err)
		assert.True(t, result)
		if !result {
			fmt.Println(actualPdfs[i])
		}

	}

}

// for zsh (escaping of [ not needed on bash https://askubuntu.com/questions/1104907/convert-single-page-from-pdf-to-jpeg-and-getting-error-no-matches-found-binde)
// convert ./tmp-delete-me/usr/exam/Practice/05-anonymous-papers/Practice-B999999.pdf null: ./expected/visual/Practice-B999999.pdf -compose Difference -layers composite -format %\[fx:mean\]\n info:

func visuallyIdenticalMultiPagePDF(pdf1, pdf2 string) (bool, error) {

	//out, err := exec.Command("compare", "-metric", "ae", pdf1, pdf2, diff).CombinedOutput()
	out, err := exec.Command("convert", pdf1, "null: ", pdf2, "-compose", "Difference", "-layers", "composite", "-format", "%[fx:mean]\\n", "info:").CombinedOutput()

	result := true

	diffs := strings.Split(string(out), "\n")

	for _, diff := range diffs {
		//fmt.Printf("%d:[%s]\n", i, diff)
		if diff != "" { //there's a blank line at the end
			if diff != "0" {
				result = false
			}
		}
	}

	return result, err
}
