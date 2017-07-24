package db

import (
	"database/sql"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/hashicorp/terraform/terraform"

	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

var db *gorm.DB

type State struct {
	gorm.Model `json:"-"`
	Path       string     `json:"path"`
	VersionId  string     `json:"version_id"`
	TFVersion  string     `json:"terraform_version"`
	Serial     int64      `json:"serial"`
	Resources  []Resource `json:"resources"`
}

type Resource struct {
	ID         int           `sql:"AUTO_INCREMENT" gorm:"primary_key" json:"-"`
	StateID    sql.NullInt64 `json:"-"`
	Type       string        `json:"type"`
	Name       string        `json:"name"`
	Attributes []Attribute   `json:"attributes"`
	Module     string        `json:"module"`
}

type Attribute struct {
	ID         int           `sql:"AUTO_INCREMENT" gorm:"primary_key" json:"-"`
	ResourceID sql.NullInt64 `json:"-"`
	Key        string        `json:"key"`
	Value      string        `json:"value"`
}

func Init() {
	var err error
	db, err = gorm.Open("sqlite3", "./db/terraboard.db")
	if err != nil {
		log.Fatal(err)
	}
	//defer db.Close()

	log.Infof("Automigrate")

	db.AutoMigrate(&State{}, &Resource{}, &Attribute{})

	db.LogMode(true)

	log.Infof("New db is %v", db)
}

func InsertState(versionId string, path string, state *terraform.State) error {
	log.Info("Inserting new state")

	st := &State{
		Path:      path,
		VersionId: versionId,
		TFVersion: state.TFVersion,
		Serial:    state.Serial,
	}

	for _, m := range state.Modules {
		mod := strings.Join(m.Path, "/")
		for n, r := range m.Resources {
			res := Resource{
				Type:   r.Type,
				Name:   n,
				Module: mod,
			}

			for k, v := range r.Primary.Attributes {
				res.Attributes = append(res.Attributes, Attribute{
					Key:   k,
					Value: v,
				})
			}

			st.Resources = append(st.Resources, res)
		}
	}

	db.Create(st)
	return nil
}

func GetState(path, versionId string) (state State) {
	log.Infof("Getting state from DB for %s/%s", path, versionId)
	db.Preload("Resources").Preload("Resources.Attributes").Find(&state, "path = ? AND version_id = ?", path, versionId)
	log.Infof("state=%v", state)
	return
}
