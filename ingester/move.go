package ingester

import (
	"os"
	"path/filepath"

	"github.com/rs/zerolog"
)

//returns true if moved
func (g *Ingester) MoveIfNewerThanDestinationInDir(source, destinationDir string, logger *zerolog.Logger) (bool, error) {

	destination := filepath.Join(destinationDir, filepath.Base(source))

	moved, err := g.MoveIfNewerThanDestination(source, destination, logger)

	return moved, err
}

func (g *Ingester) MoveToDir(source, destinationDir string) error {

	destination := filepath.Join(destinationDir, filepath.Base(source))

	return os.Rename(source, destination)
}

// if the source file is not newer, it's not an error
// we just won't move it - anything left we deal with later
// also, we delete "<file>.done" indicator
func (g *Ingester) MoveIfNewerThanDestination(source, destination string, logger *zerolog.Logger) (bool, error) {

	//check both exist
	sourceInfo, err := os.Stat(source)

	if err != nil {
		return false, err
	}

	destinationInfo, err := os.Stat(destination)

	// source newer by definition if destination does not exist
	if os.IsNotExist(err) {
		return true, os.Rename(source, destination)

	}
	if err != nil {
		//TODO  what does this case represent? an error at destination?
		return false, err
	}

	if sourceInfo.ModTime().After(destinationInfo.ModTime()) {

		err = os.Rename(source, destination)

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
