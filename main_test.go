package main

import (
	"testing"
)

func TestIsKnownStateVersion(t *testing.T) {

	statesVersions := map[string][]string{
		"fakeVersionID": []string{"myfakepath/terraform.tfstate"},
	}

	if !isKnownStateVersion(statesVersions, "fakeVersionID", "myfakepath/terraform.tfstate") {
		t.Fatalf("Expected %s, got %s", true, false)
	}
}
