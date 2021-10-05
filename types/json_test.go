package types

import (
	"encoding/json"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestPlanOutputList(t *testing.T) {
	jsonPlanOutputList := `{
		"dynamodb_table_name": {
			"actions": [
				"no-op"
			],
			"before": "terraform-state-locks",
			"after": "terraform-state-locks",
			"after_unknown": false,
			"before_sensitive": false,
			"after_sensitive": false
		},
		"s3_bucket_arn": {
			"actions": [
				"no-op"
			],
			"before": "arn:aws:s3:::terraform-learning-bucket-hbollon",
			"after": "arn:aws:s3:::terraform-learning-bucket-hbollon",
			"after_unknown": false,
			"before_sensitive": false,
			"after_sensitive": false
		}
	}`

	objPlanOutputList := planOutputList{
		{
			Name: "dynamodb_table_name",
			Change: Change{
				Actions:         "[\n\t\t\t\t\"no-op\"\n\t\t\t]",
				Before:          "\"terraform-state-locks\"",
				After:           "\"terraform-state-locks\"",
				AfterUnknown:    "false",
				BeforeSensitive: "false",
				AfterSensitive:  "false",
			},
		},
		{
			Name: "s3_bucket_arn",
			Change: Change{
				Actions:         "[\n\t\t\t\t\"no-op\"\n\t\t\t]",
				Before:          "\"arn:aws:s3:::terraform-learning-bucket-hbollon\"",
				After:           "\"arn:aws:s3:::terraform-learning-bucket-hbollon\"",
				AfterUnknown:    "false",
				BeforeSensitive: "false",
				AfterSensitive:  "false",
			},
		},
	}

	tests := []struct {
		name string
		args string
		want planOutputList
	}{
		{
			"planOutputList unmarshal",
			jsonPlanOutputList,
			objPlanOutputList,
		},
	}

	for _, tt := range tests {
		var got planOutputList
		t.Run(tt.name, func(t *testing.T) {
			if err := json.Unmarshal([]byte(tt.args), &got); err != nil {
				t.Errorf("Failed to unmarshal json: %s", err)
			}
			sort.Slice(got, func(i, j int) bool {
				cmpResult := strings.Compare(got[i].Name, got[j].Name)
				return cmpResult == -1
			})
			sort.Slice(tt.want, func(i, j int) bool {
				cmpResult := strings.Compare(tt.want[i].Name, tt.want[j].Name)
				return cmpResult == -1
			})
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf(
					"TestPlanOutputList() -> \n\ngot:\n%v,\n\nwant:\n%v",
					spew.Sdump(got),
					spew.Sdump(tt.want),
				)
			}
		})
	}
}

func TestPlanStateOutputList(t *testing.T) {
	jsonPlanStateOutputList := `{
		"s3_bucket_arn": {
			"value": "arn:aws:s3:::terraform-learning-bucket-hbollon",
			"sensitive": false
		},
		"dynamodb_table_name": {
			"value": "terraform-state-locks",
			"sensitive": false
		}
	}`

	objPlanStateOutputList := planStateOutputList{
		{
			Name:      "s3_bucket_arn",
			Sensitive: false,
			Value:     "arn:aws:s3:::terraform-learning-bucket-hbollon",
		},
		{
			Name:      "dynamodb_table_name",
			Sensitive: false,
			Value:     "terraform-state-locks",
		},
	}

	tests := []struct {
		name string
		args string
		want planStateOutputList
	}{
		{
			"planStateOutputList unmarshal",
			jsonPlanStateOutputList,
			objPlanStateOutputList,
		},
	}

	for _, tt := range tests {
		var got planStateOutputList
		t.Run(tt.name, func(t *testing.T) {
			if err := json.Unmarshal([]byte(tt.args), &got); err != nil {
				t.Errorf("Failed to unmarshal json: %s", err)
			}
			sort.Slice(got, func(i, j int) bool {
				cmpResult := strings.Compare(got[i].Name, got[j].Name)
				return cmpResult == -1
			})
			sort.Slice(tt.want, func(i, j int) bool {
				cmpResult := strings.Compare(tt.want[i].Name, tt.want[j].Name)
				return cmpResult == -1
			})
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf(
					"TestPlanStateOutputList() -> \n\ngot:\n%v,\n\nwant:\n%v",
					spew.Sdump(got),
					spew.Sdump(tt.want),
				)
			}
		})
	}
}

func TestPlanVariableList(t *testing.T) {
	jsonPlanVariableList := `{
		"s3_bucket_arn": {
			"value": "arn:aws:s3:::terraform-learning-bucket-hbollon"
		},
		"dynamodb_table_name": {
			"value": "terraform-state-locks"
		}
	}`

	objPlanVariableList := planVariableList{
		{
			Key:   "s3_bucket_arn",
			Value: "map[value:arn:aws:s3:::terraform-learning-bucket-hbollon]",
		},
		{
			Key:   "dynamodb_table_name",
			Value: "map[value:terraform-state-locks]",
		},
	}

	tests := []struct {
		name string
		args string
		want planVariableList
	}{
		{
			"planVariableList unmarshal",
			jsonPlanVariableList,
			objPlanVariableList,
		},
	}

	for _, tt := range tests {
		var got planVariableList
		t.Run(tt.name, func(t *testing.T) {
			if err := json.Unmarshal([]byte(tt.args), &got); err != nil {
				t.Errorf("Failed to unmarshal json: %s", err)
			}
			sort.Slice(got, func(i, j int) bool {
				cmpResult := strings.Compare(got[i].Key, got[j].Key)
				return cmpResult == -1
			})
			sort.Slice(tt.want, func(i, j int) bool {
				cmpResult := strings.Compare(tt.want[i].Key, tt.want[j].Key)
				return cmpResult == -1
			})
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf(
					"TestPlanVariableList() -> \n\ngot:\n%v,\n\nwant:\n%v",
					spew.Sdump(got),
					spew.Sdump(tt.want),
				)
			}
		})
	}
}

func TestPlanStateResourceAttributeList(t *testing.T) {
	jsonPlanStateResourceAttributeList := `{
		"id": "terraform-state-locks",
		"arn": "arn:aws:dynamodb:eu-west-3:456426279894:table/terraform-state-locks",
		"ttl": [
			{
				"enabled": false,
				"kms_key_arn": "",
				"attribute_name": ""
			}
		],
		"name": "terraform-state-locks"
	}`

	objPlanStateResourceAttributeList := planStateResourceAttributeList{
		{
			Key:   "id",
			Value: "terraform-state-locks",
		},
		{
			Key:   "arn",
			Value: "arn:aws:dynamodb:eu-west-3:456426279894:table/terraform-state-locks",
		},
		{
			Key:   "ttl",
			Value: "[map[attribute_name: enabled:false kms_key_arn:]]",
		},
		{
			Key:   "name",
			Value: "terraform-state-locks",
		},
	}

	tests := []struct {
		name string
		args string
		want planStateResourceAttributeList
	}{
		{
			"planStateResourceAttributeList unmarshal",
			jsonPlanStateResourceAttributeList,
			objPlanStateResourceAttributeList,
		},
	}

	for _, tt := range tests {
		var got planStateResourceAttributeList
		t.Run(tt.name, func(t *testing.T) {
			if err := json.Unmarshal([]byte(tt.args), &got); err != nil {
				t.Errorf("Failed to unmarshal json: %s", err)
			}
			sort.Slice(got, func(i, j int) bool {
				cmpResult := strings.Compare(got[i].Key, got[j].Key)
				return cmpResult == -1
			})
			sort.Slice(tt.want, func(i, j int) bool {
				cmpResult := strings.Compare(tt.want[i].Key, tt.want[j].Key)
				return cmpResult == -1
			})
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf(
					"TestPlanStateResourceAttributeList() -> \n\ngot:\n%v,\n\nwant:\n%v",
					spew.Sdump(got),
					spew.Sdump(tt.want),
				)
			}
		})
	}
}

func TestRawJSON(t *testing.T) {
	jsonRawJSON := `{
		"id": "terraform-learning-bucket-hbollon",
		"acl": "private",
		"arn": "arn:aws:s3:::terraform-learning-bucket-hbollon",
		"tags": {},
		"grant": [],
		"bucket": "terraform-learning-bucket-hbollon",
		"policy": null,
		"region": "eu-west-3",
		"logging": [],
		"website": [],
		"tags_all": {},
		"cors_rule": [],
		"versioning": [
			{
				"enabled": true,
				"mfa_delete": false
			}
		],
		"bucket_prefix": null,
		"force_destroy": false,
		"request_payer": "BucketOwner",
		"hosted_zone_id": "Z3R1K369G5AVDG",
		"lifecycle_rule": [],
		"website_domain": null,
		"website_endpoint": null,
		"bucket_domain_name": "terraform-learning-bucket-hbollon.s3.amazonaws.com",
		"acceleration_status": "",
		"object_lock_configuration": [],
		"replication_configuration": [],
		"bucket_regional_domain_name": "terraform-learning-bucket-hbollon.s3.eu-west-3.amazonaws.com",
		"server_side_encryption_configuration": [
			{
				"rule": [
					{
						"bucket_key_enabled": false,
						"apply_server_side_encryption_by_default": [
							{
								"sse_algorithm": "AES256",
								"kms_master_key_id": ""
							}
						]
					}
				]
			}
		]
	}`

	objRawJSON := "{\n\t\t\"id\": \"terraform-learning-bucket-hbollon\",\n\t\t\"acl\": \"private\",\n\t\t\"arn\": \"arn:aws:s3:::terraform-learning-bucket-hbollon\",\n\t\t\"tags\": {},\n\t\t\"grant\": [],\n\t\t\"bucket\": \"terraform-learning-bucket-hbollon\",\n\t\t\"policy\": null,\n\t\t\"region\": \"eu-west-3\",\n\t\t\"logging\": [],\n\t\t\"website\": [],\n\t\t\"tags_all\": {},\n\t\t\"cors_rule\": [],\n\t\t\"versioning\": [\n\t\t\t{\n\t\t\t\t\"enabled\": true,\n\t\t\t\t\"mfa_delete\": false\n\t\t\t}\n\t\t],\n\t\t\"bucket_prefix\": null,\n\t\t\"force_destroy\": false,\n\t\t\"request_payer\": \"BucketOwner\",\n\t\t\"hosted_zone_id\": \"Z3R1K369G5AVDG\",\n\t\t\"lifecycle_rule\": [],\n\t\t\"website_domain\": null,\n\t\t\"website_endpoint\": null,\n\t\t\"bucket_domain_name\": \"terraform-learning-bucket-hbollon.s3.amazonaws.com\",\n\t\t\"acceleration_status\": \"\",\n\t\t\"object_lock_configuration\": [],\n\t\t\"replication_configuration\": [],\n\t\t\"bucket_regional_domain_name\": \"terraform-learning-bucket-hbollon.s3.eu-west-3.amazonaws.com\",\n\t\t\"server_side_encryption_configuration\": [\n\t\t\t{\n\t\t\t\t\"rule\": [\n\t\t\t\t\t{\n\t\t\t\t\t\t\"bucket_key_enabled\": false,\n\t\t\t\t\t\t\"apply_server_side_encryption_by_default\": [\n\t\t\t\t\t\t\t{\n\t\t\t\t\t\t\t\t\"sse_algorithm\": \"AES256\",\n\t\t\t\t\t\t\t\t\"kms_master_key_id\": \"\"\n\t\t\t\t\t\t\t}\n\t\t\t\t\t\t]\n\t\t\t\t\t}\n\t\t\t\t]\n\t\t\t}\n\t\t]\n\t}"

	tests := []struct {
		name string
		args string
		want rawJSON
	}{
		{
			"rawJSON unmarshal",
			jsonRawJSON,
			rawJSON(objRawJSON),
		},
	}

	for _, tt := range tests {
		var got rawJSON
		t.Run(tt.name, func(t *testing.T) {
			if err := json.Unmarshal([]byte(tt.args), &got); err != nil {
				t.Errorf("Failed to unmarshal json: %s", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf(
					"TestRawJSON() -> \n\ngot:\n%v,\n\nwant:\n%v",
					spew.Sdump(got),
					spew.Sdump(tt.want),
				)
			}
		})
	}
}
