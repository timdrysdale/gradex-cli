package image

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/timdrysdale/gradex-cli/count"
)

func VisuallyIdenticalMultiPagePDF(pdf1, pdf2 string) (bool, error) {

	numPages, err := count.Pages(pdf1)
	if err != nil {
		return false, err
	}
	numPages2, err := count.Pages(pdf2)
	if err != nil {
		return false, err
	}
	if numPages != numPages2 {
		fmt.Printf("Page counts disagree %d != %d\n", numPages, numPages2)
		return false, nil
	}

	// convert to images
	basename := strings.TrimSuffix(filepath.Base(pdf1), filepath.Ext(pdf1))
	optname := fmt.Sprintf("%s%%04d.jpg", basename)
	jpegPath1 := filepath.Join(filepath.Dir(pdf1), optname)
	ConvertPDFToJPEGs(pdf1, filepath.Dir(pdf1), jpegPath1) //"./test"

	basename2 := strings.TrimSuffix(filepath.Base(pdf2), filepath.Ext(pdf1))
	optname2 := fmt.Sprintf("%s%%04d.jpg", basename2)
	jpegPath2 := filepath.Join(filepath.Dir(pdf2), optname2)
	ConvertPDFToJPEGs(pdf2, filepath.Dir(pdf2), jpegPath2)

	// compare each image
	for imgIdx := 1; imgIdx <= numPages; imgIdx = imgIdx + 1 {

		imgPath1 := fmt.Sprintf(jpegPath1, imgIdx)
		imgPath2 := fmt.Sprintf(jpegPath2, imgIdx)

		_, err := os.Stat(imgPath1)
		if err != nil {
			fmt.Printf("Can't find %s\n", imgPath2)
			return false, err
		}
		_, err = os.Stat(imgPath2)
		if err != nil {
			fmt.Printf("Can't find %s\n", imgPath2)
			return false, err
		}

		out, err := exec.Command("compare", "-metric", "ae", imgPath1, imgPath2, "null:").CombinedOutput()

		if err != nil {

			diffPath := filepath.Join(filepath.Dir(imgPath2),
				strings.TrimSuffix(filepath.Base(imgPath2), filepath.Ext(imgPath2))+
					"-diff"+filepath.Ext(imgPath2))
			cmd := exec.Command("compare", imgPath1, imgPath2, diffPath)
			cmd.Run()

			fmt.Printf("Images differ on page %d by %s (metric ae)\n see %s\n", imgIdx, out, diffPath)
			return false, nil
		}
	}
	return true, nil
}
