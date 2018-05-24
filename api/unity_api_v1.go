package api

import (
	"fmt"
	"errors"
	"clarin/unity-cli/http"
	"encoding/json"
	"net/url"
	"github.com/magiconair/properties"
	"strings"
	"os"
)

type UnityApiV1 struct {
	verbose bool
	client *http.UnityClient
}

func GetUnityApiV1(base_url string, verbose bool, insecure bool, username, password string, useFile bool, file string) (*UnityApiV1, error) {
	if useFile {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			fmt.Printf("Specified file (%s) does not exist\n", file)
			os.Exit(1)
		}
		key := "UNITY_ADMIN_PASSWORD"
		p := properties.MustLoadFile(file, properties.UTF8)
		v := p.GetString(key, "test")
		v = strings.Trim(v, "\"")
		password = v
	}

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
	extended := true

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

	if extended {
		//Process entity groups
		groups, err := api.GetEntityGroups(entity.Id)
		if err != nil {
			return nil, fmt.Errorf("Failed to fetch groups for entity with id: %d. Error: %v\n", entity.Id, err)
		}
		entity.Groups = groups

		//fmt.Printf("Groups: %s\n", groups)
		//Process entities attribues
		var attributes []Attribute
		//group_path := "/"
		for _, group_path := range entity.Groups {
			attrs, err := api.GetAttributes(entity_id, group_path)
			if err != nil{
				fmt.Printf("Failed to fetch attributes. Grouppath: %s, entity id: %d. Error: %v\n", group_path, entity_id, err)
			} else{
				for _, attr := range attrs{
					if attr.Name != "sys:Credential:Password credential"{
						attributes = append(attributes, attr)
					}
				}
			}
			/*
			attrs, err = api.GetAttributes(entity_id, group_path)
			if err != nil{
				fmt.Printf("Failed to fetch attributes. Grouppath: %s, entity id: %d. Error: %v\n", group_path, entity_id, err)
			} else {
				for _, attr := range attrs {
					if attr.Name != "sys:Credential:Password credential" {
						attributes = append(attributes, attr)
					}
				}
			}
			*/
		}
		entity.Attributes = attributes
	}


	return &entity, nil
}

func (api *UnityApiV1) GetEntityGroups(entity_id int64) ([]string, error) {
	api.client.SetPathByString("./entity/%d/groups", entity_id)

	response := api.client.IssueRequest()
	if response.ErrorMessage != nil {
		return nil, errors.New(*response.ErrorMessage)
	}

	// Unmarshal response
	var groups []string
	err := json.Unmarshal(response.Body, &groups)
	if err != nil {
		return nil, errors.New(fmt.Sprint("Failed to unmarshall response. Error: %v\nBody:\n%v\n", err, string(response.Body)))
	}

	return groups, nil
}

/*
Return all members and subgroups of a given group.

Request:
	@Path("/group/{groupPath}")
	@GET

Response:
	{
  	"subGroups" : [ ],
  	"members" : [ 3 ]
	}
*/
func (api *UnityApiV1) GetGroup(group_path *string) (*Group, error) {
	//Build url
	api.client.SetPathByString("./group/%s", url.QueryEscape(*group_path))

	//Issue request
	response := api.client.IssueRequest()
	if response.ErrorMessage != nil {
		return nil, errors.New(*response.ErrorMessage)
	}

	// Unmarshal response
	var group Group
	err := json.Unmarshal(response.Body, &group)
	if err != nil {
		return nil, errors.New(fmt.Sprint("Failed to unmarshall response. Error: %v\nBody:\n%v\n", err, string(response.Body)))
	}
	group.Path = *group_path

	return &group, nil
}

func (api *UnityApiV1) GetGroupRecusive(group_path *string) (*Group, error) {
	group, err := api.GetGroup(group_path)
	if err != nil {
		return group, err
	}

	for _, sub_group_path := range group.SubGroupPaths {
		sub_group, err := api.GetGroupRecusive(&sub_group_path)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		} else {
			group.SubGroups = append(group.SubGroups, *sub_group)
		}
	}

	return group, nil
}

/**
@Path("/group/{groupPath}")
@POST

Creates a new group. The created group will be empty.
 */
func (api *UnityApiV1) CreateGroup(group_path *string) (error) {
	//Build url
	api.client.SetPathByString("./group/%s", url.QueryEscape(*group_path))
	api.client.SetMethod("POST")

	//Issue request
	response := api.client.IssueRequest()
	if response.ErrorMessage != nil {
		return errors.New(*response.ErrorMessage)
	}

	return nil
}

/**
@Path("/group/{groupPath}")
@QueryParam("recursive")
@DELETE

Removes a given group. The optional recursive query parameter can be used to enforce recursive removal (off by default).
*/
func (api *UnityApiV1) DeleteGroup(group_path *string, recursive *bool) (error) {
	//Build url
	api.client.SetPathByString("./group/%s", url.QueryEscape(*group_path))
	api.client.SetMethod("DELETE")
	api.client.AddQueryParam("recursive", fmt.Sprintf("%t", *recursive))

	//Issue request
	response := api.client.IssueRequest()
	if response.ErrorMessage != nil {
		return errors.New(*response.ErrorMessage)
	}

	return nil
}

/**
@Path("/registrationRequests")
@GET

Returns an array with all registration requests which are stored in the system.
 */
func (api *UnityApiV1) GetRegistrationRequests() ([]RegistrationRequest, error) {
	var requests []RegistrationRequest

	//Intialize API call
	api.client.SetPathByString("./registrationRequests/")
	api.client.SetMethod("GET")

	//Issue request
	response := api.client.IssueRequest()
	if response.ErrorMessage != nil {
		return requests, errors.New(*response.ErrorMessage)
	}

	// Unmarshal response
	err := json.Unmarshal(response.Body, &requests)
	if err != nil {
		return requests, errors.New(fmt.Sprint("Failed to unmarshall response. Error: %v\nBody:\n%v\n", err, string(response.Body)))
	}

	return requests, nil
}

/**
@Path("/entity/{entityId}/attributes")
@QueryParam("group")
@QueryParam("effective")
@QueryParam("identityType")
@GET
 */
func (api *UnityApiV1) GetAttributes(entity_id int64, group string) ([]Attribute, error) {
	var attributes []Attribute

	api.client.SetPathByString("./entity/%d/attributes", entity_id)
	api.client.SetMethod("GET")
	api.client.AddQueryParam("group", group)

	//Issue request
	response := api.client.IssueRequest()
	if response.ErrorMessage != nil {
		return attributes, errors.New(*response.ErrorMessage)
	}

	// Unmarshal response
	err := json.Unmarshal(response.Body, &attributes)
	if err != nil {
		return attributes, errors.New(fmt.Sprint("Failed to unmarshall response. Error: %v\nBody:\n%v\n", err, string(response.Body)))
	}

	return attributes, nil
}

func (api *UnityApiV1) UpdateAttribute(entity_id int64, attribute interface{}) (error) {
	//Build request
	api.client.SetPathByString("./entity/%d/attribute", entity_id)
	api.client.SetMethod("PUT")
	if err := api.client.SetMarshallBody(attribute); err != nil {
		return fmt.Errorf("Failed to marshall request body. Error: %s", err)
	}

	//Issue request
	response := api.client.IssueRequest()
	if response.ErrorMessage != nil {
		return errors.New(*response.ErrorMessage)
	}

	return nil
}

func (api *UnityApiV1) UpdateAttributes(entity_id int64, attributes []AttributeWithConfirmationData) (error) {
	//Build request
	api.client.SetPathByString("./entity/%d/attributes", entity_id)
	api.client.SetMethod("PUT")
	if err := api.client.SetMarshallBody(attributes); err != nil {
		return fmt.Errorf("Failed to marshall request body. Error: %s", err)
	}

	//Issue request
	response := api.client.IssueRequest()
	if response.ErrorMessage != nil {
		return errors.New(*response.ErrorMessage)
	}

	// Unmarshal response
	fmt.Printf("Response:\n%v\n", response.Body)
	return nil
}

func (api *UnityApiV1) RemoveAttribute(entity_id int64, name, identity_type, group string) (error) {
	api.client.SetPathByString("./entity/%d/attribute/%s", entity_id, url.QueryEscape(name))
	api.client.SetMethod("DELETE")
	//api.client.AddQueryParam("identityType",url.QueryEscape(identity_type))
	api.client.AddQueryParam("group", group) //Seems not to work url escaped

	//Issue request
	response := api.client.IssueRequest()
	if response.ErrorMessage != nil {
		return errors.New(*response.ErrorMessage)
	}
	return nil
}

/**
@Path("/registrationRequest/{requestId}")
@GET

Returns a registration request by its id.
 */
func (api *UnityApiV1) GetRegistrationRequest(id string) {

}

/**
 * @Path("/invitations")
 * @GET
 *
 * Returns a JSON array with all existing invitations.
 */
func (api *UnityApiV1) GetInvitations() {

}

/**
 * @Path("/invitation/{code}")
 * @GET
 *
 * Returns a JSON encoded invitation with the specified code.
 */
func (api *UnityApiV1) GetInvitation(code string) {

}

/**
@Path("invitation/{code}")
@DELETE

Removes an invitation with a specified code.

@Path("invitation/{code}/send")
@POST

Triggers sending a message with an invitation. The registration form of the invitation must have an invitation template defined, and the invitation must have contact address and channel set.

@Path("invitation")
@POST
@Consumes(MediaType.APPLICATION_JSON)
@Produces(MediaType.TEXT_PLAIN)

Creates a new invitation. Returned string is the unique code of the newly created invitation. Example invitation definition:

{
  "formId" : "exForm",
  "expiration" : 1454703788,
  "contactAddress" : "someAddr@example.com",
  "channelId" : "channelId",
  "identities" : {},
  "groupSelections" : {},
  "attributes" : {}
}
Syntax of prefilled parameters, can be seen in the result of retrieving an AdminUI defined invitation via the REST GET methods.
 */