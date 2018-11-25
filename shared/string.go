package shared

import "strings"

// String utils.

// WindowsifyString converts "\n" to "\r\n" in a string.
func WindowsifyString(inp string) string {
	return strings.Replace(inp, "\n", "\r\n", -1)
}

// EscapeString, escapes a single \ by converting it to "\ \".
func EscapeString(inp string) string {
	return strings.Replace(inp, "\\", "\\\\", -1)
}
