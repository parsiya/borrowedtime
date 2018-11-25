package shared

import (
	"testing"
)

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
