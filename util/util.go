package util

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

var baseURL string

func init() {
	baseURL = os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "/"
	}
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
