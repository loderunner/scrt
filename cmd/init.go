package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes a new repository",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Initializing repository")
	},
}
