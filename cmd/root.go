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

	"github.com/spf13/cobra"
)

var (
	OpticalShrink  int
	OpticalVanilla bool
	redo           bool
	noversion      bool
	Template       string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gradex-cli",
	Short: "Grade PDF assessments like paper",
	Long: `Gradex makes marking with PDF more like marking on paper, because it protects previous steps in the marking process from editing. 
Grids at the side of the page speed up marking for those on keyboard,
and allow stylus users to use the same forms.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().BoolVar(&redo, "redo", false, "Force redo of already processed files")
	rootCmd.PersistentFlags().BoolVarP(&OpticalVanilla, "background-vanilla", "b", true, "Assume vanilla background for optical checkboxes? [default true]")
	rootCmd.PersistentFlags().IntVarP(&OpticalShrink, "box-shrink", "s", 15, "Number of pixels to shrink optical boxes to avoid false positives from boundaries [default 15]")
	rootCmd.PersistentFlags().StringVar(&Template, "layout", "layout.svg", "Use this layout [default layout.svg]")
	rootCmd.PersistentFlags().BoolVar(&noversion, "noversion", false, "don't show version")
}

func initConfig() {

}
