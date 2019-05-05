package shared

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
)

// Open with methods.

// OpenWithEditor opens paths with the editor in the config file.
// The editor must support passing one or more paths as the first parameter
// or just ignore everything after the first one.
func OpenWithEditor(editor string, paths ...string) error {
	// Only check if the first target path exists.
	exists, err := PathExists(paths[0])
	if err != nil {
		return fmt.Errorf("project.OpenWithEditor: target path check - %s", err.Error())
	}
	if !exists {
		return fmt.Errorf("project.OpenWithEditor: target path not found at %s", paths[0])
	}

	// Clean the paths. Is this needed?
	cleanPaths := paths[:0]
	for _, path := range paths {
		cleanPaths = append(cleanPaths, filepath.Clean(path))
	}

	cmd := exec.Command(editor, cleanPaths...)
	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("config.OpenWithEditor: open %s with %s editor - %s", paths, editor, err.Error())
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
