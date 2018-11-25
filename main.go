package main

import (
	"fmt"

	prompt "github.com/c-bata/go-prompt"
	"github.com/parsiya/borrowedtime/cmd"
	"github.com/parsiya/borrowedtime/config"
	"github.com/starkriedesel/prompter"
)

func main() {

	configCmd := cmd.ConfigCmd()
	deployCmd := cmd.DeployCmd()
	templateCmd := cmd.TemplateCmd()
	dataCmd := cmd.DataCmd()
	projectCmd := cmd.ProjectCmd()

	exitCmd := cmd.ExitCmd()

	comp := prompter.NewCompleter()
	err := comp.RegisterCommands(configCmd, deployCmd, templateCmd, dataCmd, projectCmd, exitCmd)
	if err != nil {
		panic(err)
	}

	// Read the config file. If it does not exist, print an error.
	_, err = config.Read()
	if err != nil {
		fmt.Println("Config could not be read, use Deploy.")
	}

	// fmt.Println(shared.StructToJSONString(cfg, true))

	p := prompt.New(
		comp.Execute,
		comp.Complete,
		prompt.OptionPrefix(">>> "),
		prompt.OptionPrefixTextColor(prompt.White),
		prompt.OptionTitle("borrowed time"),
		prompt.OptionMaxSuggestion(10), // TODO: Add this to the config file?
	)

	p.Run()
}
