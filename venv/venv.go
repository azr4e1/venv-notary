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

type Venv struct {
	Path string
}

type Notary struct {
	venvList   string
	globalVenv string
}

func (v Venv) IsVenv() bool {
	// check dir is created with correct files
	stat, err := os.Stat(v.Path)
	if err != nil {
		return false
	}
	if !stat.IsDir() {
		return false
	}
	activate_path := path.Join(v.Path, "bin/activate")
	stat, err = os.Stat(activate_path)
	if err != nil {
		return false
	}
	if !stat.Mode().IsRegular() {
		return false
	}
	python_path := path.Join(v.Path, "bin/python")
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
	_, err := os.Stat(v.Path)

	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return err
		}
	} else {
		return errors.New("Directory of file already exists with this name.")
	}
	cmd := exec.Command("python", "-m", "venv", v.Path)
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	return err
}

func (v Venv) VenvDelete() error {
	if v.IsVenv() {
		err := os.RemoveAll(v.Path)
		return err
	}
	return errors.New(fmt.Sprintf("'%s' is not a python environment!", v.Path))
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
