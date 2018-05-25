package report

import (
	"clarin/unity-cli/api"
	"encoding/json"
	"fmt"
)

type ReportEmptyAttributes struct {
	Full_name map[int64]string `json:"fullname,omitempty"`
	Member map[int64]string `json:"member,omitempty"`
	Motivation map[int64]string `json:"motivation,omitempty"`
	Language_Resource map[int64]string `json:"language_resources,omitempty"`
	Country map[int64]string `json:"country,omitempty"`
	Purpose map[int64]string `json:"purpose,omitempty"`
}

func (r *ReportEmptyAttributes) reset() {
	r.Full_name = map[int64]string{}
	r.Member = map[int64]string{}
	r.Motivation = map[int64]string{}
	r.Language_Resource = map[int64]string{}
	r.Country = map[int64]string{}
	r.Purpose = map[int64]string{}
}

func (r *ReportEmptyAttributes) Compute(entities []api.Entity, attribute_set []string) {
	r.reset()
	for _, e := range entities {
		for _, attribute_name := range attribute_set {
			switch attribute_name {
			case "clarin-full-name": r.Full_name = r.checkAttribute(e, "clarin-full-name", r.Full_name)
			case "member": r.Member = r.checkAttribute(e, "member", r.Member)
			case "clarin-motivation": r.Motivation = r.checkAttribute(e, "clarin-motivation", r.Motivation)
			case "clarin-lr-list": r.Language_Resource = r.checkAttribute(e, "clarin-lr-list", r.Language_Resource)
			case "clarin-country": r.Country = r.checkAttribute(e, "clarin-country", r.Country)
			case "clarin-purpose": r.Purpose = r.checkAttribute(e, "clarin-purpose", r.Purpose)
			}
		}
	}
}

func (r *ReportEmptyAttributes) checkAttribute(e api.Entity, attribute_name string, m map[int64]string) (map[int64]string) {
	email, _, _ := getEmail(e)
	value := e.GetAttributeValuesAsString(attribute_name)
	if value == "" {
		m[e.Id] = email
	}
	return m
}

func (r *ReportEmptyAttributes) ToJson() ([]byte) {
	json_bytes, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		fmt.Printf("Failed to marshall JSON. Error: %v\n", err)
	}
	return json_bytes
}

func (r *ReportEmptyAttributes) GetReportAsArray() ([][]string) {
	data := [][]string{}
	return data
}




type ReportEmptyAttributesGroupedByEntityId struct {
	Entities map[int64][]string
}

func (r *ReportEmptyAttributesGroupedByEntityId) Compute(entities []api.Entity, attribute_set []string) {
	r.Entities = map[int64][]string{}
	for _, e := range entities {
		r.checkAttributes(e, attribute_set)
	}
}

func (r *ReportEmptyAttributesGroupedByEntityId) checkAttributes(e api.Entity, attribute_set []string) {
	//email := getEmail(e)
	for _, attribute_name := range attribute_set {
		value := e.GetAttributeValuesAsString(attribute_name)
		if value == "" {
			r.Entities[e.Id] = append(r.Entities[e.Id], attribute_name)
		}
	}
}

func (r *ReportEmptyAttributesGroupedByEntityId) ToJson() ([]byte) {
	json_bytes, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		fmt.Printf("Failed to marshall JSON. Error: %v\n", err)
	}
	return json_bytes
}

func (r *ReportEmptyAttributesGroupedByEntityId) ToTabular() ([][]string) {
	data := [][]string{}
	return data
}