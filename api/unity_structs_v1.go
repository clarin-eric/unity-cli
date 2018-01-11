package api

import (
	"fmt"
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
	Path string
	SubGroups []string `json:"subGroups"`
	Members []int64 `json:"members"`
}

func (g *Group) Print() {
	fmt.Printf("Group:\n")
	fmt.Printf("    path         : %v\n", g.Path)
	fmt.Printf("    subGroups    : %v\n", g.SubGroups)
	fmt.Printf("    member count : %v\n", len(g.Members))
	fmt.Printf("    members      : %v\n", g.Members)
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
	Name string `json:"name`
	GroupPath string `json:"groupPath"`
	Visibility string `json:"visibility"`
}