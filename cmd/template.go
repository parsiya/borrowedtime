package cmd

import (
	"fmt"

	prompt "github.com/c-bata/go-prompt"
	"github.com/parsiya/borrowedtime/config"
	"github.com/parsiya/borrowedtime/shared"
	"github.com/starkriedesel/prompter"
)

// Template command.

// TemplateCmd returns the template command.
func TemplateCmd() prompter.Command {

	listTemplateCmd := prompter.Command{
		Name:        "list",
		Description: "list all templates",
		Executor:    listTemplateExecutor,
	}

	// SubCommands are normal commands. In this case we are creating an argument
	// to add edit as an option (but without the "-") instead of a subcommand.
	templateCmd := prompter.Command{
		Name:        "template",
		Description: "template configuration",
		Executor:    templateExecutor,
	}

	addArgument := prompter.Argument{
		Name:              "add",
		Description:       "add template",
		ArgumentCompleter: addTemplateCompleter,
	}

	editArgument := prompter.Argument{
		Name:              "edit",
		Description:       "edit template",
		ArgumentCompleter: editTemplateCompleter,
	}
	templateCmd.AddArguments(addArgument, editArgument)

	templateCmd.AddSubCommands(listTemplateCmd)

	return templateCmd
}

// listTemplateExecutor lists all templates (a.k.a. all files in the templates directory).
func listTemplateExecutor(args prompter.CmdArgs) error {
	// Get template map.
	templateMap, sortedNames, err := config.Templates()
	if err != nil {
		return err
	}

	// Print the map as "templatename: fullpath".
	var data [][]string
	for _, name := range sortedNames {
		data = append(data, []string{name, templateMap[name]})
	}
	fmt.Println(Table(data, false))
	return nil
}

// templateExecutor executes the template command, in this case we are only
// going to use it for the view argument masquerading as subcommand.
func templateExecutor(args prompter.CmdArgs) (err error) {
	fmt.Println("inside templateExecutor")
	fmt.Printf("args: %v\n", args)

	if args.Contains("view") {
		// Get the template map.
		templateMap, _, err := config.Templates()
		if err != nil {
			return err
		}
		// Get the value.
		tmpl, err := args.GetFirstValue("view")
		if err != nil {
			return err
		}
		// Read the template.
		tmplString, err := shared.ReadFileString(templateMap[tmpl])
		if err != nil {
			return err
		}
		// And print.
		fmt.Println(tmplString)
	}

	if args.Contains("add") {
		// Read the value and add the file.
		tmplFile, err := args.GetFirstValue("add")
		if err != nil {
			return err
		}
		// Read contents of the file.
		tmplContent, err := shared.ReadFileString(tmplFile)
		if err != nil {
			return err
		}
		return config.AddTemplate(tmplFile, tmplContent, false)
	}

	if args.Contains("edit") {
		// Get the template map.
		templateMap, _, err := config.Templates()
		if err != nil {
			return err
		}
		// Get the value.
		tmpl, err := args.GetFirstValue("edit")
		if err != nil {
			return err
		}
		// Get the config dir.
		cfgDir, err := config.ConfigDir()
		if err != nil {
			return err
		}
		return OpenWith(templateMap[tmpl], cfgDir)
	}
	return nil
}

// editTemplateCompleter displays all template files for "template edit".
func editTemplateCompleter(_ string, _ []string) []prompt.Suggest {
	// Create an empty list of suggestions.
	sugs := []prompt.Suggest{}
	// Get template map.
	templateMap, sortedNames, err := config.Templates()
	if err != nil {
		return sugs
	}
	// Add template names in alphabetical order to suggestions.
	for _, name := range sortedNames {
		sugs = append(sugs, prompt.Suggest{
			Text:        name,
			Description: templateMap[name], // Remove this if it takes too much space.
		})
	}
	return sugs
}

// addTemplateCompleter displays a single suggestion for the add template option.
func addTemplateCompleter(_ string, _ []string) []prompt.Suggest {
	return []prompt.Suggest{
		{Text: "sample.json", Description: "path to the template file"},
	}
}
