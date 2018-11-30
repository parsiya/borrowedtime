package cmd

import (
	"fmt"

	"github.com/parsiya/borrowedtime/config"
	"github.com/starkriedesel/prompter"
)

// Deploy command.

// DeployCmd returns the init command that initializes borrowed time.
func DeployCmd() prompter.Command {
	return prompter.Command{
		Name:        "deploy",
		Description: "deploy borrowed time and generate a config file",
		Executor:    deployExecutor,
	}
}

// deployExecutor initializes borrowed time but fails if config is already present.
func deployExecutor(args prompter.CmdArgs) error {
	fmt.Println("inside deployExecutor")
	fmt.Printf("args: %v\n", args)

	return config.Deploy()
	// err := config.Deploy()
	// if err != nil {
	// 	return err
	// }

	// // Read config.
	// cfg, err = config.Read()
	// if err != nil {
	// 	return err
	// }

	// // Enable all commands.
	// configCmd.Hidden = false
	// templateCmd.Hidden = false
	// dataCmd.Hidden = false

	// return nil
}
