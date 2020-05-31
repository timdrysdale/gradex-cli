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
)

// traceCmd represents the trace command
var traceCmd = &cobra.Command{
	Use:   "trace [stage] [step] [exam] --who=XXX",
	Short: "Verify that UUID link in pagedata of all files",
	Args:  cobra.ExactArgs(3),
	Long:  `Select [stage] and [step] to compare to 05-anonymous papers; with the option to compare in a subdirectory`,
	Run: func(cmd *cobra.Command, args []string) {
		stage := strings.ToLower(os.Args[2])
		step := strings.ToLower(os.Args[3])
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
			Str("command", "trace").
			Str("stage", stage).
			Str("step", step).
			Str("exam", exam).
			Logger()

		g, err := ingester.New(s.Root, mch, &logger)
		if err != nil {
			fmt.Printf("Failed getting New Ingester %v", err)
			os.Exit(1)
		}

		g.EnsureDirectoryStructure()
		g.SetupExamDirs(exam)

		targetDir := ""

		switch step {

		case "back":
			targetDir, err := g.FlattenProcessedPapersFromDir(exam, stage)
			if err != nil {
				fmt.Printf("Failed to find step %s in stage %s\n", step, stage)
				os.Exit(1)
			}

		case "flattened":
			targetDir, err := g.FlattenProcessedPapersToDir(exam, stage)
			if err != nil {
				fmt.Printf("Failed to find step %s in stage %s\n", step, stage)
				os.Exit(1)
			}
		case "processed":
			targetDir, err := g.MergeProcessedPapersToDir(exam, stage)
			if err != nil {
				fmt.Printf("Failed to find step %s in stage %s\n", step, stage)
				os.Exit(1)
			}
		}

		// treat who as subdir without constraining the name - for max flexibility in troubleshooting, e.g. inactive sets
		// or (!) misnamed directories (that will during dev only, obvs)
		targetDir = filepath.Join(targetDir, traceWho)

		/*

			case "badpage", "badpages", "pagebad":
				files, err := g.GetFileList(g.GetExamDir(exam, ingester.PageBad))
				if err != nil {
					return
				}
				for _, file := range files {
					fmt.Println(file)
				}

			case "tree":

				lines, err := tree.Tree(g.GetExamRoot(exam), false)

				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				fmt.Println(lines)

			case "pagetree":

				lines, err := tree.Tree(g.GetExamRoot(exam), true)

				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				fmt.Println(lines)

			case "sortcheck":

				g.SortCheck(exam)

			default:
				fmt.Printf("Unknown list type: %s\n", what)
			} // switch
		*/
	},
}

func init() {
	rootCmd.AddCommand(traceCmd)
	traceCmd.Flags().StringVarP(&traceWho, "who", "w", "", "Name of actor to which to confine the check [default is to read all files at that stage]")
	rootCmd.Flags().BoolVarP(&refreshAnon, "refresh-anon", "r", false, "Reread the pagedata from the first stage (normally only needed after new un-seen submission added) [default false]")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// traceCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// traceCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
