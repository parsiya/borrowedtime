package config

import (
	"fmt"
	"path/filepath"
	"sort"

	"github.com/parsiya/borrowedtime/shared"
)

// Template represents one template.
type Template struct {
	Name     string `json:"name"`
	FullPath string `json:"fullpath"`
}

// TemplateDir returns the templates directory.
// "homedir/borrowedtime/templates" or "configDir/templates"
// TODO: Make private.
func TemplateDir() (string, error) {
	configDir, err := configDir()
	if err != nil {
		return "", fmt.Errorf("config.TemplateDir: %s", err.Error())
	}
	return filepath.Join(configDir, "templates"), nil
}

// TemplatePath returns template path. Returns "" if template does not exist.
// Template names should be passed without the extension.
// TODO: Make private.
func TemplatePath(name string) (string, error) {
	// Remove extension if it's passed in template name.
	name = shared.RemoveExtension(name)
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

	tmplDir, _ := TemplateDir()
	tmplIndex, err := listTemplateFiles(tmplDir)
	if err != nil {
		return mp, fmt.Errorf("config.makeTemplateMap: %s", err.Error())
	}

	// Convert slice to map.
	for _, tm := range tmplIndex {
		mp[tm.Name] = tm.FullPath
	}
	return mp, nil
}

// AddTemplate adds a new template to the root of templates directory.
func AddTemplate(name, content string, overwrite bool) error {
	// Get Base.
	name = filepath.Base(name)
	// Remove extension from template name if any.
	name = shared.RemoveExtension(name)

	// Get template path to check for duplicates.
	tmplPath, err := TemplatePath(name)
	if err != nil {
		return fmt.Errorf("config.AddTemplate: %s", err.Error())
	}
	// Check for duplicates.
	if tmplPath != "" && !overwrite {
		return fmt.Errorf("config.AddTemplate: template %s already exists at %s", name, tmplPath)
	}

	// Write template to file.
	tmplDir, _ := TemplateDir()
	// Add json extension. This does not change the template name in map because
	// extensions are removed when creating the map but it helps with syntax
	// highlighting and editor support during file edit.
	name = shared.AddExtension(name, "json")
	pa := filepath.Join(tmplDir, name)
	err = shared.WriteFileString(content, pa, true)
	// Can replace with "return shared.WriteFileString" if we do not want the custom error message.
	if err != nil {
		return fmt.Errorf("config.AddTemplate: %s", err.Error())
	}
	return nil
}

// listTemplateFiles creates a list of all files regardless of hierarchy under a
// location and returns a []Template.
func listTemplateFiles(root string) (index []Template, err error) {
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

// Templates returns the template map and a sorted slice of keys to print
// them in order.
func Templates() (map[string]string, []string, error) {
	templateMap, err := makeTemplateMap()
	if err != nil {
		return templateMap, []string{}, err
	}
	// Sort the map keys and then print them in that order.
	names := make([]string, len(templateMap))
	i := 0
	for name := range templateMap {
		names[i] = name
		i++
	}
	sort.Strings(names)
	return templateMap, names, nil
}
