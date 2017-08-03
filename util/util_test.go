package util

import (
	"testing"
	"net/http"
)

func TestAddBase(t *testing.T) {
	expectedStr := "/mypath/"

	result := AddBase("mypath/")

	if result != expectedStr {
		t.Fatalf("Expected %s, got %s", expectedStr, result)
	}
}

func TestReplaceBase(t *testing.T) {
	expectedStr := "<a href=\"/\">root</a>"
	str := "<a href=\"/fakebase/\">root</a>"

	result := ReplaceBase(str, "href=\"/fakebase/\"", "href=\"%s\"")

	if result != expectedStr {
		t.Fatalf("Expected %s, got %s", expectedStr, result)
	}
}


func TestTrimBase(t *testing.T) {
	expectedStr := ""

	req, _ := http.NewRequest("GET", "/api/state", nil)
	result := TrimBase(req, "api/state")

	if result != expectedStr {
		t.Fatalf("Expected %s, got %s", expectedStr, result)
	}
}
