package util

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

var baseUrl string

func init() {
	baseUrl = os.Getenv("BASE_URL")
	if baseUrl == "" {
		baseUrl = "/"
	}
}

func ReplaceBase(str, old, new string) string {
	return strings.Replace(str, old, fmt.Sprintf(new, baseUrl), 1)
}

func AddBase(path string) string {
	return fmt.Sprintf("%s%s", baseUrl, path)
}

func TrimBase(r *http.Request, prefix string) string {
	return strings.TrimPrefix(r.URL.Path, AddBase(prefix))
}
