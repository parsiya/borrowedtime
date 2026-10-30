package main

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/parsiya/borrowedtime/command"
	"github.com/parsiya/borrowedtime/project"
)

func TestMain(tests *testing.M) {
	if os.Getenv("BORROWEDTIME_TEST_EDITOR") == "1" && len(os.Args) == 2 && filepath.IsAbs(os.Args[1]) {
		info, err := os.Stat(os.Args[1])
		if err != nil || !info.IsDir() {
			os.Exit(2)
		}
		os.Exit(0)
	}
	os.Exit(tests.Run())
}

func TestCLICommandProcess(t *testing.T) {
	for index, arg := range os.Args {
		if arg != "--borrowedtime-command" {
			continue
		}
		args := os.Args[index+1:]
		if len(args) == 1 && args[0] == "fail" {
			os.Exit(7)
		}
		if len(args) == 3 && args[0] == "write" {
			if err := os.WriteFile(args[1], []byte(args[2]), 0600); err != nil {
				os.Exit(2)
			}
			os.Exit(0)
		}
		os.Exit(2)
	}
}

func TestMainCLI(t *testing.T) {
	root := t.TempDir()
	binary := filepath.Join(root, "borrowedtime.exe")
	buildCtx, cancelBuild := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancelBuild()
	build := exec.CommandContext(buildCtx, "go", "build", "-mod=readonly", "-o", binary, ".")
	if output, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build CLI: %v\n%s", err, output)
	}
	home := filepath.Join(root, "home with spaces")
	workspace := filepath.Join(root, "workspace with spaces")
	configDir := filepath.Join(home, "borrowedtime")
	for _, dir := range []string{filepath.Join(configDir, "backups"), workspace} {
		if err := os.MkdirAll(dir, 0700); err != nil {
			t.Fatal(err)
		}
	}
	encode := func(value interface{}) []byte {
		t.Helper()
		content, err := json.Marshal(value)
		if err != nil {
			t.Fatal(err)
		}
		return content
	}
	helperArgs := []string{"-test.run=^TestCLICommandProcess$", "--", "--borrowedtime-command"}
	literal := "spaces ; & $HOME $(touch PWNED)"
	files := map[string][]byte{
		"config.json": encode(map[string]string{
			"workspace": workspace, "editor": os.Args[0], "projectstructure": "default", "customer": "Example Client",
		}),
		"templates/file/notes.md":      []byte("# {{ .ProjectName }}\nCustomer: {{ index .Config \"customer\" }}\n"),
		"templates/file/alternate.md":  []byte("Alternate: {{ .ProjectName }}\n"),
		"templates/file/broken.md":     []byte("{{"),
		"templates/file/settings.json": []byte(`{"project":"{{ .ProjectName }}","customer":"{{ index .Config "customer" }}"}`),
		"data/commands.json": encode(map[string][]command.Command{
			"setup":         {{Executable: os.Args[0], Args: append(append([]string{}, helperArgs...), "write", "hook.txt", "created")}},
			"write-literal": {{Executable: os.Args[0], Args: append(append([]string{}, helperArgs...), "write", "literal args.txt", literal)}},
			"fail": {
				{Executable: os.Args[0], Args: append(append([]string{}, helperArgs...), "fail")},
				{Executable: os.Args[0], Args: append(append([]string{}, helperArgs...), "write", "must-not-run.txt", "bad")},
			},
		}),
	}
	structure := project.Node{
		FullPath: "{{ .Workspace }}/{{ .ProjectName }}",
		Info:     &project.FileInfo{IsDir: true},
		Children: []*project.Node{
			{FullPath: "notes.md", Info: &project.FileInfo{Template: "notes"}},
			{FullPath: ".config.json", Info: &project.FileInfo{Template: "settings"}},
			{FullPath: "evidence", Info: &project.FileInfo{IsDir: true}, Children: []*project.Node{
				{FullPath: "README.txt", Info: &project.FileInfo{}},
			}},
		},
		PostCreate: []string{"setup"},
	}
	files["templates/project/default.json"] = encode(structure)
	structure.Children[0].Info.Template = "alternate"
	structure.PostCreate = nil
	files["templates/project/template with spaces.json"] = encode(structure)
	structure.Children[0].Info.Template = "broken"
	files["templates/project/broken.json"] = encode(structure)
	for relative, content := range files {
		path := filepath.Join(configDir, relative)
		if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, content, 0600); err != nil {
			t.Fatal(err)
		}
	}
	projectDir := filepath.Join(workspace, "client review")
	notes := filepath.Join(projectDir, "notes.md")
	checkFile := func(t *testing.T, path, want string) {
		t.Helper()
		content, err := os.ReadFile(path)
		if err != nil || string(content) != want {
			t.Fatalf("file %s = %q, error = %v, want %q", path, content, err, want)
		}
	}
	checkMissing := func(t *testing.T, path string) {
		t.Helper()
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			t.Fatalf("unexpected file %s: %v", path, err)
		}
	}
	for _, test := range []struct {
		name   string
		args   []string
		code   int
		stdout string
		stderr string
		before func(*testing.T)
		after  func(*testing.T)
	}{
		{name: "help", args: []string{"-h"}, stderr: "Usage: borrowedtime"},
		{name: "empty-list", args: []string{"project", "list"}, stdout: "\n"},
		{name: "create-spaced-name", args: []string{"project", "create", "client review"}, after: func(t *testing.T) {
			checkFile(t, notes, "# client review\nCustomer: Example Client\n")
			checkFile(t, filepath.Join(projectDir, "hook.txt"), "created")
			checkFile(t, filepath.Join(projectDir, "evidence", "README.txt"), "")
			content, err := os.ReadFile(filepath.Join(projectDir, ".config.json"))
			if err != nil {
				t.Fatal(err)
			}
			var settings map[string]string
			if err := json.Unmarshal(content, &settings); err != nil {
				t.Fatal(err)
			}
			if settings["project"] != "client review" || settings["customer"] != "Example Client" {
				t.Fatalf("generated project config = %v", settings)
			}
		}},
		{name: "reject-existing", args: []string{"project", "create", "client review"}, code: 1, stderr: "already exists", before: func(t *testing.T) {
			if err := os.WriteFile(notes, []byte("manual notes"), 0600); err != nil {
				t.Fatal(err)
			}
		}, after: func(t *testing.T) { checkFile(t, notes, "manual notes") }},
		{name: "overwrite-flags-first", args: []string{"project", "create", "-overwrite", "-template", "template with spaces", "client review"}, after: func(t *testing.T) {
			checkFile(t, notes, "Alternate: client review\n")
		}},
		{name: "broken-preserves-file", args: []string{"project", "create", "client review", "-template", "broken", "-overwrite"}, code: 1, stderr: "create new template", after: func(t *testing.T) {
			checkFile(t, notes, "Alternate: client review\n")
		}},
		{name: "template-after-name", args: []string{"project", "create", "second project", "-template", "template with spaces"}, after: func(t *testing.T) {
			checkFile(t, filepath.Join(workspace, "second project", "notes.md"), "Alternate: second project\n")
		}},
		{name: "dash-name", args: []string{"project", "create", "-template", "template with spaces", "--", "-dash-project"}, after: func(t *testing.T) {
			checkFile(t, filepath.Join(workspace, "-dash-project", "notes.md"), "Alternate: -dash-project\n")
		}},
		{name: "option-name", args: []string{"project", "create", "open"}},
		{name: "open-option-name", args: []string{"project", "open", "open"}},
		{name: "open-spaced-name", args: []string{"project", "open", "client review"}},
		{name: "list", args: []string{"project", "list"}, stdout: "client review"},
		{name: "list-batches", args: []string{"command", "list"}, stdout: "fail\nsetup\nwrite-literal\n"},
		{name: "literal-arguments", args: []string{"command", "run", "write-literal", "-project", "client review"}, after: func(t *testing.T) {
			checkFile(t, filepath.Join(projectDir, "literal args.txt"), literal)
			checkMissing(t, filepath.Join(projectDir, "PWNED"))
		}},
		{name: "startup-cwd", args: []string{"command", "run", "write-literal"}, after: func(t *testing.T) {
			checkFile(t, filepath.Join(root, "literal args.txt"), literal)
		}},
		{name: "stop-on-failure", args: []string{"command", "run", "fail", "-project", "client review"}, code: 1, stderr: "step 1", after: func(t *testing.T) {
			checkMissing(t, filepath.Join(projectDir, "must-not-run.txt"))
		}},
		{name: "debug", args: []string{"-debug", "command", "run", "write-literal", "-project", "client review"}, stderr: "run command"},
		{name: "backup-spaced-name", args: []string{"config", "backup", "-file", "backup with spaces"}, after: func(t *testing.T) {
			info, err := os.Stat(filepath.Join(configDir, "backups", "backup with spaces.zip"))
			if err != nil || info.Size() == 0 {
				t.Fatalf("backup missing or empty: %v", err)
			}
			if err := os.WriteFile(notes, []byte("manual notes after backup"), 0600); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(configDir, "templates", "file", "notes.md"), []byte("changed"), 0600); err != nil {
				t.Fatal(err)
			}
		}},
		{name: "restore-spaced-name", args: []string{"config", "restore", "backup with spaces.zip"}, after: func(t *testing.T) {
			checkFile(t, filepath.Join(configDir, "templates", "file", "notes.md"), "# {{ .ProjectName }}\nCustomer: {{ index .Config \"customer\" }}\n")
			checkFile(t, notes, "manual notes after backup")
		}},
		{name: "unknown", args: []string{"unknown"}, code: 1, stderr: "unknown command"},
		{name: "incomplete", args: []string{"command"}, code: 1, stderr: "incomplete command"},
		{name: "missing-name", args: []string{"project", "create"}, code: 1, stderr: "provide project name"},
		{name: "missing-value", args: []string{"project", "create", "new", "-template"}, code: 1, stderr: "requires a value"},
		{name: "unknown-option", args: []string{"project", "list", "-unknown"}, code: 1, stderr: "unknown option"},
		{name: "duplicate-option", args: []string{"project", "create", "new", "-overwrite", "-overwrite"}, code: 1, stderr: "more than once"},
		{name: "extra-positional", args: []string{"project", "create", "one", "two"}, code: 1, stderr: "unexpected argument"},
		{name: "unknown-batch", args: []string{"command", "run", "missing"}, code: 1, stderr: "not found"},
		{name: "global-option-error", args: []string{"-unknown"}, code: 2, stderr: "flag provided but not defined"},
		{name: "open-missing-value", args: []string{"project", "open"}, code: 1, stderr: "requires a value"},
	} {
		if !t.Run(test.name, func(t *testing.T) {
			if test.before != nil {
				test.before(t)
			}
			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer cancel()
			process := exec.CommandContext(ctx, binary, test.args...)
			process.Dir = root
			for _, entry := range os.Environ() {
				key := strings.ToUpper(strings.SplitN(entry, "=", 2)[0])
				if key != "HOME" && key != "USERPROFILE" && key != "BORROWEDTIME_TEST_EDITOR" {
					process.Env = append(process.Env, entry)
				}
			}
			process.Env = append(process.Env, "HOME="+home, "USERPROFILE="+home, "BORROWEDTIME_TEST_EDITOR=1")
			var stdout, stderr bytes.Buffer
			process.Stdout = &stdout
			process.Stderr = &stderr
			err := process.Run()
			if ctx.Err() != nil {
				t.Fatalf("CLI did not exit: %v", ctx.Err())
			}
			code := 0
			if err != nil {
				exitError, ok := err.(*exec.ExitError)
				if !ok {
					t.Fatal(err)
				}
				code = exitError.ExitCode()
			}
			if code != test.code {
				t.Fatalf("exit code = %d, want %d; stdout=%q stderr=%q", code, test.code, stdout.String(), stderr.String())
			}
			if !strings.Contains(stdout.String(), test.stdout) || !strings.Contains(stderr.String(), test.stderr) {
				t.Fatalf("stdout=%q stderr=%q, want %q and %q", stdout.String(), stderr.String(), test.stdout, test.stderr)
			}
			if test.stdout == "" && stdout.Len() != 0 {
				t.Fatalf("unexpected stdout: %q", stdout.String())
			}
			if test.stderr == "" && stderr.Len() != 0 {
				t.Fatalf("unexpected stderr: %q", stderr.String())
			}
			if strings.Contains(stdout.String(), ">>>") {
				t.Fatal("one-shot command started interactive prompt")
			}
			if strings.Contains(stderr.String(), literal) {
				t.Fatal("diagnostics exposed command arguments")
			}
			if test.after != nil {
				test.after(t)
			}
		}) {
			t.FailNow()
		}
	}
}
