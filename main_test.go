package main

import (
	"net/http"
	"testing"

	"github.com/camptocamp/terraboard/db"
	"github.com/jinzhu/gorm"
)

func TestIsKnownStateVersion_Match(t *testing.T) {

	statesVersions := map[string][]string{
		"fakeVersionID": []string{"myfakepath/terraform.tfstate"},
	}

	if !isKnownStateVersion(statesVersions, "fakeVersionID", "myfakepath/terraform.tfstate") {
		t.Fatalf("Expected %t, got %t", true, false)
	}
}

func TestIsKnownStateVersion_NoMatch(t *testing.T) {

	statesVersions := map[string][]string{
		"fakeVersionID": []string{"myfakepath/terraform.tfstate"},
	}

	if isKnownStateVersion(statesVersions, "VersionID", "myfakepath/terraform.tfstate") {
		t.Fatalf("Expected %t, got %t", false, true)
	}
}

func handlerWithDB(w http.ResponseWriter, r *http.Request, d *db.Database) {
}

func TestHandleWithDB(t *testing.T) {
	d := db.Database{DB: &gorm.DB{}}
	handleWithDB(handlerWithDB, &d)
}
