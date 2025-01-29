package shell

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

type shellType string

const (
	bash       shellType = "bash"
	zsh        shellType = "zsh"
	fish       shellType = "fish"
	powershell shellType = "powershell"
)

var supportedShells = []shellType{bash, zsh, fish, powershell}

// var os = runtime.GOOS

type Shell struct {
	os    string
	shell shellType
}

func (s Shell) Source(script string) error {
	var command *exec.Cmd
	switch sh := s.shell; sh {
	case bash, zsh:
		command = exec.Command(string(sh), "-c", fmt.Sprintf("source '%s'; %s", script, sh))
	case fish:
		command = exec.Command(string(sh), "--interactive", "-C", fmt.Sprintf("source '%s'", script))
	case powershell:
		command = exec.Command(string(sh), "-NoExit", "-Command", fmt.Sprintf(". '%s'", script))
	default:
		return errors.New("No shell available.")
	}
	command.Stderr = os.Stderr
	command.Stdout = os.Stdout
	command.Stdin = os.Stdin

	err := command.Run()
	return err
}

func (s Shell) OS() string {
	return s.os
}

func (s Shell) Name() string {
	switch s.shell {
	case bash:
		return "Bash"
	case zsh:
		return "Zsh"
	case fish:
		return "Fish"
	case powershell:
		return "Powershell"
	default:
		return "Unknown"
	}
}

func NewShell() Shell {
	newShell := Shell{
		os: runtime.GOOS,
	}

	var currentShell shellType
	var availableShells = []shellType{}
	for _, sh := range supportedShells {
		if hasShell(sh) {
			availableShells = append(availableShells, sh)
		}
	}

	currentShell, _ = getShellName(availableShells)

	if currentShell == "" && len(availableShells) > 0 {
		currentShell = availableShells[0]
	}

	newShell.shell = currentShell

	return newShell
}

func (s Shell) GetActivationScript() string {
	switch s.shell {
	case bash, zsh:
		return "activate"
	case fish:
		return "activate.fish"
	case powershell:
		return "Activate.ps1"
	default:
		return ""
	}
}
