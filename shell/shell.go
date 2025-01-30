package shell

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

type shellType int

const (
	bash shellType = iota
	zsh
	fish
	powershell
)

var shellExecutables = map[shellType][]string{
	bash:       []string{"bash", "bash.exe"},
	zsh:        []string{"zsh", "zsh.exe"},
	fish:       []string{"fish", "fish.exe"},
	powershell: []string{"pwsh", "powershell", "powershell.exe"},
}

type Shell struct {
	os         string
	name       shellType
	executable string
}

func (s Shell) Source(script string) error {
	var command *exec.Cmd
	switch s.name {
	case bash, zsh:
		command = exec.Command(s.executable, "-c", fmt.Sprintf("source '%s'; %s", script, s.executable))
	case fish:
		command = exec.Command(s.executable, "--interactive", "-C", fmt.Sprintf("source '%s'", script))
	case powershell:
		command = exec.Command(s.executable, "-NoExit", "-Command", fmt.Sprintf(". '%s'", script))
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
	switch s.name {
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
	currOs := runtime.GOOS

	var availableShells = []Shell{}
	for name, executables := range shellExecutables {
		for _, executable := range executables {
			sh := Shell{
				os:         currOs,
				name:       name,
				executable: executable,
			}
			if hasShell(sh) {
				availableShells = append(availableShells, sh)
			}
		}
	}

	currentShell, err := getShellName(availableShells)

	if err != nil && len(availableShells) > 0 {
		currentShell = availableShells[0]
	}

	return currentShell
}

func (s Shell) GetActivationScript() string {
	switch s.name {
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
