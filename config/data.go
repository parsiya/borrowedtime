package config

import (
	"fmt"
	"path/filepath"

	"github.com/parsiya/borrowedtime/shared"
)

// AddData adds a new data file to the root of data directory.
func AddData(name, content string, overwrite bool) error {
	// Get base file.
	name = filepath.Base(name)
	// Remove extension from data file name if any.
	name = shared.RemoveExtension(name)

	// Do not worry about duplicates.

	dataDir, _ := dataDir()

	// Add json extension.
	name = shared.AddExtension(name, "json")
	pa := filepath.Join(dataDir, name)
	err := shared.WriteFileString(content, pa, true)
	// Can replace with "return err" if we do not want the custom error message.
	if err != nil {
		return fmt.Errorf("config.AddData: %s", err.Error())
	}
	return nil
}

// DataFiles returns all files in the data directory.
func DataFiles() (fi []string, err error) {
	dataDir, err := dataDir()
	if err != nil {
		return fi, fmt.Errorf("config.DataFiles: %s", err.Error())
	}
	// List all files in the data directory.
	// Previously we used "*.*" here which filtered files without extensions.
	return shared.ListFiles(dataDir, "*")
}

// GetDataFile returns the contents of a data file using a path relative to the
// data directory.
func GetDataFile(filename string) (string, error) {
	pa, err := GetDataPath(filename)
	if err != nil {
		return "", fmt.Errorf("config.GetDataFile: %s", err.Error())
	}
	return shared.ReadFileString(pa)
}

// GetDataPath returns the full path to a data file.
func GetDataPath(filename string) (string, error) {
	// Get the path and join it with datadir.
	dataDir, err := dataDir()
	if err != nil {
		return "", fmt.Errorf("config.GetDataPath: %s", err.Error())
	}
	return filepath.Join(dataDir, filename), nil
}
