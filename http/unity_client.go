package http

import (
	"fmt"
	"net/url"
	"errors"
	"time"
	"encoding/json"
)

type Request struct {
	Url *url.URL
	Method string
}

func (r *Request) Print() {
	fmt.Printf("Request:\n")
	fmt.Printf("    Url: %s\n", r.Url.String())
	fmt.Printf("    Method: %s\n", r.Method)
}

type Response struct {
	Code int
	ErrorMessage *string
	Body []byte
	ResponseTime time.Duration
}

func (r *Response) Print() {
	fmt.Printf("Response:\n")
	fmt.Printf("    HTTP response code : HTTP %d\n", r.Code)
	fmt.Printf("    Response time      : %v\n", r.ResponseTime)
	if r.ErrorMessage != nil {
		fmt.Printf("    Error message      : %s\n", *r.ErrorMessage)
	}
	fmt.Printf("    Response length    : %d bytes\n", len(r.Body))
}

type UnityError struct {
	Message string `json:"message"`
	Error string `json:"error"`

}

type UnityClient struct {
	verbose bool
	base *url.URL
	client *HttpClient
	path *url.URL
	values *url.Values
	method string
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
		method: "GET",
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

func (c *UnityClient) SetMethod(method string) {
	c.method = method
}

func (c *UnityClient) Reset() {
	c.path = nil
	c.values = nil
	c.method = "GET"
}

func (c *UnityClient) IssueRequest() (Response) {
	var response Response

	if c.path == nil {
		msg := fmt.Sprintf("No url path specified")
		response.ErrorMessage = &msg
		return response
	}

	var request Request
	request.Method = c.method
	request.Url = c.base.ResolveReference(c.path)
	if c.values != nil {
		request.Url.RawQuery = c.values.Encode();
	}

	//Print request information when in verbose mode
	if c.verbose {
		request.Print()
		fmt.Printf("\n")
	}

	//Issue request
	switch c.method {
	case "GET":
		response = c.client.Get(request.Url)
		break
	case "POST":
		response = c.client.Post(request.Url)
		break
	case "DELETE":
		response = c.client.Delete(request.Url)
		break
	}

	//Check for any unity error
	if len(response.Body) > 0 {
		var unityError UnityError
		err := json.Unmarshal(response.Body, &unityError)
		if err == nil && len(unityError.Error) > 0 {
			msg := fmt.Sprintf("%s: %s", unityError.Error, unityError.Message)
			response.ErrorMessage = &msg
		}
	}

	//Print response information when in verbose mode
	if c.verbose {
		response.Print()
		fmt.Printf("\n")
	}

	return response
}
