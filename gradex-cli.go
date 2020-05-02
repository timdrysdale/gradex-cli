package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rs/zerolog"
	"github.com/timdrysdale/chmsg"
	"github.com/timdrysdale/gradex-cli/ingester"
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

func main() {

	verbose := true

	collectOutputs := true

	var err error

	if collectOutputs {
		err = os.RemoveAll("./example-output")
		if err != nil {
			fmt.Printf("Error %v", err)
			os.Exit(1)
		}

		err = gradexpath.EnsureDir("./example-output")
		if err != nil {
			fmt.Printf("Error %v", err)
			os.Exit(1)
		}

	}

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

	root := "/home/tim/gradex/"
	logFile := filepath.Join(root, "var/log/gradex-cli.log")
	f, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer f.Close()
	logger := zerolog.New(f).With().Timestamp().Logger()
	g, err := ingester.New(root, mch, &logger)
	if err != nil {
		fmt.Printf("Failed getting New Ingester %v", err)
		os.Exit(1)
	}

	g.EnsureDirectoryStructure()

	err = g.StageFromIngest()
	if err != nil {
		fmt.Printf("Ingest error %v", err)
		os.Exit(1)
	}

	g.ValidateNewPapers()
	if err != nil {
		fmt.Printf("Validate error %v", err)
		os.Exit(1)
	}

	exam := "PGEE11107 Solar Energy & Photovoltaic Systems (MSc) - Exam Dropbox"

	err = g.FlattenNewPapers(exam)
	if err != nil {
		fmt.Printf("Flatten  error %v", err)
		os.Exit(1)
	}

	CollectFilesFrom(g.AnonymousPapers(exam))

	marker := "tddrysdale"
	err = g.AddMarkBar(exam, marker)
	if err != nil {
		fmt.Printf("Flatten  error %v", err)
		os.Exit(1)
	}

	CollectFilesFrom(g.MarkerReady(exam, marker))

}

/*	assert.NoError(t, err)

	expectedMarker1Pdf := []string{
		"Practice Exam Drop Box-B999995-maTDD.pdf",
		"Practice Exam Drop Box-B999997-maTDD.pdf",
		"Practice Exam Drop Box-B999998-maTDD.pdf",
		"Practice Exam Drop Box-B999999-maTDD.pdf",
	}

	CollectFilesFrom(gradexpath.MarkerReady(exam, marker))
	assert.NoError(t, err)

	readyPdf, err := gradexpath.GetFileList(gradexpath.MarkerReady(exam, marker))

	assert.NoError(t, err)

	assert.Equal(t, len(expectedMarker1Pdf), len(readyPdf))

	assert.True(t, gradexpath.CopyIsComplete(expectedMarker1Pdf, readyPdf))

	pds, err = pdfpagedata.GetPageDataFromFile(readyPdf[0])
	assert.NoError(t, err)
	pd = pds[0]
	assert.Equal(t, pd[0].Questions[0].Name, "marking")

	for _, file := range readyPdf[0:2] {
		destination := filepath.Join(gradexpath.ModerateActive(exam), filepath.Base(file))
		err := gradexpath.Copy(file, destination)
		assert.NoError(t, err)
	}
	for _, file := range readyPdf[2:4] {
		destination := filepath.Join(gradexpath.ModerateInActive(exam), filepath.Base(file))
		err := gradexpath.Copy(file, destination)
		assert.NoError(t, err)
	}

	//>>>>>>>>>>>>>>>>>>>>>>>>> ADD ACTIVE MODERATE BAR  >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
	moderator := "ABC"
	err = AddModerateActiveBar(exam, moderator, mch)
	assert.NoError(t, err)

	expectedActive := []string{ //note the d is missing for convenience here
		"Practice Exam Drop Box-B999995-maTDD-moABC.pdf",
		"Practice Exam Drop Box-B999997-maTDD-moABC.pdf",
	}

	activePdf, err := gradexpath.GetFileList(gradexpath.ModeratorReady(exam, moderator))
	assert.NoError(t, err)

	assert.Equal(t, len(expectedActive), len(activePdf))

	assert.True(t, gradexpath.CopyIsComplete(expectedActive, activePdf))

	CollectFilesFrom(gradexpath.ModeratorReady(exam, moderator))
	assert.NoError(t, err)
	//>>>>>>>>>>>>>>>>>>>>>>>>> ADD INACTIVE MODERATE BAR  >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
	err = AddModerateInActiveBar(exam, mch)
	assert.NoError(t, err)

	expectedInActive := []string{ //note the d is missing for convenience here
		"Practice Exam Drop Box-B999998-maTDD-moX.pdf",
		"Practice Exam Drop Box-B999999-maTDD-moX.pdf",
	}

	inActivePdf, err := gradexpath.GetFileList(gradexpath.ModeratedInActiveBack(exam))
	assert.NoError(t, err)

	CollectFilesFrom(gradexpath.ModeratedInActiveBack(exam))
	assert.NoError(t, err)

	assert.Equal(t, len(expectedInActive), len(inActivePdf))

	assert.True(t, gradexpath.CopyIsComplete(expectedInActive, inActivePdf))

	// copy files to common area (as if have processed them - not checked here)

	for _, file := range activePdf {
		destination := filepath.Join(gradexpath.ModeratedReady(exam), filepath.Base(file))
		err := gradexpath.Copy(file, destination)
		assert.NoError(t, err)
	}

	for _, file := range inActivePdf {
		destination := filepath.Join(gradexpath.ModeratedReady(exam), filepath.Base(file))
		err := gradexpath.Copy(file, destination)
		assert.NoError(t, err)
	}

	expectedModeratedReadyPdf := []string{ //note the d is missing for convenience here
		"Practice Exam Drop Box-B999995-maTDD-moABC.pdf",
		"Practice Exam Drop Box-B999997-maTDD-moABC.pdf",
		"Practice Exam Drop Box-B999998-maTDD-moX.pdf",
		"Practice Exam Drop Box-B999999-maTDD-moX.pdf",
	}

	moderatedReadyPdf, err := gradexpath.GetFileList(gradexpath.ModeratedReady(exam))
	assert.NoError(t, err)

	assert.Equal(t, len(expectedModeratedReadyPdf), len(moderatedReadyPdf))

	assert.True(t, gradexpath.CopyIsComplete(expectedModeratedReadyPdf, moderatedReadyPdf))

	//>>>>>>>>>>>>>>>>>>>>>>>>> ADD CHECK BAR  >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>

	checker := "LD"

	err = AddCheckBar(exam, checker, mch)
	assert.NoError(t, err)
	expectedChecked := []string{ //note the d is missing for convenience here
		"Practice Exam Drop Box-B999995-maTDD-moABC-cLD.pdf",
		"Practice Exam Drop Box-B999997-maTDD-moABC-cLD.pdf",
		"Practice Exam Drop Box-B999998-maTDD-moX-cLD.pdf",
		"Practice Exam Drop Box-B999999-maTDD-moX-cLD.pdf",
	}

	checkedPdf, err := gradexpath.GetFileList(gradexpath.CheckerReady(exam, checker))
	assert.NoError(t, err)

	assert.Equal(t, len(expectedChecked), len(checkedPdf))

	assert.True(t, gradexpath.CopyIsComplete(expectedChecked, checkedPdf))
	CollectFilesFrom(gradexpath.CheckerReady(exam, checker))
	assert.NoError(t, err)
}
*/
