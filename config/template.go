package config

import (
	"fmt"
	"path/filepath"

	"github.com/parsiya/borrowedtime/shared"
)

// Template represents one template.
type Template struct {
	Name     string `json:"name"`
	FullPath string `json:"fullpath"`
}

// templateDir returns the templates directory.
// "homedir/borrowedtime/templates" or "configDir/templates"
// TODO: Remove this.
func templateDir() (string, error) {
	configDir, err := configDir()
	if err != nil {
		return "", fmt.Errorf("config.TemplateDir: %s", err.Error())
	}
	return filepath.Join(configDir, "templates"), nil
}

// fileTemplateDir returns the "templates/file" directory.
// "homedir/borrowedtime/templates/file" or "configDir/templates/file"
func fileTemplateDir() (string, error) {
	configDir, err := configDir()
	if err != nil {
		return "", fmt.Errorf("config.TemplateDir: %s", err.Error())
	}
	return filepath.Join(configDir, "templates/file"), nil
}

// projectTemplateDir returns the "templates/project`" directory.
// "homedir/borrowedtime/templates/project" or "configDir/templates/project"
func projectTemplateDir() (string, error) {
	configDir, err := configDir()
	if err != nil {
		return "", fmt.Errorf("config.TemplateDir: %s", err.Error())
	}
	return filepath.Join(configDir, "templates/project"), nil
}

// FileTemplates returns returns map[filename]fullpath of all files inside the
// file template directory.
func FileTemplates() (mp map[string]string, err error) {
	// Get the file template directory.
	dir, err := fileTemplateDir()
	if err != nil {
		return mp, err
	}
	return templateMap(dir, "*")
}

// ProjectTemplates returns returns map[filename]fullpath of all files inside
// the project template directory.
func ProjectTemplates() (mp map[string]string, err error) {
	// Get the project template directory.
	dir, err := projectTemplateDir()
	if err != nil {
		return mp, err
	}
	return templateMap(dir, "*")
}

// templateMap creates and returns a map[TemplateName]FullPath of files matching
// pattern in root. Pattern is the typical "shell file name pattern" (e.g. *.exe
// or * to list all files).
// TemplateName is the name of the template file without extensions. Duplicates
// will be overwritten but that is expected.
func templateMap(root, pattern string) (map[string]string, error) {
	mp := make(map[string]string, 0)
	// List all files matching pattern in root.
	files, err := shared.ListFiles(root, pattern)
	if err != nil {
		return mp, err
	}
	// Add them all to the map. Duplicates will overwrite but that is expected.
	for _, fi := range files {
		mp[shared.RemoveExtension(fi)] = filepath.Join(root, fi)
	}
	return mp, nil
}

// TemplatePath returns template path. Returns "" if the template does not exist.
// Template names should be passed without the extension.
// TODO: Make private.
func TemplatePath(name string) (string, error) {
	// Get template map.
	mp, err := makeTemplateMap()
	if err != nil {
		return "", fmt.Errorf("config.TemplatePath: %s", err.Error())
	}
	return mp[name], nil
}

// makeTemplateMap returns map[templatename]fullpath.
func makeTemplateMap() (map[string]string, error) {
	mp := make(map[string]string, 0)

	tmplDir, _ := templateDir()
	tmplIndex, err := listTemplates(tmplDir)
	if err != nil {
		return mp, fmt.Errorf("config.makeTemplateMap: %s", err.Error())
	}

	// Convert slice to map.
	for _, tm := range tmplIndex {
		mp[tm.Name] = tm.FullPath
	}
	return mp, nil
}

// addFileTemplate adds a new file template to "templates/files".
func addFileTemplate(name, content string, overwrite bool) error {
	// Write the template to file.
	dir, err := fileTemplateDir()
	if err != nil {
		return err
	}
	return shared.WriteFileString(filepath.Join(dir, name), content, overwrite)
}

// addTemplate adds a new template to the root of templates directory.
func addTemplate(name, content string, overwrite bool) error {

	// Write template to file.
	tmplDir, _ := templateDir()
	pa := filepath.Join(tmplDir, name)
	err := shared.WriteFileString(pa, content, true)
	// Can replace with "return shared.WriteFileString" if we do not want the custom error message.
	if err != nil {
		return fmt.Errorf("config.AddTemplate: %s", err.Error())
	}
	return nil
}

// listTemplates creates a list of all files regardless of hierarchy under a
// location and returns a []Template.
func listTemplates(root string) (index []Template, err error) {
	files, err := shared.ListFiles(root, "*")
	if err != nil {
		return index, fmt.Errorf("config.listTemplates: %s", err.Error())
	}

	// Create full path of files by joining them with root.
	for _, fi := range files {
		ix := Template{
			// Remove extension from file name when creating the map.
			Name:     shared.RemoveExtension(fi),
			FullPath: filepath.Join(root, fi),
		}
		index = append(index, ix)
	}
	return index, err
}
