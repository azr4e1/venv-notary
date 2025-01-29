package venv

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/azr4e1/venv-notary/shell"
)

type Venv struct {
	Path   string
	Name   string
	Python string
}

func (v Venv) String() string {
	return v.Path
}

func (v Venv) IsVenv() bool {
	dir := v.Path
	// check dir is created with correct files
	stat, err := os.Stat(dir)
	if err != nil {
		return false
	}
	if !stat.IsDir() {
		return false
	}
	activatePath := filepath.Join(dir, "bin/activate")
	stat, err = os.Stat(activatePath)
	if err != nil {
		return false
	}
	if !stat.Mode().IsRegular() {
		return false
	}
	pythonPath := filepath.Join(dir, "bin/python")
	stat, err = os.Stat(pythonPath)
	if err != nil {
		return false
	}
	if !stat.Mode().IsRegular() {
		return false
	}
	return true
}

func (v Venv) Create() error {
	_, err := os.Stat(v.Path)

	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return err
		}
	} else {
		return errors.New("Directory or file already exists with this name.")
	}
	executable := v.Python
	if executable == "" {
		executable = "python"
	}
	cmdEls := []string{executable, "-m", "venv"}
	if v.Name != "" {
		cmdEls = append(cmdEls, "--prompt", v.Name)
	}
	cmdEls = append(cmdEls, v.Path)
	cmd := exec.Command(cmdEls[0], cmdEls[1:]...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		os.RemoveAll(v.Path)
		return fmt.Errorf("%v. Error message: '%s'", strings.TrimSpace(err.Error()), strings.TrimSpace(string(output)))
	}
	return nil
}

func (v Venv) Delete() error {
	if v.IsActive() {
		return errors.New("environment is active. Deactivate it before deleting it.")
	}
	if v.IsVenv() {
		err := os.RemoveAll(v.Path)
		return err
	}
	return errors.New(fmt.Sprintf("'%s' is not a python environment!", v.Path))
}

func (v Venv) Activate() error {
	if !v.IsVenv() {
		return errors.New(fmt.Sprintf("'%s' is not a python environment!", v.Path))
	}
	if v.IsActive() {
		return errors.New("environment is already active!")
	}
	activeShell := shell.NewShell()
	activateScript := activeShell.GetActivationScript()
	if activateScript == "" {
		return errors.New("cannot locate activation script")
	}
	activatePath := filepath.Join(v.Path, "bin", activateScript)
	err := activeShell.Source(activatePath)
	return err
}

func (v Venv) IsActive() bool {
	venv := os.Getenv("VIRTUAL_ENV")
	if v.Path == venv {
		return true
	}
	return false
}

func (v Venv) GetPythonVersion() (string, error) {
	executable := v.Python
	if executable == "" {
		executable = "python"
	}
	return PythonVersion(executable)
}
