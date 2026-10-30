package project

import (
	"bytes"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"testing"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/parsiya/borrowedtime/command"
	"github.com/parsiya/borrowedtime/shared"
)

func TestProjectCommandProcess(t *testing.T) {
	for index, arg := range os.Args {
		if arg == "--borrowedtime-project-command" {
			args := os.Args[index+1:]
			if len(args) != 2 {
				os.Exit(2)
			}
			if err := os.WriteFile(args[0], []byte(args[1]), 0600); err != nil {
				os.Exit(2)
			}
			os.Exit(0)
		}
	}
}

func TestNodeCreate(t *testing.T) {
	root := t.TempDir()
	for _, key := range []string{"HOME", "USERPROFILE"} {
		key := key
		old, existed := os.LookupEnv(key)
		if err := os.Setenv(key, root); err != nil {
			t.Fatal(err)
		}
		defer func() {
			if existed {
				os.Setenv(key, old)
			} else {
				os.Unsetenv(key)
			}
		}()
	}
	oldDisableCache := homedir.DisableCache
	homedir.DisableCache = true
	defer func() { homedir.DisableCache = oldDisableCache }()
	templateDir := filepath.Join(root, "borrowedtime", "templates", "file")
	if err := os.MkdirAll(templateDir, 0700); err != nil {
		t.Fatal(err)
	}
	for name, content := range map[string]string{
		"valid.md":   "# {{ .ProjectName }}",
		"invalid.md": "{{",
		"failure.md": "{{ .MissingField }}",
	} {
		if err := os.WriteFile(filepath.Join(templateDir, name), []byte(content), 0600); err != nil {
			t.Fatal(err)
		}
	}
	for _, test := range []struct {
		name      string
		template  string
		exists    bool
		overwrite bool
		wantError bool
		want      string
	}{
		{"render", "valid", false, false, false, "# example"},
		{"blank", "", false, false, false, ""},
		{"missing", "missing", false, false, true, ""},
		{"missing-preserves", "missing", true, true, true, "original"},
		{"invalid-preserves", "invalid", true, true, true, "original"},
		{"failure-preserves", "failure", true, true, true, "original"},
		{"no-overwrite", "valid", true, false, true, "original"},
		{"overwrite", "valid", true, true, false, "# example"},
	} {
		t.Run(test.name, func(t *testing.T) {
			path := filepath.Join(root, test.name+".md")
			if test.exists {
				if err := os.WriteFile(path, []byte("original"), 0600); err != nil {
					t.Fatal(err)
				}
			}
			node := Node{FullPath: path, Info: &FileInfo{Template: test.template}}
			err := node.Create(Project{ProjectName: "example"}, test.overwrite)
			if (err != nil) != test.wantError {
				t.Fatalf("Create() error = %v, wantError %v", err, test.wantError)
			}
			content, readErr := os.ReadFile(path)
			if test.wantError && !test.exists {
				if !os.IsNotExist(readErr) {
					t.Fatalf("failed render created file: %v", readErr)
				}
				return
			}
			if readErr != nil || string(content) != test.want {
				t.Fatalf("content = %q, error = %v, want %q", content, readErr, test.want)
			}
		})
	}
	if err := (&Node{}).Create(Project{}, false); err == nil {
		t.Fatal("expected missing node info error")
	}
	t.Run("debug-logging", func(t *testing.T) {
		var output bytes.Buffer
		oldOutput := log.Writer()
		oldDebug := shared.Debug
		defer func() {
			log.SetOutput(oldOutput)
			shared.Debug = oldDebug
		}()
		log.SetOutput(&output)
		for _, enabled := range []bool{false, true} {
			output.Reset()
			shared.Debug = enabled
			node := Node{FullPath: filepath.Join(root, "debug.md"), Info: &FileInfo{}}
			if err := node.Create(Project{}, true); err != nil {
				t.Fatal(err)
			}
			if (output.Len() > 0) != enabled {
				t.Fatalf("debug=%t, output=%q", enabled, output.String())
			}
		}
	})
	t.Run("project-create", func(t *testing.T) {
		workspace := filepath.Join(root, "workspace")
		cfg, err := json.Marshal(map[string]string{"workspace": workspace, "custom": "client"})
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(root, "borrowedtime", "config.json"), cfg, 0600); err != nil {
			t.Fatal(err)
		}
		projectTemplateDir := filepath.Join(root, "borrowedtime", "templates", "project")
		if err := os.MkdirAll(projectTemplateDir, 0700); err != nil {
			t.Fatal(err)
		}
		structure := `{"path":"{{ .Workspace }}/{{ .ProjectName }}","info":{"isdir":true},"children":[{"path":"notes.md","info":{"template":"valid"}},{"path":".config.json","info":{"template":"settings"}}]}`
		dataDir := filepath.Join(root, "borrowedtime", "data")
		if err := os.MkdirAll(dataDir, 0700); err != nil {
			t.Fatal(err)
		}
		batches, err := json.Marshal(map[string][]command.Command{
			"setup":   {{Executable: os.Args[0], Args: []string{"-test.run=^TestProjectCommandProcess$", "--", "--borrowedtime-project-command", "hook-ran.txt", "created"}}},
			"failure": {{Executable: filepath.Join(root, "missing-program")}},
		})
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dataDir, "commands.json"), batches, 0600); err != nil {
			t.Fatal(err)
		}
		for _, test := range []struct {
			name       string
			settings   string
			wantError  bool
			hooks      []string
			wantMarker bool
		}{
			{"complete", `{"client":"{{ index .Config "custom" }}"}`, false, nil, false},
			{"bad-config", `{`, true, nil, false},
			{"with-hook", `{"client":"client"}`, false, []string{"setup"}, true},
			{"without-config", "", false, []string{"setup"}, true},
			{"hook-bad-config", `{`, true, []string{"setup"}, false},
			{"missing-hook", `{}`, true, []string{"missing"}, false},
			{"missing-later-hook", `{}`, true, []string{"setup", "missing"}, false},
			{"failed-hook", `{}`, true, []string{"failure", "setup"}, false},
		} {
			t.Run(test.name, func(t *testing.T) {
				var tree Node
				if err := json.Unmarshal([]byte(structure), &tree); err != nil {
					t.Fatal(err)
				}
				tree.PostCreate = test.hooks
				if test.settings == "" {
					tree.Children = tree.Children[:1]
				}
				templateBytes, err := json.Marshal(tree)
				if err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(filepath.Join(projectTemplateDir, "sample.json"), templateBytes, 0600); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(filepath.Join(templateDir, "settings.json"), []byte(test.settings), 0600); err != nil {
					t.Fatal(err)
				}
				project := New(test.name)
				err = project.Create("sample", false)
				if (err != nil) != test.wantError {
					t.Fatalf("Create() error = %v, wantError %v", err, test.wantError)
				}
				_, markerErr := os.Stat(filepath.Join(project.ProjectRoot, "hook-ran.txt"))
				if test.wantMarker && markerErr != nil {
					t.Fatalf("hook did not run: %v", markerErr)
				}
				if !test.wantMarker && !os.IsNotExist(markerErr) {
					t.Fatalf("unexpected hook marker: %v", markerErr)
				}
				if test.wantError {
					return
				}
				if test.settings != "" && project.ProjectConfig["client"] != "client" {
					t.Fatalf("project config = %v", project.ProjectConfig)
				}
				content, err := os.ReadFile(filepath.Join(project.ProjectRoot, "notes.md"))
				if err != nil || string(content) != "# "+test.name {
					t.Fatalf("notes = %q, error = %v", content, err)
				}
				if err := project.Create("sample", false); err == nil {
					t.Fatal("expected existing project error")
				}
			})
		}
	})
}
