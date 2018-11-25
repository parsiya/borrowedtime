package cmd

import (
	"github.com/starkriedesel/prompter"
)

// Exit Command.

// ExitCmd returns the exit command. This cannot be in main because of import cycle.
func ExitCmd() prompter.Cmd {
	return prompter.ExitCommand("exit", "exit the application")
}
