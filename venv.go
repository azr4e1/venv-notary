package venv

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
)

type Location string

const (
	GlobalLoc Location = "global"
	LocalLoc  Location = "local"
)

const (
	DATAHOMEENV = "XDG_DATA_HOME"
	DATAHOMEDIR = ".local/share"
	NotaryDir   = "venv-notary"
)

const (
	HASHLEN = 64
)

type Venv string

type Notary struct {
	venvDir  string
	venvList map[Venv]Location
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

func (v Venv) CreateWithName(name string) error {
	_, err := os.Stat(string(v))

	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return err
		}
	} else {
		return errors.New("Directory or file already exists with this name.")
	}
	cmd := exec.Command("python", "-m", "venv", "--prompt", name, string(v))
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	return err
}

func (v Venv) Delete() error {
	if v.IsActive() {
		return errors.New("environment is active. Deactivate it before deleting it.")
	}
	if v.IsVenv() {
		err := os.RemoveAll(string(v))
		return err
	}
	return errors.New(fmt.Sprintf("'%s' is not a python environment!", string(v)))
}

func (v Venv) Activate() error {
	if !v.IsVenv() {
		return errors.New(fmt.Sprintf("'%s' is not a python environment!", string(v)))
	}
	activatePath := filepath.Join(string(v), "bin/activate")
	cmd := exec.Command("bash", "-c", "source "+activatePath+"; bash")
	// cmd := exec.Command("bash")
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return err
}

func (v Venv) IsActive() bool {
	venv := os.Getenv("VIRTUAL_ENV")
	if string(v) == venv {
		return true
	}
	return false
}

func NewNotary() (Notary, error) {
	dataHome := os.Getenv(DATAHOMEENV)
	if dataHome == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return Notary{}, err
		}
		dataHome = filepath.Join(home, DATAHOMEDIR)
	}
	notaryDir := filepath.Join(dataHome, NotaryDir)
	notary := Notary{
		venvDir: notaryDir,
	}
	err := notary.SetUp()
	if err != nil {
		return Notary{}, err
	}
	err = notary.GetVenvs()
	if err != nil {
		return Notary{}, err
	}
	return notary, nil
}

func (n Notary) SetUp() error {
	globalDir := n.GlobalDir()
	localDir := n.LocalDir()
	err := os.MkdirAll(globalDir, os.ModePerm)
	if err != nil {
		return err
	}
	err = os.MkdirAll(localDir, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func (n Notary) GlobalDir() string {
	return filepath.Join(n.venvDir, string(GlobalLoc))
}

func (n Notary) LocalDir() string {
	return filepath.Join(n.venvDir, string(LocalLoc))
}

func (n *Notary) GetVenvs() error {
	globDirs, err := os.ReadDir(n.GlobalDir())
	if err != nil {
		return err
	}
	localDirs, err := os.ReadDir(n.LocalDir())
	if err != nil {
		return err
	}
	venvList := map[Venv]Location{}

	// add global venvs
	for _, g := range globDirs {
		v := Venv(filepath.Join(n.GlobalDir(), g.Name()))
		if !v.IsVenv() {
			continue
		}
		venvList[v] = GlobalLoc
	}
	// add local venvs
	for _, g := range localDirs {
		v := Venv(filepath.Join(n.LocalDir(), g.Name()))
		if !v.IsVenv() {
			continue
		}
		venvList[v] = LocalLoc
	}
	n.venvList = venvList
	return nil
}

func createLocalName() (string, error) {
	currDir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	headDir := filepath.Base(currDir)
	h := sha256.New()
	h.Write([]byte(currDir))
	venvName := fmt.Sprintf("%s-%x", headDir, h.Sum(nil))

	return venvName, nil
}

func removeHash(name string) string {
	hashLength := HASHLEN + 1
	if len(name) > hashLength {
		name = name[:len(name)-(HASHLEN+1)]
	}
	return name
}

func (n *Notary) CreateLocal() error {
	venvName, err := createLocalName()
	if err != nil {
		return err
	}
	venv := Venv(filepath.Join(n.LocalDir(), venvName))
	_, ok := n.venvList[venv]
	if ok {
		return errors.New("Environment already exists at this location.")
	}
	err = venv.CreateWithName(removeHash(venvName))
	if err != nil {
		return err
	}
	n.venvList[venv] = LocalLoc
	return nil
}

func (n *Notary) CreateGlobal(name string) error {
	venv := Venv(filepath.Join(n.GlobalDir(), name))
	_, ok := n.venvList[venv]
	if ok {
		return errors.New("Environment already exists with this name.")
	}
	err := venv.Create()
	if err != nil {
		return err
	}
	n.venvList[venv] = GlobalLoc
	return nil
}

func (n *Notary) delete(venv Venv) error {
	err := venv.Delete()
	if err != nil {
		return err
	}
	delete(n.venvList, venv)
	return nil
}

func (n *Notary) DeleteLocal() error {
	venvName, err := createLocalName()
	if err != nil {
		return err
	}
	venv := Venv(filepath.Join(n.LocalDir(), venvName))
	_, ok := n.venvList[venv]
	if ok {
		return n.delete(venv)
	} else {
		return errors.New("No environment is registered for current directory.")
	}
}

func (n *Notary) DeleteGlobal(name string) error {
	venv := Venv(filepath.Join(n.GlobalDir(), name))
	_, ok := n.venvList[venv]
	if ok {
		return n.delete(venv)
	} else {
		return fmt.Errorf("No environment with name '%s' is registered.", name)
	}
}

func (n Notary) ListGlobal() []Venv {
	venvs := []Venv{}
	for venv, loc := range n.venvList {
		if loc == GlobalLoc {
			venvs = append(venvs, venv)
		}
	}
	return venvs
}

func (n Notary) ListLocal() []Venv {
	venvs := []Venv{}
	for venv, loc := range n.venvList {
		if loc == LocalLoc {
			venvs = append(venvs, venv)
		}
	}
	return venvs
}

func (n Notary) GetGlobalVenv(name string) Venv {
	venv := Venv(filepath.Join(n.GlobalDir(), name))

	return venv
}

func (n Notary) GetLocalVenv() (Venv, error) {
	venvName, err := createLocalName()
	if err != nil {
		return Venv(""), err
	}
	venv := Venv(filepath.Join(n.LocalDir(), venvName))

	return venv, nil
}

func (n Notary) IsRegistered(venv Venv) bool {
	_, ok := n.venvList[venv]

	return ok
}

func (n Notary) ActivateGlobal(name string) error {
	venv := n.GetGlobalVenv(name)
	if !n.IsRegistered(venv) {
		return fmt.Errorf("No environment with name '%s' is registered.", name)
	}
	err := venv.Activate()
	return err
}

func (n Notary) ActivateLocal() error {
	venv, err := n.GetLocalVenv()
	if err != nil {
		return err
	}
	if !n.IsRegistered(venv) {
		return errors.New("No environment is registered for current directory.")
	}

	err = venv.Activate()
	if err != nil {
		return err
	}
	return nil
}

func (n Notary) GetActiveEnv() (Venv, error) {
	venv := Venv(os.Getenv("VIRTUAL_ENV"))
	_, ok := n.venvList[venv]
	if ok {
		return venv, nil
	}
	return "", errors.New("No active registered virtual environments.")
}
