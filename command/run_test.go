package command

import (
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestCommandProcess(t *testing.T) {
	if os.Getenv("BORROWEDTIME_TEST_PROCESS") != "1" {
		return
	}
	for index, arg := range os.Args {
		if arg != "--" {
			continue
		}
		args := os.Args[index+1:]
		if len(args) == 1 && args[0] == "fail" {
			os.Exit(7)
		}
		if len(args) == 3 && args[0] == "write" {
			if err := ioutil.WriteFile(args[1], []byte(args[2]), 0600); err != nil {
				os.Exit(2)
			}
			os.Exit(0)
		}
		os.Exit(2)
	}
	os.Exit(2)
}

func TestRun(t *testing.T) {
	root, err := ioutil.TempDir("", "borrowedtime-command-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(root)
	old, existed := os.LookupEnv("BORROWEDTIME_TEST_PROCESS")
	if err := os.Setenv("BORROWEDTIME_TEST_PROCESS", "1"); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if existed {
			os.Setenv("BORROWEDTIME_TEST_PROCESS", old)
		} else {
			os.Unsetenv("BORROWEDTIME_TEST_PROCESS")
		}
	}()
	baseArgs := []string{"-test.run=^TestCommandProcess$", "--"}
	literal := "spaces ; & $HOME $(echo unsafe)"
	write := Command{Executable: os.Args[0], Args: append(append([]string{}, baseArgs...), "write", "result.txt", literal)}
	if err := Run([]Command{write}, root); err != nil {
		t.Fatal(err)
	}
	content, err := ioutil.ReadFile(filepath.Join(root, "result.txt"))
	if err != nil || string(content) != literal {
		t.Fatalf("content = %q, error = %v", content, err)
	}
	second := Command{Executable: os.Args[0], Args: append(append([]string{}, baseArgs...), "write", "result.txt", "second")}
	if err := Run([]Command{write, second}, root); err != nil {
		t.Fatal(err)
	}
	content, err = ioutil.ReadFile(filepath.Join(root, "result.txt"))
	if err != nil || string(content) != "second" {
		t.Fatalf("commands ran out of order: content = %q, error = %v", content, err)
	}
	if err := os.Remove(filepath.Join(root, "result.txt")); err != nil {
		t.Fatal(err)
	}
	failure := Command{Executable: os.Args[0], Args: append(append([]string{}, baseArgs...), "fail")}
	err = Run([]Command{failure, write}, root)
	var exitError *exec.ExitError
	if !errors.As(err, &exitError) || exitError.ExitCode() != 7 {
		t.Fatalf("expected exit code 7, got %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, "result.txt")); !os.IsNotExist(err) {
		t.Fatalf("later command ran after failure: %v", err)
	}
	if err := Run([]Command{write, {}}, root); err == nil {
		t.Fatal("expected missing executable error")
	}
	if _, err := os.Stat(filepath.Join(root, "result.txt")); !os.IsNotExist(err) {
		t.Fatalf("commands ran before validation: %v", err)
	}
	if err := Run([]Command{{Executable: filepath.Join(root, "missing-executable")}}, root); err == nil {
		t.Fatal("expected missing program error")
	}
}
