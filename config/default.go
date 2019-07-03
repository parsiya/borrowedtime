package config

// Default config.
const (
	// TODO: Convert this to template and pass a config struct.
	defaultConfig = `
{
	"editor": "",
	"workspace":"",
	"projectstructure":"",
	"burppath":"",
	"yourname":""
}`

	// Default workspace config template.
	defaultConfigFilename = "config.json"
)

// Default templates.
// Create content as a const string to "defaultTemplates.go" and add the names
// to this map to get them created in initConfig.

// Default file templates.
var defaultFileTemplates = map[string]string{
	"notes.md":          defaultNoteTemplate,
	"creds.md":          defaultCredsTemplate,
	"project-config.md": defaultProjectConfig,
	"todo.md":           defaultTODOTemplate,
}

// Default project templates.
var defaultProjectTemplates = map[string]string{
	"project-structure.json": defaultProjectStructure,
}
