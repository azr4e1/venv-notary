package cmd

import (
	"os"
	"path/filepath"

	venv "github.com/azr4e1/venv-notary"
	"github.com/spf13/cobra"
)

var (
	globalVenvName string
	localVenv      bool
	globalVenv     bool
	pythonVersion  string
	namePattern    string
	rootCmd        = &cobra.Command{
		Use:     "vn",
		Short:   "A wrapper for python-venv",
		Long:    `venv-notary is an application that makes it easy to manage global and local virtual environments for Python.`,
		Version: "0.7.1",
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
	rootCmd.AddCommand(cleanCmd)
}

func initConfig() {
}

func venvCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	notary, err := venv.NewNotary()
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	globalVenvs := []string{}
	for _, v := range notary.ListGlobal() {
		name, _ := venv.ExtractVersion(filepath.Base(v))
		globalVenvs = append(globalVenvs, name)
	}
	return globalVenvs, cobra.ShellCompDirectiveNoFileComp
}
