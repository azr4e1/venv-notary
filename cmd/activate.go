package cmd

import (
	"errors"

	"github.com/azr4e1/venv-notary/venv"
	"github.com/spf13/cobra"
)

var (
	activateCmd = &cobra.Command{
		Use:   "activate",
		Short: "activate a local or global environment in a new shell",
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
		err = notary.ActivateLocal()
		if err != nil {
			return err
		}
	} else if len(args) > 0 {
		err = notary.ActivateGlobal(args[0])
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
