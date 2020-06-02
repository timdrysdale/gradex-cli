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
	"github.com/timdrysdale/gradex-cli/pagedata"
	"github.com/timdrysdale/gradex-cli/tree"
	"github.com/timdrysdale/gradex-cli/util"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list [what] [exam/file]",
	Args:  cobra.ExactArgs(2),
	Short: "show a list of [what] for [exam/file]",
	Long: `Shows selected information about an exam

for example

badpages - pages marked with badpage
tree - tree diagram of exam folder with file/done counts
pagetree - as above but with page counts
sortcheck - checks the sort was ok
pagedata - read and prettyprint the pagedata from a file

For example:

gradex-cli list badpages 'PGEE00000 A B D Exam'

`,
	Run: func(cmd *cobra.Command, args []string) {
		what := os.Args[2]
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
			Str("command", "list").
			Str("what", what).
			Str("exam", exam).
			Logger()

		g, err := ingester.New(s.Root, mch, &logger)
		if err != nil {
			fmt.Printf("Failed getting New Ingester %v", err)
			os.Exit(1)
		}

		g.EnsureDirectoryStructure()

		switch what {

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

		case "pagedata":

			pageDataMap, err := pagedata.UnMarshalAllFromFile(exam)

			if err != nil {
				msg := fmt.Sprintf("Skipping (%s): error obtaining pagedata\n", exam)
				logger.Error().
					Str("file", exam).
					Str("error", err.Error()).
					Msg(msg)
				fmt.Println(msg)
				os.Exit(1)
			}

			if pagedata.GetLen(pageDataMap) < 1 {
				msg := fmt.Sprintf("Skipping (%s): no pagedata in file\n", exam)
				logger.Error().
					Str("file", exam).
					Msg(msg)
				fmt.Println(msg)
				os.Exit(1)
			}

			util.PrettyPrintStruct(pageDataMap)
		default:
			fmt.Printf("Unknown list type: %s\n", what)
		} // switch
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
