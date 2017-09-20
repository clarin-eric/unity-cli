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
		fmt.Printf("Base: %v\n", urlBase.ResolveReference(urlPath))
	}

	return &UnityClient{
		verbose: verbose,
		base: urlBase.ResolveReference(urlPath),
		client: InitHttpClient(verbose, insecure, &username, &password),
	}, nil
}

func (c *UnityClient) SetPathByString(path string) (error) {
	urlPath, err := url.Parse(path)
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
	if c.verbose {
		fmt.Printf("\tUrl:%s\n", u.String())
	}

	//Issue request
	return c.client.Get(u)
}
