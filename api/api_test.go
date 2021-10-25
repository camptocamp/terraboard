package api

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/camptocamp/terraboard/db"
)

func TestJSONError(t *testing.T) {
	buf := httptest.NewRecorder()
	JSONError(buf, "test", errors.New("test error"))

	if buf.Body.String() != "{\"details\":\"test error\",\"error\":\"test\"}" {
		t.Errorf("JSONError returned unexpected body: %s", buf.Body.String())
	}
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

	db := &db.Database{
		DB: gormDB,
	}

	buf := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/lineages/tfversion/count?orderBy=version", nil)
	ListTerraformVersionsWithCount(buf, req, db)

	if buf.Body.String() != `[{"count":"1","name":"1.0.0"},{"count":"1","name":"1.0.1"}]` {
		t.Errorf("TestListTerraformVersionsWithCount returned unexpected body: %s", buf.Body.String())
	}
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

	db := &db.Database{
		DB: gormDB,
	}

	buf := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/lineages/stats?page=1", nil)
	ListStateStats(buf, req, db)

	if buf.Body.String() != `{"page":1,"states":[{"path":"foo","lineage_value":"","terraform_version":"","serial":0,"version_id":"","last_modified":"0001-01-01T00:00:00Z","resource_count":0},{"path":"bar","lineage_value":"","terraform_version":"","serial":0,"version_id":"","last_modified":"0001-01-01T00:00:00Z","resource_count":0},{"path":"baz","lineage_value":"","terraform_version":"","serial":0,"version_id":"","last_modified":"0001-01-01T00:00:00Z","resource_count":0}],"total":3}` {
		t.Errorf("TestListStateStats returned unexpected body: %s", buf.Body.String())
	}
}

func TestGetState(t *testing.T) {
	fakeDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer fakeDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: fakeDB,
	}))
	assert.Nil(t, err)

	mock.ExpectQuery(`^SELECT (.+) FROM "states" (.+)`).
		WithArgs("123456789", "foo").
		WillReturnRows(sqlmock.NewRows([]string{"id", "path"}).AddRow(1, `path`))
	mock.ExpectQuery(`^SELECT (.+) FROM "modules" (.+)`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	db := &db.Database{
		DB: gormDB,
	}

	buf := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/lineages/123456789?versionid=foo", nil)
	// Hack to fake gorilla/mux vars
	vars := map[string]string{
		"lineage": "123456789",
	}
	req = mux.SetURLVars(req, vars)
	GetState(buf, req, db)

	if buf.Body.String() != `{"path":"path","version":{"version_id":"","last_modified":"0001-01-01T00:00:00Z"},"terraform_version":"","serial":0,"modules":[]}` {
		t.Errorf("TestGetState returned unexpected body: %s", buf.Body.String())
	}
}

func TestGetLineageActivity(t *testing.T) {
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
		WithArgs("123456789").
		WillReturnRows(sqlmock.NewRows([]string{"path", "version_id"}).AddRow("path", "foo"))

	db := &db.Database{
		DB: gormDB,
	}

	buf := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/lineages/123456789/activity", nil)
	// Hack to fake gorilla/mux vars
	vars := map[string]string{
		"lineage": "123456789",
	}
	req = mux.SetURLVars(req, vars)
	GetLineageActivity(buf, req, db)

	if buf.Body.String() != `[{"path":"path","lineage_value":"","terraform_version":"","serial":0,"version_id":"foo","last_modified":"0001-01-01T00:00:00Z","resource_count":0}]` {
		t.Errorf("TestGetLineageActivity returned unexpected body: %s", buf.Body.String())
	}
}

func TestStateCompare(t *testing.T) {
	fakeDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer fakeDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: fakeDB,
	}))
	assert.Nil(t, err)

	mock.ExpectQuery(`^SELECT (.+) FROM "states" (.+)`).
		WithArgs("123456789", "123").
		WillReturnRows(sqlmock.NewRows([]string{"id", "path"}).AddRow(1, `path`))
	mock.ExpectQuery(`^SELECT (.+) FROM "modules" (.+)`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	mock.ExpectQuery(`^SELECT (.+) FROM "states" (.+)`).
		WithArgs("123456789", "456").
		WillReturnRows(sqlmock.NewRows([]string{"id", "path"}).AddRow(2, `path2`))
	mock.ExpectQuery(`^SELECT (.+) FROM "modules" (.+)`).
		WithArgs(2).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	db := &db.Database{
		DB: gormDB,
	}

	buf := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/lineages/123456789/compare?from=123&to=456", nil)
	// Hack to fake gorilla/mux vars
	vars := map[string]string{
		"lineage": "123456789",
	}
	req = mux.SetURLVars(req, vars)
	StateCompare(buf, req, db)

	if buf.Body.String() != `{"stats":{"from":{"path":"path","version_id":"","resource_count":0,"terraform_version":"","serial":0},"to":{"path":"path2","version_id":"","resource_count":0,"terraform_version":"","serial":0}},"differences":{"only_in_old":{},"only_in_new":{},"in_both":null,"resource_diff":{}}}` {
		t.Errorf("TestStateCompare returned unexpected body: %s", buf.Body.String())
	}
}

func TestGetLocks(t *testing.T) {
	// TODO: Test with state provider
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

	mock.ExpectQuery("^SELECT count(.+)").
		WillReturnRows(sqlmock.NewRows([]string{"total"}).AddRow(1))
	mock.ExpectQuery("^SELECT (.+)").
		WithArgs("%test_thing%", "%baz%", "%woozles%", `%"confuzles"%`, 20).
		WillReturnRows(sqlmock.NewRows([]string{"path", "version_id", "tf_version"}).AddRow("path", "foo", "1.0.0"))

	db := &db.Database{
		DB: gormDB,
	}

	buf := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, `/search/attribute?name=baz&type=test_thing&key=woozles&value="confuzles"&tf_version=1.0.0`, nil)
	SearchAttribute(buf, req, db)

	if buf.Body.String() != `{"page":1,"results":[{"path":"path","version_id":"foo","tf_version":"1.0.0","serial":0,"lineage_value":"","module_path":"","resource_type":"","resource_name":"","resource_index":"","attribute_key":"","attribute_value":""}],"total":1}` {
		t.Errorf("TestSearchAttribute returned unexpected body: %s", buf.Body.String())
	}
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

	db := &db.Database{
		DB: gormDB,
	}

	buf := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, `/resource/types`, nil)
	ListResourceTypes(buf, req, db)

	if buf.Body.String() != `["foo","bar"]` {
		t.Errorf("TestListResourceTypes returned unexpected body: %s", buf.Body.String())
	}
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

	db := &db.Database{
		DB: gormDB,
	}

	buf := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, `/resource/types/count`, nil)
	ListResourceTypesWithCount(buf, req, db)

	if buf.Body.String() != `[{"count":"2","name":"foo"},{"count":"4","name":"bar"}]` {
		t.Errorf("TestListResourceTypesWithCount returned unexpected body: %s", buf.Body.String())
	}
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

	db := &db.Database{
		DB: gormDB,
	}

	buf := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, `/resource/names`, nil)
	ListResourceNames(buf, req, db)

	if buf.Body.String() != `["foo","bar"]` {
		t.Errorf("TestListResourceNames returned unexpected body: %s", buf.Body.String())
	}
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

	db := &db.Database{
		DB: gormDB,
	}

	buf := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, `/attribute/keys`, nil)
	ListAttributeKeys(buf, req, db)

	if buf.Body.String() != `["bar"]` {
		t.Errorf("TestListAttributeKeys returned unexpected body: %s", buf.Body.String())
	}
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

	db := &db.Database{
		DB: gormDB,
	}

	buf := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, `/tf_versions`, nil)
	ListTfVersions(buf, req, db)

	if buf.Body.String() != `["1.0.0","1.1.0"]` {
		t.Errorf("TestListTfVersions returned unexpected body: %s", buf.Body.String())
	}
}

func TestGetUser(t *testing.T) {
	buf := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, `/tf_versions"`, nil)
	req.Header.Set("X-Forwarded-User", "testUser")
	req.Header.Set("X-Forwarded-Email", "testUser@gmail.com")
	GetUser(buf, req)

	if buf.Body.String() != `{"name":"testUser","avatar_url":"http://www.gravatar.com/avatar/15847e15e9f672649d5e3199e34f7ad9","logout_url":""}` {
		t.Errorf("TestGetUser returned unexpected body: %s", buf.Body.String())
	}
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

	db := &db.Database{
		DB: gormDB,
	}

	buf := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, `/plans/summary?lineage=lineage_value&limit=10&page=1`, nil)
	GetPlansSummary(buf, req, db)

	if buf.Body.String() != `{"page":1,"plans":[{"ID":1,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"lineage_data":{"ID":0,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"lineage":"","states":null,"plans":null},"terraform_version":"","git_remote":"","git_commit":"","ci_url":"","source":"","exit_code":0,"parsed_plan":{"ID":0,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"planned_values":{"ID":0,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"root_module":{"ID":0,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null}},"prior_state":{"ID":0,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"values":{"ID":0,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"root_module":{"ID":0,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null}}}},"plan_json":null},{"ID":2,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"lineage_data":{"ID":0,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"lineage":"","states":null,"plans":null},"terraform_version":"","git_remote":"","git_commit":"","ci_url":"","source":"","exit_code":0,"parsed_plan":{"ID":0,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"planned_values":{"ID":0,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"root_module":{"ID":0,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null}},"prior_state":{"ID":0,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"values":{"ID":0,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"root_module":{"ID":0,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null}}}},"plan_json":null},{"ID":3,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"lineage_data":{"ID":0,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"lineage":"","states":null,"plans":null},"terraform_version":"","git_remote":"","git_commit":"","ci_url":"","source":"","exit_code":0,"parsed_plan":{"ID":0,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"planned_values":{"ID":0,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"root_module":{"ID":0,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null}},"prior_state":{"ID":0,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"values":{"ID":0,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"root_module":{"ID":0,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null}}}},"plan_json":null}],"total":3}` {
		t.Errorf("TestGetPlansSummary returned unexpected body: %s", buf.Body.String())
	}
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

	db := &db.Database{
		DB: gormDB,
	}

	buf := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, `/plans?planid=1`, nil)
	ManagePlans(buf, req, db)

	if buf.Body.String() != `{"ID":1,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"lineage_data":{"ID":0,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"lineage":"","states":null,"plans":null},"terraform_version":"1.0.0","git_remote":"","git_commit":"","ci_url":"","source":"","exit_code":0,"parsed_plan":{"ID":0,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"planned_values":{"ID":0,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"root_module":{"ID":0,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null}},"prior_state":{"ID":0,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"values":{"ID":0,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"root_module":{"ID":0,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null}}}},"plan_json":null}` {
		t.Errorf("TestGetPlan returned unexpected body: %s", buf.Body.String())
	}
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

	db := &db.Database{
		DB: gormDB,
	}

	buf := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, `/plans?lineage=lineage_value&limit=10&page=1`, nil)
	ManagePlans(buf, req, db)

	if buf.Body.String() != `{"page":1,"plans":[{"ID":1,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"lineage_data":{"ID":0,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"lineage":"","states":null,"plans":null},"terraform_version":"","git_remote":"","git_commit":"","ci_url":"","source":"","exit_code":0,"parsed_plan":{"ID":0,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"planned_values":{"ID":0,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"root_module":{"ID":0,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null}},"prior_state":{"ID":0,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"values":{"ID":0,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"root_module":{"ID":0,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null}}}},"plan_json":null},{"ID":2,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"lineage_data":{"ID":0,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"lineage":"","states":null,"plans":null},"terraform_version":"","git_remote":"","git_commit":"","ci_url":"","source":"","exit_code":0,"parsed_plan":{"ID":0,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"planned_values":{"ID":0,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"root_module":{"ID":0,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null}},"prior_state":{"ID":0,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"values":{"ID":0,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"root_module":{"ID":0,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null}}}},"plan_json":null},{"ID":3,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"lineage_data":{"ID":0,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"lineage":"","states":null,"plans":null},"terraform_version":"","git_remote":"","git_commit":"","ci_url":"","source":"","exit_code":0,"parsed_plan":{"ID":0,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"planned_values":{"ID":0,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"root_module":{"ID":0,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null}},"prior_state":{"ID":0,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"values":{"ID":0,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"root_module":{"ID":0,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null}}}},"plan_json":null}],"total":3}` {
		t.Errorf("TestGetPlans returned unexpected body: %s", buf.Body.String())
	}
}

func TestSubmitPlan(t *testing.T) {
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

	db := &db.Database{
		DB: gormDB,
	}

	buf := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, `/plans`, bytes.NewReader([]byte(`{"lineage":"lineage_value","terraform_version":"1.0.0","git_remote":"foo.com","git_commit":"#12345","ci_url":"","source":"","plan_json":{"format_version":"0.1","terraform_version":"0.12.6","planned_values":{"root_module":{"resources":[{"address":"aws_autoscaling_group.my_asg","mode":"managed","type":"aws_autoscaling_group","name":"my_asg","provider_name":"aws","schema_version":0,"values":{"availability_zones":["us-west-1a"],"desired_capacity":4,"enabled_metrics":null,"force_delete":true,"health_check_grace_period":300,"health_check_type":"ELB","initial_lifecycle_hook":[],"launch_configuration":"my_web_config","launch_template":[],"max_size":5,"metrics_granularity":"1Minute","min_elb_capacity":null,"min_size":1,"mixed_instances_policy":[],"name":"my_asg","name_prefix":null,"placement_group":null,"protect_from_scale_in":false,"suspended_processes":null,"tag":[],"tags":null,"termination_policies":null,"timeouts":null,"wait_for_capacity_timeout":"10m","wait_for_elb_capacity":null}},{"address":"aws_instance.web","mode":"managed","type":"aws_instance","name":"web","provider_name":"aws","schema_version":1,"values":{"ami":"ami-09b4b74c","credit_specification":[],"disable_api_termination":null,"ebs_optimized":null,"get_password_data":false,"iam_instance_profile":null,"instance_initiated_shutdown_behavior":null,"instance_type":"t2.micro","monitoring":null,"source_dest_check":true,"tags":null,"timeouts":null,"user_data":null,"user_data_base64":null}},{"address":"aws_launch_configuration.my_web_config","mode":"managed","type":"aws_launch_configuration","name":"my_web_config","provider_name":"aws","schema_version":0,"values":{"associate_public_ip_address":false,"enable_monitoring":true,"ephemeral_block_device":[],"iam_instance_profile":null,"image_id":"ami-09b4b74c","instance_type":"t2.micro","name":"my_web_config","name_prefix":null,"placement_tenancy":null,"security_groups":null,"spot_price":null,"user_data":null,"user_data_base64":null,"vpc_classic_link_id":null,"vpc_classic_link_security_groups":null}}]}},"resource_changes":[{"address":"aws_autoscaling_group.my_asg","mode":"managed","type":"aws_autoscaling_group","name":"my_asg","provider_name":"aws","change":{"actions":["create"],"before":null,"after":{"availability_zones":["us-west-1a"],"desired_capacity":4,"enabled_metrics":null,"force_delete":true,"health_check_grace_period":300,"health_check_type":"ELB","initial_lifecycle_hook":[],"launch_configuration":"my_web_config","launch_template":[],"max_size":5,"metrics_granularity":"1Minute","min_elb_capacity":null,"min_size":1,"mixed_instances_policy":[],"name":"my_asg","name_prefix":null,"placement_group":null,"protect_from_scale_in":false,"suspended_processes":null,"tag":[],"tags":null,"termination_policies":null,"timeouts":null,"wait_for_capacity_timeout":"10m","wait_for_elb_capacity":null},"after_unknown":{"arn":true,"availability_zones":[false],"default_cooldown":true,"id":true,"initial_lifecycle_hook":[],"launch_template":[],"load_balancers":true,"mixed_instances_policy":[],"service_linked_role_arn":true,"tag":[],"target_group_arns":true,"vpc_zone_identifier":true}}},{"address":"aws_instance.web","mode":"managed","type":"aws_instance","name":"web","provider_name":"aws","change":{"actions":["create"],"before":null,"after":{"ami":"ami-09b4b74c","credit_specification":[],"disable_api_termination":null,"ebs_optimized":null,"get_password_data":false,"iam_instance_profile":null,"instance_initiated_shutdown_behavior":null,"instance_type":"t2.micro","monitoring":null,"source_dest_check":true,"tags":null,"timeouts":null,"user_data":null,"user_data_base64":null},"after_unknown":{"arn":true,"associate_public_ip_address":true,"availability_zone":true,"cpu_core_count":true,"cpu_threads_per_core":true,"credit_specification":[],"ebs_block_device":true,"ephemeral_block_device":true,"host_id":true,"id":true,"instance_state":true,"ipv6_address_count":true,"ipv6_addresses":true,"key_name":true,"network_interface":true,"network_interface_id":true,"password_data":true,"placement_group":true,"primary_network_interface_id":true,"private_dns":true,"private_ip":true,"public_dns":true,"public_ip":true,"root_block_device":true,"security_groups":true,"subnet_id":true,"tenancy":true,"volume_tags":true,"vpc_security_group_ids":true}}},{"address":"aws_launch_configuration.my_web_config","mode":"managed","type":"aws_launch_configuration","name":"my_web_config","provider_name":"aws","change":{"actions":["create"],"before":null,"after":{"associate_public_ip_address":false,"enable_monitoring":true,"ephemeral_block_device":[],"iam_instance_profile":null,"image_id":"ami-09b4b74c","instance_type":"t2.micro","name":"my_web_config","name_prefix":null,"placement_tenancy":null,"security_groups":null,"spot_price":null,"user_data":null,"user_data_base64":null,"vpc_classic_link_id":null,"vpc_classic_link_security_groups":null},"after_unknown":{"ebs_block_device":true,"ebs_optimized":true,"ephemeral_block_device":[],"id":true,"key_name":true,"root_block_device":true}}}],"configuration":{"provider_config":{"aws":{"name":"aws","expressions":{"region":{"constant_value":"us-west-1"}}}},"root_module":{"resources":[{"address":"aws_autoscaling_group.my_asg","mode":"managed","type":"aws_autoscaling_group","name":"my_asg","provider_config_key":"aws","expressions":{"availability_zones":{"constant_value":["us-west-1a"]},"desired_capacity":{"constant_value":4},"force_delete":{"constant_value":true},"health_check_grace_period":{"constant_value":300},"health_check_type":{"constant_value":"ELB"},"launch_configuration":{"constant_value":"my_web_config"},"max_size":{"constant_value":5},"min_size":{"constant_value":1},"name":{"constant_value":"my_asg"}},"schema_version":0},{"address":"aws_instance.web","mode":"managed","type":"aws_instance","name":"web","provider_config_key":"aws","expressions":{"ami":{"constant_value":"ami-09b4b74c"},"instance_type":{"constant_value":"t2.micro"}},"schema_version":1},{"address":"aws_launch_configuration.my_web_config","mode":"managed","type":"aws_launch_configuration","name":"my_web_config","provider_config_key":"aws","expressions":{"image_id":{"constant_value":"ami-09b4b74c"},"instance_type":{"constant_value":"t2.micro"},"name":{"constant_value":"my_web_config"}},"schema_version":0}]}}}}`)))
	ManagePlans(buf, req, db)

	if buf.Body.String() != `` {
		t.Errorf("TestSubmitPlan returned unexpected body: %s", buf.Body.String())
	}
}

func TestManagePlansMethodError(t *testing.T) {
	buf := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, `/plans`, nil)
	ManagePlans(buf, req, nil)

	if buf.Body.String() != "Invalid request method." && buf.Code != http.StatusMethodNotAllowed {
		t.Errorf("TestManagePlansMethodError returned unexpected body: %s", buf.Body.String())
	}
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

	db := &db.Database{
		DB: gormDB,
	}

	buf := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, `/lineages?limit=10`, nil)
	GetLineages(buf, req, db)

	if buf.Body.String() != `[{"ID":1,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"lineage":"","states":null,"plans":null},{"ID":2,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"lineage":"","states":null,"plans":null},{"ID":3,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"lineage":"","states":null,"plans":null}]` {
		t.Errorf("TestGetLineages returned unexpected body: %s", buf.Body.String())
	}
}
