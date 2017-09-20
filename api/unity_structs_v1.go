package api

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