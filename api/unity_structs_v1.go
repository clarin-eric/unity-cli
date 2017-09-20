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

type Identity struct {
	TypeId string `json:""`
	Value string `json:""`
	Target string `json:""`
	Realm string `json:""`
	Local bool `json:""`
	EntityId int64 `json:"entityId"`
	ComparableValue string `json:"comparableValue"`
}

type CredentialInfo struct {
	CredentialRequirementId string `json:"credentialRequirementId"`
	CredentialsState interface{} `json:"credentialsState"` //TODO: properly implement
}