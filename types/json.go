package types

import (
	"encoding/json"
	"fmt"
)

/*********************************************
 * Json object types allowing to define custom unmarshalling
 * rules while maintaining compatibility with Gorm
 *********************************************/

type planOutputList []PlanOutput
type planStateOutputList []PlanStateOutput
type planStateResourceAttributeList []PlanStateResourceAttribute
type rawJSON string

func (p *planOutputList) UnmarshalJSON(b []byte) error {
	tmp := map[string]Change{}
	err := json.Unmarshal(b, &tmp)
	if err != nil {
		return err
	}

	var list planOutputList
	for key, value := range tmp {
		output := PlanOutput{
			Name:   key,
			Change: value,
		}
		list = append(list, output)
	}

	*p = list
	return nil
}

func (p *planStateOutputList) UnmarshalJSON(b []byte) error {
	tmp := map[string]PlanStateOutput{}
	err := json.Unmarshal(b, &tmp)
	if err != nil {
		return err
	}

	var list planStateOutputList
	for key, value := range tmp {
		value.Name = key
		list = append(list, value)
	}

	*p = list
	return nil
}

func (p *planStateResourceAttributeList) UnmarshalJSON(b []byte) error {
	var tmp map[string]interface{}
	err := json.Unmarshal(b, &tmp)
	if err != nil {
		return err
	}

	var list planStateResourceAttributeList
	for key, value := range tmp {
		list = append(list, PlanStateResourceAttribute{
			Key:   key,
			Value: fmt.Sprintf("%v", value),
		})
	}

	*p = list
	return nil
}

func (p *rawJSON) UnmarshalJSON(b []byte) error {
	*p = rawJSON(string(b))
	return nil
}
