package cmd

import (
	"path/filepath"

	venv "github.com/azr4e1/venv-notary"
	"github.com/spf13/cobra"
)

var (
	cleanCmd = &cobra.Command{
		Use:   "clean",
		Short: "Delete all local or global environments. You can filter by Python version.",
		RunE:  cleanCobraFunc,
		Args:  cobra.NoArgs,
	}
)

func cleanCobraFunc(cmd *cobra.Command, args []string) error {
	notary, err := venv.NewNotary()
	if err != nil {
		return err
	}
	if localVenv {
		err = deleteVenv(notary.ListLocal(), pythonVersion)
		if err != nil {
			return err
		}
	}
	if globalVenv {
		err = deleteVenv(notary.ListGlobal(), pythonVersion)
		if err != nil {
			return err
		}
	}
	return nil
}

func deleteVenv(vPath []string, python string) error {
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
		_, venvVersion := venv.ExtractVersion(filepath.Base(venvPath))
		if version != "" && version != venvVersion {
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
	cleanCmd.MarkFlagsOneRequired("local", "global")
}
