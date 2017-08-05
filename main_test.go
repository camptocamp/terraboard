package main

import (
	"net/http"
	"testing"

	"github.com/camptocamp/terraboard/db"
	"github.com/jinzhu/gorm"
)

func TestIsKnownStateVersion(t *testing.T) {

	statesVersions := map[string][]string{
		"fakeVersionID": []string{"myfakepath/terraform.tfstate"},
	}

	if !isKnownStateVersion(statesVersions, "fakeVersionID", "myfakepath/terraform.tfstate") {
		t.Fatalf("Expected %s, got %s", true, false)
	}
}

func handlerWithDB(w http.ResponseWriter, r *http.Request, d *db.Database) {
}

func TestHandleWithDB(t *testing.T) {
	d := db.Database{DB: &gorm.DB{}}
	handleWithDB(handlerWithDB, &d)
}
