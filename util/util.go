package util

import (
	"fmt"
	"net/http"
	"strings"
)

var baseURL string

// UpdateBase replaces baseURL with a new one
func UpdateBase(new string) {
	baseURL = new
}

// ReplaceBase replaces a pattern in a string, injecting baseURL into it
func ReplaceBase(str, old, new string) string {
	return strings.Replace(str, old, fmt.Sprintf(new, baseURL), 1)
}

// AddBase preprends baseURL to a string
func AddBase(path string) string {
	return fmt.Sprintf("%s%s", baseURL, path)
}

// TrimBase removes baseURL from the beginning of a string
func TrimBase(r *http.Request, prefix string) string {
	return strings.TrimPrefix(r.URL.Path, AddBase(prefix))
}
