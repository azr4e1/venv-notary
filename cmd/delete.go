package cmd

import (
	"errors"

	"github.com/azr4e1/venv-notary/venv"
	"github.com/spf13/cobra"
)

var (
	deleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "delete a local or global virtual environment.",
		Long: `'delete' will delete a named virtual environment if choosing to use the --global flag,
or a virtual environment corresponding to the current directory if choosing to use the --local flag.`,
		RunE: deleteCobraFunc,
		Args: cobra.MaximumNArgs(1),
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
		err = notary.DeleteLocal()
		if err != nil {
			return err
		}
	} else if len(args) > 0 {
		err = notary.DeleteGlobal(args[0])
		if err != nil {
			return err
		}
	} else {
		return errors.New("you need to either delete a local or global venv.")
	}
	return nil
}

func init() {
	// createCmd.Flags().StringVarP(&globalVenvName, "global", "g", "", "name of the global venv.")
	deleteCmd.Flags().BoolVarP(&localVenv, "local", "l", false, "delete a local venv.")
	// createCmd.MarkFlagsMutuallyExclusive("global", "local")
}
