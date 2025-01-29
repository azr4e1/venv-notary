package venv

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
)

type Notary struct {
	venvDir  string
	venvList map[string]Location
}

type Location string

const (
	GlobalLoc Location = "global"
	LocalLoc  Location = "local"
)

const (
	NotaryDir     = "venv-notary"
	VersionPrefix = "py"
)

func NewNotary() (Notary, error) {
	dataHome, err := getDataHome()
	if err != nil {
		return Notary{}, err
	}
	notaryDir := filepath.Join(dataHome, NotaryDir)
	notary := Notary{
		venvDir: notaryDir,
	}
	err = notary.SetUp()
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

func (n *Notary) CreateLocal(python string) error {
	venvName, err := createLocalName()
	if err != nil {
		return err
	}
	// change to empty dir so there is no clash with python script for venv
	err = SafeDir(func() error {
		venv := Venv{Path: filepath.Join(n.LocalDir(), venvName), Name: RemoveHash(venvName), Python: python}
		venv, err = addVersion(venv)
		if err != nil {
			return err
		}
		_, ok := n.venvList[venv.Path]
		if ok {
			return errors.New("Environment already exists at this location and with this Python version.")
		}
		err = venv.Create()
		if err != nil {
			return err
		}
		n.venvList[venv.Path] = LocalLoc
		return nil
	})
	return err
}

func (n *Notary) CreateGlobal(name, python string) error {
	name = NormalizeName(name)
	if name == "" {
		return errors.New("Invalid venv name. Please use a name that contains only letters, digits, '_' and '-'.")
	}
	// change to empty dir so there is no clash with python script for venv
	var err error
	err = SafeDir(func() error {
		venv := Venv{Path: filepath.Join(n.GlobalDir(), name), Python: python, Name: name}
		venv, err = addVersion(venv)
		if err != nil {
			return err
		}
		_, ok := n.venvList[venv.Path]
		if ok {
			return errors.New("Environment already exists with this name and this Python version.")
		}
		err = venv.Create()
		if err != nil {
			return err
		}
		n.venvList[venv.Path] = GlobalLoc
		return nil

	})
	return err
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
	venv, err := n.GetLocalVenv(python)
	if err != nil {
		return err
	}
	if venvs := n.GetRegisteredVersionsOfVenv(venv, true); python == "" && len(venvs) >= 1 {
		if len(venvs) == 1 {
			registeredVenv := Venv{Path: venvs[0]}
			return n.delete(registeredVenv)
		}
		return MultipleVersionsError{"Multiple Python versions associated with this environment. Select one Python version."}
	}
	venv, err = addVersion(venv)
	if err != nil {
		return err
	}
	_, ok := n.venvList[venv.Path]
	if ok {
		return n.delete(venv)
	} else {
		return errors.New("No environment is registered for current directory with this Python version.")
	}
}

func (n *Notary) DeleteGlobal(name, python string) error {
	venv, err := n.GetGlobalVenv(name, python)
	if err != nil {
		return err
	}
	if venvs := n.GetRegisteredVersionsOfVenv(venv, false); python == "" && len(venvs) >= 1 {
		if len(venvs) == 1 {
			registeredVenv := Venv{Path: venvs[0]}
			return n.delete(registeredVenv)
		}
		return MultipleVersionsError{"Multiple Python versions associated with this environment. Select one Python version."}
	}
	venv, err = addVersion(venv)
	if err != nil {
		return err
	}
	_, ok := n.venvList[venv.Path]
	if ok {
		return n.delete(venv)
	} else {
		return fmt.Errorf("No environment with name '%s' is registered with this Python version.", name)
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
	name = NormalizeName(name)
	if name == "" {
		return Venv{}, errors.New("Invalid venv name. Please use a name that contains only letters, digits, '_' and '-'.")
	}
	venv := Venv{Path: filepath.Join(n.GlobalDir(), name), Name: name, Python: python}

	return venv, nil
}

func (n Notary) GetLocalVenv(python string) (Venv, error) {
	venvName, err := createLocalName()
	if err != nil {
		return Venv{}, err
	}
	venv := Venv{Path: filepath.Join(n.LocalDir(), venvName), Name: venvName, Python: python}

	return venv, nil
}

func (n Notary) IsRegistered(venv Venv) bool {
	_, ok := n.venvList[venv.Path]

	return ok
}

func (n Notary) IsRegisteredNoVersion(venv Venv, isLocal bool) bool {
	match, _ := ExtractVersion(venv.Path)
	for v, l := range n.venvList {
		if isLocal {
			if l != LocalLoc {
				continue
			}
		} else {
			if l != GlobalLoc {
				continue
			}
		}
		name, _ := ExtractVersion(v)
		if name == match {
			return true
		}
	}
	return false
}

func (n Notary) GetRegisteredVersionsOfVenv(venv Venv, isLocal bool) []string {
	match, _ := ExtractVersion(venv.Path)
	registeredVenvs := []string{}
	for v, l := range n.venvList {
		if isLocal {
			if l != LocalLoc {
				continue
			}
		} else {
			if l != GlobalLoc {
				continue
			}
		}
		name, _ := ExtractVersion(v)
		if name == match {
			registeredVenvs = append(registeredVenvs, v)
		}
	}
	return registeredVenvs
}

func (n Notary) GetHighestVersionRegistered(venv Venv, isLocal bool) Venv {
	registeredVersions := map[string]string{}
	versionSlice := []string{}
	match, _ := ExtractVersion(venv.Path)
	for v, l := range n.venvList {
		if isLocal {
			if l != LocalLoc {
				continue
			}
		} else {
			if l != GlobalLoc {
				continue
			}
		}
		name, version := ExtractVersion(v)
		if name == match {
			// remove "py" in front of version
			version = version[len(VersionPrefix):]
			registeredVersions[version] = v
			versionSlice = append(versionSlice, version)
		}
	}
	if len(registeredVersions) == 0 {
		return Venv{}
	}
	slices.SortFunc(versionSlice, SemanticVersioningSort)
	highestVersion := versionSlice[len(versionSlice)-1]
	venvPath := registeredVersions[highestVersion]
	return Venv{Path: venvPath}
}

func (n Notary) ActivateGlobal(name, python string) error {
	venv, err := n.GetGlobalVenv(name, python)
	if err != nil {
		return err
	}
	if venvs := n.GetRegisteredVersionsOfVenv(venv, false); python == "" && len(venvs) >= 1 {
		if len(venvs) == 1 {
			registeredVenv := Venv{Path: venvs[0]}
			return registeredVenv.Activate()
		}
		return MultipleVersionsError{"Multiple Python versions associated with this environment. Select one Python version."}
	}
	venv, err = addVersion(venv)
	if err != nil {
		return err
	}
	if !n.IsRegistered(venv) {
		return VenvNotRegisteredError{Message: fmt.Sprintf("No environment with name '%s' is registered with this Python version.", name)}
	}
	err = venv.Activate()
	return err
}

func (n Notary) ActivateLocal(python string) error {
	venv, err := n.GetLocalVenv(python)
	if err != nil {
		return err
	}
	if venvs := n.GetRegisteredVersionsOfVenv(venv, true); python == "" && len(venvs) >= 1 {
		if len(venvs) == 1 {
			registeredVenv := Venv{Path: venvs[0]}
			return registeredVenv.Activate()
		}
		return MultipleVersionsError{"Multiple Python versions associated with this environment. Select one Python version."}
	}
	venv, err = addVersion(venv)
	if err != nil {
		return err
	}
	if !n.IsRegistered(venv) {
		return VenvNotRegisteredError{Message: "No environment is registered for current directory with this Python version."}
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
