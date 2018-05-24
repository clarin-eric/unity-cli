package api

import (
	"fmt"
	"strings"
	"encoding/json"
)

/*
{
  "id" : 3,
  "state" : "valid",
  "identities" : [ {
    "typeId" : "userName",
    "value" : "tested",
    "target" : null,
    "realm" : null,
    "local" : true,
    "entityId" : 3,
    "comparableValue" : "tested"
  }, {
    "typeId" : "persistent",
    "value" : "129ffe63-63b9-4467-ae24-6bc889327b0d",
    "target" : null,
    "realm" : null,
    "local" : true,
    "entityId" : 3,
    "comparableValue" : "129ffe63-63b9-4467-ae24-6bc889327b0d"
  } ],
  "credentialInfo" : {
    "credentialRequirementId" : "cr-pass",
    "credentialsState" : {
      "credential1" : {
        "state" : "notSet",
        "extraInformation" : ""
      }
    }
  }
}
 */


type Entity struct {
	Id int64 `json:"id"`
	State string `json:"state"`
	Identities []Identity `json:"identities"`
	CredentialInfo CredentialInfo `json:"credentialInfo"`
	Groups []string
	Attributes []Attribute `json:"-"`
}

func (e *Entity) Print() {
	fmt.Printf("Entity:\n")
	fmt.Printf("    id    : %d\n", e.Id)
	fmt.Printf("    state : %s\n", e.State)
	fmt.Printf("    Identities:\n")
	for _, id := range e.Identities {
		fmt.Printf("        Type id : %s\n", id.TypeId)
		fmt.Printf("        Value   : %s\n", id.Value)
		fmt.Printf("\n")
	}
	fmt.Printf("    Credential info:\n")
	fmt.Printf("        Not supported\n")
	fmt.Printf("\n")

	fmt.Printf("    Groups:\n")
	for _, group := range e.Groups {
		fmt.Printf("        %s\n", group)
	}
	fmt.Printf("\n")

	fmt.Printf("    Attributes:\n")
	for _, attr := range e.Attributes {
		fmt.Printf("        [%10s] %s=%s\n", attr.GroupPath, attr.Name, attr.Values)
	}
	fmt.Printf("\n")
}

func (e *Entity) PrintCsv() {
	fmt.Printf("Entity:\n")
	fmt.Printf("    id    : %d\n", e.Id)
	fmt.Printf("    state : %s\n", e.State)
	fmt.Printf("    Identities:\n")
	for _, id := range e.Identities {
		fmt.Printf("        Type id : %s\n", id.TypeId)
		fmt.Printf("        Value   : %s\n", id.Value)
		fmt.Printf("\n")
	}
	fmt.Printf("    Credential info:\n")
	fmt.Printf("        Not supported\n")
}

func (e *Entity) GetAttributeValues(name string) ([]string) {
	for _, attr := range e.Attributes {
		if attr.Name == name {
			//fmt.Printf("%s=%v\n", attr.Name, attr.Values)
			return attr.Values
		}
	}
	return []string{}
}

func (e *Entity) GetAttributeValuesAsString(name string) (string) {
	result := ""
	for _, attr := range e.GetAttributeValues(name) {
		cleaned_attr := strings.Replace(attr, "\n", " ", -1)
		cleaned_attr = strings.Replace(cleaned_attr, "\r", " ", -1)

		if result == "" {
			result = cleaned_attr
		} else {
			result = fmt.Sprintf("%s,%s", result, cleaned_attr)
		}
	}
	return result
}

func (e *Entity) UpdateAttributeGroupPath(name, old_path, new_path string, client *UnityApiV1) (error) {
	for _, attr := range e.Attributes {
		if attr.GroupPath == old_path && attr.Name == name {
			new_attr := Attribute{
				Name:       attr.Name,
				GroupPath:  new_path,
				Visibility: attr.Visibility,
				Values:     attr.Values,
			}
			err := client.UpdateAttribute(e.Id, new_attr)
			if err != nil {
				return err
 			}
			return nil //success
		}
	}
	return fmt.Errorf("Failed to update attribute. Error: name=%s and path=%s not found\n", name, old_path)
}

func (e *Entity) UpdateAttributeValue(name, old_path string, new_values []string, client *UnityApiV1) (error) {
	for _, attr := range e.Attributes {
		if attr.GroupPath == old_path && attr.Name == name {
			new_attr := Attribute{
				Name:       attr.Name,
				GroupPath:  attr.GroupPath,
				Visibility: attr.Visibility,
				Values:     new_values,
			}
			err := client.UpdateAttribute(e.Id, new_attr)
			if err != nil {
				return err
			}
			return nil //success
		}
	}
	return fmt.Errorf("Failed to update attribute. Error: name=%s and path=%s not found\n", name, old_path)
}

type AttributeRepairResult struct {
	Id int64 `json:"id"`
	Updated bool `json:"updated"`
	Updates []AttributeUpdate `json:"updates,omitempty"`
}

type AttributeUpdate struct {
	Old Attribute `json:"old"`
	New Attribute `json:"new"`
	Updated bool `json:"updated"`
	Error *string `json:"error,omitempty"`
}

func (e *Entity) RepairAttributeSet(client *UnityApiV1) (*AttributeRepairResult) {
	//fmt.Printf("Processing %d\n", e.Id)

	result := AttributeRepairResult{Id: e.Id, Updated: false, Updates: []AttributeUpdate{}}
	for _, attr := range e.Attributes {
		if isClarinAttribute(attr.Name) {
			needs_update := false
			needs_group_path_update := false
			//Clone
			new_attr := Attribute{
				Name:       attr.Name,
				GroupPath:  attr.GroupPath,
				Visibility: attr.Visibility,
				Values:     attr.Values,
			}

			if new_attr.GroupPath != "/clarin" {
				new_attr.GroupPath = "/clarin"
				needs_group_path_update = true
			}

			if isRequiredClarinAttribute(new_attr.Name) && isEmpty(attr.Values){
				new_attr.Values = []string{"Unspecified"}
				needs_update = true
			}

			if needs_group_path_update || needs_update {
				attr_update := AttributeUpdate{Old: attr, New: new_attr, Updated: true}

				err := client.UpdateAttribute(e.Id, new_attr)
				if err != nil {
					//return nil, err
					attr_update.Updated = false
					error_msg := fmt.Sprintf("%s", err)
					attr_update.Error = &error_msg
				}
//				fmt.Printf("  Updated %s\n", new_attr.Name)
				//Remove the old attribute if the group path was updated
				if needs_group_path_update {
//					fmt.Printf("  Removing old attribute %s @ %s\n", attr.Name, attr.GroupPath)
					err = client.RemoveAttribute(e.Id, attr.Name, e.Identities[0].TypeId, attr.GroupPath)
					if err != nil {
						//return nil, fmt.Errorf("  Error: %v\n", err)
						attr_update.Updated = false
						error_msg := fmt.Sprintf("%s", err)
						attr_update.Error = &error_msg
					}
				}

				result.Updates = append(result.Updates, attr_update)
			}
		}
	}
	result.Updated = len(result.Updates) > 0

	return &result
}



func isClarinAttribute(name string) (bool) {
	clarin_attribute_set := []string{"clarin-full-name", "member", "clarin-motivation", "clarin-lr-list", "clarin-country", "clarin-purpose"}
	for _, clarin_name := range clarin_attribute_set {
		if name == clarin_name {
			return true
		}
	}
	return false
}

func isRequiredClarinAttribute(name string) (bool) {
	clarin_attribute_set := []string{"clarin-full-name", "member", "clarin-motivation"}
	for _, clarin_name := range clarin_attribute_set {
		if name == clarin_name {
			return true
		}
	}
	return false
}

func isEmpty(values []string) (bool) {
	if values == nil {
		return true
	}
	if len(values) == 0 {
		return true
	}
	if values[0] == "" {
		return true
	}
	return false
}

/**
 * {
 *   "typeId": "email",
 *   "value": "willem@clarin.eu",
 *   "target": null,
 *   "realm": null,
 *   "translationProfile": null,
 *   "remoteIdp": null,
 *   "confirmationInfo": {
 *     "confirmed": false,
 *     "confirmationDate": 0,
 *     "sentRequestAmount": 1,
 *     "serializedConfiguration": "{\"confirmed\":false,\"confirmationDate\":0,\"sentRequestAmount\":1}"
 *   },
 *   "metadata": null
 * }
 */
type Identity struct {
	TypeId string `json:"typeId"`
	Value string `json:"value"`
	Target string `json:"target"`
	Realm string `json:"realm"`
	Local bool `json:""`
	EntityId int64 `json:"entityId,omitempty"`
	ComparableValue string `json:"comparableValue"`
}

type CredentialInfo struct {
	CredentialRequirementId string `json:"credentialRequirementId"`
	CredentialsState interface{} `json:"credentialsState"` //TODO: properly implement
}

type Group struct {
	Path string `json:"path"`
	SubGroupPaths []string `json:"subGroups"`
	SubGroups []Group
	Members []int64 `json:"members"`

}

func (g *Group) ToJson() ([]byte) {
	json_bytes, err := json.MarshalIndent(g, "", "  ")
	if err != nil {
		fmt.Printf("Failed to marshall JSON. Error: %v\n", err)
	}
	return json_bytes
}

func (g *Group) Print() {
	g.printIndented(0)
}

func (g *Group) printIndented(indents int) {
	format := fmt.Sprintf("%%%ds", indents)
	indent := fmt.Sprintf(format, " ")

	fmt.Printf("%sGroup:\n", indent)
	fmt.Printf("%s    path           : %v\n", indent, g.Path)
	fmt.Printf("%s    subgroup names : %v\n", indent, g.SubGroupPaths)
	fmt.Printf("%s    member count   : %v\n", indent, len(g.Members))
	fmt.Printf("%s    members        : %v\n", indent, g.Members)
	fmt.Printf("%s    subgroups:\n", indent)
	for _, sub_group := range g.SubGroups {
		sub_group.printIndented(indents+4)
	}
}

/*
{
	"requestId": "3cc7435b-8a1a-4390-8eea-a571abad9ffc",
	"timestamp": 1482155419266,
	"adminComments": [],
	"status": "accepted",
	"createdEntityId": 133
	"request": {...},
	"registrationContext": {...},
}
 */
 type RegistrationRequest struct {
	 RequestId string `json:"requestId"`
	 Timestamp int64 `json:"timestamp"`
	 //AdminComments []string `json:"adminComments"`
	 Status string `json:"status"`
	 CreatedEntityId int64 `json:"createdEntityId"`
	 Content *RegistrationRequestContent `json:"request"`
	 Context *RegistrationContext `json:"registrationContext"`
 }

/**
{
	   "formId": "CLARIN Identity Registration",
	   "identities": [...],
	   "credentials": [{
		   "credentialId": "Password credential",
		   "secrets": "{\"passwords\":[{\"hash\":\"uMAStxibhd4N3qsyaiWlOvvrtGuNpYp7MfzDBVE7ksE=\",\"salt\":\"4634714456602327934\",\"time\":1482155419279,\"rehashNumber\":1}],\"outdated\":false}"
	   }],
	   "groupSelections": [{
		   "selected": false,
		   "externalIdp": null,
		   "translationProfile": null
	   }],
	   "agreements": [],
	   "comments": null,
	   "userLocale": "en",
	   "attributes": [...],
	   "registrationCode": null
   }
 */
type RegistrationRequestContent struct {
	FormId string `json:"formId"`
	RegistrationCode string `json:"registrationCode"`
	Identities []*Identity `json:identities"`
	Attributes []*Attribute `json:"attributes"`
}

/**
 * {
 *   "tryAutoAccept": true,
 *   "isOnIdpEndpoint": false,
 *   "triggeringMode": "manualAtLogin"
 * }
 */
type RegistrationContext struct {
	TryAutoAccept bool `json:"tryAutoAccept"`
	IsOnIdpEndpoint bool `json:"isOnIdpEndpoint"`
	TriggeringMode string `json:"triggeringMode"`
}

/**
 * {
 *   "values": ["Netherlands, the"],
 *   "name": "country",
 *   "groupPath": "/clarin",
 *   "visibility": "full"
 * }
 */
type Attribute struct {
	Values []string `json:"values"`
	Name string `json:"name""`
	GroupPath string `json:"groupPath"`
	Visibility string `json:"visibility"`
}

type AttributeWithConfirmationData struct {
	Values []AttributeValueWithConfirmationData `json:"values"`
	Name string `json:"name"`
	GroupPath string `json:"groupPath"`
	Visibility string `json:"visibility"`
}

type AttributeValueWithConfirmationData struct {
	Value string `json:"value"`
	ConfirmationData AttributeConfirmationData `json:"confirmationData"`
}

type AttributeConfirmationData struct {
	Confirmed bool `json:"confirmed"`
	ConfirmationDate int64 `json:"confirmationDate"`
	SentRequestAmount int64 `json:"sentRequestAmount"`
}