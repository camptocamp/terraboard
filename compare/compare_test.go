package compare

import (
	"github.com/camptocamp/terraboard/types"
	"reflect"
	"testing"
	"time"
)

var fakeAttribute = types.Attribute{
	Key:   "fakeKey",
	Value: "fakeValue",
}

var fakeResource = types.Resource{
	Type:       "fakeType",
	Name:       "fakeName",
	Attributes: []types.Attribute{fakeAttribute},
}

var fakeModule = types.Module{
	Path:      "root",
	Resources: []types.Resource{fakeResource},
}

var fakeState = types.State{
	Path: "myfakepath/terraform.tfstate",
	Version: types.Version{
		VersionID:    "h8qExjo2Blk3S37CiWm7ljKxEJPuYCZw",
		LastModified: time.Unix(1501782443, 0),
	},
	TFVersion: "0.9.8",
	Serial:    182,
	Modules:   []types.Module{fakeModule},
}

var fakeStateInfo = types.StateInfo{
	Path:          "myfakepath/terraform.tfstate",
	VersionID:     "h8qExjo2Blk3S37CiWm7ljKxEJPuYCZw",
	ResourceCount: 1,
	TFVersion:     "0.9.8",
	Serial:        182,
}

func TestStateResources(t *testing.T) {
	expectedResult := []string{"root.fakeType.fakeName"}

	result := stateResources(fakeState)

	if !reflect.DeepEqual(result, expectedResult) {
		t.Fatalf("Expected %s, got %s", expectedResult, result)
	}
}

func TestSliceDiff(t *testing.T) {
	expectedResult := []string{"apple", "orange"}
	s1 := []string{"apple", "banana", "orange", "melon"}
	s2 := []string{"melon", "banana", "lemon"}

	result := sliceDiff(s1, s2)

	if !reflect.DeepEqual(result, expectedResult) {
		t.Fatalf("Expected %s, got %s", expectedResult, result)
	}
}

func TestSliceInter(t *testing.T) {
	expectedResult := []string{"banana", "melon"}
	s1 := []string{"apple", "banana", "orange", "melon"}
	s2 := []string{"melon", "banana", "lemon"}

	result := sliceInter(s1, s2)

	if !reflect.DeepEqual(result, expectedResult) {
		t.Fatalf("Expected %s, got %s", expectedResult, result)
	}
}

func TestGetResource(t *testing.T) {
	expectedResult := types.Resource{
		Type:       "fakeType",
		Name:       "fakeName",
		Attributes: []types.Attribute{fakeAttribute},
	}

	result := getResource(fakeState, "root.fakeType.fakeName")

	if !reflect.DeepEqual(result, expectedResult) {
		t.Fatalf("Expected %s, got %s", expectedResult, result)
	}
}

func TestResourceAttributes(t *testing.T) {
	expectedResult := []string{"fakeKey"}

	result := resourceAttributes(fakeResource)

	if !reflect.DeepEqual(result, expectedResult) {
		t.Fatalf("Expected %s, got %s", expectedResult, result)
	}
}

func TestFormatResource(t *testing.T) {
	expectedResult := `resource "fakeType" "fakeName" {
  fakeKey = "fakeValue"
}
`

	result := formatResource(fakeResource)

	if result != expectedResult {
		t.Fatalf("Expected %s, got %s", expectedResult, result)
	}
}

func TestStateInfo(t *testing.T) {
	expectedResult := "myfakepath/terraform.tfstate (2017-08-03 19:47:23 +0200 CEST)"

	result := stateInfo(fakeState)

	if result != expectedResult {
		t.Fatalf("Expected %s, got %s", expectedResult, result)
	}
}

func TestCompareResource(t *testing.T) {
	expectedResult := types.ResourceDiff{
		OnlyInOld: map[string]string{"fakeKey": "fakeValue"},
		OnlyInNew: map[string]string{"fakeNewKey": "fakeNewValue"},
		UnifiedDiff: `--- myfakepath/terraform.tfstate (2017-08-03 19:47:23 +0200 CEST)
+++ myfakepath/terraform.tfstate (2017-08-03 19:47:23 +0200 CEST)
@@ -1,4 +1,4 @@
 resource "fakeType" "fakeName" {
-  fakeKey = "fakeValue"
+  fakeNewKey = "fakeNewValue"
 }
 
`,
	}

	fakeNewAttribute := types.Attribute{
		Key:   "fakeNewKey",
		Value: "fakeNewValue",
	}

	fakeNewResource := types.Resource{
		Type:       "fakeType",
		Name:       "fakeName",
		Attributes: []types.Attribute{fakeNewAttribute},
	}

	fakeNewModule := types.Module{
		Path:      "root",
		Resources: []types.Resource{fakeNewResource},
	}

	fakeNewState := types.State{
		Path: "myfakepath/terraform.tfstate",
		Version: types.Version{
			VersionID:    "h8qExjo2Blk3S37CiWm7ljKxEJPuYCZw",
			LastModified: time.Unix(1501782443, 0),
		},
		TFVersion: "0.9.8",
		Serial:    182,
		Modules:   []types.Module{fakeNewModule},
	}

	result := compareResource(fakeState, fakeNewState, "root.fakeType.fakeName")

	if !reflect.DeepEqual(result, expectedResult) {
		t.Fatalf("Expected %s, got %s", expectedResult, result)
	}
}

func TestCompare_Result(t *testing.T) {

	fromStateInfo := types.StateInfo{
		Path:          "myfakepath/terraform.tfstate",
		VersionID:     "h8qExjo2Blk3SUYGBniWm7ljKxEJPuYCZw",
		ResourceCount: 1,
		TFVersion:     "0.9.8",
		Serial:        183,
	}

	expectedResult := types.StateCompare{
		Stats: struct {
			From types.StateInfo `json:"from"`
			To   types.StateInfo `json:"to"`
		}{
			From: fromStateInfo,
			To:   fakeStateInfo,
		},
		Differences: struct {
			OnlyInOld    map[string]string             `json:"only_in_old"`
			OnlyInNew    map[string]string             `json:"only_in_new"`
			InBoth       []string                      `json:"in_both"`
			ResourceDiff map[string]types.ResourceDiff `json:"resource_diff"`
		}{
			OnlyInOld: map[string]string{},
			OnlyInNew: map[string]string{},
			InBoth:    []string{"root.fakeType.fakeName"},
			ResourceDiff: map[string]types.ResourceDiff{
				"root.fakeType.fakeName": types.ResourceDiff{
					OnlyInOld: map[string]string{"fakeKey": "fakeValue"},
					OnlyInNew: map[string]string{"fakeNewKey": "fakeNewValue"},
					UnifiedDiff: `--- myfakepath/terraform.tfstate (2017-08-03 19:47:23 +0200 CEST)
+++ myfakepath/terraform.tfstate (2017-08-03 19:47:23 +0200 CEST)
@@ -1,4 +1,4 @@
 resource "fakeType" "fakeName" {
-  fakeKey = "fakeValue"
+  fakeNewKey = "fakeNewValue"
 }
 
`,
				},
			},
		},
	}

	fakeNewAttribute := types.Attribute{
		Key:   "fakeNewKey",
		Value: "fakeNewValue",
	}

	fakeNewResource := types.Resource{
		Type:       "fakeType",
		Name:       "fakeName",
		Attributes: []types.Attribute{fakeNewAttribute},
	}

	fakeNewModule := types.Module{
		Path:      "root",
		Resources: []types.Resource{fakeNewResource},
	}

	fakeNewState := types.State{
		Path: "myfakepath/terraform.tfstate",
		Version: types.Version{
			VersionID:    "h8qExjo2Blk3SUYGBniWm7ljKxEJPuYCZw",
			LastModified: time.Unix(1501782443, 0),
		},
		TFVersion: "0.9.8",
		Serial:    183,
		Modules:   []types.Module{fakeNewModule},
	}

	result, _ := Compare(fakeNewState, fakeState)

	if !reflect.DeepEqual(result, expectedResult) {
		t.Fatalf("Expected %s, got %s", expectedResult, result)
	}
}
