/*
Copyright © 2020 Tim Drysdale <timothy.d.drysdale@gmail.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/timdrysdale/chmsg"
	"github.com/timdrysdale/gradex-cli/ingester"
)

// exportCmd represents the export command
var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Put files into the export directory",
	Args:  cobra.ExactArgs(3),
	Long: `A helper command to take files from key points in the process, and put them in the export 
directory where they are easier to find. Example usage

gradex-cli export readyToMark tdd exam

Exported files are usually flagged in some way, e.g. being moved to a "sent" folder internally.`,
	Run: func(cmd *cobra.Command, args []string) {
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

		/*
			QuestionReady  = "questionReady"
			QuestionSent   = "questionSent"
			MarkerReady    = "markerReady"
			MarkerSent     = "markerSent"
			ModeratorReady = "moderatorReady"
			ModeratorSent  = "moderatorSent"
			CheckerReady   = "checkerReady"
			CheckerSent    = "checkerSent"
			RemarkerReady  = "remarkerReady"
			RemarkerSent   = "remarkerSent"
			RecheckerReady = "recheckerReady"
			RecheckerSent  = "recheckerSent"

		*/

	},
}

func init() {
	rootCmd.AddCommand(exportCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// exportCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// exportCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
