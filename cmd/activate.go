package cmd

import (
	"errors"

	"github.com/azr4e1/venv-notary/venv"
	"github.com/spf13/cobra"
)

var (
	activateCmd = &cobra.Command{
		Use:   "activate",
		Short: "Activate a local or global environment in a new shell",
		Args:  cobra.MaximumNArgs(1),
		RunE:  activateCobraFunction,
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
		venv, err := notary.GetLocalVenv()
		if err != nil {
			return err
		}
		if !notary.IsRegistered(venv) {
			err = notary.CreateLocal()
			if err != nil {
				return err
			}
		}
		err = notary.ActivateLocal()
		if err != nil {
			return err
		}
	} else if len(args) > 0 {
		name := args[0]
		venv := notary.GetGlobalVenv(name)
		if !notary.IsRegistered(venv) {
			err = notary.CreateGlobal(name)
			if err != nil {
				return err
			}
		}
		err = notary.ActivateGlobal(name)
		if err != nil {
			return err
		}
	} else {
		return errors.New("you need to either activate a local or global venv.")
	}
	return nil
}

func init() {
	activateCmd.Flags().BoolVarP(&localVenv, "local", "l", false, "activate local venv.")
}
