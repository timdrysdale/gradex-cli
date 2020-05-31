package ingester

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog"
	"github.com/timdrysdale/gradex-cli/copy"
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
		return copy.Copy(source, destination, 32*1024, false)
	} else {
		return copy.Copy(source, destination, 1024*1024, false)
	}
}

func CopyOverWrite(source, destination string) error {
	// last param is buffer size ...
	info, err := os.Stat(source)
	if err != nil {
		return err
	}
	if info.Size() > 1024*1024 {
		return copy.Copy(source, destination, 32*1024, true)
	} else {
		return copy.Copy(source, destination, 1024*1024, true)
	}
}

func (g *Ingester) CopyToDir(source, destinationDir string) error {

	err := g.EnsureDirAll(destinationDir)

	if err != nil {
		g.logger.Error().
			Str("source", source).
			Str("destinationDir", destinationDir).
			Str("error", err.Error()).
			Msg("could not ensureDirAll for CopyRoDir")
		return err
	}

	destination := filepath.Join(destinationDir, filepath.Base(source))

	return Copy(source, destination)
}

func BareFile(name string) string {
	return strings.TrimSuffix(filepath.Base(name), filepath.Ext(name))
}

//returns true if moved
func (g *Ingester) CopyIfNewerThanDestinationInDir(source, destinationDir string, logger *zerolog.Logger) (bool, error) {

	destination := filepath.Join(destinationDir, filepath.Base(source))

	copied, err := g.CopyIfNewerThanDestination(source, destination, logger)

	return copied, err
}

// if the source file is not newer, it's not an error
// we just won't move it - anything left we deal with later
// also, we delete "<file>.done" indicator
func (g *Ingester) CopyIfNewerThanDestination(source, destination string, logger *zerolog.Logger) (bool, error) {

	//check both exist
	sourceInfo, err := os.Stat(source)

	if err != nil {
		return false, err
	}

	destinationInfo, err := os.Stat(destination)

	// source newer by definition if destination does not exist
	if os.IsNotExist(err) {
		return true, CopyOverWrite(source, destination)
	}

	//for any other case, let's try and do the copy

	if sourceInfo.ModTime().After(destinationInfo.ModTime()) {

		err = CopyOverWrite(source, destination)

		if err == nil {
			doneFile := doneFilePath(destination)
			_, err := os.Stat(doneFile)

			if err == nil {

				err = os.Remove(doneFile)

				if err != nil {
					logger.Error().
						Str("file", doneFile).
						Str("error", err.Error()).
						Msg("Error removing done file")
					return true, err
				}

			} // no done file to remove
			return true, nil
		} else {
			return false, err
		}
	} else {
		return false, nil //too old
	}
}
