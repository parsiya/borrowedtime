package shared

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
)

// Open with methods.

// OpenWithEditor opens a path (file or directory) with the editor in the config file.
// The editor must support passing the path as the first parameter.
func OpenWithEditor(editor, path string) error {
	// Check if target path exists.
	exists, err := PathExists(path)
	if err != nil {
		return fmt.Errorf("project.OpenWithEditor: target path check - %s", err.Error())
	}
	if !exists {
		return fmt.Errorf("project.OpenWithEditor: target path not found at %s", path)
	}

	// Clean the path. Is this needed?
	path = filepath.Clean(path)

	cmd := exec.Command(editor, path)
	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("config.OpenWithEditor: open %s with %s editor - %s", path, editor, err.Error())
	}
	return nil
}

// OpenWithDefaultEditor opens a file with the default editor depdening on the OS.
// Only Windows is supported.
// TODO: Add other OS default editors.
func OpenWithDefaultEditor(path string) error {
	// Check if target path exists.
	exists, err := PathExists(path)
	if err != nil {
		return fmt.Errorf("project.OpenWithDefaultEditor: target path check - %s", err.Error())
	}
	if !exists {
		return fmt.Errorf("project.OpenWithDefaultEditor: target path not found at %s", path)
	}

	if runtime.GOOS == "windows" {
		return OpenWithEditor("notepad", path)
	}
	return nil
}
