package cmd

import (
	"fmt"
	"os"
	"sort"

	prompt "github.com/c-bata/go-prompt"
	"github.com/parsiya/borrowedtime/command"
	"github.com/parsiya/borrowedtime/config"
	"github.com/starkriedesel/prompter"
)

// CommandCmd defines the interactive command-batch menu.
// Example: CommandCmd() defines the command list and command run subcommands.
func CommandCmd() prompter.Command {
	run := prompter.Command{
		Name: "run", Description: "run a named command batch", Executor: runCommandExecutor,
	}
	run.AddArguments(
		prompter.Argument{Name: " ", Description: "command batch name", ArgumentCompleter: commandBatchCompleter},
		prompter.Argument{Name: "-project", Description: "(optional) project working directory", ArgumentCompleter: openProjectCompleter},
	)
	menu := prompter.Command{Name: "command", Description: "run configured command batches"}
	menu.AddSubCommands(
		prompter.Command{Name: "list", Description: "list configured command batches", Executor: listCommandExecutor},
		run,
	)
	return menu
}

// runCommandExecutor runs one batch in the current directory or a selected project.
// Example: runCommandExecutor(prompter.CmdArgs{"_": {"setup"}}).
func runCommandExecutor(args prompter.CmdArgs) error {
	if len(args["_"]) != 1 || args["_"][0] == "" {
		return fmt.Errorf("command run: provide one command batch name")
	}
	name := args["_"][0]
	batches, err := config.Commands()
	if err != nil {
		return err
	}
	steps, exists := batches[name]
	if !exists {
		return fmt.Errorf("command run: batch %q not found", name)
	}
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	if args.Contains("-project") {
		name, err := args.GetFirstValue("-project")
		if err != nil {
			return err
		}
		if name == "" {
			return fmt.Errorf("command run: project name is empty")
		}
		dir, err = projectPath(name)
		if err != nil {
			return err
		}
	}
	return command.Run(steps, dir)
}

// listCommandExecutor prints the available batches in alphabetical order.
// Example: listCommandExecutor(nil) lists configured batches.
func listCommandExecutor(_ prompter.CmdArgs) error {
	batches, err := config.Commands()
	if err != nil {
		return err
	}
	var names []string
	for name := range batches {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		fmt.Println(name)
	}
	return nil
}

// commandBatchCompleter suggests configured batch names in alphabetical order.
// Example: commandBatchCompleter("", nil) suggests all batches.
func commandBatchCompleter(_ string, _ []string) []prompt.Suggest {
	batches, err := config.Commands()
	if err != nil {
		return nil
	}
	var names []string
	for name := range batches {
		names = append(names, name)
	}
	sort.Strings(names)
	var suggestions []prompt.Suggest
	for _, name := range names {
		suggestions = append(suggestions, prompt.Suggest{Text: name, Description: fmt.Sprintf("%d steps", len(batches[name]))})
	}
	return suggestions
}
