package cmd

import (
	"errors"

	venv "github.com/azr4e1/venv-notary"
	"github.com/azr4e1/venv-notary/graphics"
	"github.com/spf13/cobra"
)

var (
	activateCmd = &cobra.Command{
		Use:   "activate",
		Short: "Activate a local or global environment in a new shell (default local)",
		Args:  cobra.NoArgs,
		RunE:  activateCobraFunction,
	}
)

func activateCobraFunction(cmd *cobra.Command, args []string) error {
	notary, err := venv.NewNotary()
	if err != nil {
		return err
	}
	if globalVenvName != "" {
		err = activateGlobal(notary, cmd, args)
		if err != nil {
			return err
		}
	} else {
		err = activateLocal(notary, cmd, args)
		if err != nil {
			return err
		}
	}
	return nil
}

func activateGlobal(notary venv.Notary, cmd *cobra.Command, args []string) error {
	err := notary.ActivateGlobal(globalVenvName, pythonVersion)
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
			err = notary.ActivateGlobal(globalVenvName, pythonVersion)
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
	activateCmd.Flags().StringVarP(&globalVenvName, "global", "g", "", "activate global venv.")
	activateCmd.Flags().StringVarP(&pythonVersion, "python", "p", "", "use this python version.")
	activateCmd.RegisterFlagCompletionFunc("global", venvCompletion)
}
