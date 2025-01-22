package cmd

import (
	"errors"

	venv "github.com/azr4e1/venv-notary"
	"github.com/azr4e1/venv-notary/graphics"
	"github.com/spf13/cobra"
)

var (
	createCmd = &cobra.Command{
		Use:               "create",
		Short:             "Create a local or global virtual environment",
		RunE:              graphics.StatusMain("Creating environment...", "Environment successfully created.", createAction),
		Args:              cobra.MaximumNArgs(1),
		ValidArgsFunction: venvCompletion,
	}
)

// func createCobraFunc(cmd *cobra.Command, args []string) error {

func createAction(cmd *cobra.Command, args []string) func() error {
	return func() error {
		if len(args) > 0 && localVenv {
			return errors.New("you cannot create a global venv and a local venv at the same time.")
		}
		notary, err := venv.NewNotary()
		if err != nil {
			return err
		}
		if localVenv {
			err = notary.CreateLocal(pythonVersion)
			if err != nil {
				return err
			}
		} else if len(args) > 0 {
			err = notary.CreateGlobal(args[0], pythonVersion)
			if err != nil {
				return err
			}
		} else {
			return errors.New("you need to either create a local or global venv.")
		}
		return nil
	}
}

func init() {
	createCmd.Flags().BoolVarP(&localVenv, "local", "l", false, "create a local venv.")
	createCmd.Flags().StringVarP(&pythonVersion, "python", "p", "", "use this python version.")
}
