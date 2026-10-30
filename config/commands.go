package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/parsiya/borrowedtime/command"
	"github.com/parsiya/borrowedtime/shared"
)

// Commands loads named command batches from data/commands.json.
// Example: Commands() loads the current deployment's batches.
func Commands() (map[string][]command.Command, error) {
	batches := make(map[string][]command.Command)
	dir, err := dataDir()
	if err != nil {
		return nil, err
	}
	content, err := shared.ReadFileByte(filepath.Join(dir, "commands.json"))
	if os.IsNotExist(err) {
		return batches, nil
	}
	if err != nil {
		return nil, fmt.Errorf("config.Commands: %w", err)
	}
	if err := json.Unmarshal(content, &batches); err != nil {
		return nil, fmt.Errorf("config.Commands: %w", err)
	}
	return batches, nil
}
