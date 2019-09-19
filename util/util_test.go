package util

import (
	"net/http"
	"testing"
)

func TestSetBasePath(t *testing.T) {
	newBase := "/test/"
	expectedStr := newBase

	SetBasePath(newBase)

	result := GetFullPath("")

	if result != expectedStr {
		t.Fatalf("Expected %s, got %s", expectedStr, result)
	}
}

func TestGetFullPath(t *testing.T) {
	expectedStr := "/mypath/"

	SetBasePath("/")
	result := GetFullPath("mypath/")

	if result != expectedStr {
		t.Fatalf("Expected %s, got %s", expectedStr, result)
	}
}

func TestReplaceBasePath(t *testing.T) {
	expectedStr := "<a href=\"/\">root</a>"
	str := "<a href=\"/fakebase/\">root</a>"

	SetBasePath("/")
	result := ReplaceBasePath(str, "href=\"/fakebase/\"", "href=\"%s\"")

	if result != expectedStr {
		t.Fatalf("Expected %s, got %s", expectedStr, result)
	}
}

func TestTrimBasePath(t *testing.T) {
	expectedStr := ""

	SetBasePath("/")
	req, _ := http.NewRequest("GET", "/api/state", nil)
	result := TrimBasePath(req, "api/state")

	if result != expectedStr {
		t.Fatalf("Expected %s, got %s", expectedStr, result)
	}
}
