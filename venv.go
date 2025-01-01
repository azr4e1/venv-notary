package venv

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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

type Venv struct {
	Path   string
	Name   string
	Python string
}

func (v Venv) String() string {
	return v.Path
}

type Notary struct {
	venvDir  string
	venvList map[string]Location
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
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	return err
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
	activatePath := filepath.Join(v.Path, "bin/activate")
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
	if v.Path == venv {
		return true
	}
	return false
}

func getMinorVersion(version string) string {
	parts := strings.Split(version, ".")

	if len(parts) != 3 {
		return version
	}

	return strings.Join(parts[:2], ".")
}

func (v Venv) GetPythonVersion() (string, error) {
	executable := v.Python
	if executable == "" {
		executable = "python"
	}
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
	venvList := map[string]Location{}

	// add global venvs
	for _, g := range globDirs {
		// fullName := g.Name()
		v := Venv{Path: filepath.Join(n.GlobalDir(), g.Name())}
		if !v.IsVenv() {
			continue
		}
		venvList[v.Path] = GlobalLoc
	}
	// add local venvs
	for _, g := range localDirs {
		v := Venv{Path: filepath.Join(n.LocalDir(), g.Name())}
		if !v.IsVenv() {
			continue
		}
		venvList[v.Path] = LocalLoc
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
	version := ""
	length := len(parts)
	if length == 1 {
		return name, version
	}
	return strings.Join(parts[:length-1], separator), parts[length-1]
}

func (n *Notary) CreateLocal(python string) error {
	venvName, err := createLocalName()
	if err != nil {
		return err
	}
	venv := Venv{Path: filepath.Join(n.LocalDir(), venvName), Name: RemoveHash(venvName), Python: python}
	venv, err = addVersion(venv)
	if err != nil {
		return err
	}
	_, ok := n.venvList[venv.Path]
	if ok {
		return errors.New("Environment already exists at this location.")
	}
	err = venv.Create()
	if err != nil {
		return err
	}
	n.venvList[venv.Path] = LocalLoc
	return nil
}

func (n *Notary) CreateGlobal(name, python string) error {
	venv := Venv{Path: filepath.Join(n.GlobalDir(), name), Python: python, Name: name}
	venv, err := addVersion(venv)
	if err != nil {
		return err
	}
	_, ok := n.venvList[venv.Path]
	if ok {
		return errors.New("Environment already exists with this name.")
	}
	err = venv.Create()
	if err != nil {
		return err
	}
	n.venvList[venv.Path] = GlobalLoc
	return nil
}

func (n *Notary) delete(venv Venv) error {
	err := venv.Delete()
	if err != nil {
		return err
	}
	delete(n.venvList, venv.Path)
	return nil
}

func (n *Notary) DeleteLocal(python string) error {
	venvName, err := createLocalName()
	if err != nil {
		return err
	}
	venv := Venv{Path: filepath.Join(n.LocalDir(), venvName), Name: venvName, Python: python}
	venv, err = addVersion(venv)
	if err != nil {
		return err
	}
	_, ok := n.venvList[venv.Path]
	if ok {
		return n.delete(venv)
	} else {
		return errors.New("No environment is registered for current directory.")
	}
}

func (n *Notary) DeleteGlobal(name, python string) error {
	venv := Venv{Path: filepath.Join(n.GlobalDir(), name), Name: name, Python: python}
	venv, err := addVersion(venv)
	if err != nil {
		return err
	}
	_, ok := n.venvList[venv.Path]
	if ok {
		return n.delete(venv)
	} else {
		return fmt.Errorf("No environment with name '%s' is registered.", name)
	}
}

func (n Notary) ListGlobal() []string {
	venvs := []string{}
	for venv, loc := range n.venvList {
		if loc == GlobalLoc {
			venvs = append(venvs, venv)
		}
	}
	return venvs
}

func (n Notary) ListLocal() []string {
	venvs := []string{}
	for venv, loc := range n.venvList {
		if loc == LocalLoc {
			venvs = append(venvs, venv)
		}
	}
	return venvs
}

func (n Notary) GetGlobalVenv(name, python string) (Venv, error) {
	venv := Venv{Path: filepath.Join(n.GlobalDir(), name), Name: name, Python: python}
	venv, err := addVersion(venv)
	if err != nil {
		return Venv{}, err
	}

	return venv, nil
}

func (n Notary) GetLocalVenv(python string) (Venv, error) {
	venvName, err := createLocalName()
	if err != nil {
		return Venv{}, err
	}
	venv := Venv{Path: filepath.Join(n.LocalDir(), venvName), Name: venvName, Python: python}
	venv, err = addVersion(venv)
	if err != nil {
		return Venv{}, err
	}

	return venv, nil
}

func (n Notary) IsRegistered(venv Venv) bool {
	_, ok := n.venvList[venv.Path]

	return ok
}

func (n Notary) ActivateGlobal(name, python string) error {
	venv, err := n.GetGlobalVenv(name, python)
	if err != nil {
		return err
	}
	if !n.IsRegistered(venv) {
		return fmt.Errorf("No environment with name '%s' is registered.", name)
	}
	err = venv.Activate()
	return err
}

func (n Notary) ActivateLocal(python string) error {
	venv, err := n.GetLocalVenv(python)
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
	venv := Venv{Path: os.Getenv("VIRTUAL_ENV")}
	_, ok := n.venvList[venv.Path]
	if ok {
		return venv, nil
	}
	return Venv{}, errors.New("No active registered virtual environments.")
}
