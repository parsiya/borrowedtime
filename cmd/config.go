package cmd

import (
	"fmt"
	"time"

	prompt "github.com/c-bata/go-prompt"
	"github.com/parsiya/borrowedtime/config"
	"github.com/starkriedesel/prompter"
)

// Config command.

// ConfigCmd returns the config command.
func ConfigCmd() prompter.Cmd {

	resetCmd := prompter.SubCommand(
		"reset",
		"reset borrowed time and overwrite current configuration",
		resetExecutor,
	)
	resetCmd.AddOption("-file", "(optional) backup file before reset", false, resetCompleter)

	backupCmd := prompter.SubCommand(
		"backup",
		"backup configuration and templates",
		backupExecutor,
	)
	backupCmd.AddOption("-file", "(optional) backup file", false, resetCompleter)

	editConfigCmd := prompter.SubCommand(
		"edit",
		"edit configuration file",
		editConfigExecutor,
	)

	viewConfigCmd := prompter.SubCommand(
		"view",
		"view contents of configuration file",
		viewConfigExecutor,
	)

	// SubCommands are normal commands. In this case we are creating an argument
	// to add view as an option (but without the "-") instead of a subcommand.
	configCmd := prompter.SubCommand(
		"config",
		"configure workspace",
		configExecutor,
	)
	configCmd.AddOption(
		"restore",
		"restore configuration and templates",
		false,
		restoreCompleter,
	)
	configCmd.AddSubCommands(resetCmd, backupCmd, editConfigCmd, viewConfigCmd)

	return configCmd
}

// resetExecutor resets borrowed time and optionally creates a backup.
func resetExecutor(args prompter.CmdArgs) (err error) {
	fmt.Println("inside resetExecutor")
	fmt.Printf("args: %v\n", args)
	// If backup is specified, try and get the file name from first value.
	backupFile := ""
	createBackup := false
	if args.Contains("-file") {
		backupFile, err = args.GetFirstValue("-file")
		createBackup = true
		if err != nil {
			// If value does not exist, assume it was not provided.
			backupFile = ""
		}
	}
	return config.Reset(createBackup, backupFile)
}

// backupExecutor creates a backup. If no filename is specified, one with the
// current timestamp is created.
func backupExecutor(args prompter.CmdArgs) error {
	fmt.Println("inside backupExecutor")
	fmt.Printf("args: %v\n", args)
	if args.Contains("-file") {
		filename, err := args.GetFirstValue("-file")
		if err != nil {
			return err
		}
		return config.Backup(filename)
	}
	return config.Backup("")
}

// resetCompleter displays backup file suggestions for backup and reset commands.
func resetCompleter(optName string, _ []string) []prompt.Suggest {
	// Create an empty list of suggestions.
	sugs := []prompt.Suggest{}
	switch optName {
	case "-file":
		backupTimestamp := time.Now().Format("2006-01-02-15-04-05")
		sugs = append(sugs,
			prompt.Suggest{Text: backupTimestamp, Description: "current time"})
	}
	return sugs
}

// templateExecutor restores a backup and secretly creates one.
func configExecutor(args prompter.CmdArgs) error {
	fmt.Println("inside templateExecutor")
	fmt.Printf("deploy args: %v\n", args)

	if args.Contains("restore") {
		restoreFile, err := args.GetFirstValue("restore")
		if err != nil {
			return err
		}
		// Create a backup with the current timestamp.
		err = config.Backup("")
		if err != nil {
			return err
		}
		return config.Restore(restoreFile)
	}
	return nil
}

// restoreCompleter returns the list of files in the backup directory as suggestions.
func restoreCompleter(_ string, _ []string) []prompt.Suggest {
	// Create an empty list of suggestions.
	sugs := []prompt.Suggest{}
	// Get a list of all files in the backup directory.
	files, err := config.BackupFiles()
	// Print error and do not display any suggestions.
	if err != nil {
		return sugs
	}
	// Sort the files by name and add them to suggestions.
	// sort.Strings(files)
	// Add files to suggestions.
	for _, fi := range files {
		sugs = append(sugs, prompt.Suggest{Text: fi})
	}
	return sugs
}

// editConfigExecutor opens the config file with the workspace editor.
func editConfigExecutor(args prompter.CmdArgs) error {
	fmt.Println("inside editConfigExecutor")
	fmt.Printf("deploy args: %v\n", args)

	// Open config file with default editor.
	cfg, err := config.Read()
	if err != nil {
		return err
	}

	// // Check if editor is set.
	// if hasEditor := cfg.Has("editor"); !hasEditor {
	// 	return fmt.Errorf("editor is not set in config, use Deploy")
	// }
	return config.Edit(cfg.Key("editor"))
}

// viewConfigExecutor prints the config file to console.
func viewConfigExecutor(args prompter.CmdArgs) error {
	fmt.Println("inside viewConfigExecutor")
	fmt.Printf("deploy args: %v\n", args)

	cfg, err := config.Read()
	if err != nil {
		return err
	}

	data := [][]string{}
	for k, v := range cfg {
		data = append(data, []string{k, v})
	}
	fmt.Println(Table(data, false))
	return nil
}
