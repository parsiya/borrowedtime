package main

import (
	"flag"
	"fmt"

	prompt "github.com/c-bata/go-prompt"
	"github.com/parsiya/borrowedtime/cmd"
	"github.com/parsiya/borrowedtime/config"
	"github.com/parsiya/borrowedtime/shared"
	"github.com/starkriedesel/prompter"
)

func main() {
	flag.BoolVar(&shared.Debug, "debug", false, "enable diagnostic logging to stderr")
	flag.Parse()

	configCmd := cmd.ConfigCmd()
	deployCmd := cmd.DeployCmd()
	projectCmd := cmd.ProjectCmd()
	commandCmd := cmd.CommandCmd()
	exitCmd := cmd.ExitCmd()

	comp := prompter.NewCompleter()
	err := comp.RegisterCommands(configCmd, deployCmd, projectCmd, commandCmd, exitCmd)
	if err != nil {
		panic(err)
	}

	// Read the config file. If it does not exist, print an error.
	_, err = config.Read()
	if err != nil {
		fmt.Println("Cannot read config, use Deploy.")
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
