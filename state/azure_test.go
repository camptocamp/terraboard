package state

import (
	"github.com/camptocamp/terraboard/config"

	"testing"
)

func TestNewAzure(t *testing.T) {
	azInstance := NewAzure(
		config.AzureConfig{
			StorageAccount: "storage-account",
			Container:      "container",
		},
	)

	if azInstance == nil {
		t.Error("Azure instance is nil")
	}
}

func TestNewAzureKey(t *testing.T) {
	azInstance := NewAzure(
		config.AzureConfig{
			StorageAccount: "storage-account",
			Container:      "container",
			AccountKey:     "VGhpcyBpcyBhIGZha2UgYmFzZTY0IHN0cmluZw==",
		},
	)

	if azInstance == nil {
		t.Error("Azure instance is nil")
	}
}
