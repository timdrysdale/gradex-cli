package ingester

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

// >>>>>>>>>>>>>>>>> SINGLE USERS >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
// create or time from https://golangbyexample.com/touch-file-golang/
func setDone(path string, logger *zerolog.Logger) {

	doneFile := doneFilePath(path)

	_, err := os.Stat(doneFile)
	if os.IsNotExist(err) {
		file, err := os.Create(doneFile)
		if err != nil {
			logger.Error().
				Str("file", doneFile).
				Str("error", err.Error()).
				Msg("Could not write done file")
		}
		defer file.Close()
	} else {
		currentTime := time.Now().Local()
		err = os.Chtimes(doneFile, currentTime, currentTime)
		if err != nil {
			logger.Error().
				Str("file", doneFile).
				Str("error", err.Error()).
				Msg("Could not update time of done file")
		}
	}
}

func SetDone(path string, logger *zerolog.Logger) {
	setDone(path, logger)
}
func GetDone(path string) bool {
	return getDone(path)
}

func getDone(path string) bool {

	donefile := doneFilePath(path)

	_, err := os.Stat(donefile)

	if err == nil {
		return true //done file exists, is done
	} else {
		return false
	}
}

func doneFilePath(path string) string {

	base := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))

	return filepath.Join(filepath.Dir(path), "."+base+".done")
}

// >>>>>>>>>>>> MULTIPLE USERS >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>

// create or time from https://golangbyexample.com/touch-file-golang/
func setDoneFor(path, who string, logger *zerolog.Logger) {

	doneFile := doneFilePathFor(path, who)

	_, err := os.Stat(doneFile)
	if os.IsNotExist(err) {
		file, err := os.Create(doneFile)
		if err != nil {
			logger.Error().
				Str("file", doneFile).
				Str("error", err.Error()).
				Msg("Could not write done file")
		}
		defer file.Close()
	} else {
		currentTime := time.Now().Local()
		err = os.Chtimes(doneFile, currentTime, currentTime)
		if err != nil {
			logger.Error().
				Str("file", doneFile).
				Str("error", err.Error()).
				Msg("Could not update time of done file")
		}
	}
}

func getDoneFor(path, who string) bool {

	donefile := doneFilePathFor(path, who)

	_, err := os.Stat(donefile)

	if err == nil {
		return true //done file exists, is done
	} else {
		return false
	}
}

func doneFilePathFor(path, who string) string {
	who = strings.TrimPrefix(who, "-")
	base := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))

	return filepath.Join(filepath.Dir(path), "."+who+"."+base+".done")
}
