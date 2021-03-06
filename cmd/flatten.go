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
	"strings"

	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/timdrysdale/chmsg"
	"github.com/timdrysdale/gradex-cli/ingester"
)

var (
	ensureAncestor bool
)

// flattenCmd represents the flatten command
var flattenCmd = &cobra.Command{
	Use:   "flatten [exam] [stage]",
	Short: "Validates and anonymises PDF files submitted via Learn",
	Args:  cobra.ExactArgs(2),
	Long: `You must specify the exam and the stage for which you wish to flatten files. You should have already ingested the files.
Example usage:

gradex-cli flatten SomeExam new

Possible stages to flatten

new
marked
remarked
remoderated
checked
rechecked`,
	Run: func(cmd *cobra.Command, args []string) {

		exam := os.Args[3]

		stage := os.Args[2]

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
		logger := zerolog.New(f).With().Timestamp().Logger()
		g, err := ingester.New(s.Root, mch, &logger)
		if err != nil {
			fmt.Printf("Failed getting New Ingester %v", err)
			os.Exit(1)
		}

		g.EnsureDirectoryStructure()
		g.SetupExamDirs(exam)

		if Template != "" {
			err := g.SetOverlayTemplatePath(Template)
			if err != nil {
				fmt.Printf("Overlay not usable because %s\n", err.Error())
				os.Exit(1)
			}
		}

		g.Redo = redo

		switch {

		case strings.ToLower(stage) == "new":

			err = g.FlattenNewPapers(exam)

		case ingester.ValidStageForProcessedPapers(stage):

			g.SetBackgroundIsVanilla(OpticalVanilla)
			g.SetOpticalShrink(OpticalShrink)
			g.SetChangeAncestor(ensureAncestor)

			err = g.FlattenProcessedPapers(exam, stage)

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			if stage != ingester.Checked {
				err = g.MergeProcessedPapers(exam, stage)
			}

		default:

			fmt.Printf("Stage [%s] not known\n", stage)

		}

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		os.Exit(0)

	},
}

func init() {
	rootCmd.AddCommand(flattenCmd)
	flattenCmd.Flags().BoolVar(&ensureAncestor, "ensure-ancestor", false, "Force rewriting of pagedata to ensure link to ancestor in anonPapers [default false]")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// flattenCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// flattenCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
