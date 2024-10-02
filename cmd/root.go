package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "rck",
	Short: "Rocket 🚀 An easier way to write queries.",
}

func Execute() {
	rootCmd.Execute()
}
