package shared

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestListFiles(t *testing.T) {
	root := t.TempDir()
	if err := os.Mkdir(filepath.Join(root, "nested"), 0700); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"notes.md", "config.json", filepath.Join("nested", "todo.md")} {
		if err := os.WriteFile(filepath.Join(root, name), nil, 0600); err != nil {
			t.Fatal(err)
		}
	}
	for _, test := range []struct {
		pattern string
		want    []string
	}{
		{"*", []string{"config.json", filepath.Join("nested", "todo.md"), "notes.md"}},
		{"*.md", []string{filepath.Join("nested", "todo.md"), "notes.md"}},
	} {
		t.Run(test.pattern, func(t *testing.T) {
			got, err := ListFiles(root, test.pattern)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Fatalf("ListFiles() = %v, want %v", got, test.want)
			}
		})
	}
	if _, err := ListFiles(root, "["); err == nil {
		t.Fatal("expected invalid pattern error")
	}
	if _, err := ListFiles(filepath.Join(root, "missing"), "*"); err == nil {
		t.Fatal("expected missing directory error")
	}
}

func TestRemoveExtension(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"normal-file", args{"filename.ext"}, "filename"},
		{"long-extension", args{"filename.ext213123123123123"}, "filename"},
		{"no-extension", args{"filename"}, "filename"},
		{"empty-filename", args{".ext"}, ""},
		{"empty-string", args{""}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RemoveExtension(tt.args.str); got != tt.want {
				t.Errorf("RemoveExtension() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAddExtension(t *testing.T) {
	type args struct {
		str string
		ext string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"normal", args{"filename", "ext"}, "filename.ext"},
		{"long-filename", args{"filename33333333333", "ext"}, "filename33333333333.ext"},
		{"long-extension", args{"filename", "ext1111111111111111"}, "filename.ext1111111111111111"},
		{"no-extension", args{"filename", ""}, "filename"},
		{"no-filename", args{"", "ext"}, ".ext"},
		{"empty-string", args{"", ""}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AddExtension(tt.args.str, tt.args.ext); got != tt.want {
				t.Errorf("AddExtension() = %v, want %v", got, tt.want)
			}
		})
	}
}
