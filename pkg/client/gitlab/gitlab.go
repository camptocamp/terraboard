package gitlab

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/machinebox/graphql"
)

// Client ..
type Client struct {
	GraphQL  *graphql.Client
	Endpoint string
	Token    string
}

// TerraformState ..
type TerraformState struct {
	ID                       string
	Name                     string
	ProjectPathWithNamespace string
	LatestVersion            struct {
		CreatedAt time.Time
		CreatedBy string
		Serial    int
	}
	Lock *TerraformStateLock
}

// TerraformStateLock ..
type TerraformStateLock struct {
	CreatedAt time.Time
	CreatedBy string
}

// TerraformStates ..
type TerraformStates []TerraformState

// Project ..
type Project struct {
	ID                string
	PathWithNamespace string
	TerraformStates   TerraformStates
}

// Projects ..
type Projects []Project

const projectsQuery string = `
query($first: Int, $after: String!) {
	projects(first: $first, after: $after) {
		pageInfo {
			endCursor
			hasNextPage
		}
		nodes {
			id
			fullPath
			terraformStates {
				count
			}
		}
	}
}`

const projectTerraformStatesQuery string = `
query($fullPath: ID!, $first: Int, $after: String!){
	project(fullPath: $fullPath) {
    terraformStates(first: $first, after: $after) {
      pageInfo {
        endCursor
        hasNextPage
      }
      nodes {
				id
				name
        lockedAt
        lockedByUser {
          publicEmail
        }
        latestVersion {
          serial
          createdAt
          createdByUser {
            publicEmail
          }
        }
      }
    }
	}
}`

// ProjectsResponse ..
type ProjectsResponse struct {
	Projects struct {
		PageInfo struct {
			EndCursor   string `json:"endCursor"`
			HasNextPage bool   `json:"hasNextPage"`
		} `json:"pageInfo"`
		Nodes []struct {
			ID              string `json:"id"`
			FullPath        string `json:"fullPath"`
			TerraformStates struct {
				Count int `json:"count"`
			} `json:"terraformStates"`
		} `json:"nodes"`
	} `json:"projects"`
}

// ProjectTerraformStatesResponse ..
type ProjectTerraformStatesResponse struct {
	Project struct {
		TerraformStates struct {
			PageInfo struct {
				EndCursor   string `json:"endCursor"`
				HasNextPage bool   `json:"hasNextPage"`
			} `json:"pageInfo"`
			Nodes []struct {
				ID           string     `json:"id"`
				Name         string     `json:"name"`
				LockedAt     *time.Time `json:"lockedAt"`
				LockedByUser *struct {
					PublicEmail string `json:"publicEmail"`
				} `json:"lockedByUser"`
				LatestVersion struct {
					Serial        int       `json:"serial"`
					CreatedAt     time.Time `json:"createdAt"`
					CreatedByUser struct {
						PublicEmail string `json:"publicEmail"`
					} `json:"createdByUser"`
				} `json:"latestVersion"`
			} `json:"nodes"`
		} `json:"terraformStates"`
	} `json:"project"`
}

// NewClient returns a new Client
func NewClient(endpoint, token string) Client {
	return Client{
		GraphQL:  graphql.NewClient(fmt.Sprintf("%s/api/graphql", endpoint)),
		Endpoint: endpoint,
		Token:    token,
	}
}

// GetProjectsWithTerraformStates ..
func (c *Client) GetProjectsWithTerraformStates() (projects Projects, err error) {
	resp := ProjectsResponse{}
	vars := map[string]interface{}{
		"first": 50,
		"after": "",
	}

	for {
		if err = c.Query(projectsQuery, &resp, vars); err != nil {
			return
		}

		for _, project := range resp.Projects.Nodes {
			if project.TerraformStates.Count > 0 {
				p := Project{
					ID:                project.ID,
					PathWithNamespace: project.FullPath,
				}

				p.TerraformStates, err = c.GetProjectTerraformStates(project.FullPath)
				if err != nil {
					return
				}

				projects = append(projects, p)
			}
		}

		if resp.Projects.PageInfo.HasNextPage {
			vars["after"] = resp.Projects.PageInfo.EndCursor
			continue
		}

		break
	}

	return
}

// GetProjectTerraformStates ..
func (c *Client) GetProjectTerraformStates(pathWithNamespace string) (terraformStates TerraformStates, err error) {
	resp := ProjectTerraformStatesResponse{}
	vars := map[string]interface{}{
		"fullPath": pathWithNamespace,
		"first":    50,
		"after":    "",
	}

	for {
		if err = c.Query(projectTerraformStatesQuery, &resp, vars); err != nil {
			return
		}

		for _, state := range resp.Project.TerraformStates.Nodes {
			terraformState := TerraformState{
				ID:                       state.ID,
				Name:                     state.Name,
				ProjectPathWithNamespace: pathWithNamespace,
			}
			terraformState.LatestVersion.CreatedAt = state.LatestVersion.CreatedAt
			terraformState.LatestVersion.CreatedBy = state.LatestVersion.CreatedByUser.PublicEmail
			terraformState.LatestVersion.Serial = state.LatestVersion.Serial

			if state.LockedAt != nil {
				terraformState.Lock = &TerraformStateLock{
					CreatedAt: *state.LockedAt,
				}

				if state.LockedByUser != nil {
					terraformState.Lock.CreatedBy = state.LockedByUser.PublicEmail
				}
			}

			terraformStates = append(terraformStates, terraformState)
		}

		if resp.Project.TerraformStates.PageInfo.HasNextPage {
			vars["after"] = resp.Project.TerraformStates.PageInfo.EndCursor
			continue
		}

		break
	}
	return
}

// GlobalPath ..
func (s *TerraformState) GlobalPath() string {
	return fmt.Sprintf("[%s] %s", s.ProjectPathWithNamespace, s.Name)
}

// GetState ..
func (c *Client) GetState(projectID, stateName, version string) (state []byte, err error) {
	var req *http.Request
	var resp *http.Response
	req, err = http.NewRequest("GET", fmt.Sprintf("%s/api/v4/projects/%s/terraform/state/%s/versions/%s", c.Endpoint, projectID, stateName, version), nil)
	if err != nil {
		return
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.Token))

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	state, err = ioutil.ReadAll(resp.Body)
	return
}

// Query ..
func (c *Client) Query(request string, response interface{}, vars map[string]interface{}) error {
	req := graphql.NewRequest(request)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	for k, v := range vars {
		req.Var(k, v)
	}
	return c.GraphQL.Run(context.TODO(), req, response)
}
