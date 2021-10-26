package api

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/camptocamp/terraboard/auth"
	"github.com/camptocamp/terraboard/compare"
	"github.com/camptocamp/terraboard/db"
	"github.com/camptocamp/terraboard/state"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"gorm.io/datatypes"
)

// Terraform plan payload structure usedfor swagger documentation
type planPayload struct {
	Lineage   string         `json:"lineage"`
	TFVersion string         `json:"terraform_version"`
	GitRemote string         `json:"git_remote"`
	GitCommit string         `json:"git_commit"`
	CiURL     string         `json:"ci_url"`
	Source    string         `json:"source"`
	ExitCode  int            `json:"exit_code"`
	PlanJSON  datatypes.JSON `json:"plan_json" swaggertype:"object"`
}

var _ *planPayload = nil // Avoid deadcode warning for planPayload

// JSONError is a wrapper function for errors
// which prints them to the http.ResponseWriter as a JSON response
func JSONError(w http.ResponseWriter, message string, err error) {
	errObj := make(map[string]string)
	errObj["error"] = message
	errObj["details"] = fmt.Sprintf("%v", err)
	j, _ := json.Marshal(errObj)
	if _, err := io.WriteString(w, string(j)); err != nil {
		log.Error(err.Error())
	}
}

// ListTerraformVersionsWithCount lists Terraform versions with their associated
// counts, sorted by the 'orderBy' parameter (version by default)
// @Summary Lists Terraform versions with counts
// @Description Get terraform version with their associated counts, sorted by the 'orderBy' parameter (version by default)
// @ID list-terraform-versions-with-count
// @Produce  json
// @Param   orderBy      query   string     false  "Order by constraint"
// @Success 200 {string} string	"ok"
// @Router /lineages/tfversion/count [get]
func ListTerraformVersionsWithCount(w http.ResponseWriter, r *http.Request, d *db.Database) {
	query := r.URL.Query()
	versions, _ := d.ListTerraformVersionsWithCount(query)

	j, err := json.Marshal(versions)
	if err != nil {
		JSONError(w, "Failed to marshal states", err)
		return
	}
	if _, err := io.WriteString(w, string(j)); err != nil {
		log.Error(err.Error())
	}
}

// ListStateStats returns State information for a given path as parameter
// @Summary Get Lineage states stats
// @Description Returns Lineage states stats along with paging information
// @ID list-state-stats
// @Produce  json
// @Param   page      query   integer     false  "Current page for pagination"
// @Success 200 {string} string	"ok"
// @Router /lineages/stats [get]
func ListStateStats(w http.ResponseWriter, r *http.Request, d *db.Database) {
	query := r.URL.Query()
	states, page, total := d.ListStateStats(query)

	// Build response object
	response := make(map[string]interface{})
	response["states"] = states
	response["page"] = page
	response["total"] = total
	j, err := json.Marshal(response)
	if err != nil {
		JSONError(w, "Failed to marshal states", err)
		return
	}
	if _, err := io.WriteString(w, string(j)); err != nil {
		log.Error(err.Error())
	}
}

// GetState provides information on a State
// @Summary Provides information on a State
// @Description Retrieves a State from the database by its lineage and versionID
// @ID get-state
// @Produce  json
// @Param   versionid      query   string     false  "Version ID"
// @Param   lineage      path   string     true  "Lineage"
// @Success 200 {string} string	"ok"
// @Router /lineages/{lineage} [get]
func GetState(w http.ResponseWriter, r *http.Request, d *db.Database) {
	params := mux.Vars(r)
	versionID := r.URL.Query().Get("versionid")
	var err error
	if versionID == "" {
		versionID, err = d.DefaultVersion(params["lineage"])
		if err != nil {
			JSONError(w, "Failed to retrieve default version", err)
			return
		}
	}
	state := d.GetState(params["lineage"], versionID)

	j, err := json.Marshal(state)
	if err != nil {
		JSONError(w, "Failed to marshal state", err)
		return
	}
	if _, err := io.WriteString(w, string(j)); err != nil {
		log.Error(err.Error())
	}
}

// GetLineageActivity returns the activity (version history) of a Lineage
// @Summary Get Lineage activity
// @Description Retrieves the activity (version history) of a Lineage
// @ID get-lineage-activity
// @Produce  json
// @Param   lineage      path   string     true  "Lineage"
// @Success 200 {string} string	"ok"
// @Router /lineages/{lineage}/activity [get]
func GetLineageActivity(w http.ResponseWriter, r *http.Request, d *db.Database) {
	params := mux.Vars(r)
	activity := d.GetLineageActivity(params["lineage"])

	j, err := json.Marshal(activity)
	if err != nil {
		JSONError(w, "Failed to marshal state activity", err)
		return
	}
	if _, err := io.WriteString(w, string(j)); err != nil {
		log.Error(err.Error())
	}
}

// StateCompare compares two versions ('from' and 'to') of a State
// @Summary Compares two versions of a State
// @Description Compares two versions ('from' and 'to') of a State
// @ID state-compare
// @Produce  json
// @Param   lineage      path   string     true  "Lineage"
// @Param   from      query   string     true  "Version from"
// @Param   to      query   string     true  "Version to"
// @Success 200 {string} string	"ok"
// @Router /lineages/{lineage}/compare [get]
func StateCompare(w http.ResponseWriter, r *http.Request, d *db.Database) {
	params := mux.Vars(r)
	query := r.URL.Query()
	fromVersion := query.Get("from")
	toVersion := query.Get("to")

	from := d.GetState(params["lineage"], fromVersion)
	to := d.GetState(params["lineage"], toVersion)
	compare, err := compare.Compare(from, to)
	if err != nil {
		JSONError(w, "Failed to compare state versions", err)
		return
	}

	j, err := json.Marshal(compare)
	if err != nil {
		JSONError(w, "Failed to marshal state compare", err)
		return
	}
	if _, err := io.WriteString(w, string(j)); err != nil {
		log.Error(err.Error())
	}
}

// GetLocks returns information on locked States
// @Summary Get locked states information
// @Description Returns information on locked States
// @ID get-locks
// @Produce  json
// @Success 200 {string} string	"ok"
// @Router /locks [get]
func GetLocks(w http.ResponseWriter, _ *http.Request, sps []state.Provider) {
	allLocks := make(map[string]state.LockInfo)
	for _, sp := range sps {
		locks, err := sp.GetLocks()
		if err != nil {
			JSONError(w, "Failed to get locks on a provider", err)
			return
		}
		for k, v := range locks {
			allLocks[k] = v
		}
	}

	j, err := json.Marshal(allLocks)
	if err != nil {
		JSONError(w, "Failed to marshal locks", err)
		return
	}
	if _, err := io.WriteString(w, string(j)); err != nil {
		log.Error(err.Error())
	}
}

// SearchAttribute performs a search on Resource Attributes
// by various parameters
// @Summary Search Resource Attributes
// @Description Performs a search on Resource Attributes by various parameters
// @ID search-attribute
// @Produce  json
// @Param   versionid      query   string     false  "Version ID"
// @Param   type      query   string     false  "Ressource type"
// @Param   name      query   string     false  "Resource ID"
// @Param   key      query   string     false  "Attribute Key"
// @Param   value      query   string     false  "Attribute Value"
// @Param   tf_version      query   string     false  "Terraform Version"
// @Param   lineage_value      query   string     false  "Lineage"
// @Success 200 {string} string	"ok"
// @Router /search/attribute [get]
func SearchAttribute(w http.ResponseWriter, r *http.Request, d *db.Database) {
	query := r.URL.Query()
	result, page, total := d.SearchAttribute(query)

	// Build response object
	response := make(map[string]interface{})
	response["results"] = result
	response["page"] = page
	response["total"] = total

	j, err := json.Marshal(response)
	if err != nil {
		JSONError(w, "Failed to marshal json", err)
		return
	}
	if _, err := io.WriteString(w, string(j)); err != nil {
		log.Error(err.Error())
	}
}

// ListResourceTypes lists all Resource types
// @Summary Get Resource types
// @Description Lists all Resource types
// @ID list-resource-types
// @Produce  json
// @Success 200 {string} string	"ok"
// @Router /resource/types [get]
func ListResourceTypes(w http.ResponseWriter, _ *http.Request, d *db.Database) {
	result, _ := d.ListResourceTypes()
	j, err := json.Marshal(result)
	if err != nil {
		JSONError(w, "Failed to marshal json", err)
		return
	}
	if _, err := io.WriteString(w, string(j)); err != nil {
		log.Error(err.Error())
	}
}

// ListResourceTypesWithCount lists all Resource types with their associated count
// @Summary Get resource types with count
// @Description Lists all resource types with their associated count
// @ID list-resource-types-with-count
// @Produce  json
// @Success 200 {string} string	"ok"
// @Router /resource/types/count [get]
func ListResourceTypesWithCount(w http.ResponseWriter, _ *http.Request, d *db.Database) {
	result, _ := d.ListResourceTypesWithCount()
	j, err := json.Marshal(result)
	if err != nil {
		JSONError(w, "Failed to marshal json", err)
		return
	}
	if _, err := io.WriteString(w, string(j)); err != nil {
		log.Error(err.Error())
	}
}

// ListResourceNames lists all Resource names
// @Summary Get resource names
// @Description Lists all resource names
// @ID list-resource-names
// @Produce  json
// @Success 200 {string} string	"ok"
// @Router /resource/names [get]
func ListResourceNames(w http.ResponseWriter, _ *http.Request, d *db.Database) {
	result, _ := d.ListResourceNames()
	j, err := json.Marshal(result)
	if err != nil {
		JSONError(w, "Failed to marshal json", err)
		return
	}
	if _, err := io.WriteString(w, string(j)); err != nil {
		log.Error(err.Error())
	}
}

// ListAttributeKeys lists all Resource Attribute Keys,
// optionally filtered by resource_type
// @Summary Get resource attribute keys
// @Description Lists all resource attribute keys, optionally filtered by resource_type
// @ID list-attribute-keys
// @Produce  json
// @Param   resource_type      query   string     false  "Resource Type"
// @Success 200 {string} string	"ok"
// @Router /attribute/keys [get]
func ListAttributeKeys(w http.ResponseWriter, r *http.Request, d *db.Database) {
	resourceType := r.URL.Query().Get("resource_type")
	result, _ := d.ListAttributeKeys(resourceType)
	j, err := json.Marshal(result)
	if err != nil {
		JSONError(w, "Failed to marshal json", err)
		return
	}
	if _, err := io.WriteString(w, string(j)); err != nil {
		log.Error(err.Error())
	}
}

// ListTfVersions lists all Terraform versions
// @Summary Get terraform versions
// @Description Lists all terraform versions
// @ID list-tf-versions
// @Produce  json
// @Success 200 {string} string	"ok"
// @Router /tf_versions [get]
func ListTfVersions(w http.ResponseWriter, _ *http.Request, d *db.Database) {
	result, _ := d.ListTfVersions()
	j, err := json.Marshal(result)
	if err != nil {
		JSONError(w, "Failed to marshal json", err)
		return
	}
	if _, err := io.WriteString(w, string(j)); err != nil {
		log.Error(err.Error())
	}
}

// GetUser returns information about the logged user
// @Summary Get logged user information
// @Description Returns information about the logged user
// @ID get-user
// @Produce  json
// @Success 200 {string} string	"ok"
// @Router /user [get]
func GetUser(w http.ResponseWriter, r *http.Request) {
	name := r.Header.Get("X-Forwarded-User")
	email := r.Header.Get("X-Forwarded-Email")

	user := auth.UserInfo(name, email)

	j, err := json.Marshal(user)
	if err != nil {
		JSONError(w, "Failed to marshal user information", err)
		return
	}
	if _, err := io.WriteString(w, string(j)); err != nil {
		log.Error(err.Error())
	}
}

// SubmitPlan inserts a new Terraform plan in the database.
// /api/plans POST endpoint callback
// @Summary Submit a new plan
// @Description Submits and inserts a new Terraform plan in the database.
// @ID submit-plan
// @Accept  json
// @Param   plan      body   api.planPayload     false  "Wrapped plan"
// @Router /plans [post]
func SubmitPlan(w http.ResponseWriter, r *http.Request, db *db.Database) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorf("Failed to read body: %v", err)
		JSONError(w, "Failed to read body during plan submit", err)
		return
	}

	if err = db.InsertPlan(body); err != nil {
		log.Errorf("Failed to insert plan to db: %v", err)
		JSONError(w, "Failed to insert plan to db", err)
		return
	}
}

// GetPlansSummary provides summary of all Plan by lineage (only metadata added by the wrapper).
// Optional "&limit=X" parameter to limit requested quantity of plans.
// Optional "&page=X" parameter to add an offset to the query and enable pagination.
// Sorted by most recent to oldest.
// /api/plans/summary GET endpoint callback
// Also return pagination informations (current page ans total items count in database)
// @Summary Get summary of all Plan by lineage
// @Description Provides summary of all Plan by lineage (only metadata added by the wrapper). Sorted by most recent to oldest. Returns also paging informations (current page ans total items count in database)
// @ID get-plans-summary
// @Produce  json
// @Param   lineage      query   string     false  "Lineage"
// @Param   page      query   integer     false  "Page"
// @Param   limit      query   integer     false  "Limit"
// @Success 200 {string} string	"ok"
// @Router /plans/summary [get]
func GetPlansSummary(w http.ResponseWriter, r *http.Request, db *db.Database) {
	lineage := r.URL.Query().Get("lineage")
	limit := r.URL.Query().Get("limit")
	page := r.URL.Query().Get("page")
	plans, currentPage, total := db.GetPlansSummary(lineage, limit, page)

	response := make(map[string]interface{})
	response["plans"] = plans
	response["page"] = currentPage
	response["total"] = total
	j, err := json.Marshal(response)
	if err != nil {
		log.Errorf("Failed to marshal plans: %v", err)
		JSONError(w, "Failed to marshal plans", err)
		return
	}
	if _, err := io.WriteString(w, string(j)); err != nil {
		log.Error(err.Error())
	}
}

// GetPlan provides a specific Plan of a lineage using ID.
// /api/plans GET endpoint callback on request with ?plan_id=X parameter
// @Summary Get plans
// @Description Provides a specific Plan of a lineage using ID or all plans if no ID is provided
// @ID get-plans
// @Produce  json
// @Param   planid      query   string     false  "Plan's ID"
// @Param   page      query   integer     false  "Page"
// @Param   limit      query   integer     false  "Limit"
// @Success 200 {string} string	"ok"
// @Router /plans [get]
func GetPlan(w http.ResponseWriter, r *http.Request, db *db.Database) {
	id := r.URL.Query().Get("planid")
	plan := db.GetPlan(id)

	j, err := json.Marshal(plan)
	if err != nil {
		log.Errorf("Failed to marshal plan: %v", err)
		JSONError(w, "Failed to marshal plan", err)
		return
	}
	if _, err := io.WriteString(w, string(j)); err != nil {
		log.Error(err.Error())
	}
}

// GetPlans provides all Plan by lineage.
// Optional "&limit=X" parameter to limit requested quantity of plans.
// Optional "&page=X" parameter to add an offset to the query and enable pagination.
// Sorted by most recent to oldest.
// /api/plans GET endpoint callback
// Also return pagination informations (current page ans total items count in database)
func GetPlans(w http.ResponseWriter, r *http.Request, db *db.Database) {
	lineage := r.URL.Query().Get("lineage")
	limit := r.URL.Query().Get("limit")
	page := r.URL.Query().Get("page")
	plans, currentPage, total := db.GetPlans(lineage, limit, page)

	response := make(map[string]interface{})
	response["plans"] = plans
	response["page"] = currentPage
	response["total"] = total
	j, err := json.Marshal(response)
	if err != nil {
		log.Errorf("Failed to marshal plans: %v", err)
		JSONError(w, "Failed to marshal plans", err)
		return
	}
	if _, err := io.WriteString(w, string(j)); err != nil {
		log.Error(err.Error())
	}
}

// ManagePlans is used to route the request to the appropriated handler function
// on /api/plans request
func ManagePlans(w http.ResponseWriter, r *http.Request, db *db.Database) {
	if r.Method == "GET" {
		if r.URL.Query().Get("planid") != "" {
			GetPlan(w, r, db)
		} else {
			GetPlans(w, r, db)
		}
	} else if r.Method == "POST" {
		SubmitPlan(w, r, db)
	} else {
		http.Error(w, "Invalid request method.", 405)
	}
}

// GetLineages recover all Lineage from db.
// Optional "&limit=X" parameter to limit requested quantity of them.
// Sorted by most recent to oldest.
// @Summary Get lineages
// @Description List all existing lineages
// @ID get-lineages
// @Produce  json
// @Param   limit      query   integer     false  "Limit"
// @Success 200 {string} string	"ok"
// @Router /lineages [get]
func GetLineages(w http.ResponseWriter, r *http.Request, db *db.Database) {
	limit := r.URL.Query().Get("limit")
	lineages := db.GetLineages(limit)

	j, err := json.Marshal(lineages)
	if err != nil {
		log.Errorf("Failed to marshal lineages: %v", err)
		JSONError(w, "Failed to marshal lineages", err)
		return
	}
	if _, err := io.WriteString(w, string(j)); err != nil {
		log.Error(err.Error())
	}
}
