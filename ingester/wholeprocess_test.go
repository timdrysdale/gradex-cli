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
)

func CollectFilesFrom(path string) error {
	files, err := GetFileList(path)
	if err != nil {
		return err
	}
	for _, file := range files {
		destination := filepath.Join("./example-output", filepath.Base(file))
		err := Copy(file, destination)
		if err != nil {
			fmt.Printf("ERROR COPYING FILES %v %s %s\n", err, file, destination)
		}
	}
	return err //only tracking last error for this out of convenience
}

func checkPageData(t *testing.T, path string) {

	pdMap, err := pagedata.UnMarshalAllFromFile(path)

	assert.NoError(t, err)

	dataFields := pdMap[1].Current.Data

	mergeMessageCount := 0
	pageOKCount := 0
	pageBadCount := 0
	pageOKOpticalCount := 0
	pageBadOpticalCount := 0

	for _, item := range dataFields {
		if item.Key == "merge-message" {
			mergeMessageCount++
		}
		if item.Key == "tf-page-bad" {
			pageBadCount++
		}
		if item.Key == "tf-page-ok" {
			pageOKCount++
		}
		if item.Key == "tf-page-bad-optical" {
			pageBadOpticalCount++
		}
		if item.Key == "tf-page-ok-optical" {
			pageOKOpticalCount++
		}

	}

	assert.Equal(t, 1, mergeMessageCount)
	assert.Equal(t, 1, pageOKCount)
	assert.Equal(t, 1, pageBadCount)
	assert.Equal(t, 1, pageOKOpticalCount)
	assert.Equal(t, 1, pageBadOpticalCount)
}

func checkNoMergePageData(t *testing.T, path string) {

	pdMap, err := pagedata.UnMarshalAllFromFile(path)

	assert.NoError(t, err)

	dataFields := pdMap[1].Current.Data

	mergeMessageCount := 0
	pageOKCount := 0
	pageBadCount := 0
	pageOKOpticalCount := 0
	pageBadOpticalCount := 0

	for _, item := range dataFields {
		if item.Key == "merge-message" {
			mergeMessageCount++
		}
		if item.Key == "tf-page-bad" {
			pageBadCount++
		}
		if item.Key == "tf-page-ok" {
			pageOKCount++
		}
		if item.Key == "tf-page-bad-optical" {
			pageBadOpticalCount++
		}
		if item.Key == "tf-page-ok-optical" {
			pageOKOpticalCount++
		}

	}

	assert.Equal(t, 0, mergeMessageCount)
	assert.Equal(t, 1, pageOKCount)
	assert.Equal(t, 1, pageBadCount)
	assert.Equal(t, 1, pageOKOpticalCount)
	assert.Equal(t, 1, pageBadOpticalCount)
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

	marker := "tddrysdale" //this gets shortend to TDD
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

	//>>>>>>>>>>>>> export files for "action" (albeit none will be taken) >>>>>>>>>>>>>>>>>>>
	// ingest moves these into the appropriate <stage>Back directory
	stage := marking
	actor := marker
	err = g.ExportFiles(exam, stage, actor)
	assert.NoError(t, err)
	err = g.MoveExportedFilesToIngest(exam, stage, actor)
	assert.NoError(t, err)
	err = g.StageFromIngest()
	assert.NoError(t, err)

	// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> FLATTEN MARKED >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>

	markerBackPdf, err := g.GetFileList(g.GetExamDir(exam, markerBack))
	assert.Equal(t, 4, len(markerBackPdf))

	stage = marked

	err = g.FlattenProcessedPapers(exam, stage)
	assert.NoError(t, err)

	err = g.MergeProcessedPapers(exam, stage) //cmd combines this with flatten
	assert.NoError(t, err)

	// >>>>>>>>>>>>>>>>>>>>>>>> MODERATE SPLIT (MANUAL FOR TEST) >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
	// We do a manual split because we don't have enough files to force the moderation split (>10 files)
	// and that logic has been separately tested so does not need testing here

	processedPdf, err := g.GetFileList(g.GetExamDir(exam, markerProcessed))
	assert.Equal(t, 4, len(processedPdf))

	for _, file := range processedPdf[0:2] {
		destination := filepath.Join(g.GetExamDir(exam, moderatorActive), filepath.Base(file))
		err := Copy(file, destination)
		assert.NoError(t, err)
	}
	for _, file := range processedPdf[2:4] {
		destination := filepath.Join(g.GetExamDir(exam, moderatorInactive), filepath.Base(file))
		err := Copy(file, destination)
		assert.NoError(t, err)
	}

	//>>>>>>>>>>>>>>> Check PageData has (one) merge message >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
	checkPageData(t, processedPdf[0])

	//>>>>>>>>>>>>>>>>>>>>>>>>> ADD ACTIVE MODERATE BAR  >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
	actor = "ABC"
	err = g.AddModerateActiveBar(exam, actor)
	assert.NoError(t, err)

	expectedActive := []string{ //note the d is missing for convenience here
		"Practice-B999995-merge-moABC.pdf",
		"Practice-B999997-merge-moABC.pdf",
	}

	activePdf, err := g.GetFileList(g.GetExamDirNamed(exam, moderatorReady, actor))
	assert.NoError(t, err)

	assert.Equal(t, len(expectedActive), len(activePdf))

	assert.True(t, CopyIsComplete(expectedActive, activePdf))

	CollectFilesFrom(g.GetExamDirNamed(exam, moderatorReady, actor))
	assert.NoError(t, err)

	//>>>>>>>>>>>>> export files for "action" (albeit none will be taken) >>>>>>>>>>>>>>>>>>>
	// ingest moves these into the appropriate <stage>Back directory
	stage = moderating
	readyDir, sentDir, exportDir, err := g.GetExportDirs(exam, stage, actor)
	readyPdf, err = g.GetFileList(readyDir)
	assert.NoError(t, err)
	assert.True(t, CopyIsComplete(expectedActive, readyPdf))

	err = g.ExportFiles(exam, stage, actor)
	assert.NoError(t, err)

	_, sentDir, exportDir, err = g.GetExportDirs(exam, stage, actor)
	assert.NoError(t, err)

	exportedPdf, err := g.GetFileList(exportDir)
	assert.NoError(t, err)
	assert.True(t, CopyIsComplete(expectedActive, exportedPdf))

	sentPdf, err := g.GetFileList(sentDir)
	assert.NoError(t, err)
	assert.True(t, CopyIsComplete(expectedActive, sentPdf))

	err = g.MoveExportedFilesToIngest(exam, stage, actor)
	assert.NoError(t, err)

	err = g.StageFromIngest()
	assert.NoError(t, err)

	//>>>>>>>>>>>>>>>>>>>>>>>>> ADD INACTIVE MODERATE BAR  >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
	err = g.AddModerateInActiveBar(exam)
	assert.NoError(t, err)

	expectedInActive := []string{ //note the d is missing for convenience here
		"Practice-B999998-merge-moX.pdf",
		"Practice-B999999-merge-moX.pdf",
	}

	inActivePdf, err := g.GetFileList(g.GetExamDirSub(exam, moderatorBack, inactive))
	assert.NoError(t, err)

	CollectFilesFrom(g.GetExamDirSub(exam, moderatorBack, inactive))
	assert.NoError(t, err)

	assert.Equal(t, len(expectedInActive), len(inActivePdf))

	assert.True(t, CopyIsComplete(expectedInActive, inActivePdf))

	stage = moderated
	g.FlattenProcessedPapers(exam, stage)
	assert.NoError(t, err)
	g.MergeProcessedPapers(exam, stage)
	assert.NoError(t, err)

	expectedModeratedProcessedPdf := []string{ //note the d is missing for convenience here
		"Practice-B999995-merge.pdf",
		"Practice-B999997-merge.pdf",
		"Practice-B999998-merge.pdf",
		"Practice-B999999-merge.pdf",
	}

	moderatedProcessedPdf, err := g.GetFileList(g.GetExamDir(exam, moderatorProcessed))
	assert.NoError(t, err)

	assert.Equal(t, len(expectedModeratedProcessedPdf), len(moderatedProcessedPdf))

	assert.True(t, CopyIsComplete(expectedModeratedProcessedPdf, moderatedProcessedPdf))

	//>>>>>>>>>>>>>>>>>>>> Check page data for one merge message >>>>>>>>>>>>>>>>>>>>>>>>
	checkPageData(t, moderatedProcessedPdf[0])

	//>>>>>>>>>>>>>>>>>>>>>>>> ADD ENTER BARS >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>

	for _, path := range processedPdf[0:2] {
		err := g.CopyToDir(path, g.GetExamDir(exam, enterActive))
		assert.NoError(t, err)
	}
	for _, path := range processedPdf[2:4] {
		err := g.CopyToDir(path, g.GetExamDir(exam, enterInactive))
		assert.NoError(t, err)
	}
	//for _, path := range readyPdf[2:4] {
	//	err := g.CopyToDir(path, g.GetExamDirSub(exam, enterProcessed, inactive))
	//	assert.NoError(t, err)
	//}

	err = g.AddEnterInactiveBar(exam)
	assert.NoError(t, err)
	expectedInactiveEnter := []string{ //note the d is missing for convenience here
		"Practice-B999998-merge-enX.pdf",
		"Practice-B999999-merge-enX.pdf",
	}

	inactiveEnterPdf, err := g.GetFileList(g.GetExamDirSub(exam, enterBack, inactive))
	assert.NoError(t, err)

	assert.Equal(t, len(expectedInactiveEnter), len(inactiveEnterPdf))

	assert.True(t, CopyIsComplete(expectedInactiveEnter, inactiveEnterPdf))

	//>>>>>>>>>>>>>>>>>>>>>>> ADD ENTER ACTIVE BAR >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>

	actor = "JM"

	err = g.AddEnterActiveBar(exam, actor)

	assert.NoError(t, err)

	expectedActiveEnter := []string{ //note the d is missing for convenience here
		"Practice-B999995-merge-enJM.pdf",
		"Practice-B999997-merge-enJM.pdf",
	}

	activeEnterPdf, err := g.GetFileList(g.GetExamDirNamed(exam, enterReady, actor))
	assert.NoError(t, err)

	assert.Equal(t, len(expectedActiveEnter), len(activeEnterPdf))

	assert.True(t, CopyIsComplete(expectedActiveEnter, activeEnterPdf))

	//>>>>>>>>>>>>> export files for "action" (albeit none will be taken) >>>>>>>>>>>>>>>>>>>
	// ingest moves these into the appropriate <stage>Back directory
	stage = entering
	expectedExport := expectedActiveEnter

	readyDir, sentDir, exportDir, err = g.GetExportDirs(exam, stage, actor)
	readyPdf, err = g.GetFileList(readyDir)
	assert.NoError(t, err)
	assert.True(t, CopyIsComplete(expectedExport, readyPdf))

	err = g.ExportFiles(exam, stage, actor)
	assert.NoError(t, err)

	_, sentDir, exportDir, err = g.GetExportDirs(exam, stage, actor)
	assert.NoError(t, err)

	exportedPdf, err = g.GetFileList(exportDir)
	assert.NoError(t, err)
	assert.True(t, CopyIsComplete(expectedExport, exportedPdf))

	sentPdf, err = g.GetFileList(sentDir)
	assert.NoError(t, err)
	assert.True(t, CopyIsComplete(expectedExport, sentPdf))

	err = g.MoveExportedFilesToIngest(exam, stage, actor)
	assert.NoError(t, err)

	err = g.StageFromIngest()
	assert.NoError(t, err)

	//>>>>>>>>>>>>>>>>>>>>>>>> FLATTEN ACTIVE ENTER >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>

	err = g.FlattenProcessedPapers(exam, "entered")
	assert.NoError(t, err)

	err = g.MergeProcessedPapers(exam, "entered")
	assert.NoError(t, err)

	expectedProcessedEnter := []string{
		"Practice-B999995-merge.pdf",
		"Practice-B999997-merge.pdf",
		"Practice-B999998-merge.pdf",
		"Practice-B999999-merge.pdf",
	}

	processedEnterPdf, err := g.GetFileList(g.GetExamDir(exam, enterProcessed))
	assert.NoError(t, err)

	assert.Equal(t, len(expectedProcessedEnter), len(processedEnterPdf))

	assert.True(t, CopyIsComplete(expectedProcessedEnter, processedEnterPdf))

	//>>>>>>>>>>> Check for one merge message >>>>>>>>>>>>>>
	checkPageData(t, processedEnterPdf[0])

	//>>>>>>>>>>>>>>>>>>>>>>>>> ADD CHECK BAR  >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>

	actor = "LD"

	err = g.AddCheckBar(exam, actor)
	assert.NoError(t, err)
	expectedCheckerReady := []string{ //note the d is missing for convenience here
		"Practice-B999995-merge-chLD.pdf",
		"Practice-B999997-merge-chLD.pdf",
		"Practice-B999998-merge-chLD.pdf",
		"Practice-B999999-merge-chLD.pdf",
	}

	checkerReadyPdf, err := g.GetFileList(g.GetExamDirNamed(exam, checkerReady, actor))
	assert.NoError(t, err)

	assert.Equal(t, len(expectedCheckerReady), len(checkerReadyPdf))

	assert.True(t, CopyIsComplete(expectedCheckerReady, checkerReadyPdf))
	CollectFilesFrom(g.GetExamDirNamed(exam, checkerReady, actor))
	assert.NoError(t, err)

	//>>>>>>>>>>>>> export files for "action" (albeit none will be taken) >>>>>>>>>>>>>>>>>>>
	// ingest moves these into the appropriate <stage>Back directory
	stage = checking
	expectedExport = expectedCheckerReady

	readyDir, sentDir, exportDir, err = g.GetExportDirs(exam, stage, actor)
	readyPdf, err = g.GetFileList(readyDir)
	assert.NoError(t, err)
	assert.True(t, CopyIsComplete(expectedExport, readyPdf))

	err = g.ExportFiles(exam, stage, actor)
	assert.NoError(t, err)

	_, sentDir, exportDir, err = g.GetExportDirs(exam, stage, actor)
	assert.NoError(t, err)

	exportedPdf, err = g.GetFileList(exportDir)
	assert.NoError(t, err)
	assert.True(t, CopyIsComplete(expectedExport, exportedPdf))

	sentPdf, err = g.GetFileList(sentDir)
	assert.NoError(t, err)
	assert.True(t, CopyIsComplete(expectedExport, sentPdf))

	err = g.MoveExportedFilesToIngest(exam, stage, actor)
	assert.NoError(t, err)

	err = g.StageFromIngest()
	assert.NoError(t, err)

	// >>>>>>>>>>>>>>>>>>>>> Processed checked papers >>>>>>>>>>>>>>>>>>>>>>>>>>>
	stage = checked

	err = g.FlattenProcessedPapers(exam, stage)
	assert.NoError(t, err)

	expectedCheckerProcessed := []string{ //note the d is missing for convenience here
		"Practice-B999995-merge-chLD.pdf",
		"Practice-B999997-merge-chLD.pdf",
		"Practice-B999998-merge-chLD.pdf",
		"Practice-B999999-merge-chLD.pdf",
	}

	checkerProcessedPdf, err := g.GetFileList(g.GetExamDir(exam, checkerProcessed))
	assert.NoError(t, err)

	assert.Equal(t, len(expectedCheckerProcessed), len(checkerProcessedPdf))

	assert.True(t, CopyIsComplete(expectedCheckerProcessed, checkerProcessedPdf))
	// we don't collect because have same name, since not merging, so throw errors

	//>>>>>>>>>>> Check for one merge message >>>>>>>>>>>>>>
	checkNoMergePageData(t, checkerProcessedPdf[0])

	// Now do visual checks

	actualPdfs := []string{
		"./tmp-delete-me/usr/exam/Practice/05-anonymous-papers/Practice-B999999.pdf",
		"./tmp-delete-me/usr/exam/Practice/21-marker-sent/TDD/Practice-B999999-maTDD.pdf",
		"./tmp-delete-me/usr/exam/Practice/31-moderator-sent/ABC/Practice-B999995-merge-moABC.pdf",
		"./tmp-delete-me/usr/exam/Practice/32-moderator-back/inactive/Practice-B999999-merge-moX.pdf",
		"./tmp-delete-me/usr/exam/Practice/41-enter-sent/JM/Practice-B999995-merge-enJM.pdf",
		"./tmp-delete-me/usr/exam/Practice/42-enter-back/inactive/Practice-B999999-merge-enX.pdf",
		"./tmp-delete-me/usr/exam/Practice/51-checker-sent/LD/Practice-B999995-merge-chLD.pdf",
		"./tmp-delete-me/usr/exam/Practice/51-checker-sent/LD/Practice-B999999-merge-chLD.pdf",
		"./tmp-delete-me/usr/exam/Practice/54-checker-processed/Practice-B999995-merge-chLD.pdf",
		"./tmp-delete-me/usr/exam/Practice/54-checker-processed/Practice-B999999-merge-chLD.pdf",
	}

	expectedPdfs := []string{
		"./expected/visual/Practice-B999999.pdf",
		"./expected/visual/Practice-B999999-maTDD.pdf",
		"./expected/visual/Practice-B999995-merge-moABC.pdf",
		"./expected/visual/Practice-B999999-merge-moX.pdf",
		"./expected/visual/Practice-B999995-merge-enJM.pdf",
		"./expected/visual/Practice-B999999-merge-enX.pdf",
		"./expected/visual/Practice-B999995-merge-chLD.pdf",
		"./expected/visual/Practice-B999999-merge-chLD.pdf",
		"./expected/visual/Practice-B999995-merge-checkerProcessed.pdf",
		"./expected/visual/Practice-B999999-merge-checkerProcessed.pdf",
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
