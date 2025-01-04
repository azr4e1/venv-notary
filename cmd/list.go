package cmd

import (
	"fmt"
	"path/filepath"
	"slices"
	"strings"

	venv "github.com/azr4e1/venv-notary"
	"github.com/spf13/cobra"
)

var (
	listCmd = &cobra.Command{
		Use:   "list",
		Short: "List registered environments",
		Args:  cobra.NoArgs,
		RunE:  listCobraFunc,
	}
)

func listCobraFunc(cmd *cobra.Command, args []string) error {
	finalStr := ""
	if !localVenv {
		globalStr, err := printVenvs(false, pythonVersion)
		if err != nil {
			return err
		}
		finalStr += fmt.Sprintln(globalStr)
	}
	if !globalVenv {
		localStr, err := printVenvs(true, pythonVersion)
		if err != nil {
			return err
		}
		finalStr += fmt.Sprintln(localStr)
	}
	fmt.Println(strings.TrimSpace(finalStr))
	return nil
}

func printVenvs(isLocal bool, python string) (string, error) {
	notary, err := venv.NewNotary()
	if err != nil {
		return "", err
	}
	var str string
	var placeholder string
	var sortedVenvs []string
	var version string
	if python != "" {
		version, err = venv.PythonVersion(python)
		if err != nil {
			return "", err
		}
	}
	if isLocal {
		sortedVenvs = notary.ListLocal()
		str = "Local Environments\n"
	} else {
		sortedVenvs = notary.ListGlobal()
		str = "Global Environments\n"
	}
	slices.SortFunc(sortedVenvs, venv.AlphanumericSort)
	for _, venvPath := range sortedVenvs {
		v := venv.Venv{Path: venvPath}
		_, venvVersion := venv.ExtractVersion(filepath.Base(venvPath))
		if version != "" && version != venvVersion {
			continue
		}
		if v.IsActive() {
			placeholder = "*"
		} else {
			placeholder = " "
		}
		str += fmt.Sprintf("  %s %s\n", placeholder, filepath.Base(v.Path))
	}
	return str, nil
}

func init() {
	listCmd.Flags().BoolVarP(&globalVenv, "global", "g", false, "list only global venvs.")
	listCmd.Flags().BoolVarP(&localVenv, "local", "l", false, "list only local venvs.")
	listCmd.Flags().StringVarP(&pythonVersion, "python", "p", "", "filter by python version.")
	listCmd.MarkFlagsMutuallyExclusive("global", "local")
}
