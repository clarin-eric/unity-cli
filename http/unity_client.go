package http

import (
	"fmt"
	"net/url"
	"errors"
)

type Response struct {
	Code int
	ErrorMessage *string
	Body []byte
}

type UnityClient struct {
	verbose bool
	base *url.URL
	client *HttpClient
	path *url.URL
	values *url.Values
}

func GetNewUnityClient(base_url string, verbose bool, insecure bool, username, password string) (*UnityClient, error) {
	urlBase, err := url.Parse(base_url)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to parse base url. Error: %v", err))
	}
	path := "/rest-admin/v1/"
	urlPath, err := url.Parse(path)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to parse url path (%s). Error: %v", path, err))
	}

	if verbose {
		fmt.Printf("    API root url : %v\n", urlBase.ResolveReference(urlPath))
	}

	return &UnityClient{
		verbose: verbose,
		base: urlBase.ResolveReference(urlPath),
		client: InitHttpClient(verbose, insecure, &username, &password),
	}, nil
}

func (c *UnityClient) SetPathByString(_fmt string, args ...interface{}) (error) {
	_url := fmt.Sprintf(_fmt, args...)
	urlPath, err := url.Parse(_url)
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to parse url path. Error: %v", err))
	}
	c.SetPath(urlPath)
	return nil
}

func (c *UnityClient) SetPath(path *url.URL) {
	if c.path != nil {
		c.Reset()
	}
	c.path = path
}

func (c *UnityClient) AddQueryParam(key, value string) {
	if c.values == nil {
		v := url.Values{}
		c.values = &v
	}
	c.values.Add(key, value)
}

func (c *UnityClient) Reset() {
	c.path = nil
	c.values = nil
}

func (c *UnityClient) IssueRequest() (Response) {
	var response Response

	if c.path == nil {
		msg := fmt.Sprintf("No url path specified")
		response.ErrorMessage = &msg
		return response
	}
	u := c.base.ResolveReference(c.path)
	if c.values != nil {
		u.RawQuery = c.values.Encode();
	}

	//
	if c.verbose {
		fmt.Printf("Request:\n")
		fmt.Printf("    Url: %s\n", u.String())
		fmt.Printf("\n")
	}

	//Issue request
	response = c.client.Get(u)

	//Print response information when in verbose mode
	if c.verbose {
		fmt.Printf("Response:\n")
		fmt.Printf("    HTTP response code: %d\n", response.Code)
		if response.ErrorMessage != nil {
			fmt.Printf("    Error message: %s\n", response.ErrorMessage)
		}
		//response.Body
		fmt.Printf("\n")
	}

	return response
}
