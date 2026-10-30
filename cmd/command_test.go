package cmd

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/parsiya/borrowedtime/command"
	"github.com/parsiya/borrowedtime/config"
	"github.com/starkriedesel/prompter"
)

func TestCommandProcess(t *testing.T) {
	for _, arg := range os.Args {
		if arg == "--borrowedtime-command-test" {
			if err := ioutil.WriteFile("command-ran.txt", []byte("done"), 0600); err != nil {
				os.Exit(2)
			}
			os.Exit(0)
		}
	}
}

func TestCommandMenu(t *testing.T) {
	root, err := ioutil.TempDir("", "borrowedtime-menu-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(root)
	oldHome, hadHome := os.LookupEnv("HOME")
	oldCache := homedir.DisableCache
	if err := os.Setenv("HOME", root); err != nil {
		t.Fatal(err)
	}
	homedir.DisableCache = true
	defer func() {
		homedir.DisableCache = oldCache
		if hadHome {
			os.Setenv("HOME", oldHome)
		} else {
			os.Unsetenv("HOME")
		}
	}()
	if batches, err := config.Commands(); err != nil || len(batches) != 0 {
		t.Fatalf("optional command config = %v, error = %v", batches, err)
	}
	dataDir := filepath.Join(root, "borrowedtime", "data")
	projectDir := filepath.Join(root, "projects", "example")
	for _, dir := range []string{dataDir, projectDir} {
		if err := os.MkdirAll(dir, 0700); err != nil {
			t.Fatal(err)
		}
	}
	if err := config.Write(config.ConfigMap{"workspace": filepath.Dir(projectDir)}); err != nil {
		t.Fatal(err)
	}
	batches, err := json.Marshal(map[string][]command.Command{
		"setup": {{Executable: os.Args[0], Args: []string{"-test.run=^TestCommandProcess$", "--", "--borrowedtime-command-test"}}},
		"empty": {},
	})
	if err != nil {
		t.Fatal(err)
	}
	configPath := filepath.Join(dataDir, "commands.json")
	if err := ioutil.WriteFile(configPath, batches, 0600); err != nil {
		t.Fatal(err)
	}
	completer := prompter.NewCompleter()
	if err := completer.RegisterCommands(CommandCmd()); err != nil {
		t.Fatal(err)
	}
	suggestions := commandBatchCompleter("", nil)
	if len(suggestions) != 2 || suggestions[0].Text != "empty" || suggestions[1].Text != "setup" {
		t.Fatalf("suggestions = %v", suggestions)
	}
	if err := listCommandExecutor(nil); err != nil {
		t.Fatal(err)
	}
	completer.Execute("command run setup -project example")
	if _, err := os.Stat(filepath.Join(projectDir, "command-ran.txt")); err != nil {
		t.Fatalf("command did not run in project: %v", err)
	}
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(oldDir)
	completer.Execute("command run setup")
	if _, err := os.Stat(filepath.Join(root, "command-ran.txt")); err != nil {
		t.Fatalf("command did not run in current directory: %v", err)
	}
	for _, args := range []prompter.CmdArgs{
		nil,
		{"_": {"missing"}},
		{"_": {"setup", "extra"}},
		{"_": {"setup"}, "-project": {}},
	} {
		if err := runCommandExecutor(args); err == nil {
			t.Fatalf("expected argument error for %v", args)
		}
	}
	if err := ioutil.WriteFile(configPath, []byte("{"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := listCommandExecutor(nil); err == nil {
		t.Fatal("expected malformed config error")
	}
}
