package shell

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
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
	procName, err := getProcessName()
	if err != nil {
		return "", err
	}

	for _, sh := range availableShells {
		if sh == shellType(procName) {
			return sh, nil
		}
	}
	return "", errors.New("couldn't detect current shell among the supported ones")
}

func getProcessName() (string, error) {
	ppid := os.Getppid()
	var exePath string
	switch runtime.GOOS {
	case "linux":
		if path, err := os.Readlink(fmt.Sprintf("/proc/%d/exe", ppid)); err == nil {
			exePath = path
		}
	}

	if exePath == "" {
		return "", errors.New("couldn't detect parent process")
	}

	shellName := filepath.Base(exePath)

	return shellName, nil
}
