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

var UseFullAssignmentName bool

// ingestCmd represents the ingest command
var ingestCmd = &cobra.Command{
	Use:   "ingest",
	Args:  cobra.ExactArgs(0),
	Short: "Ingest files into exams for further processing",
	Long: `This command works on all files in the ingest directory for your exam processing system. 
You MUST set the root of this system as an environment variable 

eg. on linux 
export $GRADEX_CLI_ROOT=/usr/local/gradex

Then you can issue the ingest command as

gradex-cli ingest

If you are chopping and changing between test and production systems, you might wish to use a local "one time setting" of the environment
variable, e.g. on linux

GRADEX_CLI_ROOT=/some/test/gradex; gradex-cli ingest

`,
	Run: func(cmd *cobra.Command, args []string) {
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

		if UseFullAssignmentName {
			logger.Info().Msg("Using long names for assignments")
			g.SetUseFullAssignmentName()
		} else {
			logger.Info().Msg("Using short names for assignments")
		}

		g.EnsureDirectoryStructure()

		err = g.StageFromIngest()

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		g.ValidateNewPapers()
		if err != nil {
			fmt.Printf("Validate error %v", err)
			os.Exit(1)
		}

		os.Exit(0)

	},
}

func init() {
	rootCmd.AddCommand(ingestCmd)
	ingestCmd.Flags().BoolVarP(&UseFullAssignmentName, "use-long-name", "l", false, "Use long name for assignment [default false]")
}
