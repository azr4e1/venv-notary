package venv

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path"
)

const (
	DATAHOMEENV = "XDG_DATA_HOME"
	DATAHOMEDIR = ".local/share"
	NotaryDir   = "venv-notary"
	VenvList    = "venv-list.txt"
	VenvDir     = "global-venv"
)

type Venv string

type Notary struct {
	venvList   string
	globalVenv string
}

func (v Venv) IsVenv() bool {
	dir := string(v)
	// check dir is created with correct files
	stat, err := os.Stat(dir)
	if err != nil {
		return false
	}
	if !stat.IsDir() {
		return false
	}
	activate_path := filepath.Join(dir, "bin/activate")
	stat, err = os.Stat(activate_path)
	if err != nil {
		return false
	}
	if !stat.Mode().IsRegular() {
		return false
	}
	python_path := filepath.Join(dir, "bin/python")
	stat, err = os.Stat(python_path)
	if err != nil {
		return false
	}
	if !stat.Mode().IsRegular() {
		return false
	}
	return true
}

func (v Venv) Create() error {
	_, err := os.Stat(string(v))

	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return err
		}
	} else {
		return errors.New("Directory or file already exists with this name.")
	}
	cmd := exec.Command("python", "-m", "venv", string(v))
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	return err
}

func (v Venv) Delete() error {
	if v.IsVenv() {
		err := os.RemoveAll(string(v))
		return err
	}
	return errors.New(fmt.Sprintf("'%s' is not a python environment!", string(v)))
}

func NewNotary() Notary {
	dataHome := os.Getenv(DATAHOMEENV)
	if dataHome == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			panic(err)
		}
		dataHome = path.Join(home, DATAHOMEDIR)
	}
	notaryDir := path.Join(dataHome, NotaryDir)
	return Notary{
		venvList:   path.Join(notaryDir, VenvList),
		globalVenv: path.Join(notaryDir, VenvDir),
	}
}
