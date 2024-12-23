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
	_ = VenvCreate(dir)
	if !IsVenv(dir) {
		t.Error("Command doesn't exist.")
	}
}
