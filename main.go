package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path"
)

type Venv struct {
	path string
}

func (v Venv) IsVenv() bool {
	// check dir is created with correct files
	stat, err := os.Stat(v.path)
	if err != nil {
		return false
	}
	if !stat.IsDir() {
		return false
	}
	activate_path := path.Join(v.path, "bin/activate")
	stat, err = os.Stat(activate_path)
	if err != nil {
		return false
	}
	if !stat.Mode().IsRegular() {
		return false
	}
	python_path := path.Join(v.path, "bin/python")
	stat, err = os.Stat(python_path)
	if err != nil {
		return false
	}
	if !stat.Mode().IsRegular() {
		return false
	}
	return true
}

func (v Venv) VenvCreate() error {
	_, err := os.Stat(v.path)

	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return err
		}
	} else {
		return errors.New("Directory of file already exists with this name.")
	}
	cmd := exec.Command("python", "-m", "venv", v.path)
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	return err
}

func (v Venv) VenvDelete() error {
	if v.IsVenv() {
		err := os.RemoveAll(v.path)
		return err
	}
	return errors.New(fmt.Sprintf("'%s' is not a python environment!", v.path))
}

func main() {
	Venv{"venv"}.VenvCreate()
}
