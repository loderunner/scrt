package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print scrt version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("scrt v0.0.0")
	},
}
