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
	traceWho    string
	refreshAnon bool
	showOK      bool
)

// traceCmd represents the trace command
var traceCmd = &cobra.Command{
	Use:   "trace [stage] [exam] [--verbose=true]",
	Short: "Trace and report issues with files at specified step",
	Args:  cobra.ExactArgs(2),
	Long:  `The verbose flag shows OK files too, if true`,
	Run: func(cmd *cobra.Command, args []string) {
		stage := strings.ToLower(os.Args[2])
		exam := os.Args[3]

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
			Str("command", "trace").
			Str("stage", stage).
			Str("exam", exam).
			Logger()

		g, err := ingester.New(s.Root, mch, &logger)
		if err != nil {
			fmt.Printf("Failed getting New Ingester %v", err)
			os.Exit(1)
		}

		// setup dirs for later when writing report
		g.EnsureDirectoryStructure()
		g.SetupExamDirs(exam)

		dir, err := g.MergeProcessedPapersToDir(exam, stage)

		fmt.Println(dir)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		tokens, err := g.ReportOnProcessedDir(dir, showOK)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		for _, token := range tokens {
			fmt.Println(token)
		}

	},
}

func init() {
	rootCmd.AddCommand(traceCmd)
	traceCmd.Flags().StringVarP(&traceWho, "who", "w", "", "Name of actor to which to confine the check [default is to read all files at that stage]")
	traceCmd.Flags().BoolVarP(&refreshAnon, "refresh-anon", "r", false, "Reread the pagedata from the first stage (normally only needed after new un-seen submission added) [default false]")
	traceCmd.Flags().BoolVarP(&showOK, "verbose", "v", false, "Show the OK files as well [default false]")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// traceCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// traceCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
