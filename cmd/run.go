package cmd

import (
	venv "github.com/azr4e1/venv-notary"
	"github.com/spf13/cobra"
)

var (
	runCmd = &cobra.Command{
		Use:   "run",
		Short: "Run a command in the virtual environment (default local)",
		RunE:  runCobraFunction,
	}
)

func runCobraFunction(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return cmd.Help()
	}
	comm := args[0]
	commArgs := args[1:]
	notary, err := venv.NewNotary()
	if err != nil {
		return err
	}
	if globalVenvName != "" {
		err = notary.RunGlobal(globalVenvName, pythonVersion, comm, commArgs...)
		if err != nil {
			return err
		}
	} else {
		err = notary.RunLocal(pythonVersion, comm, commArgs...)
		if err != nil {
			return err
		}
	}
	return nil
}

func init() {
	runCmd.Flags().StringVarP(&globalVenvName, "global", "g", "", "run in global venv")
	runCmd.Flags().StringVarP(&pythonVersion, "python", "p", "", "use this python version")
	runCmd.RegisterFlagCompletionFunc("global", venvCompletion)
}
