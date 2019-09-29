package io

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

//go:generate mockgen -destination mock_cui/io_mock.go github.com/tevino/the-clean-architecture-demo/todo/cui/io IO

// IO abstracts input/output functions tied to the OS.
type IO interface {
	GetInputByLaunchingEditor() (string, error)
}

// UnixLikeIO implements IO for *nix platform.
type UnixLikeIO struct{}

func (UnixLikeIO) GetInputByLaunchingEditor() (string, error) {
	filePath := tempFilePath("task")
	defer os.Remove(filePath)
	err := launchEditor(filePath)
	if err != nil {
		return "", fmt.Errorf("launching editor: %w", err)
	}
	fl, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("reading temp file: %w", err)
	}
	buf, err := ioutil.ReadAll(fl)
	if err != nil {
		return "", fmt.Errorf("reading file: %w", err)
	}
	return string(buf), nil
}

const defaultEditor = "vi"

func getEditor() string {
	var editor string

	if editor == "" {
		editor = os.Getenv("VISUAL")
	}
	if editor == "" {
		editor = os.Getenv("EDITOR")
	}
	if editor == "" {
		editor = defaultEditor
	}

	return editor
}

func launchEditor(path string) error {
	editor := getEditor()

	cmd := exec.Command(editor, path)
	cmd.Env = os.Environ()
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("running editor: %w", err)
	}
	return nil
}

func tempFilePath(pattern string) string {
	fl, err := ioutil.TempFile("", pattern)
	if err != nil {
		panic(err)
	}
	defer fl.Close()
	return fl.Name()
}
