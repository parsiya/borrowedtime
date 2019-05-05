package cmd

import (
	"fmt"

	prompt "github.com/c-bata/go-prompt"
	"github.com/parsiya/borrowedtime/config"
	"github.com/parsiya/borrowedtime/shared"
	"github.com/starkriedesel/prompter"
)

// Data command.

// DataCmd returns the data command.
func DataCmd() prompter.Command {

	listDataCmd := prompter.Command{
		Name:        "list",
		Description: "list all data files",
		Executor:    listDataExecutor,
	}

	dataCmd := prompter.Command{
		Name:        "data",
		Description: "data files configuration",
		Executor:    dataExecutor,
	}

	addArgument := prompter.Argument{
		Name:              "add",
		Description:       "add data file",
		ArgumentCompleter: addDataCompleter,
	}

	editArgument := prompter.Argument{
		Name:              "edit",
		Description:       "edit data file",
		ArgumentCompleter: editDataCompleter,
	}

	// Add all arguments to the command.
	dataCmd.AddArguments(addArgument, editArgument)

	dataCmd.AddSubCommands(listDataCmd)

	return dataCmd
}

// listExecutor lists all data files (a.k.a. all files in the data directory).
func listDataExecutor(args prompter.CmdArgs) error {
	// List all data files.
	files, err := config.DataFiles()
	if err != nil {
		return err
	}
	var data [][]string
	for _, fi := range files {
		data = append(data, []string{fi})
	}
	fmt.Println(Table(data, false))
	return nil
}

// dataExecutor executes the data command.
func dataExecutor(args prompter.CmdArgs) (err error) {
	fmt.Println("inside dataExecutor")
	fmt.Printf("args: %v\n", args)

	if args.Contains("view") {
		dataFile, err := args.GetFirstValue("view")
		if err != nil {
			return err
		}
		content, err := config.GetDataFile(dataFile)
		if err != nil {
			return err
		}
		fmt.Println(content)
		return nil
	}

	if args.Contains("add") {
		// Read the value and add the file.
		dataFile, err := args.GetFirstValue("add")
		if err != nil {
			return err
		}
		// Read contents of the file.
		dataContent, err := shared.ReadFileString(dataFile)
		if err != nil {
			return err
		}
		return config.AddData(dataFile, dataContent, false)
	}

	if args.Contains("edit") {
		dataFile, err := args.GetFirstValue("edit")
		if err != nil {
			return err
		}
		dataPath, err := config.GetDataPath(dataFile)
		if err != nil {
			return err
		}
		cfgDir, err := config.ConfigDir()
		if err != nil {
			return err
		}
		return OpenWith(dataPath, cfgDir)
	}

	return nil
}

// editDataCompleter displays all data files for the "data edit" command.
func editDataCompleter(_ string, _ []string) []prompt.Suggest {
	// Create an empty list of suggestions.
	sugs := []prompt.Suggest{}
	// Get data files.
	files, err := config.DataFiles()
	if err != nil {
		return sugs
	}

	// Add data files to suggestions.
	for _, fi := range files {
		sugs = append(sugs, prompt.Suggest{
			Text: fi,
		})
	}
	return sugs
}

// addCompleter displays a single suggestion for the add data option.
func addDataCompleter(_ string, _ []string) []prompt.Suggest {
	return []prompt.Suggest{
		{Text: "data.json", Description: "path to the data file"},
	}
}
