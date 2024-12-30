package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	globalVenvName string
	localVenv      bool
	globalVenv     bool
	rootCmd        = &cobra.Command{
		Use:   "vn",
		Short: "A wrapper for python-venv",
		Long: `venv-notary is an application that makes it easy to manage global and local virtual environments for Python.
It is a wrapper around python-venv. It does not manage Python installations!`,
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(activateCmd)
	rootCmd.AddCommand(deleteCmd)
}

func initConfig() {
}