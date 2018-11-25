package cmd

import (
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"

	xj "github.com/basgys/goxml2json"
	"github.com/olekukonko/tablewriter"
	"github.com/parsiya/borrowedtime/config"
	"github.com/parsiya/borrowedtime/shared"
)

// Utilities.

// Table returns the [][]string in a tablewriter table.
func Table(data [][]string, border bool) string {
	var sb strings.Builder
	table := tablewriter.NewWriter(&sb)
	table.SetBorder(border) // Remove borders.
	if !border {
		table.SetColumnSeparator("") // Remove column.
	}
	table.AppendBulk(data)
	table.Render()
	return sb.String()
}

// TopDirs returns the name and full path of top-level directories of root in this format:
// [][]string{name, fullpath}.
func TopDirs(root string) (dirs [][]string, err error) {
	// Use ioutil.ReadDir to get a slice of []os.FileInfo.
	files, err := ioutil.ReadDir(root)
	if err != nil {
		return dirs, err
	}
	for _, file := range files {
		if file.IsDir() {
			dirs = append(dirs, []string{file.Name(), filepath.Join(root, file.Name())})
		}
	}
	return dirs, nil
}

// OpenWith opens the file or directory with the editor specified in the config.
func OpenWith(filename string) error {
	// Get the config.
	cfg, err := config.Read()
	if err != nil {
		return err
	}
	// Open the path.
	return shared.OpenWithEditor(cfg.Key("editor"), filename)
}

// OpenProject opens projectName using the editor specified in the config file.
func OpenProject(projectName string) error {
	prjPath, err := projectPath(projectName)
	if err != nil {
		return err
	}
	return OpenWith(prjPath)
}

// projectPath returns the path to projectName's directory.
func projectPath(projectName string) (string, error) {
	workspace, err := workspacePath()
	if err != nil {
		return "", err
	}
	projectPath := path.Join(workspace, projectName)
	return projectPath, nil
}

// workspacePath reads the config file and returns the workspace path.
func workspacePath() (string, error) {
	// Get config.
	cfg, err := config.Read()
	if err != nil {
		return "", err
	}
	// Check if workspace is set.
	if !cfg.Has("workspace") {
		return "", fmt.Errorf("workspace is missing from the config file")
	}
	return cfg.Key("workspace"), nil
}

// XMLToJSON converts an xml file to JSON.
func XMLToJSON(xmlString string) (string, error) {
	xml := strings.NewReader(xmlString)
	json, err := xj.Convert(xml)
	if err != nil {
		return "", fmt.Errorf("convert.XMLToJSON: %s", err.Error())
	}
	return json.String(), nil
}
