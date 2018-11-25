package shared

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// JSON utils.

// StructToJSONString converts a struct into a JSON string.
func StructToJSONString(v interface{}, indent bool) (string, error) {
	var sb strings.Builder
	enc := json.NewEncoder(&sb)
	if indent {
		enc.SetIndent("", "    ")
	}
	// Do not escape < > & in HTML.
	enc.SetEscapeHTML(false)
	err := enc.Encode(v)
	return sb.String(), err
}

// PrettyPrintJSON gets a string and pretty prints it to io.Writer.
func PrettyPrintJSON(in string, w io.Writer) error {
	var content bytes.Buffer
	if err := json.Indent(&content, []byte(in), "", "\t"); err != nil {
		return fmt.Errorf("shared.PrettyPrintJSON: %s", err.Error())
	}
	_, err := w.Write(content.Bytes())
	if err != nil {
		return fmt.Errorf("shared.PrettyPrintJSON: %s", err.Error())
	}
	return nil
}
