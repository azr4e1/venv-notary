package cmd

import (
	"path/filepath"
	"strings"

	venv "github.com/azr4e1/venv-notary"
	"github.com/azr4e1/venv-notary/graphics"
	"github.com/spf13/cobra"
)

var (
	cleanCmd = &cobra.Command{
		Use:   "clean",
		Short: "Delete all local or global environments. You can filter by Python version.",
		RunE:  graphics.StatusMain("Cleaning up environments...", "All environments deleted.", cleanAction),
		Args:  cobra.NoArgs,
	}
)

func cleanAction(cmd *cobra.Command, args []string) func() error {
	return func() error {
		notary, err := venv.NewNotary()
		if err != nil {
			return err
		}
		if localVenv {
			err = deleteVenv(notary.ListLocal(), pythonVersion, namePattern)
			if err != nil {
				return err
			}
		}
		if globalVenv {
			err = deleteVenv(notary.ListGlobal(), pythonVersion, namePattern)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

func deleteVenv(vPath []string, python, namePattern string) error {
	var version string
	var err error
	if python != "" {
		version, err = venv.PythonVersion(python)
		if err != nil {
			return err
		}
	}
	for _, venvPath := range vPath {
		v := venv.Venv{Path: venvPath}
		name, venvVersion := venv.ExtractVersion(filepath.Base(venvPath))
		if version != "" && version != venvVersion {
			continue
		}
		if namePattern != "" && !strings.Contains(name, namePattern) {
			continue
		}
		err = v.Delete()
		if err != nil {
			return err
		}
	}
	return nil
}

func init() {
	cleanCmd.Flags().BoolVarP(&localVenv, "local", "l", false, "delete all local venvs.")
	cleanCmd.Flags().BoolVarP(&globalVenv, "global", "g", false, "delete all global venvs.")
	cleanCmd.Flags().StringVarP(&pythonVersion, "python", "p", "", "delete venvs with this python version")
	cleanCmd.Flags().StringVarP(&namePattern, "name", "n", "", "delete venvs with this name pattern")
	cleanCmd.MarkFlagsOneRequired("local", "global")
}
