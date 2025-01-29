package shell

import (
	"errors"
	ps "github.com/mitchellh/go-ps"
	"os"
	"os/exec"
)

var shellVariables = map[shellType]string{
	bash:       "BASH_VERSION",
	zsh:        "ZSH_NAME",
	fish:       "FISH_VERSION",
	powershell: "PSEdition",
}

func hasShell(shellName shellType) bool {
	command := exec.Command(string(shellName), "--version")
	_, err := command.CombinedOutput()
	if err != nil {
		return false
	}
	return true
}

func getShellName(availableShells []shellType) (shellType, error) {
	ppid := os.Getppid()
	proc, err := ps.FindProcess(ppid)
	if err != nil {
		return "", err
	}

	for _, sh := range availableShells {
		if sh == shellType(proc.Executable()) {
			return sh, nil
		}
	}
	return "", errors.New("couldn't detect current shell among the supported ones")
}
