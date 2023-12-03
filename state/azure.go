package state

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/camptocamp/terraboard/config"
	"github.com/camptocamp/terraboard/internal/terraform/states/statefile"
	log "github.com/sirupsen/logrus"
	"time"
)

// Azure is a state provider type, leveraging S3 and DynamoDB
type Azure struct {
	blobClient *azblob.Client
	container  string
}

// NewAzure creates an Azure object
func NewAzure(az config.AzureConfig) *Azure {
	client := azblobClient(az.StorageAccount, az.AccountKey)
	return &Azure{
		blobClient: client,
		container:  az.Container,
	}
}

// NewAzureCollection instantiate all needed Azure objects configurated by the user and return a slice
func NewAzureCollection(c *config.Config) []*Azure {
	var azInstances []*Azure
	for _, az := range c.Azure {
		if azInstance := NewAzure(az); azInstance != nil {
			azInstances = append(azInstances, azInstance)
		}
	}

	return azInstances
}

func (a Azure) GetLocks() (locks map[string]LockInfo, err error) {
	locks = make(map[string]LockInfo)

	pager := a.blobClient.NewListBlobsFlatPager(a.container, nil)

	for pager.More() {
		resp, err := pager.NextPage(context.TODO())
		handleError(err)

		currentTime := time.Now()

		for _, blob := range resp.Segment.BlobItems {
			if *blob.Properties.ContentType == "application/json" && "locked" == *blob.Properties.LeaseStatus {
				info := LockInfo{
					ID:        "N/A",
					Operation: "N/A",
					Info:      "N/A",
					Who:       "N/A",
					Version:   "N/A",
					Created:   &currentTime,
					Path:      *blob.Name,
				}

				locks[*blob.Name] = info
			}
		}
	}

	return locks, nil
}

func (a Azure) GetVersions(state string) (versions []Version, err error) {
	versions = []Version{}
	versions = append(versions, Version{
		ID:           state,
		LastModified: time.Now(),
	})

	return
}

func (a Azure) GetStates() (states []string, err error) {
	var keys []string

	pager := a.blobClient.NewListBlobsFlatPager(a.container, nil)

	for pager.More() {
		resp, err := pager.NextPage(context.TODO())
		handleError(err)

		for _, blob := range resp.Segment.BlobItems {
			if *blob.Properties.ContentType == "application/json" {
				keys = append(keys, *blob.Name)
			}
		}
	}
	states = keys
	log.WithFields(log.Fields{
		"states": len(states),
	}).Debug("Found states from Azure Storage Account")

	return states, nil

}

func (a Azure) GetState(st string, _ string) (sf *statefile.File, er error) {
	ctx := context.Background()

	// Download the blob
	get, err := a.blobClient.DownloadStream(ctx, a.container, st, nil)
	handleError(err)

	retryReader := get.NewRetryReader(ctx, &azblob.RetryReaderOptions{})
	handleError(err)

	sf, err = statefile.Read(retryReader)
	if sf == nil || err != nil {
		return sf, fmt.Errorf("Failed to find state: %v", err)
	}

	err = retryReader.Close()
	handleError(err)

	return
}

func azblobClient(st string, t string) (client *azblob.Client) {
	url := "https://" + st + ".blob.core.windows.net/"

	if t != "" {
		credential, err := azblob.NewSharedKeyCredential(st, t)
		handleError(err)

		client, err = azblob.NewClientWithSharedKeyCredential(url, credential, nil)
		handleError(err)
	} else {
		credential, err := azidentity.NewDefaultAzureCredential(nil)
		handleError(err)

		client, err = azblob.NewClient(url, credential, nil)
		handleError(err)
	}

	return client
}

func handleError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}
