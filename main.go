package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path"
)

func IsVenv(dir string) bool {
	// check dir is created with correct files
	stat, err := os.Stat(dir)
	if err != nil {
		return false
	}
	if !stat.IsDir() {
		return false
	}
	activate_path := path.Join(dir, "bin/activate")
	stat, err = os.Stat(activate_path)
	if err != nil {
		return false
	}
	if !stat.Mode().IsRegular() {
		return false
	}
	python_path := path.Join(dir, "bin/python")
	stat, err = os.Stat(python_path)
	if err != nil {
		return false
	}
	if !stat.Mode().IsRegular() {
		return false
	}
	return true
}

func VenvCreate(path string) error {
	_, err := os.Stat(path)

	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return err
		}
	} else {
		return errors.New("Directory of file already exists with this name.")
	}
	cmd := exec.Command("python", "-m", "venv", path)
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	return err
}

func VenvDelete(path string) error {
	if IsVenv(path) {
		err := os.RemoveAll(path)
		return err
	}
	return errors.New(fmt.Sprintf("'%s' is not a python environment!", path))
}

func main() {
	VenvCreate("venv")
}
