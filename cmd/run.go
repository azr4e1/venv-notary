package cmd

import (
	"errors"

	venv "github.com/azr4e1/venv-notary"
	"github.com/azr4e1/venv-notary/graphics"
	"github.com/spf13/cobra"
)

var (
	runCmd = &cobra.Command{
		Use:   "run",
		Short: "Run a command inside the given venv",
		RunE:  activateCobraFunction,
	}
)
