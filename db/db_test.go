package db

import (
	"database/sql"
	"net/url"
	"reflect"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/hashicorp/go-version"
	"github.com/stretchr/testify/assert"
	"github.com/zclconf/go-cty/cty"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/camptocamp/terraboard/internal/terraform/addrs"
	"github.com/camptocamp/terraboard/internal/terraform/states"
	"github.com/camptocamp/terraboard/internal/terraform/states/statefile"
	"github.com/camptocamp/terraboard/state"
	"github.com/camptocamp/terraboard/types"
)

func TestGetResourceIndex(t *testing.T) {
	tests := []struct {
		name string
		args addrs.InstanceKey
		want string
	}{
		{
			"StringKey",
			addrs.StringKey("module.bar"),
			"[\"module.bar\"]",
		},
		{
			"NoKey",
			addrs.NoKey,
			"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getResourceIndex(tt.args)
			if got != tt.want {
				t.Errorf(
					"TestGetResourceIndex() -> \n\ngot:\n%v,\n\nwant:\n%v",
					got,
					tt.want,
				)
			}
		})
	}
}

func TestMarshalAttributeValues(t *testing.T) {
	tests := []struct {
		name string
		args *states.ResourceInstanceObjectSrc
		want []types.Attribute
	}{
		{
			"Nil src",
			nil,
			nil,
		},
		{
			"Empty AttrsFlat",
			&states.ResourceInstanceObjectSrc{
				AttrsJSON: []byte(`{"ami":"bar"}`),
				Status:    states.ObjectReady,
			},
			[]types.Attribute{
				{
					Key:   "ami",
					Value: "\"bar\"",
				},
			},
		},
		{
			"Empty AttrsFlat with bad AttrsJSON JSON format",
			&states.ResourceInstanceObjectSrc{
				AttrsJSON: []byte(`"bar"`),
				Status:    states.ObjectReady,
			},
			nil,
		},
		{
			"With valid AttrsFlat",
			&states.ResourceInstanceObjectSrc{
				AttrsJSON: []byte(`{"ami":"bar"}`),
				AttrsFlat: map[string]string{
					"ami": "bar",
				},
				Status: states.ObjectReady,
			},
			[]types.Attribute{
				{
					Key:   "ami",
					Value: "\"bar\"",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := marshalAttributeValues(tt.args)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf(
					"TestMarshalAttributeValues() -> \n\ngot:\n%v,\n\nwant:\n%v",
					got,
					tt.want,
				)
			}
		})
	}

}

func TestInsertState(t *testing.T) {
	fakeDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer fakeDB.Close()

	gormDB, err := gorm.Open(postgres.New(
		postgres.Config{
			Conn: fakeDB,
		},
	))
	assert.Nil(t, err)

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "versions"
		 WHERE "versions"."version_id" = $1
		 ORDER BY "versions"."id"
		 LIMIT 1`,
	)).
		WithArgs("foo").
		WillReturnRows(sqlmock.NewRows([]string{"version_id"}))

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "lineages" 
		 WHERE "lineages"."value" = $1 AND "lineages"."deleted_at" IS NULL 
		 ORDER BY "lineages"."id" 
		 LIMIT 1`,
	)).WithArgs("lineage").WillReturnRows(sqlmock.NewRows([]string{"id"}))

	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT (.+)").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, "lineage").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT (.+)").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, "path", nil, "1.0.0", 2, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	// Following queries have multiples args in a random order so we can't use real args here
	mock.ExpectQuery("^INSERT (.+)").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectQuery("^INSERT (.+)").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectQuery("^INSERT (.+)").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectQuery("^INSERT (.+)").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	db := &Database{
		DB: gormDB,
	}

	version, _ := version.NewSemver("v1.0.0")
	err = db.InsertState("path", "foo", &statefile.File{
		TerraformVersion: version,
		Serial:           2,
		Lineage:          "lineage",
		State: &states.State{
			Modules: map[string]*states.Module{
				"": {
					Addr: addrs.RootModuleInstance,
					LocalValues: map[string]cty.Value{
						"foo": cty.StringVal("foo value"),
					},
					OutputValues: map[string]*states.OutputValue{
						"bar": {
							Addr: addrs.AbsOutputValue{
								OutputValue: addrs.OutputValue{
									Name: "bar",
								},
							},
							Value:     cty.StringVal("bar value"),
							Sensitive: false,
						},
						"secret": {
							Addr: addrs.AbsOutputValue{
								OutputValue: addrs.OutputValue{
									Name: "secret",
								},
							},
							Value:     cty.StringVal("secret value"),
							Sensitive: true,
						},
					},
					Resources: map[string]*states.Resource{
						"test_thing.baz": {
							Addr: addrs.Resource{
								Mode: addrs.ManagedResourceMode,
								Type: "test_thing",
								Name: "baz",
							}.Absolute(addrs.RootModuleInstance),

							Instances: map[addrs.InstanceKey]*states.ResourceInstance{
								addrs.IntKey(0): {
									Current: &states.ResourceInstanceObjectSrc{
										SchemaVersion: 1,
										Status:        states.ObjectReady,
										AttrsJSON:     []byte(`{"woozles":"confuzles"}`),
									},
									Deposed: map[states.DeposedKey]*states.ResourceInstanceObjectSrc{},
								},
							},
							ProviderConfig: addrs.AbsProviderConfig{
								Provider: addrs.NewDefaultProvider("test"),
								Module:   addrs.RootModule,
							},
						},
					},
				},
				"module.child": {
					Addr:        addrs.RootModuleInstance.Child("child", addrs.NoKey),
					LocalValues: map[string]cty.Value{},
					OutputValues: map[string]*states.OutputValue{
						"pizza": {
							Addr: addrs.AbsOutputValue{
								Module: addrs.RootModuleInstance.Child("child", addrs.NoKey),
								OutputValue: addrs.OutputValue{
									Name: "pizza",
								},
							},
							Value:     cty.StringVal("hawaiian"),
							Sensitive: false,
						},
					},
					Resources: map[string]*states.Resource{},
				},
			},
		},
	})
	assert.Nil(t, err)

	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestUpdateStateWithLineage(t *testing.T) {
	testState := types.State{
		Model: gorm.Model{
			ID: 1,
		},
		Path:      "foo",
		TFVersion: "bar",
		Serial:    1,
		VersionID: sql.NullInt64{Int64: 1, Valid: true},
	}

	fakeDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer fakeDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: fakeDB}))
	assert.Nil(t, err)

	mock.ExpectQuery(
		"SELECT (.+)",
	).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"lineage"}).AddRow("lineage_value"))

	mock.ExpectQuery(
		"SELECT (.+)").
		WithArgs("lineage_value").
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT (.+)").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), "lineage_value").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE (.+)").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), "foo", 1, "bar", 1, 1, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	db := &Database{
		DB: gormDB,
	}

	err = db.UpdateState(testState)
	assert.Nil(t, err)

	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestUpdateStateWithoutLineage(t *testing.T) {
	testState := types.State{
		Model: gorm.Model{
			ID: 1,
		},
		Path:      "foo",
		TFVersion: "bar",
		Serial:    1,
		VersionID: sql.NullInt64{Int64: 1, Valid: true},
	}

	fakeDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer fakeDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: fakeDB}))
	assert.Nil(t, err)

	mock.ExpectQuery(
		"SELECT (.+)",
	).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"lineage"}))

	mock.ExpectQuery(
		"SELECT (.+)").
		WithArgs("foo").
		WillReturnRows(sqlmock.NewRows([]string{"lineage"}).AddRow("lineage_value"))

	mock.ExpectQuery(
		"SELECT (.+)").
		WithArgs("lineage_value").
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT (.+)").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), "lineage_value").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE (.+)").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), "foo", 1, "bar", 1, 1, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	db := &Database{
		DB: gormDB,
	}

	err = db.UpdateState(testState)
	assert.Nil(t, err)

	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestUpdateStateFail(t *testing.T) {
	testState := types.State{
		Model: gorm.Model{
			ID: 1,
		},
		Path:      "foo",
		TFVersion: "bar",
		Serial:    1,
		VersionID: sql.NullInt64{Int64: 1, Valid: true},
	}

	fakeDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer fakeDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: fakeDB}))
	assert.Nil(t, err)

	mock.ExpectQuery(
		"SELECT (.+)",
	).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"lineage"}))

	mock.ExpectQuery(
		"SELECT (.+)").
		WithArgs("foo").
		WillReturnRows(sqlmock.NewRows([]string{"lineage"}).AddRow(""))

	db := &Database{
		DB: gormDB,
	}

	err = db.UpdateState(testState)
	assert.Nil(t, err)

	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestGetState(t *testing.T) {
	fakeDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer fakeDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: fakeDB}))
	assert.Nil(t, err)

	// State retrieval
	mock.ExpectQuery(`^SELECT (.+) FROM "states" (.+)`).
		WithArgs("lineage", "foo").
		WillReturnRows(sqlmock.NewRows([]string{"id", "path"}).AddRow(1, `path`))
	mock.ExpectQuery(`^SELECT (.+) FROM "modules" (.+)`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	db := &Database{
		DB: gormDB,
	}

	state := db.GetState("lineage", "foo")
	assert.NotNil(t, state)
	assert.Equal(t, uint(1), state.ID)
	assert.Equal(t, "path", state.Path)

	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestGetLineageActivity(t *testing.T) {
	fakeDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer fakeDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: fakeDB}))
	assert.Nil(t, err)

	// Lineage activity retrieval
	mock.ExpectQuery("^SELECT (.+)").
		WithArgs("lineage").
		WillReturnRows(sqlmock.NewRows([]string{"path", "version_id"}).AddRow("path", "foo"))

	db := &Database{
		DB: gormDB,
	}

	states := db.GetLineageActivity("lineage")
	assert.NotNil(t, states)
	assert.Equal(t, "path", states[0].Path)
	assert.Equal(t, "foo", states[0].VersionID)

	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestInsertVersionCreated(t *testing.T) {
	fakeDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer fakeDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: fakeDB}))
	assert.Nil(t, err)

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "versions"
		 WHERE "versions"."version_id" = $1
		 ORDER BY "versions"."id"
		 LIMIT 1`)).
		WithArgs("foo").
		WillReturnRows(sqlmock.NewRows([]string{"version_id"}))

	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT (.+)").
		WithArgs("foo", time.Time{}).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	db := &Database{
		DB: gormDB,
	}
	err = db.InsertVersion(&state.Version{
		ID: "foo",
	})
	assert.Nil(t, err)
	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestKnownVersions(t *testing.T) {
	fakeDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer fakeDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: fakeDB}))
	assert.Nil(t, err)

	// Versions insertion
	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "versions"
		 WHERE "versions"."version_id" = $1
		 ORDER BY "versions"."id"
		 LIMIT 1`)).
		WithArgs("foo").
		WillReturnRows(sqlmock.NewRows([]string{"version_id"}))

	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT (.+)").
		WithArgs("foo", time.Time{}).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	// Known versions retrieval
	mock.ExpectQuery("^SELECT (.+)").
		WillReturnRows(sqlmock.NewRows([]string{"version_id"}).AddRow("foo"))

	db := &Database{
		DB: gormDB,
	}
	err = db.InsertVersion(&state.Version{
		ID: "foo",
	})
	assert.Nil(t, err)

	versions := db.KnownVersions()
	assert.Equal(t, 1, len(versions))
	assert.Equal(t, "foo", versions[0])

	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestSearchAttribute(t *testing.T) {
	fakeDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer fakeDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: fakeDB,
	}))
	assert.Nil(t, err)

	// Search attribute
	mock.ExpectQuery("^SELECT count(.+)").
		WillReturnRows(sqlmock.NewRows([]string{"total"}).AddRow(1))
	mock.ExpectQuery("^SELECT (.+)").
		WithArgs("%test_thing%", "%baz%", "%woozles%", `%"confuzles"%`, `%1.0.0%`, 20).
		WillReturnRows(sqlmock.NewRows([]string{"path", "version_id", "tf_version"}).AddRow("path", "foo", "1.0.0"))

	db := &Database{
		DB: gormDB,
	}

	params := url.Values{}
	params.Add("name", "baz")
	params.Add("type", "test_thing")
	params.Add("key", "woozles")
	params.Add("value", `"confuzles"`)
	params.Add("tf_version", "1.0.0")

	results, page, total := db.SearchAttribute(params)
	assert.Equal(t, 1, len(results))
	assert.Equal(t, 1, page)
	assert.Equal(t, 1, total)
	assert.Equal(t, "path", results[0].Path)
	assert.Equal(t, "1.0.0", results[0].TFVersion)
	assert.Equal(t, "foo", results[0].VersionID)

	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestListStatesVersions(t *testing.T) {
	fakeDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer fakeDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: fakeDB,
	}))
	assert.Nil(t, err)

	mock.ExpectQuery("^SELECT (.+)").
		WillReturnRows(sqlmock.NewRows([]string{"states.path", "versions.version_id"}).
			AddRow("foo", "bar").
			AddRow("baz", "bar"))

	db := &Database{
		DB: gormDB,
	}

	statesVer := db.ListStatesVersions()
	assert.NotNil(t, statesVer)
	assert.Equal(t, []string{"foo", "baz"}, statesVer["bar"])

	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestListTerraformVersionsWithCount(t *testing.T) {
	fakeDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer fakeDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: fakeDB,
	}))
	assert.Nil(t, err)

	mock.ExpectQuery("^SELECT (.+)").
		WillReturnRows(sqlmock.NewRows([]string{"tf_version", "count"}).
			AddRow("1.0.0", 1).
			AddRow("1.0.1", 1))

	db := &Database{
		DB: gormDB,
	}

	params := url.Values{}
	params.Add("orderBy", "version")

	tfVersions, err := db.ListTerraformVersionsWithCount(params)
	assert.NotNil(t, tfVersions)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(tfVersions))

	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestListStateStats(t *testing.T) {
	fakeDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer fakeDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: fakeDB,
	}))
	assert.Nil(t, err)

	mock.ExpectQuery("^SELECT (.+)").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).
			AddRow(3))

	mock.ExpectQuery("^SELECT (.+)").
		WithArgs(0).
		WillReturnRows(sqlmock.NewRows([]string{"path"}).
			AddRow("foo").
			AddRow("bar").
			AddRow("baz"))

	db := &Database{
		DB: gormDB,
	}

	params := url.Values{}
	params.Add("page", "1")

	states, page, total := db.ListStateStats(params)
	assert.NotNil(t, states)
	assert.Equal(t, 3, len(states))
	assert.Equal(t, 1, page)
	assert.Equal(t, 3, total)

	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestListResourceTypes(t *testing.T) {
	fakeDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer fakeDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: fakeDB,
	}))
	assert.Nil(t, err)

	mock.ExpectQuery("^SELECT (.+)").
		WillReturnRows(sqlmock.NewRows([]string{"resource_type"}).
			AddRow("foo").
			AddRow("bar"))

	db := &Database{
		DB: gormDB,
	}

	resourceTypes, err := db.ListResourceTypes()
	assert.NotNil(t, resourceTypes)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(resourceTypes))

	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestListResourceTypesWithCount(t *testing.T) {
	fakeDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer fakeDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: fakeDB,
	}))
	assert.Nil(t, err)

	mock.ExpectQuery("^SELECT (.+)").
		WillReturnRows(sqlmock.NewRows([]string{"resource_type", "count"}).
			AddRow("foo", 2).
			AddRow("bar", 4))

	db := &Database{
		DB: gormDB,
	}

	resourceTypes, err := db.ListResourceTypesWithCount()
	assert.NotNil(t, resourceTypes)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(resourceTypes))
	assert.Equal(t, map[string]string{
		"name":  "foo",
		"count": "2",
	}, resourceTypes[0])
	assert.Equal(t, map[string]string{
		"name":  "bar",
		"count": "4",
	}, resourceTypes[1])

	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestListResourceNames(t *testing.T) {
	fakeDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer fakeDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: fakeDB,
	}))
	assert.Nil(t, err)

	mock.ExpectQuery("^SELECT (.+)").
		WillReturnRows(sqlmock.NewRows([]string{"resource_name"}).
			AddRow("foo").
			AddRow("bar"))

	db := &Database{
		DB: gormDB,
	}

	resourceTypes, err := db.ListResourceNames()
	assert.NotNil(t, resourceTypes)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(resourceTypes))

	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestListTfVersions(t *testing.T) {
	fakeDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer fakeDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: fakeDB,
	}))
	assert.Nil(t, err)

	mock.ExpectQuery("^SELECT (.+)").
		WillReturnRows(sqlmock.NewRows([]string{"tf_version"}).
			AddRow("1.0.0").
			AddRow("1.1.0"))

	db := &Database{
		DB: gormDB,
	}

	resourceTypes, err := db.ListTfVersions()
	assert.NotNil(t, resourceTypes)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(resourceTypes))

	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestListAttributeKeys(t *testing.T) {
	fakeDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer fakeDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: fakeDB,
	}))
	assert.Nil(t, err)

	mock.ExpectQuery("^SELECT (.+)").
		WillReturnRows(sqlmock.NewRows([]string{"tf_version"}).
			AddRow("bar"))

	db := &Database{
		DB: gormDB,
	}

	attrs, err := db.ListAttributeKeys("foo")
	assert.NotNil(t, attrs)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(attrs))

	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestInsertPlan(t *testing.T) {
	fakeDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer fakeDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: fakeDB,
	}))
	assert.Nil(t, err)

	mock.ExpectQuery("^SELECT (.+)").
		WillReturnRows(sqlmock.NewRows([]string{"id", "lineage"}).
			AddRow(1, "lineage_value"))

	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT (.+)").
		WillReturnRows(sqlmock.NewRows([]string{"id"}))
	mock.ExpectQuery("^INSERT (.+)").
		WillReturnRows(sqlmock.NewRows([]string{"id"}))
	mock.ExpectQuery("^INSERT (.+)").
		WillReturnRows(sqlmock.NewRows([]string{"id"}))
	mock.ExpectQuery("^INSERT (.+)").
		WillReturnRows(sqlmock.NewRows([]string{"id"}))
	mock.ExpectQuery("^INSERT (.+)").
		WillReturnRows(sqlmock.NewRows([]string{"id"}))
	mock.ExpectQuery("^INSERT (.+)").
		WillReturnRows(sqlmock.NewRows([]string{"id"}))
	mock.ExpectQuery("^INSERT (.+)").
		WillReturnRows(sqlmock.NewRows([]string{"id"}))
	mock.ExpectQuery("^INSERT (.+)").
		WillReturnRows(sqlmock.NewRows([]string{"id"}))
	mock.ExpectCommit()

	db := &Database{
		DB: gormDB,
	}

	err = db.InsertPlan([]byte(`{"lineage":"lineage_value","terraform_version":"1.0.0","git_remote":"foo.com","git_commit":"#12345","ci_url":"","source":"","plan_json":{"format_version":"0.1","terraform_version":"0.12.6","planned_values":{"root_module":{"resources":[{"address":"aws_autoscaling_group.my_asg","mode":"managed","type":"aws_autoscaling_group","name":"my_asg","provider_name":"aws","schema_version":0,"values":{"availability_zones":["us-west-1a"],"desired_capacity":4,"enabled_metrics":null,"force_delete":true,"health_check_grace_period":300,"health_check_type":"ELB","initial_lifecycle_hook":[],"launch_configuration":"my_web_config","launch_template":[],"max_size":5,"metrics_granularity":"1Minute","min_elb_capacity":null,"min_size":1,"mixed_instances_policy":[],"name":"my_asg","name_prefix":null,"placement_group":null,"protect_from_scale_in":false,"suspended_processes":null,"tag":[],"tags":null,"termination_policies":null,"timeouts":null,"wait_for_capacity_timeout":"10m","wait_for_elb_capacity":null}},{"address":"aws_instance.web","mode":"managed","type":"aws_instance","name":"web","provider_name":"aws","schema_version":1,"values":{"ami":"ami-09b4b74c","credit_specification":[],"disable_api_termination":null,"ebs_optimized":null,"get_password_data":false,"iam_instance_profile":null,"instance_initiated_shutdown_behavior":null,"instance_type":"t2.micro","monitoring":null,"source_dest_check":true,"tags":null,"timeouts":null,"user_data":null,"user_data_base64":null}},{"address":"aws_launch_configuration.my_web_config","mode":"managed","type":"aws_launch_configuration","name":"my_web_config","provider_name":"aws","schema_version":0,"values":{"associate_public_ip_address":false,"enable_monitoring":true,"ephemeral_block_device":[],"iam_instance_profile":null,"image_id":"ami-09b4b74c","instance_type":"t2.micro","name":"my_web_config","name_prefix":null,"placement_tenancy":null,"security_groups":null,"spot_price":null,"user_data":null,"user_data_base64":null,"vpc_classic_link_id":null,"vpc_classic_link_security_groups":null}}]}},"resource_changes":[{"address":"aws_autoscaling_group.my_asg","mode":"managed","type":"aws_autoscaling_group","name":"my_asg","provider_name":"aws","change":{"actions":["create"],"before":null,"after":{"availability_zones":["us-west-1a"],"desired_capacity":4,"enabled_metrics":null,"force_delete":true,"health_check_grace_period":300,"health_check_type":"ELB","initial_lifecycle_hook":[],"launch_configuration":"my_web_config","launch_template":[],"max_size":5,"metrics_granularity":"1Minute","min_elb_capacity":null,"min_size":1,"mixed_instances_policy":[],"name":"my_asg","name_prefix":null,"placement_group":null,"protect_from_scale_in":false,"suspended_processes":null,"tag":[],"tags":null,"termination_policies":null,"timeouts":null,"wait_for_capacity_timeout":"10m","wait_for_elb_capacity":null},"after_unknown":{"arn":true,"availability_zones":[false],"default_cooldown":true,"id":true,"initial_lifecycle_hook":[],"launch_template":[],"load_balancers":true,"mixed_instances_policy":[],"service_linked_role_arn":true,"tag":[],"target_group_arns":true,"vpc_zone_identifier":true}}},{"address":"aws_instance.web","mode":"managed","type":"aws_instance","name":"web","provider_name":"aws","change":{"actions":["create"],"before":null,"after":{"ami":"ami-09b4b74c","credit_specification":[],"disable_api_termination":null,"ebs_optimized":null,"get_password_data":false,"iam_instance_profile":null,"instance_initiated_shutdown_behavior":null,"instance_type":"t2.micro","monitoring":null,"source_dest_check":true,"tags":null,"timeouts":null,"user_data":null,"user_data_base64":null},"after_unknown":{"arn":true,"associate_public_ip_address":true,"availability_zone":true,"cpu_core_count":true,"cpu_threads_per_core":true,"credit_specification":[],"ebs_block_device":true,"ephemeral_block_device":true,"host_id":true,"id":true,"instance_state":true,"ipv6_address_count":true,"ipv6_addresses":true,"key_name":true,"network_interface":true,"network_interface_id":true,"password_data":true,"placement_group":true,"primary_network_interface_id":true,"private_dns":true,"private_ip":true,"public_dns":true,"public_ip":true,"root_block_device":true,"security_groups":true,"subnet_id":true,"tenancy":true,"volume_tags":true,"vpc_security_group_ids":true}}},{"address":"aws_launch_configuration.my_web_config","mode":"managed","type":"aws_launch_configuration","name":"my_web_config","provider_name":"aws","change":{"actions":["create"],"before":null,"after":{"associate_public_ip_address":false,"enable_monitoring":true,"ephemeral_block_device":[],"iam_instance_profile":null,"image_id":"ami-09b4b74c","instance_type":"t2.micro","name":"my_web_config","name_prefix":null,"placement_tenancy":null,"security_groups":null,"spot_price":null,"user_data":null,"user_data_base64":null,"vpc_classic_link_id":null,"vpc_classic_link_security_groups":null},"after_unknown":{"ebs_block_device":true,"ebs_optimized":true,"ephemeral_block_device":[],"id":true,"key_name":true,"root_block_device":true}}}],"configuration":{"provider_config":{"aws":{"name":"aws","expressions":{"region":{"constant_value":"us-west-1"}}}},"root_module":{"resources":[{"address":"aws_autoscaling_group.my_asg","mode":"managed","type":"aws_autoscaling_group","name":"my_asg","provider_config_key":"aws","expressions":{"availability_zones":{"constant_value":["us-west-1a"]},"desired_capacity":{"constant_value":4},"force_delete":{"constant_value":true},"health_check_grace_period":{"constant_value":300},"health_check_type":{"constant_value":"ELB"},"launch_configuration":{"constant_value":"my_web_config"},"max_size":{"constant_value":5},"min_size":{"constant_value":1},"name":{"constant_value":"my_asg"}},"schema_version":0},{"address":"aws_instance.web","mode":"managed","type":"aws_instance","name":"web","provider_config_key":"aws","expressions":{"ami":{"constant_value":"ami-09b4b74c"},"instance_type":{"constant_value":"t2.micro"}},"schema_version":1},{"address":"aws_launch_configuration.my_web_config","mode":"managed","type":"aws_launch_configuration","name":"my_web_config","provider_config_key":"aws","expressions":{"image_id":{"constant_value":"ami-09b4b74c"},"instance_type":{"constant_value":"t2.micro"},"name":{"constant_value":"my_web_config"}},"schema_version":0}]}}}}`))
	assert.Nil(t, err)

	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestGetPlansSummary(t *testing.T) {
	fakeDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer fakeDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: fakeDB,
	}))
	assert.Nil(t, err)

	mock.ExpectQuery("^SELECT (.+)").
		WithArgs("lineage_value").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).
			AddRow(3))

	mock.ExpectQuery("^SELECT (.+)").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).
			AddRow(1).AddRow(2).AddRow(3))

	db := &Database{
		DB: gormDB,
	}

	plans, page, total := db.GetPlansSummary("lineage_value", "10", "1")
	assert.NotNil(t, plans)
	assert.Equal(t, 3, len(plans))
	assert.Equal(t, 1, page)
	assert.Equal(t, 3, total)

	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestGetPlan(t *testing.T) {
	fakeDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer fakeDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: fakeDB,
	}))
	assert.Nil(t, err)

	mock.ExpectQuery("^SELECT (.+)").
		WithArgs("1").
		WillReturnRows(sqlmock.NewRows([]string{"id", "tf_version"}).
			AddRow(1, "1.0.0"))

	db := &Database{
		DB: gormDB,
	}

	plan := db.GetPlan("1")
	assert.NotNil(t, plan)
	assert.Equal(t, uint(1), plan.ID)
	assert.Equal(t, "1.0.0", plan.TFVersion)

	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestGetPlans(t *testing.T) {
	fakeDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer fakeDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: fakeDB,
	}))
	assert.Nil(t, err)

	mock.ExpectQuery("^SELECT (.+)").
		WithArgs("lineage_value").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).
			AddRow(3))

	mock.ExpectQuery("^SELECT (.+)").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).
			AddRow(1).AddRow(2).AddRow(3))

	db := &Database{
		DB: gormDB,
	}

	plans, page, total := db.GetPlans("lineage_value", "10", "1")
	assert.NotNil(t, plans)
	assert.Equal(t, 3, len(plans))
	assert.Equal(t, 1, page)
	assert.Equal(t, 3, total)

	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestGetLineages(t *testing.T) {
	fakeDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer fakeDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: fakeDB,
	}))
	assert.Nil(t, err)

	mock.ExpectQuery("^SELECT (.+)").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).
			AddRow(1).AddRow(2).AddRow(3))

	db := &Database{
		DB: gormDB,
	}

	lineages := db.GetLineages("10")
	assert.NotNil(t, lineages)
	assert.Equal(t, 3, len(lineages))

	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestDefaultVersion(t *testing.T) {
	fakeDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer fakeDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: fakeDB,
	}))
	assert.Nil(t, err)

	mock.ExpectQuery("^SELECT (.+)").
		WillReturnRows(sqlmock.NewRows([]string{"version_id"}).
			AddRow("foo"))

	db := &Database{
		DB: gormDB,
	}

	version, err := db.DefaultVersion("lineage_value")
	assert.NotNil(t, version)
	assert.Equal(t, "foo", version)
	assert.Nil(t, err)

	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}

func TestClose(t *testing.T) {
	fakeDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer fakeDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: fakeDB,
	}))
	assert.Nil(t, err)

	mock.ExpectClose()

	db := &Database{
		DB: gormDB,
	}

	db.Close()

	err = mock.ExpectationsWereMet()
	assert.Nil(t, err)
}
