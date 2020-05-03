package ingester

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/timdrysdale/copy"
)

func CopyIsComplete(source, dest []string) bool {

	sourceBase := BaseList(source)
	destBase := BaseList(dest)

	for _, item := range sourceBase {

		if !ItemExists(destBase, item) {
			return false
		}
	}

	return true

}

func Copy(source, destination string) error {
	// last param is buffer size ...
	info, err := os.Stat(source)
	if err != nil {
		return err
	}
	if info.Size() > 1024*1024 {
		return copy.Copy(source, destination, 32*1024)
	} else {
		return copy.Copy(source, destination, 1024*1024)
	}
}

func (g *Ingester) CopyToDir(source, destinationDir string) error {

	destination := filepath.Join(destinationDir, filepath.Base(source))

	return Copy(source, destination)
}

func BareFile(name string) string {
	return strings.TrimSuffix(filepath.Base(name), filepath.Ext(name))
}
