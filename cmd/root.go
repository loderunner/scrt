package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "scrt",
	Short: "A secret manager for the command-line",
}

func Init() {
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(versionCmd)
}

func Execute() error {
	return rootCmd.Execute()
}
