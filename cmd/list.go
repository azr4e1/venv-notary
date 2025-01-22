package cmd

import (
	"os"

	ui "github.com/azr4e1/venv-notary/graphics"
	"github.com/spf13/cobra"
)

var (
	listCmd = &cobra.Command{
		Use:   "list",
		Short: "List registered environments",
		Args:  cobra.NoArgs,
		RunE:  ui.ListMain(&localVenv, &globalVenv, &pythonVersion, os.Stdout),
	}
)

func init() {
	listCmd.Flags().BoolVarP(&globalVenv, "global", "g", false, "list only global venvs.")
	listCmd.Flags().BoolVarP(&localVenv, "local", "l", false, "list only local venvs.")
	listCmd.Flags().StringVarP(&pythonVersion, "python", "p", "", "filter by python version.")
	listCmd.MarkFlagsMutuallyExclusive("local", "global")
}
