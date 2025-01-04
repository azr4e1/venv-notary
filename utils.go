package venv

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	HASHLEN = 64
)

func getMinorVersion(version string) string {
	parts := strings.Split(version, ".")

	if len(parts) != 3 {
		return version
	}

	return strings.Join(parts[:2], ".")
}

func PythonVersion(executable string) (string, error) {
	cmd := exec.Command(executable, "-V")
	version, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	versionEls := strings.Fields(string(version))
	if len(versionEls) != 2 {
		return "", fmt.Errorf("Something went wrong in fetching python version: '%s'", string(version))
	}
	if versionEls[0] != "Python" {
		return "", errors.New("Executable is not Python binary.")
	}
	versionNo := getMinorVersion(versionEls[1])
	return fmt.Sprintf("py%s", versionNo), nil
}

func createLocalName() (string, error) {
	currDir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	headDir := filepath.Base(currDir)
	h := sha256.New()
	_, err = h.Write([]byte(currDir))
	if err != nil {
		return "", err
	}
	venvName := fmt.Sprintf("%s-%x", headDir, h.Sum(nil))

	return venvName, nil
}

func addVersion(venv Venv) (Venv, error) {
	version, err := venv.GetPythonVersion()
	if err != nil {
		return venv, err
	}
	venv.Path = fmt.Sprintf("%s-%s", venv.Path, version)
	venv.Name = fmt.Sprintf("%s-%s", venv.Name, version)

	return venv, nil
}

func RemoveHash(name string) string {
	hashLength := HASHLEN + 1
	if len(name) > hashLength {
		name = name[:len(name)-(HASHLEN+1)]
	}
	return name
}

func ExtractVersion(name string) (string, string) {
	separator := "-py"
	parts := strings.Split(name, separator)
	length := len(parts)
	if length == 1 {
		return name, ""
	}
	version := "py" + parts[length-1]
	if name[len(name)-len(version):] != version {
		return name, ""
	}
	return strings.Join(parts[:length-1], separator), "py" + parts[length-1]
}

func SafeDir(f func() error) error {
	dir, err := os.MkdirTemp("", "*")
	if err != nil {
		return err
	}
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}
	err = os.Chdir(dir)
	if err != nil {
		return err
	}
	err = f()
	if err != nil {
		return err
	}
	err = os.Chdir(currentDir)
	return err
}
