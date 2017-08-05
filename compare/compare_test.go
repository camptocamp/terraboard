package compare

import (
	"reflect"
	"testing"
	"time"

	"github.com/camptocamp/terraboard/types"
)

var fakeAttribute = types.Attribute{
	Key:   "fakeKey",
	Value: "fakeValue",
}

var fakeAttribute2 = types.Attribute{
	Key:   "fakeKey2",
	Value: "fakeValue2",
}

var fakeResource = types.Resource{
	Type:       "fakeType",
	Name:       "fakeName",
	Attributes: []types.Attribute{fakeAttribute, fakeAttribute2},
}

var fakeModule = types.Module{
	Path:      "root",
	Resources: []types.Resource{fakeResource},
}

var fakeModule2 = types.Module{
	Path:      "root/foo",
	Resources: []types.Resource{},
}

var fakeState = types.State{
	Path: "myfakepath/terraform.tfstate",
	Version: types.Version{
		VersionID:    "h8qExjo2Blk3S37CiWm7ljKxEJPuYCZw",
		LastModified: time.Unix(1501782443, 0).UTC(),
	},
	TFVersion: "0.9.8",
	Serial:    182,
	Modules:   []types.Module{fakeModule, fakeModule2},
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
		Attributes: []types.Attribute{fakeAttribute, fakeAttribute2},
	}

	result, err := getResource(fakeState, "root.fakeType.fakeName")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !reflect.DeepEqual(result, expectedResult) {
		t.Fatalf("Expected %s, got %s", expectedResult, result)
	}
}

func TestGetResource_nomatch(t *testing.T) {
	_, err := getResource(fakeState, "root.fakeType.wrongName")

	if err == nil {
		t.Fatalf("Expected an error, got nil")
	}

	expectedError := "Could not find resource with key root.fakeType.wrongName in state myfakepath/terraform.tfstate"

	if err.Error() != expectedError {
		t.Fatalf("Expected %s, got %s", expectedError, err.Error())
	}
}

func TestResourceAttributes(t *testing.T) {
	expectedResult := []string{"fakeKey", "fakeKey2"}

	result := resourceAttributes(fakeResource)

	if !reflect.DeepEqual(result, expectedResult) {
		t.Fatalf("Expected %s, got %s", expectedResult, result)
	}
}

func TestGetResourceAttribute_nomatch(t *testing.T) {
	expectedError := "Could not find attribute foo for resource fakeType.fakeName"

	_, err := getResourceAttribute(fakeResource, "foo")

	if err == nil {
		t.Fatalf("Expected error, got nil")
	}

	if err.Error() != expectedError {
		t.Fatalf("Expected %s, got %s", expectedError, err.Error())
	}
}

func TestFormatResource(t *testing.T) {
	expectedResult := `resource "fakeType" "fakeName" {
  fakeKey = "fakeValue"
  fakeKey2 = "fakeValue2"
}
`

	result := formatResource(fakeResource)

	if result != expectedResult {
		t.Fatalf("Expected %s, got %s", expectedResult, result)
	}
}

func TestStateInfo(t *testing.T) {
	expectedResult := "myfakepath/terraform.tfstate (2017-08-03 17:47:23 +0000 UTC)"

	result := stateInfo(fakeState)

	if result != expectedResult {
		t.Fatalf("Expected %s, got %s", expectedResult, result)
	}
}

func TestCompareResource(t *testing.T) {
	expectedResult := types.ResourceDiff{
		OnlyInOld: map[string]string{"fakeKey": "fakeValue", "fakeKey2": "fakeValue2"},
		OnlyInNew: map[string]string{"fakeNewKey": "fakeNewValue", "fakeNewKey2": "fakeNewValue2"},
		UnifiedDiff: `--- myfakepath/terraform.tfstate (2017-08-03 17:47:23 +0000 UTC)
+++ myfakepath/terraform.tfstate (2017-08-03 17:47:23 +0000 UTC)
@@ -1,5 +1,5 @@
 resource "fakeType" "fakeName" {
-  fakeKey = "fakeValue"
-  fakeKey2 = "fakeValue2"
+  fakeNewKey = "fakeNewValue"
+  fakeNewKey2 = "fakeNewValue2"
 }
 
`,
	}

	fakeNewAttribute := types.Attribute{
		Key:   "fakeNewKey",
		Value: "fakeNewValue",
	}

	fakeNewAttribute2 := types.Attribute{
		Key:   "fakeNewKey2",
		Value: "fakeNewValue2",
	}

	fakeNewResource := types.Resource{
		Type:       "fakeType",
		Name:       "fakeName",
		Attributes: []types.Attribute{fakeNewAttribute, fakeNewAttribute2},
	}

	fakeNewModule := types.Module{
		Path:      "root",
		Resources: []types.Resource{fakeNewResource},
	}

	fakeNewState := types.State{
		Path: "myfakepath/terraform.tfstate",
		Version: types.Version{
			VersionID:    "h8qExjo2Blk3S37CiWm7ljKxEJPuYCZw",
			LastModified: time.Unix(1501782443, 0).UTC(),
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
					OnlyInOld: map[string]string{"fakeKey": "fakeValue", "fakeKey2": "fakeValue2"},
					OnlyInNew: map[string]string{"fakeNewKey": "fakeNewValue"},
					UnifiedDiff: `--- myfakepath/terraform.tfstate (2017-08-03 17:47:23 +0000 UTC)
+++ myfakepath/terraform.tfstate (2017-08-03 17:47:23 +0000 UTC)
@@ -1,5 +1,4 @@
 resource "fakeType" "fakeName" {
-  fakeKey = "fakeValue"
-  fakeKey2 = "fakeValue2"
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
			LastModified: time.Unix(1501782443, 0).UTC(),
		},
		TFVersion: "0.9.8",
		Serial:    183,
		Modules:   []types.Module{fakeNewModule},
	}

	result, err := Compare(fakeNewState, fakeState)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !reflect.DeepEqual(result, expectedResult) {
		t.Fatalf("Expected %s, got %s", expectedResult, result)
	}
}

func TestCompare_nofrom(t *testing.T) {
	expectedError := "from version is unknown"

	_, err := Compare(types.State{}, types.State{})

	if err == nil {
		t.Fatalf("Expected error, got nil")
	}

	if err.Error() != expectedError {
		t.Fatalf("Expected %s, got %s", expectedError, err.Error())
	}
}

func TestCompare_noto(t *testing.T) {
	expectedError := "to version is unknown"

	_, err := Compare(types.State{Path: "path/to/foo.tfstate"}, types.State{})

	if err == nil {
		t.Fatalf("Expected error, got nil")
	}

	if err.Error() != expectedError {
		t.Fatalf("Expected %s, got %s", expectedError, err.Error())
	}
}
