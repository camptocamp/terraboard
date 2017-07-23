package db

import (
	"database/sql"
	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/hashicorp/terraform/terraform"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func Init() {
	var err error
	db, err = sql.Open("sqlite3", "./db/terraboard.db")
	if err != nil {
		log.Fatal(err)
	}
	//defer db.Close()

	log.Infof("New db is %v", db)

	sqlStmt := `
	CREATE TABLE IF NOT EXISTS states (id INTEGER PRIMARY KEY AUTOINCREMENT, version_id TEXT, path TEXT, tf_version TEXT, serial INTEGER);
	CREATE TABLE IF NOT EXISTS modules (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT, UNIQUE(name));
	CREATE TABLE IF NOT EXISTS resources (id INTEGER PRIMARY KEY AUTOINCREMENT, state_id INTEGER, module_id INTEGER, type TEXT, name TEXT);
	CREATE TABLE IF NOT EXISTS attributes (id INTEGER PRIMARY KEY AUTOINCREMENT, resource_id INTEGER, name TEXT, value TEXT);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}
}

func InsertState(versionId string, path string, state *terraform.State) error {
	log.Info("Inserting new state")
	stmt, err := db.Prepare("INSERT INTO states(version_id, path, tf_version, serial) VALUES(?, ?, ?, ?)")
	if err != nil {
		return err
	}

	res, err := stmt.Exec(versionId, path, state.TFVersion, state.Serial)
	if err != nil {
		return err
	}
	sId, _ := res.LastInsertId()

	for _, m := range state.Modules {
		stmt, err = db.Prepare("INSERT OR IGNORE INTO modules(name) VALUES(?)")
		if err != nil {
			return err
		}
		mName := strings.Join(m.Path, "/")
		_, err := stmt.Exec(mName)
		if err != nil {
			return err
		}

		rows, err := db.Query(fmt.Sprintf("SELECT id FROM modules WHERE name = '%s'", mName))
		if err != nil {
			return err
		}
		var mId int
		rows.Next()
		err = rows.Scan(&mId)
		if err != nil {
			return err
		}
		rows.Close()

		for n, r := range m.Resources {
			stmt, err = db.Prepare("INSERT INTO resources(state_id, module_id, type, name) VALUES(?, ?, ?, ?)")
			if err != nil {
				return err
			}
			res, err := stmt.Exec(sId, mId, r.Type, n)
			if err != nil {
				return err
			}

			rId, _ := res.LastInsertId()

			for k, v := range r.Primary.Attributes {
				stmt, err = db.Prepare("INSERT INTO attributes(resource_id, name, value) VALUES(?, ?, ?)")
				if err != nil {
					return err
				}
				_, err := stmt.Exec(rId, k, v)
				if err != nil {
					return err
				}
			}
		}
	}

	log.Infof("res: %v", res)

	return nil
}

func GetState(path, versionId string) (state *terraform.State, err error) {
	rows, err := db.Query(fmt.Sprintf("SELECT (id, tf_version, serial) FROM states WHERE path = '%s' AND version_id = '%s'", path, versionId))
	if err != nil {
		return state, err
	}
	var mId int
	var tfVersion string
	var serial int64
	rows.Next()
	err = rows.Scan(&mId, &tfVersion, &serial)
	if err != nil {
		return state, err
	}
	rows.Close()

	return &terraform.State{
		TFVersion: tfVersion,
		Serial:    serial,
	}, nil
}
