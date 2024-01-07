package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "license-tool",
	Short: "A tool for doing things with licenses.",
	Long:  `Longer description to be put here.`,
}

var Path string
var Verbosity int

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.license-tool.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.PersistentFlags().StringVarP(&Path, "path", "p", ".", "Working path directory")
	rootCmd.PersistentFlags().CountVarP(&Verbosity, "verbosity", "v", "Log level")
}
