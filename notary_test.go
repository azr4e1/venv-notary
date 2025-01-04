package venv

import (
	"os"
	"path"
	"testing"
)

func TestCreatesAVirtualEnv(t *testing.T) {
	t.Parallel()

	// setup
	dir, err := os.MkdirTemp("", "*")
	if err != nil {
		t.Fatal(err)
	}
	os.Remove(dir)
	err = Venv{Path: dir}.Create()
	if err != nil {
		t.Fatal(err)
	}
	// check dir is created with correct files
	stat, err := os.Stat(dir)
	if err != nil {
		t.Fatal(err)
	}
	if !stat.IsDir() {
		t.Error("did not create a dir")
	}
	activate_path := path.Join(dir, "bin/activate")
	stat, err = os.Stat(activate_path)
	if err != nil {
		t.Fatal(err)
	}
	if !stat.Mode().IsRegular() {
		t.Error("did not create activate file")
	}
	python_path := path.Join(dir, "bin/python")
	stat, err = os.Stat(python_path)
	if err != nil {
		t.Fatal(err)
	}
	if !stat.Mode().IsRegular() {
		t.Error("did not create python file")
	}
}

func TestFailsAtCreatingAVirtualEnv(t *testing.T) {
	t.Parallel()

	// setup
	dir, err := os.MkdirTemp("", "*")
	if err != nil {
		t.Fatal(err)
	}
	err = Venv{Path: dir}.Create()
	if err == nil {
		t.Error("Creates virtual env even when directory already exists")
	}

	file, err := os.CreateTemp("", "*")
	if err != nil {
		t.Fatal(err)
	}
	err = Venv{Path: file.Name()}.Create()
	if err == nil {
		t.Error("Creates virtual env even when file already exists")
	}
}

func TestCheckIsVirtualEnv(t *testing.T) {
	t.Parallel()
	dir, err := os.MkdirTemp("", "*")
	if err != nil {
		t.Fatal(err)
	}
	err = os.Remove(dir)
	if err != nil {
		t.Fatal(err)
	}
	v := Venv{Path: dir}
	err = v.Create()
	if err != nil {
		t.Fatal(err)
	}
	if !v.IsVenv() {
		t.Error("directory is an environment!")
	}
}

func TestCheckIsNotVirtualEnv(t *testing.T) {
	t.Parallel()
	dir, err := os.MkdirTemp("", "*")
	if err != nil {
		t.Fatal(err)
	}
	err = os.Remove(dir)
	if err != nil {
		t.Fatal(err)
	}
	v := Venv{Path: dir}
	err = v.Create()
	if err != nil {
		t.Fatal(err)
	}
	err = os.Remove(path.Join(dir, "bin/activate"))
	if v.IsVenv() {
		t.Error("directory is not an environment!")
	}
	file, err := os.CreateTemp("", "*")
	if err != nil {
		t.Fatal(err)
	}
	v = Venv{Path: file.Name()}
	if v.IsVenv() {
		t.Error("file is not an environment!")
	}
}

func TestDeleteVenv_DeletesTheVenv(t *testing.T) {
	t.Parallel()
	dir, err := os.MkdirTemp("", "*")
	if err != nil {
		t.Fatal(err)
	}
	err = os.Remove(dir)
	if err != nil {
		t.Fatal(err)
	}
	v := Venv{Path: dir}
	err = v.Create()
	if err != nil {
		t.Fatal(err)
	}
	err = v.Delete()
	if err != nil {
		t.Fatal(err)
	}
	if v.IsVenv() {
		t.Error("Failed at deleting venv")
	}
}

func TestDeleteVenv_DoesNotDeleteVenv(t *testing.T) {
	t.Parallel()
	dir, err := os.MkdirTemp("", "*")
	if err != nil {
		t.Fatal(err)
	}
	v := Venv{Path: dir}
	err = v.Delete()
	if err == nil {
		t.Error("should return error")
	}
	file, err := os.CreateTemp("", "*")
	if err != nil {
		t.Fatal(err)
	}
	v = Venv{Path: file.Name()}
	err = v.Delete()
	if err == nil {
		t.Error("should return error")
	}
}

func TestNewNotary_GetsHomeDirCorrectly(t *testing.T) {
	t.Parallel()
	dir, err := os.MkdirTemp("", "*")
	if err != nil {
		t.Fatal(err)
	}
	os.Setenv("XDG_DATA_HOME", dir)
	notary, err := NewNotary()
	if err != nil {
		t.Fatal(err)
	}
	wantDir := path.Join(dir, "venv-notary")
	gotDir := notary.venvDir
	if wantDir != gotDir {
		t.Errorf("want global dir path '%s', got '%s'", wantDir, gotDir)
	}
	wantGlobalDir := path.Join(dir, "venv-notary/global")
	gotGlobalDir := notary.GlobalDir()
	if wantGlobalDir != gotGlobalDir {
		t.Errorf("want global dir path '%s', got '%s'", wantDir, gotDir)
	}
	wantLocalDir := path.Join(dir, "venv-notary/local")
	gotLocalDir := notary.LocalDir()
	if wantLocalDir != gotLocalDir {
		t.Errorf("want global dir path '%s', got '%s'", wantDir, gotDir)
	}
}

func TestNewNotary_GetsPathCorrectlyIfHomeDirNotSet(t *testing.T) {
	t.Parallel()
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatal(err)
	}
	os.Setenv("XDG_DATA_HOME", "")
	notary, err := NewNotary()
	if err != nil {
		t.Fatal(err)
	}
	wantDir := path.Join(home, ".local/share/venv-notary")
	gotDir := notary.venvDir
	if wantDir != gotDir {
		t.Errorf("want global dir path '%s', got '%s'", wantDir, gotDir)
	}
	wantGlobalDir := path.Join(home, ".local/share/venv-notary/global")
	gotGlobalDir := notary.GlobalDir()
	if wantGlobalDir != gotGlobalDir {
		t.Errorf("want global dir path '%s', got '%s'", wantDir, gotDir)
	}
	wantLocalDir := path.Join(home, ".local/share/venv-notary/local")
	gotLocalDir := notary.LocalDir()
	if wantLocalDir != gotLocalDir {
		t.Errorf("want global dir path '%s', got '%s'", wantDir, gotDir)
	}
}

func TestNewNotary_SetsUpDirCorrectly(t *testing.T) {
	t.Parallel()
	dir, err := os.MkdirTemp("", "*")
	if err != nil {
		t.Fatal(err)
	}
	os.Setenv("XDG_DATA_HOME", dir)
	notary, err := NewNotary()
	if err != nil {
		t.Fatal(err)
	}
	homeDir, err := os.Stat(notary.venvDir)
	if err != nil {
		t.Fatal(err)
	}
	globalDir, err := os.Stat(notary.GlobalDir())
	if err != nil {
		t.Fatal(err)
	}
	localDir, err := os.Stat(notary.LocalDir())
	if err != nil {
		t.Fatal(err)
	}
	if !homeDir.IsDir() {
		t.Error("main dir has not been created")
	}
	if !globalDir.IsDir() {
		t.Error("global env dir has not been created")
	}
	if !localDir.IsDir() {
		t.Error("local env dir has not been created")
	}
}

func TestExtractVersion(t *testing.T) {
	t.Parallel()
	type testCase struct {
		Question string
		Answer1  string
		Answer2  string
	}
	testCases := []testCase{
		{
			Question: "ciaocomeva-py3.1.43",
			Answer1:  "ciaocomeva",
			Answer2:  "py3.1.43",
		},
		{
			Question: "-py3.1.43",
			Answer1:  "",
			Answer2:  "py3.1.43",
		},
		{
			Question: "ciaocomeva",
			Answer1:  "ciaocomeva",
			Answer2:  "",
		},
		{
			Question: "-py4.2.5ciaocomeva",
			Answer1:  "",
			Answer2:  "py4.2.5ciaocomeva",
		},
	}
	for i, tc := range testCases {

		noVersion, version := tc.Answer1, tc.Answer2
		noVersionGot, versionGot := ExtractVersion(tc.Question)
		if noVersion != noVersionGot || version != versionGot {
			t.Errorf("for testcase %d: name: %s, name got: %s; version: %s, version got: %s", i, noVersion, noVersionGot, version, versionGot)
		}
	}
}
