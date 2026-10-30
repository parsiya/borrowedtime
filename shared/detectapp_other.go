//go:build !windows
// +build !windows

package shared

// DetectApp leaves Windows Start Menu detection disabled on other platforms.
// Example: DetectApp("code.exe") yields an empty string.
func DetectApp(executable string) string {
	return ""
}
