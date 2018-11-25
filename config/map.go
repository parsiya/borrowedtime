package config

import (
	"strings"
)

// ConfigMap represents the contents of the config file.
type ConfigMap map[string]string

// Set assigns the value to the key in the config. It overwrites any previous
// values so if needed, check with Has first.
func (v ConfigMap) Set(key, value string) {
	key = strings.ToLower(key)
	v[key] = value
}

// Key returns the value of a key or "" if it does not exist in the config.
func (v ConfigMap) Key(key string) string {
	key = strings.ToLower(key)
	return v[key]
}

// Has returns true if a key exists in the config.
func (v ConfigMap) Has(key string) bool {
	key = strings.ToLower(key)
	_, exists := v[key]
	return exists
}
