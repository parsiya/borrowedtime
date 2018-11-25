package config

// Add individual file templates here.

const (
	defaultNoteTemplate = `# {{ .ProjectName }} Notes

## `

	defaultCredsTemplate = `# {{ .ProjectName }} Credentials

`

	defaultTODOTemplate = `# {{ .ProjectName }} TODO

## `

	defaultProjectConfig = `
{
	"pix": "{{ .Workspace }}\\{{ .ProjectName }}\\@pix",
	"findings":"{{ .Workspace }}\\{{ .ProjectName }}\\@findings.md",
	"notes":"{{ .Workspace }}\\{{ .ProjectName }}\\@notes.md",
	"root":"{{ .Workspace }}\\{{ .ProjectName }}",
	"report":"{{ .Workspace }}\\{{ .ProjectName }}\\@report",
	"reportconfig":"{{ .Workspace }}\\{{ .ProjectName }}\\@report\\report.json"
}
`
)
