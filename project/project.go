package project

import (
	"encoding/json"
	"fmt"
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
	tmpl, err := p.generateTemplate(templateName, true)
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
		// If there is no project config in the template, this will return an
		// error. Print an error and return instead.
		fmt.Printf("project.Project.Create: read project config - %s", err.Error())
		return nil
		// return fmt.Errorf("project.Project.Create: read project config - %s", err.Error())
	}
	err = json.Unmarshal(cfgBytes, &p.ProjectConfig)
	if err != nil {
		// If the project config file is empty (no template or wrong template)
		// then we will get an error here that we can print and return.
		fmt.Printf("project.Project.Create: unmarshal project config - %s", err.Error())
		return nil
		// return fmt.Errorf("project.Project.Create: unmarshal project config - %s", err.Error())
	}
	return nil
}

// generateTemplate creates a template using the provided template string and
// project info. If isProject is set to true then we are generating a project,
// otherwise we are generating a file.
func (p Project) generateTemplate(templateName string, isProject bool) (string, error) {
	return genTemplate(p, templateName, isProject)
}

// genTemplate creates a template using the provided template string and project
// info. If isProject is set to true then we are generating a project, otherwise
// we are generating a file
func genTemplate(p Project, tmplName string, isProject bool) (string, error) {

	pth := ""
	// Remove extension from tmplName if any.
	tmplName = shared.RemoveExtension(tmplName)

	if isProject {
		// Get project templates.
		prjTmpls, err := config.ProjectTemplates()
		if err != nil {
			return "", err
		}
		// Get project template path if it exists.
		pth = prjTmpls[tmplName]
	} else {
		// Get file templates.
		fileTmpls, err := config.FileTemplates()
		if err != nil {
			return "", err
		}
		pth = fileTmpls[tmplName]
	}

	// If template is not found.
	if pth == "" {
		return "", fmt.Errorf("project.genTemplate: template %s not found", tmplName)
	}

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
