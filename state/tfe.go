package state

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/camptocamp/terraboard/config"
	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/terraform/states/statefile"
)

// TFE is a state provider type, leveraging Terraform Enterprise
type TFE struct {
	*tfe.Client
	org string
	ctx *context.Context
}

// NewTFE creates a new TFE object
func NewTFE(c *config.Config) (TFE, error) {
	config := &tfe.Config{
		Address: c.TFE.Address,
		Token:   c.TFE.Token,
	}

	client, err := tfe.NewClient(config)
	if err != nil {
		return TFE{}, err
	}

	ctx := context.Background()

	return TFE{
		Client: client,
		org:    c.TFE.Organization,
		ctx:    &ctx,
	}, nil
}

// GetLocks returns a map of locks by State path
func (t *TFE) GetLocks() (locks map[string]LockInfo, err error) {
	locks = make(map[string]LockInfo)

	options := tfe.WorkspaceListOptions{
		ListOptions: tfe.ListOptions{
			PageNumber: 1,
			PageSize:   50,
		},
	}

	for {
		resp, err := t.Workspaces.List(*t.ctx, t.org, options)
		if err != nil {
			return locks, err
		}

		now := time.Now()
		for _, workspace := range resp.Items {
			if workspace.Locked {
				locks[workspace.Name] = LockInfo{
					ID:        "N/A",
					Operation: "N/A",
					Info:      "N/A",
					Who:       "N/A",
					Version:   workspace.TerraformVersion,
					Created:   &now,
					Path:      workspace.Name,
				}
			}
		}

		if resp.Pagination.CurrentPage >= resp.Pagination.TotalPages {
			break
		}

		options.PageNumber = resp.Pagination.NextPage
	}

	return
}

// GetStates returns a slice of all found workspaces
func (t *TFE) GetStates() (states []string, err error) {
	options := tfe.WorkspaceListOptions{
		ListOptions: tfe.ListOptions{
			PageNumber: 1,
			PageSize:   50,
		},
	}

	for {
		resp, err := t.Workspaces.List(*t.ctx, t.org, options)
		if err != nil {
			return states, err
		}

		for _, workspace := range resp.Items {
			states = append(states, workspace.Name)
		}

		if resp.Pagination.CurrentPage >= resp.Pagination.TotalPages {
			break
		}

		options.PageNumber = resp.Pagination.NextPage
	}

	return
}

// GetVersions returns a slice of Version objects
func (t *TFE) GetVersions(state string) (versions []Version, err error) {
	options := tfe.StateVersionListOptions{
		ListOptions: tfe.ListOptions{
			PageNumber: 1,
			PageSize:   50,
		},
		Organization: &t.org,
		Workspace:    &state,
	}

	for {
		resp, err := t.StateVersions.List(*t.ctx, options)
		if err != nil {
			return versions, err
		}

		for _, version := range resp.Items {
			versions = append(versions, Version{
				ID:           version.ID,
				LastModified: version.CreatedAt,
			})
		}

		if resp.Pagination.CurrentPage >= resp.Pagination.TotalPages {
			break
		}

		options.PageNumber = resp.Pagination.NextPage
	}

	return
}

// GetState retrieves a single State from the S3 bucket
func (t *TFE) GetState(st, versionID string) (sf *statefile.File, err error) {
	// Fetch the version metadata
	version, err := t.StateVersions.Read(*t.ctx, versionID)
	if err != nil {
		return nil, err
	}

	// Download the statefile
	state, err := t.StateVersions.Download(*t.ctx, version.DownloadURL)
	if err != nil {
		return nil, err
	}

	// Parse the statefile
	sf, err = statefile.Read(bytes.NewReader(state))
	if sf == nil {
		return nil, fmt.Errorf("Unable to parse the statefile for workspace %s version %s", st, versionID)
	}

	return
}
