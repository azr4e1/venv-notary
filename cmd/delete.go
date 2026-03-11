package cmd

import (
	"os"

	venv "github.com/azr4e1/venv-notary"
	"github.com/azr4e1/venv-notary/graphics"
	"github.com/spf13/cobra"
)

var (
	deleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete a local or global virtual environment (default local)",
		RunE:  graphics.StatusMain("Deleting environment...", "Environment successfully deleted.", deleteAction, setupAction),
		Args:  cobra.NoArgs,
	}
)

// override error system of status line; goal is to get cobra style error output if multiple versions are available. This unfortunately means code repetition, but oh well.
func setupAction(cmd *cobra.Command, args []string) error {
	n, err := venv.NewNotary()
	if err != nil {
		return err
	}
	if globalVenvName != "" {
		vn, err := n.GetGlobalVenv(globalVenvName, pythonVersion)
		if err != nil {
			return err
		}
		if venvs := n.GetRegisteredVersionsOfVenv(vn, false); pythonVersion == "" && len(venvs) > 1 {
			return venv.MultipleVersionsError{Message: "Multiple Python versions associated with this environment. Select one Python version."}
		}
	} else {
		currDir, err := os.Getwd()
		if err != nil {
			return err
		}
		vn, err := n.GetLocalVenv(currDir, pythonVersion)
		if err != nil {
			return err
		}
		if venvs := n.GetRegisteredVersionsOfVenv(vn, true); pythonVersion == "" && len(venvs) > 1 {
			return venv.MultipleVersionsError{Message: "Multiple Python versions associated with this environment. Select one Python version."}
		}
	}
	return nil
}

func deleteAction(cmd *cobra.Command, args []string) func() error {
	return func() error {
		notary, err := venv.NewNotary()
		if err != nil {
			return err
		}
		if globalVenvName != "" {
			err = notary.DeleteGlobal(globalVenvName, pythonVersion)
			if err != nil {
				return err
			}
		} else {
			err = notary.DeleteLocal(pythonVersion)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

func init() {
	deleteCmd.Flags().StringVarP(&globalVenvName, "global", "g", "", "delete a global venv.")
	deleteCmd.Flags().StringVarP(&pythonVersion, "python", "p", "", "delete venv with this python version")
	deleteCmd.RegisterFlagCompletionFunc("global", venvCompletion)
}
