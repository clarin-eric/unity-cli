package api

import (
	"fmt"
	"errors"
	"unity-client/http"
	"encoding/json"
)

type UnityApiV1 struct {
	verbose bool
	client *http.UnityClient
}

func GetUnityApiV1(base_url string, verbose bool, insecure bool, username, password string) (*UnityApiV1, error) {
	if verbose {
		fmt.Printf("Global flags:\n")
		fmt.Printf("    Verbose      : %t\n", verbose)
		fmt.Printf("    API base     : %s\n", base_url)
		fmt.Printf("    Insecure     : %t\n", insecure)
		fmt.Printf("    Username     : %s\n", username)
		fmt.Printf("    Password     : %s\n", "******")
	}

	client,err := http.GetNewUnityClient(base_url, verbose, insecure, username, password)
	if err != nil {
		return nil, err
	}

	if verbose {
		fmt.Printf("\n")
	}
	return &UnityApiV1{
		verbose: verbose,
		client: client,
	}, nil
}

/*
@Path("/resolve/{identityType}/{identityValue}")
@GET
 */
func (api *UnityApiV1) Resolve(identity_type, identity_value string) (*Entity, error) {
	//Build url
	api.client.SetPathByString("./resolve/%s/%s", identity_type, identity_value)

	//Issue request
	response := api.client.IssueRequest()
	if response.ErrorMessage != nil {
		return nil, errors.New(*response.ErrorMessage)
	}

	// Unmarshal response
	var entity Entity
	err := json.Unmarshal(response.Body, &entity)
	if err != nil {
		return nil, errors.New(fmt.Sprint("Failed to unmarshall response. Error: %v\nBody:\n%v\n", err, string(response.Body)))
	}

	return &entity, nil
}

/**
@Path("/entity/{entityId}")
@QueryParam("identityType")
@GET
 */
func (api *UnityApiV1) Entity(entity_id int64, identity_type *string) (*Entity, error) {
	api.client.SetPathByString("./entity/%d", entity_id)
	if identity_type != nil {
		api.client.AddQueryParam("identityType", *identity_type)
	}

	response := api.client.IssueRequest()
	if response.ErrorMessage != nil {
		return nil, errors.New(*response.ErrorMessage)
	}

	// Unmarshal response
	var entity Entity
	err := json.Unmarshal(response.Body, &entity)
	if err != nil {
		return nil, errors.New(fmt.Sprint("Failed to unmarshall response. Error: %v\nBody:\n%v\n", err, string(response.Body)))
	}

	return &entity, nil
}