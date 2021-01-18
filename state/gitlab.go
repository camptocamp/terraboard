package state

import (
	"bytes"
	"fmt"
	"net/url"
	"regexp"
	"strconv"

	"github.com/camptocamp/terraboard/config"
	"github.com/camptocamp/terraboard/pkg/client/gitlab"
	"github.com/hashicorp/terraform/states/statefile"
)

// Gitlab is a state provider type, leveraging GitLab
type Gitlab struct {
	Client gitlab.Client
}

// NewGitlab creates a new Gitlab object
func NewGitlab(c *config.Config) *Gitlab {
	return &Gitlab{
		Client: gitlab.NewClient(c.Gitlab.Address, c.Gitlab.Token),
	}
}

// GetLocks returns a map of locks by State path
func (g *Gitlab) GetLocks() (locks map[string]LockInfo, err error) {
	locks = make(map[string]LockInfo)
	var projects gitlab.Projects
	projects, err = g.Client.GetProjectsWithTerraformStates()
	if err != nil {
		return
	}

	for _, project := range projects {
		for _, state := range project.TerraformStates {
			if state.Lock != nil {
				locks[state.GlobalPath()] = LockInfo{
					ID:        "N/A",
					Operation: "N/A",
					Info:      "N/A",
					Who:       state.Lock.CreatedBy,
					Version:   "N/A",
					Created:   &state.Lock.CreatedAt,
					Path:      state.GlobalPath(),
				}
			}
		}
	}

	return
}

// GetStates returns a slice of all found workspaces
func (g *Gitlab) GetStates() (states []string, err error) {
	var projects gitlab.Projects
	projects, err = g.Client.GetProjectsWithTerraformStates()
	if err != nil {
		return
	}

	for _, project := range projects {
		for _, state := range project.TerraformStates {
			states = append(states, state.GlobalPath())
		}
	}

	return
}

// GetVersions returns a slice of Version objects
func (g *Gitlab) GetVersions(state string) (versions []Version, err error) {
	var projects gitlab.Projects
	projects, err = g.Client.GetProjectsWithTerraformStates()
	if err != nil {
		return
	}

	// TODO: Highly unoptimized: whether implement a GraphQL query to fetch the correct project only
	// TODO: or cache the values locally
	for _, project := range projects {
		for _, s := range project.TerraformStates {
			if state != s.GlobalPath() {
				continue
			}

			for i := s.LatestVersion.Serial; i >= 0; i-- {
				versions = append(versions, Version{
					ID: strconv.Itoa(i),
					// TODO: Fix/implement once https://gitlab.com/gitlab-org/gitlab/-/merge_requests/45851 will be released
					// somehow it seems to be working correctly though, not sure from which place it manages to find the correct date
					LastModified: s.LatestVersion.CreatedAt,
				})
			}
		}
	}

	return
}

// GetState retrieves a single state file from the GitLab API
func (g *Gitlab) GetState(path, version string) (sf *statefile.File, err error) {
	re := regexp.MustCompile(`^\[(.*)] (.*)$`)
	stateInfo := re.FindStringSubmatch(path)
	if len(stateInfo) != 3 {
		return nil, fmt.Errorf("invalid state path: %s", path)
	}

	var state []byte
	state, err = g.Client.GetState(url.PathEscape(stateInfo[1]), url.PathEscape(stateInfo[2]), version)
	if err != nil {
		return
	}

	// Parse the statefile
	sf, err = statefile.Read(bytes.NewReader(state))
	if sf == nil {
		return nil, fmt.Errorf("Unable to parse the statefile for workspace %s version %s", path, version)
	}

	return
}
