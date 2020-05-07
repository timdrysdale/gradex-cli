package ingester

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

func (g *Ingester) ExportForLabelling(exam, labeller string) {

	source := g.QuestionReady(exam, labeller)

	g.logger.Info().Msg("Exporting")

	files, err := GetFileList(source)

	if err != nil {

		g.logger.Error().
			Str("error", err.Error()).
			Msg("could not get file list for exporting")
	}

	g.logger.Info().Msg("Exporting")

	numErrors := 0

	destination := g.ExportLabelling(exam, labeller)

	for _, file := range files {

		if !g.IsPDF(file) {
			continue
		}

		destination = g.ExportLabelling(exam, labeller)

		err := g.CopyToDir(file, destination)
		if err == nil {

			destination = g.QuestionSent(exam, labeller)

			err = g.MoveToDir(file, destination)

			if err != nil {
				g.logger.Error().
					Str("file", file).
					Str("destination", destination).
					Str("error", err.Error()).
					Msg("could not copy file to LabellerSent")

			}

		} else {
			numErrors++
			g.logger.Error().
				Str("file", file).
				Str("destination", destination).
				Str("error", err.Error()).
				Msg("could not copy file to export it")

		}

	}
	if numErrors == 0 {
		g.logger.Info().
			Int("count", len(files)).
			Str("destination", destination).
			Msg(fmt.Sprintf("Exported %d files to %s", len(files), destination))
	}

}

func (g *Ingester) ExportForMarking(exam, marker string) {

	source := g.MarkerReady(exam, marker)

	files, err := GetFileList(source)

	if err != nil {

		g.logger.Error().
			Str("error", err.Error()).
			Msg("could not get file list for exporting")
	}

	g.logger.Info().Msg("Exporting")

	numErrors := 0

	destination := g.ExportMarking(exam, marker)

	for _, file := range files {

		if !g.IsPDF(file) {
			continue
		}

		destination = g.ExportMarking(exam, marker)

		err := g.CopyToDir(file, destination)
		if err == nil {

			destination = g.MarkerSent(exam, marker)

			err = g.MoveToDir(file, destination)

			if err != nil {
				g.logger.Error().
					Str("file", file).
					Str("destination", destination).
					Str("error", err.Error()).
					Msg("could not copy file to MarkerSent")

			}

		} else {
			numErrors++
			g.logger.Error().
				Str("file", file).
				Str("destination", destination).
				Str("error", err.Error()).
				Msg("could not copy file to export it")

		}

	}
	if numErrors == 0 {
		g.logger.Info().
			Int("count", len(files)).
			Str("destination", destination).
			Msg(fmt.Sprintf("Exported %d files to %s", len(files), destination))
	}

}

func (g *Ingester) ExportForModerating(exam, moderator string) {

	files, err := GetFileList(g.ModeratorReady(exam, moderator))

	if err != nil {

		g.logger.Error().
			Str("error", err.Error()).
			Msg("could not get file list for exporting")
	}
	g.logger.Info().Msg("Exporting")
	numErrors := 0

	destination := g.ExportModerating(exam, moderator)

	for _, file := range files {

		if !g.IsPDF(file) {
			continue
		}

		destination = g.ExportModerating(exam, moderator)

		err := g.CopyToDir(file, destination)
		if err == nil {

			destination = g.ModeratorSent(exam, moderator)
			err = g.MoveToDir(file, destination)
			if err != nil {
				g.logger.Error().
					Str("file", file).
					Str("destination", destination).
					Str("error", err.Error()).
					Msg("could not copy file to ModeratorSent")

			}

		} else {
			numErrors++
			g.logger.Error().
				Str("file", file).
				Str("destination", destination).
				Str("error", err.Error()).
				Msg("could not copy file to export it")

		}
	}
	if numErrors == 0 {
		g.logger.Info().
			Int("count", len(files)).
			Str("destination", destination).
			Msg(fmt.Sprintf("Exported %d files to %s", len(files), destination))
	}

}

func (g *Ingester) ExportForChecking(exam, checker string) {

	files, err := GetFileList(g.CheckerReady(exam, checker))
	if err != nil {

		g.logger.Error().
			Str("error", err.Error()).
			Msg("could not get file list for exporting")
	}

	g.logger.Info().Msg("Exporting")

	numErrors := 0

	destination := g.ExportChecking(exam, checker)

	for _, file := range files {

		if !g.IsPDF(file) {
			continue
		}
		destination = g.ExportChecking(exam, checker)
		err := g.CopyToDir(file, destination)
		if err == nil {

			destination = g.CheckerSent(exam, checker)
			err = g.MoveToDir(file, destination)
			if err != nil {
				g.logger.Error().
					Str("file", file).
					Str("destination", destination).
					Str("error", err.Error()).
					Msg("could not copy file to CheckerSent")

			}

		} else {
			numErrors++
			g.logger.Error().
				Str("file", file).
				Str("destination", destination).
				Str("error", err.Error()).
				Msg("could not copy file to export it")

		}
	}
	if numErrors == 0 {
		g.logger.Info().
			Int("count", len(files)).
			Str("destination", destination).
			Msg(fmt.Sprintf("Exported %d files to %s", len(files), destination))
	}

}

func (g *Ingester) ExportForReMarking(exam, marker string) {

	files, err := GetFileList(g.ReMarkerReady(exam, marker))
	if err != nil {

		log.Error().
			Str("error", err.Error()).
			Msg("could not get file list for exporting")
	}
	g.logger.Info().Msg("Exporting")
	numErrors := 0

	destination := g.ExportReMarking(exam, marker)

	for _, file := range files {

		if !g.IsPDF(file) {
			continue
		}
		destination = g.ExportReMarking(exam, marker)
		err := g.CopyToDir(file, destination)
		if err == nil {
			destination = g.ReMarkerSent(exam, marker)

			err = g.MoveToDir(file, destination)
			if err != nil {
				g.logger.Error().
					Str("file", file).
					Str("destination", destination).
					Str("error", err.Error()).
					Msg("could not copy file to ReMarkerSent")

			}

		} else {
			numErrors++
			g.logger.Error().
				Str("file", file).
				Str("destination", destination).
				Str("error", err.Error()).
				Msg("could not copy file to export it")

		}
	}
	if numErrors == 0 {
		g.logger.Info().
			Int("count", len(files)).
			Str("destination", destination).
			Msg(fmt.Sprintf("Exported %d files to %s", len(files), destination))
	}

}

func (g *Ingester) ExportForReChecking(exam, checker string) {

	files, err := GetFileList(g.ReCheckerReady(exam, checker))
	if err != nil {

		log.Error().
			Str("error", err.Error()).
			Msg("could not get file list for exporting")
	}
	g.logger.Info().Msg("Exporting")
	numErrors := 0

	destination := g.ExportReChecking(exam, checker)

	for _, file := range files {

		if !g.IsPDF(file) {
			continue
		}
		destination = g.ExportReChecking(exam, checker)

		err := g.CopyToDir(file, destination)
		if err == nil {

			destination = g.ReCheckerSent(exam, checker)
			err = g.MoveToDir(file, destination)
			if err != nil {
				g.logger.Error().
					Str("file", file).
					Str("destination", destination).
					Str("error", err.Error()).
					Msg("could not copy file to CheckerSent")

			}

		} else {
			numErrors++
			g.logger.Error().
				Str("file", file).
				Str("destination", destination).
				Str("error", err.Error()).
				Msg("could not copy file to export it")

		}
	}

	if numErrors == 0 {
		g.logger.Info().
			Int("count", len(files)).
			Str("destination", destination).
			Msg(fmt.Sprintf("Exported %d files to %s", len(files), destination))
	}

}

/*


/*
func(cmd *cobra.Command, args []string) {
		which := os.Args[2]
		who := os.Args[3]
		exam := os.Args[4]

		var s Specification
		// load configuration from environment variables GRADEX_CLI_<var>
		if err := envconfig.Process("gradex_cli", &s); err != nil {
			fmt.Println("Configuration Failed")
			os.Exit(1)
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
					if s.Verbose {
						fmt.Printf("MC:%s\n", msg.Message)
					}
				}

			}
		}()

		logFile := filepath.Join(s.Root, "var/log/gradex-cli.log")
		ingester.EnsureDirAll(filepath.Dir(logFile))
		f, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer f.Close()

		logger := zerolog.
			New(f).
			With().
			Timestamp().
			Str("command", "export").
			Str("which", which).
			Str("who", who).
			Str("exam", exam).
			Logger()

		g, err := ingester.New(s.Root, mch, &logger)
		if err != nil {
			fmt.Printf("Failed getting New Ingester %v", err)
			os.Exit(1)
		}

		g.EnsureDirectoryStructure()
		g.SetupExamPaths(exam)

		switch which {
		case ingester.QuestionReady:
			files, err := g.GetFileList(g.QuestionReady(exam, who))
			if err != nil {
				logger.Error().
					Str("Error", err.Error()).
					Msg("Can't get files to export")
				return
			}
			for _, file := range files {

				err = g.CopyToDir(file, g.ExportLabelling(exam, who))

				if err != nil {
					logger.Error().
						Str("file", file).
						Str("Error", err.Error()).
						Msg("Can't copy file to export them")
					return
				}
				err = g.MoveToDir(file, g.QuestionSent(exam, who))
				if err != nil {
					logger.Error().
						Str("file", file).
						Str("Error", err.Error()).
						Msg("Can't move file to sent, after exported")
					return
				}
			}
		case ingester.MarkerReady:
			files, err := g.GetFileList(g.MarkerReady(exam, who))
			if err != nil {
				logger.Error().
					Str("Error", err.Error()).
					Msg("Can't get files to export")
				return
			}
			for _, file := range files {

				err = g.CopyToDir(file, g.ExportMarking(exam, who))

				if err != nil {
					logger.Error().
						Str("file", file).
						Str("Error", err.Error()).
						Msg("Can't copy file to export them")
					return
				}
				err = g.MoveToDir(file, g.MarkerSent(exam, who))
				if err != nil {
					logger.Error().
						Str("file", file).
						Str("Error", err.Error()).
						Msg("Can't move file to sent, after exported")
					return
				}
			}
		case ingester.ModeratorReady:
			files, err := g.GetFileList(g.ModeratorReady(exam, who))
			if err != nil {
				logger.Error().
					Str("Error", err.Error()).
					Msg("Can't get files to export")
				return
			}
			for _, file := range files {

				err = g.CopyToDir(file, g.ExportModerating(exam, who))

				if err != nil {
					logger.Error().
						Str("file", file).
						Str("Error", err.Error()).
						Msg("Can't copy file to export them")
					return
				}
				err = g.MoveToDir(file, g.ModeratorSent(exam, who))
				if err != nil {
					logger.Error().
						Str("file", file).
						Str("Error", err.Error()).
						Msg("Can't move file to sent, after exported")
					return
				}
			}
		case ingester.CheckerReady:
			files, err := g.GetFileList(g.CheckerReady(exam, who))
			if err != nil {
				logger.Error().
					Str("Error", err.Error()).
					Msg("Can't get files to export")
				return
			}
			for _, file := range files {

				err = g.CopyToDir(file, g.ExportChecking(exam, who))

				if err != nil {
					logger.Error().
						Str("file", file).
						Str("Error", err.Error()).
						Msg("Can't copy file to export them")
					return
				}
				err = g.MoveToDir(file, g.CheckerSent(exam, who))
				if err != nil {
					logger.Error().
						Str("file", file).
						Str("Error", err.Error()).
						Msg("Can't move file to sent, after exported")
					return
				}
			}

		case ingester.ReMarkerReady:
			files, err := g.GetFileList(g.ReMarkerReady(exam, who))
			if err != nil {
				logger.Error().
					Str("Error", err.Error()).
					Msg("Can't get files to export")
				return
			}
			for _, file := range files {

				err = g.CopyToDir(file, g.ExportReMarking(exam, who))

				if err != nil {
					logger.Error().
						Str("file", file).
						Str("Error", err.Error()).
						Msg("Can't copy file to export them")
					return
				}
				err = g.MoveToDir(file, g.ReMarkerSent(exam, who))
				if err != nil {
					logger.Error().
						Str("file", file).
						Str("Error", err.Error()).
						Msg("Can't move file to sent, after exported")
					return
				}
			}
		case ingester.ReCheckerReady:
			files, err := g.GetFileList(g.ReCheckerReady(exam, who))
			if err != nil {
				logger.Error().
					Str("Error", err.Error()).
					Msg("Can't get files to export")
				return
			}
			for _, file := range files {

				err = g.CopyToDir(file, g.ExportReChecking(exam, who))

				if err != nil {
					logger.Error().
						Str("file", file).
						Str("Error", err.Error()).
						Msg("Can't copy file to export them")
					return
				}
				err = g.MoveToDir(file, g.ReCheckerSent(exam, who))
				if err != nil {
					logger.Error().
						Str("file", file).
						Str("Error", err.Error()).
						Msg("Can't move file to sent, after exported")
					return
				}
			}

		} // switch
	}
}

*/
