package cmd

import (
	"errors"
	"os"

	venv "github.com/azr4e1/venv-notary"
	"github.com/spf13/cobra"
)

var (
	runCmd = &cobra.Command{
		Use:   "run",
		Short: "Run a command in the virtual environment",
		RunE:  runCobraFunction,
	}
)

func runCobraFunction(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		cmd.Help()
		os.Exit(0)
	}
	comm := args[0]
	commArgs := args[1:]
	notary, err := venv.NewNotary()
	if err != nil {
		return err
	}
	if localVenv {
		err = notary.RunLocal(pythonVersion, comm, commArgs...)
		if err != nil {
			return err
		}
	} else {
		if globalVenvName == "" {
			return errors.New("you need to either activate a local or global venv.")
		}
		err = notary.RunGlobal(globalVenvName, pythonVersion, comm, commArgs...)
		if err != nil {
			return err
		}
	}
	return nil
}

func init() {
	runCmd.Flags().BoolVarP(&localVenv, "local", "l", false, "activate local venv.")
	runCmd.Flags().StringVarP(&globalVenvName, "global", "g", "", "activate global venv.")
	runCmd.Flags().StringVarP(&pythonVersion, "python", "p", "", "use this python version.")
	runCmd.MarkFlagsMutuallyExclusive("global", "local")
	runCmd.RegisterFlagCompletionFunc("global", venvCompletion)
}
