package cmd

import (
	"errors"

	venv "github.com/azr4e1/venv-notary"
	"github.com/azr4e1/venv-notary/graphics"
	"github.com/spf13/cobra"
)

var (
	cloneCmd = &cobra.Command{
		Use:   "clone",
		Short: "clone a local or global environment",
		Args:  cobra.NoArgs,
		RunE:  cloneCobraFunction,
	}
)

func cloneCobraFunction(cmd *cobra.Command, args []string) error {
	return nil
}
