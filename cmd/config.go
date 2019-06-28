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
func ConfigCmd() prompter.Command {

	resetCmd := prompter.Command{
		Name:        "reset",
		Description: "reset borrowed time and overwrite current configuration",
		Executor:    resetExecutor,
	}
	resetCmd.AddArguments(prompter.Argument{
		Name:              "-file",
		Description:       "(optional) backup file before reset",
		ArgumentCompleter: resetCompleter,
	})

	backupCmd := prompter.Command{
		Name:        "backup",
		Description: "backup configuration and templates",
		Executor:    backupExecutor,
	}
	backupCmd.AddArguments(prompter.Argument{
		Name:              "-file",
		Description:       "(optional) backup file",
		ArgumentCompleter: resetCompleter,
	})

	editConfigCmd := prompter.Command{
		Name:        "edit",
		Description: "edit configuration file",
		Executor:    editConfigExecutor,
	}

	configCmd := prompter.Command{
		Name:        "config",
		Description: "configure workspace",
		Executor:    configExecutor,
	}
	configCmd.AddArguments(prompter.Argument{
		Name:              "restore",
		Description:       "restore configuration and templates",
		ArgumentCompleter: restoreCompleter,
	})
	configCmd.AddSubCommands(resetCmd, backupCmd, editConfigCmd)

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
	fmt.Printf("config args: %v\n", args)

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
	fmt.Printf("config args: %v\n", args)

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
