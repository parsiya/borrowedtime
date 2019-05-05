package config

// Default config.
const (
	// TODO: Convert this to template and pass a config struct.
	defaultWorkspaceConfig = `
{
	"editor": "",
	"workspace":"",
	"projectstructure":"",
	"burppath":"",
	"yourname":""
}`

	// Default workspace config template.
	defaultWorkspaceConfigFilename = "config.json"
)

// Default templates.
// Create content as a const string to "defaultTemplates.go" and add the names
// to this map to get them created in initConfig.
var defaultTemplates = map[string]string{
	"notes.json":             defaultNoteTemplate,
	"creds.json":             defaultCredsTemplate,
	"project-config.json":    defaultProjectConfig,
	"project-structure.json": defaultProjectStructure,
	"todo.json":              defaultTODOTemplate,
}

// Start menu path
var startMenuAllUsers = "C:/ProgramData/Microsoft/Windows/Start Menu/Programs"
