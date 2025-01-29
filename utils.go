package venv

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

const (
	HASHLEN = 64
)

// allowed characters: a-z, 0-9, _, -
func NormalizeName(name string) string {
	name = strings.TrimSpace(strings.ToLower(name))
	name = strings.Join(strings.Fields(name), "_")
	normalizedName := ""
	for _, r := range name {
		switch {
		case r >= 'a' && r <= 'z':
			normalizedName += string(r)
		case r >= '0' && r <= '9':
			normalizedName += string(r)
		case r == '_' || r == '-':
			normalizedName += string(r)
		}
	}
	return normalizedName
}

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
	headDir := NormalizeName(filepath.Base(currDir))
	if headDir == "" {
		return "", errors.New("Invalid venv name. Please use a name that contains only letters, digits, '_' and '-'.")
	}
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
	separator := "-" + VersionPrefix
	parts := strings.Split(name, separator)
	length := len(parts)
	if length == 1 {
		return name, ""
	}
	version := VersionPrefix + parts[length-1]
	if name[len(name)-len(version):] != version {
		return name, ""
	}
	return strings.Join(parts[:length-1], separator), version
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

func AlphanumericSort(a, b string) int {
	byteA := []byte(a)
	byteB := []byte(b)
	var i int
	var e byte
	for i, e = range byteA {
		if i >= len(byteB) || e > byteB[i] {
			return 1
		}
		if e < byteB[i] {
			return -1
		}
	}
	if len(byteA) == len(byteB) {
		return 0
	}
	return -1
}

func SemanticVersioningSort(a, b string) int {
	partsA := strings.Split(a, ".")
	partsB := strings.Split(b, ".")
	for i, e := range partsA {
		if i >= len(partsB) {
			return -1
		}
		valA, err := strconv.Atoi(e)
		if err != nil {
			return -1
		}
		valB, err := strconv.Atoi(partsB[i])
		if err != nil {
			return 1
		}
		if valA > valB {
			return 1
		}
		if valA < valB {
			return -1
		}
	}
	if len(partsA) == len(partsB) {
		return 0
	}
	return 1
}

func getDataHome() (string, error) {
	var dataHome string
	switch runtime.GOOS {
	case "linux":
		dataHome = os.Getenv("XDG_DATA_HOME")
		if dataHome == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				return "", err
			}
			dataHome = filepath.Join(home, ".local/share")
		}
	case "darwin":
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		dataHome = filepath.Join(home, "Library")
	case "windows":
		dataHome = os.Getenv("LOCALAPPDATA")
		if dataHome == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				return "", err
			}
			dataHome = filepath.Join(home, "AppData/Local")
		}
	}
	return dataHome, nil
}

func getVenvExecDir() string {
	var execDir string
	switch runtime.GOOS {
	case "linux", "darwin":
		execDir = "bin"
	case "windows":
		execDir = "Scripts"
	}

	return execDir
}
