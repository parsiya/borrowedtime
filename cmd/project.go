package cmd

import (
	"fmt"

	prompt "github.com/c-bata/go-prompt"
	"github.com/parsiya/borrowedtime/config"
	"github.com/parsiya/borrowedtime/project"
	"github.com/starkriedesel/prompter"
)

// Project command.

// ProjectCmd returns the project command.
func ProjectCmd() prompter.Cmd {

	listProjectsCmd := prompter.SubCommand(
		"list",
		"list all projects in the workspace",
		listProjectExecutor,
	)

	projectCmd := prompter.SubCommand(
		"project",
		"project configuration",
		projectExecutor,
	)
	projectCmd.AddOption("open", "open project with default editor", false, openProjectCompleter)

	createProjectsCmd := prompter.SubCommand(
		"create",
		"create a new project in the workspace",
		createProjectExecutor,
	)
	createProjectsCmd.AddOption("name", "unique name of the new project", false, createProjectCompleter)
	createProjectsCmd.AddOption("template", "(optional) project template name", false, createProjectCompleter)

	projectCmd.AddSubCommands(listProjectsCmd, createProjectsCmd)
	return projectCmd
}

// projectExecutor executes the project command.
func projectExecutor(args prompter.CmdArgs) error {
	fmt.Println("inside projectExecutor")
	fmt.Printf("args: %v\n", args)

	if args.Contains("open") {
		project, err := args.GetFirstValue("open")
		if err != nil {
			return err
		}
		return OpenProject(project)
	}

	// if args.Contains("add") {
	// 	// Read the value and add the file.
	// 	dataFile, err := args.GetFirstValue("add")
	// 	if err != nil {
	// 		return err
	// 	}
	// 	// Read contents of the file.
	// 	dataContent, err := shared.ReadFileString(dataFile)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	return config.AddData(dataFile, dataContent, false)
	// }
	return nil
}

// openProjectCompleter lists all top-level directories in the workspace.
func openProjectCompleter(_ string, _ []string) []prompt.Suggest {
	// Create an empty list of suggestions.
	sugs := []prompt.Suggest{}
	// Get workspace path.
	workspace, err := workspacePath()
	if err != nil {
		return sugs
	}
	// Get top-level directories with ioutil.ReadDir.
	dirs, err := TopDirs(workspace)
	if err != nil {
		return sugs
	}
	for _, dir := range dirs {
		sugs = append(sugs, prompt.Suggest{
			// Text will be directory name and description is full path.
			Text:        dir[0],
			Description: dir[1],
		})
	}
	return sugs
}

// listProjectExecutor lists all projects.
func listProjectExecutor(args prompter.CmdArgs) error {
	// Get workspace path.
	workspace, err := workspacePath()
	if err != nil {
		return err
	}
	dirs, err := TopDirs(workspace)
	if err != nil {
		return err
	}
	fmt.Println(Table(dirs, false))
	return nil
}

// createProjectCompleter shows sample suggestions for the create project command.
func createProjectCompleter(optName string, _ []string) []prompt.Suggest {
	// Create an empty list of suggestions.
	sugs := []prompt.Suggest{}
	switch optName {
	case "name":
		sugs = append(sugs,
			prompt.Suggest{Text: "", Description: "should be unique"})
		return sugs
	case "template":
		// Show all templates. This will have non-project templates, but we can
		// delegate that responsibility to the user.

		// Taken from viewTemplateCompleter in template.go.
		// Get template map.
		templateMap, sortedNames, err := config.Templates()
		if err != nil {
			return sugs
		}
		// Add template names in alphabetical order to suggestions.
		for _, name := range sortedNames {
			sugs = append(sugs, prompt.Suggest{
				Text:        name,
				Description: templateMap[name],
			})
		}
	}
	return sugs
}

// createProjectExecutor is the executor for the project command. Creates a new
// project name with using an optional template.
func createProjectExecutor(args prompter.CmdArgs) (err error) {
	fmt.Println("inside createProjectExecutor")
	fmt.Printf("args: %v\n", args)

	// Read config.
	cfg, err := config.Read()
	if err != nil {
		return err
	}

	projectName := ""
	templateName := ""

	if args.Contains("name") {
		projectName, err = args.GetFirstValue("name")
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("please provide project name")
	}

	if args.Contains("template") {
		templateName, err = args.GetFirstValue("name")
		if err != nil {
			return err
		}
	} else {
		// If not provided, use the default project structure in the config.
		// Check if the config has the
		templateName = cfg.Key("projectstructure")
	}

	// Create project.
	prj := project.New(projectName)
	// Do not overwrite by default.
	// TODO: We can add an overwrite flag here but I want to keep it simple.
	err = prj.Create(templateName, false)
	if err != nil {
		return err
	}
	return OpenProject(projectName)

	// if args.Contains("add") {
	// 	// Read the value and add the file.
	// 	dataFile, err := args.GetFirstValue("add")
	// 	if err != nil {
	// 		return err
	// 	}
	// 	// Read contents of the file.
	// 	dataContent, err := shared.ReadFileString(dataFile)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	return config.AddData(dataFile, dataContent, false)
	// }
}
