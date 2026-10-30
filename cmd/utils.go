package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	xj "github.com/basgys/goxml2json"
	"github.com/olekukonko/tablewriter"
	"github.com/parsiya/borrowedtime/config"
	"github.com/parsiya/borrowedtime/shared"
	"github.com/starkriedesel/prompter"
)

// Utilities.

// ExecuteCLI dispatches shell arguments through the registered command tree.
// Example: ExecuteCLI(comp, []string{"project", "create", "example"}).
func ExecuteCLI(comp *prompter.Completer, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("provide a command")
	}
	var selected *prompter.Command
	for len(args) > 0 {
		next, exists := comp.Commands[args[0]]
		if !exists {
			break
		}
		selected = next
		comp = next.Completer
		args = args[1:]
	}
	if selected == nil {
		return fmt.Errorf("unknown command %q", args[0])
	}
	if selected.Executor == nil || (len(args) == 0 && len(comp.Commands) > 0) {
		return fmt.Errorf("incomplete command %q", selected.Name)
	}
	parsed := prompter.CmdArgs{}
	positionalOnly := false
	for index := 0; index < len(args); index++ {
		arg := args[index]
		if arg == "--" && !positionalOnly {
			positionalOnly = true
			continue
		}
		option, exists := comp.Options[arg]
		if exists && arg != " " && !positionalOnly {
			if parsed.Contains(arg) && !option.Repeatable {
				return fmt.Errorf("option %q specified more than once", arg)
			}
			if arg == "-overwrite" {
				parsed[arg] = nil
				continue
			}
			if index+1 == len(args) {
				return fmt.Errorf("option %q requires a value", arg)
			}
			value := args[index+1]
			if _, isOption := comp.Options[value]; (isOption && strings.HasPrefix(value, "-")) || value == "--" {
				return fmt.Errorf("option %q requires a value", arg)
			}
			parsed[arg] = append(parsed[arg], value)
			index++
			continue
		}
		if strings.HasPrefix(arg, "-") && !positionalOnly {
			return fmt.Errorf("unknown option %q", arg)
		}
		if _, allowed := comp.Options[" "]; !allowed || len(parsed["_"]) != 0 {
			return fmt.Errorf("unexpected argument %q", arg)
		}
		parsed["_"] = append(parsed["_"], arg)
	}
	return selected.Executor(parsed)
}

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
	files, err := os.ReadDir(root)
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

// OpenWith opens paths with the editor specified in the config.
func OpenWith(paths ...string) error {
	// Get the config.
	cfg, err := config.Read()
	if err != nil {
		return err
	}
	// Open the paths.
	return shared.OpenWithEditor(cfg.Key("editor"), paths...)
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
	projectPath := filepath.Join(workspace, projectName)
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
