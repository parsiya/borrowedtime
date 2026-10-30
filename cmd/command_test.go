package cmd

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/parsiya/borrowedtime/command"
	"github.com/parsiya/borrowedtime/config"
	"github.com/starkriedesel/prompter"
)

func TestExecuteCLI(t *testing.T) {
	failure := errors.New("handler failed")
	var received prompter.CmdArgs
	create := prompter.Command{
		Name: "create",
		Executor: func(args prompter.CmdArgs) error {
			received = args
			return nil
		},
	}
	create.AddArguments(
		prompter.Argument{Name: " "},
		prompter.Argument{Name: "-template"},
		prompter.Argument{Name: "-overwrite"},
	)
	project := prompter.Command{Name: "project", Executor: func(args prompter.CmdArgs) error {
		received = args
		return nil
	}}
	project.AddSubCommands(create)
	project.AddArguments(prompter.Argument{Name: "open"})
	completer := prompter.NewCompleter()
	if err := completer.RegisterCommands(project, prompter.Command{
		Name: "fail", Executor: func(prompter.CmdArgs) error { return failure },
	}); err != nil {
		t.Fatal(err)
	}
	for _, test := range []struct {
		name string
		args []string
		want prompter.CmdArgs
	}{
		{"spaces", []string{"project", "create", "two words", "-template", "custom template", "-overwrite"}, prompter.CmdArgs{"_": {"two words"}, "-template": {"custom template"}, "-overwrite": nil}},
		{"flags-first", []string{"project", "create", "-overwrite", "-template", "custom", "example"}, prompter.CmdArgs{"_": {"example"}, "-template": {"custom"}, "-overwrite": nil}},
		{"literal", []string{"project", "create", "spaces ; & $HOME"}, prompter.CmdArgs{"_": {"spaces ; & $HOME"}}},
		{"terminator", []string{"project", "create", "--", "-example"}, prompter.CmdArgs{"_": {"-example"}}},
		{"named-option", []string{"project", "open", "two words"}, prompter.CmdArgs{"open": {"two words"}}},
		{"option-as-value", []string{"project", "open", "open"}, prompter.CmdArgs{"open": {"open"}}},
	} {
		t.Run(test.name, func(t *testing.T) {
			received = nil
			if err := ExecuteCLI(completer, test.args); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(received, test.want) {
				t.Fatalf("args = %v, want %v", received, test.want)
			}
		})
	}
	for _, args := range [][]string{
		nil,
		{"unknown"},
		{"project"},
		{"project", "unknown"},
		{"project", "create", "example", "-unknown"},
		{"project", "create", "example", "-template"},
		{"project", "create", "example", "-template", "-overwrite"},
		{"project", "create", "example", "-overwrite", "-overwrite"},
		{"project", "create", "example", "extra"},
		{"project", "open"},
	} {
		if err := ExecuteCLI(completer, args); err == nil {
			t.Fatalf("expected error for %v", args)
		}
	}
	if err := ExecuteCLI(completer, []string{"fail"}); !errors.Is(err, failure) {
		t.Fatalf("handler error lost: %v", err)
	}
}

func TestCommandProcess(t *testing.T) {
	for _, arg := range os.Args {
		if arg == "--borrowedtime-command-test" {
			if err := os.WriteFile("command-ran.txt", []byte("done"), 0600); err != nil {
				os.Exit(2)
			}
			os.Exit(0)
		}
	}
}

func TestCommandMenu(t *testing.T) {
	root := t.TempDir()
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
	if err := os.WriteFile(configPath, batches, 0600); err != nil {
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
	defer func() {
		if err := os.Chdir(oldDir); err != nil {
			t.Errorf("restore working directory: %v", err)
		}
	}()
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
	if err := os.WriteFile(configPath, []byte("{"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := listCommandExecutor(nil); err == nil {
		t.Fatal("expected malformed config error")
	}
}
