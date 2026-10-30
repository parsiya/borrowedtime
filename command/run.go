// Package command executes ordered commands without an implicit shell.
package command

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/parsiya/borrowedtime/shared"
)

// Command describes an executable and its literal arguments.
type Command struct {
	Executable string   `json:"executable"`
	Args       []string `json:"args,omitempty"`
}

// Run executes commands sequentially in dir and stops at the first failure.
// Example: Run([]Command{{Executable: "git", Args: []string{"init"}}}, "/projects/example").
func Run(commands []Command, dir string) error {
	for index, step := range commands {
		if step.Executable == "" {
			return fmt.Errorf("command.Run: step %d has no executable", index+1)
		}
	}
	for index, step := range commands {
		if shared.Debug {
			log.Printf("run command %d: %q in %q", index+1, step.Executable, dir)
		}
		process := exec.Command(step.Executable, step.Args...)
		process.Dir = dir
		process.Stdin = os.Stdin
		process.Stdout = os.Stdout
		process.Stderr = os.Stderr
		if err := process.Run(); err != nil {
			return fmt.Errorf("command.Run: step %d (%s) - %w", index+1, step.Executable, err)
		}
	}
	return nil
}
