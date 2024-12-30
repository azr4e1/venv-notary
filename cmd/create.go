package cmd

import (
	"errors"

	"github.com/azr4e1/venv-notary/venv"
	"github.com/spf13/cobra"
)

var (
	createCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a local or global virtual environment.",
		Long: `'create' will create a named virtual environment if choosing to use the --global flag,
or a virtual environment corresponding to the current directory if choosing to use the --local flag.
In both instances, the environment is created at $XDG_DATA_HOME/venv-dir.`,
		RunE: createCobraFunc,
		Args: cobra.MaximumNArgs(1),
	}
)

func createCobraFunc(cmd *cobra.Command, args []string) error {
	if len(args) > 0 && localVenv {
		return errors.New("you cannot create a global venv and a local venv at the same time.")
	}
	notary, err := venv.NewNotary()
	if err != nil {
		return err
	}
	if localVenv {
		err = notary.CreateLocal()
		if err != nil {
			return err
		}
	} else if len(args) > 0 {
		err = notary.CreateGlobal(args[0])
		if err != nil {
			return err
		}
	} else {
		return errors.New("you need to either create a local or global venv.")
	}
	return nil
}

func init() {
	// createCmd.Flags().StringVarP(&globalVenvName, "global", "g", "", "name of the global venv.")
	createCmd.Flags().BoolVarP(&localVenv, "local", "l", false, "create a local venv.")
	// createCmd.MarkFlagsMutuallyExclusive("global", "local")
}