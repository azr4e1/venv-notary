package main

import (
	"os"
	"path"
	"testing"
)

func checkVirtualEnv(dir string, t *testing.T) {
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

func TestCreatesAVirtualEnv(t *testing.T) {
	t.Parallel()

	// setup
	dir, err := os.MkdirTemp("", "*")
	if err != nil {
		t.Fatal(err)
	}
	os.Remove(dir)
	err = Venv{dir}.VenvCreate()
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
	err = Venv{dir}.VenvCreate()
	if err == nil {
		t.Error("Creates virtual env even when directory already exists")
	}

	file, err := os.CreateTemp("", "*")
	if err != nil {
		t.Fatal(err)
	}
	err = Venv{file.Name()}.VenvCreate()
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
	v := Venv{dir}
	err = v.VenvCreate()
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
	v := Venv{dir}
	err = v.VenvCreate()
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
	v = Venv{file.Name()}
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
	v := Venv{dir}
	err = v.VenvCreate()
	if err != nil {
		t.Fatal(err)
	}
	err = v.VenvDelete()
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
	v := Venv{dir}
	err = v.VenvDelete()
	if err == nil {
		t.Error("should return error")
	}
	file, err := os.CreateTemp("", "*")
	if err != nil {
		t.Fatal(err)
	}
	v = Venv{file.Name()}
	err = v.VenvDelete()
	if err == nil {
		t.Error("should return error")
	}
}
