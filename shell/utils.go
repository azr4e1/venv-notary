package shell

import (
	"errors"
	"os"
	"os/exec"

	ps "github.com/mitchellh/go-ps"
)

var shellVariables = map[shellType]string{
	bash:       "BASH_VERSION",
	zsh:        "ZSH_NAME",
	fish:       "FISH_VERSION",
	powershell: "PSEdition",
}

func hasShell(shellName Shell) bool {
	var command *exec.Cmd
	switch shellName.name {
	case bash, fish, zsh:
		command = exec.Command(shellName.executable, "--version")
	case powershell:
		command = exec.Command(shellName.executable, "-H")
	}
	_, err := command.CombinedOutput()
	if err != nil {
		return false
	}
	return true
}

func getShellName(availableShells []Shell) (Shell, error) {
	ppid := os.Getppid()
	proc, err := ps.FindProcess(ppid)
	if err != nil {
		return Shell{}, err
	}

	executable := proc.Executable()
	for _, sh := range availableShells {
		if executable == sh.executable {
			return sh, nil
		}
	}
	return Shell{}, errors.New("couldn't detect current shell among the supported ones")
}
