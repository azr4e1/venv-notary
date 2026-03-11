package cmd

import (
	venv "github.com/azr4e1/venv-notary"
	"github.com/azr4e1/venv-notary/graphics"
	"github.com/spf13/cobra"
)

var (
	createCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a local or global virtual environment (default local)",
		RunE:  graphics.StatusMain("Creating environment...", "Environment successfully created.", createAction, nil),
		Args:  cobra.NoArgs,
	}
)

func createAction(cmd *cobra.Command, args []string) func() error {
	return func() error {
		notary, err := venv.NewNotary()
		if err != nil {
			return err
		}
		if globalVenvName != "" {
			err = notary.CreateGlobal(globalVenvName, pythonVersion)
			if err != nil {
				return err
			}
		} else {
			err = notary.CreateLocal(pythonVersion)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

func init() {
	createCmd.Flags().StringVarP(&globalVenvName, "global", "g", "", "create a global venv")
	createCmd.Flags().StringVarP(&pythonVersion, "python", "p", "", "use this python version")
	createCmd.RegisterFlagCompletionFunc("global", venvCompletion)
}
