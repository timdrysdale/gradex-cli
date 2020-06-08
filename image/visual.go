package image

import (
	"os/exec"
	"strings"
)

func VisuallyIdenticalMultiPagePDFByConvert(pdf1, pdf2 string) (bool, error) {

	out, err := exec.Command("convert", pdf1, "null: ", pdf2, "-compose", "Difference", "-layers", "composite", "-format", "%[fx:mean]\\n", "info:").CombinedOutput()

	result := true

	diffs := strings.Split(string(out), "\n")

	for _, diff := range diffs {
		if diff != "" { //there's a blank line at the end
			if diff != "0" {
				result = false
			}
		}
	}

	return result, err
}
