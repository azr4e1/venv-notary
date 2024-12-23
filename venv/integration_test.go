//go:build integration

package venv

import (
	"os"
	"os/exec"
	"testing"
)

func TestPythonAvailable(t *testing.T) {
	t.Parallel()
	cmd := exec.Command("python", "--version")
	err := cmd.Run()
	if err != nil {
		t.Error("python is not installed on the system.")
	}
}

func TestVenvCommandExists(t *testing.T) {
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
	_ = v.VenvCreate()
	if !v.IsVenv() {
		t.Error("python venv module is not installed on the system.")
	}
}
