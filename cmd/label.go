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

// labelCmd represents the label command
var labelCmd = &cobra.Command{
	Use:   "label",
	Short: "A brief description of your command",
	Args:  ExactArgs(1),
	Long: `Add labelling bars to all flattened scripts, decorating the path with the labeller name, for example

gradex-cli label x demo-exam

this will produce a bunch of files in the questionReady folder, e.g

$GRADEX_CLI_ROOT/usr/demo-exam/08-question-ready/X/<original-filename>-laTDD.pdf

Note that the exam argument is the relative path to the exam in $GRADEX_CLI_ROOT/usr/exam/

`,
	Run: func(cmd *cobra.Command, args []string) {

		exam := os.Args[2]

		labeller := "x"

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
		g.SetupExamPaths(exam)

		err = g.AddLabelBar(exam, labeller)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(labelCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// labelCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// labelCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
