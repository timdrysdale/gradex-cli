package tree

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/disiqueira/gotree"
	"github.com/timdrysdale/gradexpath"

	"github.com/timdrysdale/gradex-cli/count"
)

func Tree(path string, doPageCount bool) (string, error) {

	treemap := make(map[string]gotree.Tree)
	treemap["."] = gotree.New(filepath.Base(path))

	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {

		if err != nil {
			return err
		}
		if info.IsDir() && strings.Contains(info.Name(), "temp") {
			return filepath.SkipDir
		}
		if info.IsDir() {

			filelist, _ := gradexpath.GetFileList(path)
			numPdf := 0
			numPages := 0
			numDone := 0
			for _, file := range filelist {
				if gradexpath.IsPdf(file) {
					if doPageCount {
						thisCount, err := count.Pages(file)
						if err == nil {
							numPages = numPages + thisCount
						}
					}
					numPdf++
				}
				if gradexpath.IsTxt(file) {
					numPdf++
				}

				if filepath.Ext(file) == ".done" {
					numDone++
				}
			}

			var label string
			if doPageCount {
				label = fmt.Sprintf("%3d/%3d %4d %s", numPdf, numDone, numPages, filepath.Base(path))
				if numPdf == 0 {
					label = fmt.Sprintf("  ─    ─ %s", filepath.Base(path)) //U+00B7 ·
				}
			} else {
				label = fmt.Sprintf("%3d/%3d %s", numPdf, numDone, filepath.Base(path))
				if numPdf == 0 {
					label = fmt.Sprintf("  ─ %s", filepath.Base(path)) //U+00B7 ·
				}

			}

			parent := filepath.Dir(path)

			if parent != "" {

				if _, ok := treemap[parent]; !ok {

					treemap[parent] = (treemap["."]).Add(label)

				} else {
					treemap[path] = (treemap[parent]).Add(fmt.Sprintf(label))
				}
			} else {
				treemap[path] = (treemap["."]).Add(path)
			}
		}

		return nil
	})
	if err != nil {
		return "", err

	}

	return treemap["."].Print(), nil
}
