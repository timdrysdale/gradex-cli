package image

import (
	"fmt"
	"os/exec"
	"path/filepath"
)

// simplified https://github.com/catherinelu/evangelist/blob/master/server.go

func ConvertPDFToJPEGs(pdfPath string, jpegPath string, outputFile string) error {

	outputFileOption := fmt.Sprintf("-sOutputFile=%s", outputFile)

	cmd := exec.Command("gs", "-dNOPAUSE", "-sDEVICE=jpeg", outputFileOption, "-dJPEGQ=90", "-r175", "-q", pdfPath,
		"-c", "quit")

	err := cmd.Run()
	if err != nil {
		fmt.Printf("gs command failed: %s\n", err.Error())
		return err
	}

	return nil
}

func CropToQuestion(inputPath, outputPath string) error {
	//350x220+0+110
	cmd := exec.Command("convert", inputPath, "-crop", " 700x500+0+0", filepath.Join(outputPath))

	err := cmd.Run()

	if err != nil {
		fmt.Printf("convert command failed: %s\n", err.Error())
		return err
	}

	return nil
}

// Known good command line
// gs -dNOPAUSE -sDEVICE=jpeg -sOutputFile=edited-%d.jpg -dJPEGQ=95 -r300 -q edited5-covered.pdf -c quit
