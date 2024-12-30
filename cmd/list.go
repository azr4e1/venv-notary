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
		globalStr, err := printVenvs(false)
		if err != nil {
			return err
		}
		finalStr += fmt.Sprintln(globalStr)
	}
	if !globalVenv {
		localStr, err := printVenvs(true)
		if err != nil {
			return err
		}
		finalStr += fmt.Sprintln(localStr)
	}
	fmt.Println(strings.TrimSpace(finalStr))
	return nil
}

func init() {
	listCmd.Flags().BoolVarP(&globalVenv, "global", "g", false, "list only global venvs.")
	listCmd.Flags().BoolVarP(&localVenv, "local", "l", false, "list only local venvs.")
	listCmd.MarkFlagsMutuallyExclusive("global", "local")
}

func printVenvs(isLocal bool) (string, error) {
	notary, err := venv.NewNotary()
	if err != nil {
		return "", err
	}
	var str string
	var placeholder string
	var sortedVenvs []venv.Venv
	if isLocal {
		sortedVenvs = notary.ListLocal()
		str = "Local Environments\n"
	} else {
		sortedVenvs = notary.ListGlobal()
		str = "Global Environments\n"
	}
	slices.Sort(sortedVenvs)
	for _, v := range sortedVenvs {
		if v.IsActive() {
			placeholder = "*"
		} else {
			placeholder = " "
		}
		str += fmt.Sprintf(" %s  %s\n", placeholder, filepath.Base(string(v)))
	}
	return str, nil

}
