package project

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/parsiya/borrowedtime/config"
	"github.com/parsiya/borrowedtime/shared"
)

// Project represents one project.
type Project struct {
	// ProjectName is the name of the project directory inside the workspace.
	ProjectName string `json:"projectname"`
	// Workspace is the path to the current workspace.
	Workspace string `json:"workspace"`
	// ProjectRoot is "workspace/projectname"
	ProjectRoot string `json:"projectroot"`
	// WorkspaceConfig is a copy of the workspace configuration.
	WorkspaceConfig map[string]string `json:"workspaceconfig"`
	// ProjectConfig contains project specific configuration.
	ProjectConfig map[string]string `json:"projectconfig"`
}

// New creates a new project.
func New(name string) *Project {
	cfg, _ := config.Read()
	return &Project{
		ProjectName:     shared.EscapeString(name),
		Workspace:       shared.EscapeString(cfg.Key("workspace")),
		ProjectRoot:     shared.EscapeString(filepath.Join(cfg.Key("workspace"), name)),
		WorkspaceConfig: cfg,
	}
}

// Create creates a project according to a specific directory structure template.
func (p *Project) Create(templateName string, overwrite bool) error {
	if p.ProjectName == "" || p.Workspace == "" {
		return fmt.Errorf("project.Project.Create: empty project")
	}
	// Generate template.
	tmpl, err := p.generateTemplate(templateName)
	if err != nil {
		return fmt.Errorf("project.Project.Create: %s", err.Error())
	}
	// Execute template.
	err = p.executeTemplate(tmpl, overwrite)
	if err != nil {
		return fmt.Errorf("project.Project.Create: %s", err.Error())
	}
	// Update project with newly created project config.
	// Project config is at "projectRoot/.config.json".
	configPath := filepath.Join(p.ProjectRoot, ".config.json")
	cfgBytes, err := shared.ReadFileByte(configPath)
	if err != nil {
		return fmt.Errorf("project.Project.Create: read project config - %s", err.Error())
	}
	err = json.Unmarshal(cfgBytes, &p.ProjectConfig)
	if err != nil {
		return fmt.Errorf("project.Project.Create: unmarshal project config - %s", err.Error())
	}
	return nil
}

// generateTemplate creates a template using the provided template string and project info.
func (p Project) generateTemplate(templateName string) (string, error) {
	return genTemplate(p, templateName)
}

// genTemplate creates a template using the provided template string and project info.
func genTemplate(p Project, tmplName string) (string, error) {
	pth, err := config.TemplatePath(tmplName)
	if err != nil {
		return "", err
	}
	// If template is not found.
	if pth == "" {
		return "", fmt.Errorf("project.genTemplate: template %s not found", tmplName)
	}
	// fmt.Println(pth)
	// Read file, apparently ParseFiles does not work.
	tmplStr, err := shared.ReadFileString(pth)
	if err != nil {
		return "", err
	}
	// Read the template and execute it.
	tmpl, err := template.New(tmplName).Parse(tmplStr)
	if err != nil {
		return "", fmt.Errorf("project.genTemplate: create new template - %s", err.Error())
	}

	var tplResult strings.Builder
	// fmt.Println("-----------------")
	if err := tmpl.Execute(&tplResult, p); err != nil {
		return "", fmt.Errorf("project.genTemplate: execute template - %s", err.Error())
	}
	return tplResult.String(), nil
}

// executeTemplate creates the directory structure according to a generated template.
func (p Project) executeTemplate(tmpl string, overwrite bool) error {
	return execProjectTemplate(p, tmpl, overwrite)
}

// execProjectTemplate creates the directory structure according to a generated template.
func execProjectTemplate(p Project, tmpl string, overwrite bool) error {
	root := &Node{}
	// Unmarshal template.
	if err := json.Unmarshal([]byte(tmpl), root); err != nil {
		return fmt.Errorf("project.executeTemplate: unmarshal template - %s", err.Error())
	}
	// Create directory structure.
	if err := root.Create(p, overwrite); err != nil {
		return fmt.Errorf("project.executeTemplate: create directory structure - %s", err.Error())
	}
	return nil
}

// Create creates the file or directory represented by the node and its children.
func (n *Node) Create(p Project, overwrite bool) error {

	exists, err := shared.PathExists(n.FullPath)
	// Check if we have access to the path.
	if err != nil {
		return fmt.Errorf("project.Node.Create: %s", err.Error())
	}
	// Check if path exists and return an error if overwrite is not set.
	if exists && !overwrite {
		// TODO: See how we can change this to support overwriting partial parts
		// of the project.
		return fmt.Errorf("project.Node.Create: %s already exists, use overwrite", n.FullPath)
	}
	// If node is a directory, create the directory. No need to create if it already exists.
	if n.Info.IsDir && !exists {
		if err := os.MkdirAll(n.FullPath, os.ModePerm); err != nil {
			return err
		}
	}
	// If file, create the file.
	if !n.Info.IsDir {
		// Node is a file, create the file.
		f, err := os.Create(n.FullPath)
		if err != nil {
			return err
		}
		defer f.Close()

		// generateTemplate will return an error if template is not found, which
		// means the path is empty.
		tmplString, _ := p.generateTemplate(n.Info.Template)
		// If template string is not empty, it has been generated.
		// We could also check for lack of error here, if err == nil.
		if tmplString != "" {
			num, err := f.WriteString(tmplString)
			if err != nil || num != len(tmplString) {
				return fmt.Errorf("project.Node.Create: error writing %s - %s", n.FullPath, err)
			}
		}
	}
	// Create all node's children.
	for _, child := range n.Children {
		if err := child.Create(p, overwrite); err != nil {
			return err
		}
	}
	return nil
}
