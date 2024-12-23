//go:build integration

package main

import (
	"os"
	"testing"
)

func TestCommandExists(t *testing.T) {
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
		t.Error("Command doesn't exist.")
	}
}
