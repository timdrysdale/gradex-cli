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

// moderateCmd represents the moderate command
var enterCmd = &cobra.Command{
	Use:   "enter [enterer] [exam]",
	Short: "Add moderation bars to an exam",
	Args:  cobra.ExactArgs(2),
	Long: `Add enter bar to flattened, marked scripts, decorating the path with the enterer name, for example

gradex-cli enter abc demo-exam

this will produce a bunch of files in the enterReady directory, which can be exported.

`,
	Run: func(cmd *cobra.Command, args []string) {

		enterer := os.Args[2]
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
		logger := zerolog.New(f).With().Timestamp().Str("command", "enter").Logger()
		g, err := ingester.New(s.Root, mch, &logger)
		if err != nil {
			fmt.Printf("Failed getting New Ingester %v", err)
			os.Exit(1)
		}

		// if we've added new steps to the process, this prepares
		// the directories - is idempotent so doesn't matter
		// if we call it on structure already setup
		// these functions MUST not delete anything!
		g.EnsureDirectoryStructure()
		g.SetupExamDirs(exam)

		// TODO handling redo flag to redo the split is starting to get into
		// unclear territory - do you remove all files from moderation sets?
		// what if they have been sent?
		// do you block if any moderation files have been exported?
		// better to handle these cases manually
		err = g.SplitForEnter(exam)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		g.SetBackgroundIsVanilla(OpticalVanilla)
		g.SetOpticalShrink(OpticalShrink)

		err = g.AddEnterActiveBar(exam, enterer)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		err = g.AddEnterInactiveBar(exam)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(enterCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// moderateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// moderateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
