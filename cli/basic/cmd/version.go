package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of basic",
	Long:  `All software has versions`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("version v0.1 -- HEAD")
	},
}
