package shared

import (
	"testing"
)

func TestWindowsifyString(t *testing.T) {
	type args struct {
		inp string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"end-of-line", args{"input\n"}, "input\r\n"},
		{"double-end-of-line", args{"input\n\n"}, "input\r\n\r\n"},
		{"middle", args{"input\n\nnextline"}, "input\r\n\r\nnextline"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WindowsifyString(tt.args.inp); got != tt.want {
				t.Errorf("WindowsifyString() = %v, want %v", got, tt.want)
			}
		})
	}
}
