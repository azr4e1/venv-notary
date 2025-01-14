package cmd

import (
	"errors"

	venv "github.com/azr4e1/venv-notary"
	"github.com/spf13/cobra"
)

var (
	deleteCmd = &cobra.Command{
		Use:               "delete",
		Short:             "Delete a local or global virtual environment",
		RunE:              deleteCobraFunc,
		Args:              cobra.MaximumNArgs(1),
		ValidArgsFunction: venvCompletion,
	}
)

func deleteCobraFunc(cmd *cobra.Command, args []string) error {
	if len(args) > 0 && localVenv {
		return errors.New("you cannot delete a global venv and a local venv at the same time.")
	}
	notary, err := venv.NewNotary()
	if err != nil {
		return err
	}
	if localVenv {
		err = notary.DeleteLocal(pythonVersion)
		if err != nil {
			return err
		}
	} else if len(args) > 0 {
		err = notary.DeleteGlobal(args[0], pythonVersion)
		if err != nil {
			return err
		}
	} else {
		return errors.New("you need to either delete a local or global venv.")
	}
	return nil
}

func init() {
	deleteCmd.Flags().BoolVarP(&localVenv, "local", "l", false, "delete a local venv.")
	deleteCmd.Flags().StringVarP(&pythonVersion, "python", "p", "", "delete venv with this python version")
}
