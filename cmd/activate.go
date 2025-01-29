package cmd

import (
	"errors"

	venv "github.com/azr4e1/venv-notary"
	"github.com/azr4e1/venv-notary/graphics"
	"github.com/spf13/cobra"
)

var (
	activateCmd = &cobra.Command{
		Use:               "activate",
		Short:             "Activate a local or global environment in a new shell",
		Args:              cobra.MaximumNArgs(1),
		RunE:              activateCobraFunction,
		ValidArgsFunction: venvCompletion,
	}
)

func activateCobraFunction(cmd *cobra.Command, args []string) error {
	if len(args) > 0 && localVenv {
		return errors.New("you cannot activate a global venv and a local venv at the same time.")
	}
	notary, err := venv.NewNotary()
	if err != nil {
		return err
	}
	if localVenv {
		err = activateLocal(notary, cmd, args)
		if err != nil {
			return err
		}
	} else if len(args) > 0 {
		err = activateGlobal(notary, cmd, args)
		if err != nil {
			return err
		}
	} else {
		return errors.New("you need to either activate a local or global venv.")
	}
	return nil
}

func activateGlobal(notary venv.Notary, cmd *cobra.Command, args []string) error {
	name := args[0]
	err := notary.ActivateGlobal(name, pythonVersion)
	if err != nil {
		if errors.As(err, &venv.VenvNotRegisteredError{}) {
			// err = notary.CreateGlobal(name, pythonVersion)
			err = graphics.StatusMain("No environment registered with this name and this Python version. Creating it now...", "Environment successfully created.", createAction, nil)(cmd, args)
			if err != nil {
				return nil
			}
			err = notary.GetVenvs()
			if err != nil {
				return err
			}
			err = notary.ActivateGlobal(name, pythonVersion)
			return err
		}
		return err
	}
	return nil
}

func activateLocal(notary venv.Notary, cmd *cobra.Command, args []string) error {
	err := notary.ActivateLocal(pythonVersion)
	if err != nil {
		if errors.As(err, &venv.VenvNotRegisteredError{}) {
			// err = notary.CreateLocal(pythonVersion)
			err = graphics.StatusMain("No environment registered at this location and with this Python version. Creating it now...", "Environment successfully created.", createAction, nil)(cmd, args)
			if err != nil {
				return nil
			}
			err = notary.GetVenvs()
			if err != nil {
				return err
			}
			err = notary.ActivateLocal(pythonVersion)
			return err
		}
		return err
	}
	return nil
}

func init() {
	activateCmd.Flags().BoolVarP(&localVenv, "local", "l", false, "activate local venv.")
	activateCmd.Flags().StringVarP(&pythonVersion, "python", "p", "", "use this python version.")
}
