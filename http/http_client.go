package http

import (
	"net/url"
	net_http "net/http"
	"time"
	"crypto/tls"
	"io/ioutil"
	//"errors"
	"fmt"
	"bytes"
	"io"
)

type HttpClient struct {
	client *net_http.Client
	username *string
	password *string
	verbose bool
}

func InitHttpClient(verbose, insecure bool, username, password *string) (*HttpClient) {
	http_timeout := 30

	//Initialize client
	tr := &net_http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure},
	}
	httpTimeout := time.Duration(time.Duration(http_timeout) * time.Second)
	httpClient := net_http.Client{
		Timeout: httpTimeout,
		Transport: tr,
	}

	return &HttpClient {
		client: &httpClient,
		username: username,
		password: password,
		verbose: verbose,
	}
}

func (hc *HttpClient) Get(url *url.URL, body []byte) (Response) {
	return hc.Request(url, "GET", body)
}

func (hc *HttpClient) Put(url *url.URL, body []byte) (Response) {
	return hc.Request(url, "PUT", body)
}

func (hc *HttpClient) Post(url *url.URL, body []byte) (Response) {
	return hc.Request(url, "POST", body)
}

func (hc *HttpClient) Delete(url *url.URL, body []byte) (Response) {
	return hc.Request(url, "DELETE", body)
}

func (hc *HttpClient) Request(url *url.URL, method string, body []byte) (Response) {
	var response Response

	t1 := time.Now()
	var rdr io.Reader
	if body != nil && len(body) > 0 {
		rdr = bytes.NewReader(body)
	}

	//Build request
	req, err := net_http.NewRequest(method, url.String(), rdr)
	if body != nil && len(body) > 0 {
		req.Header.Add("Content-Type", "application/json")
	}
	if err != nil {
		msg := fmt.Sprintf("%v", err)
		response.ErrorMessage = &msg
		response.ResponseTime = time.Since(t1)
		return response
	}

	if hc.username != nil && hc.password != nil {
		req.SetBasicAuth(*hc.username, *hc.password)
	}

	//Execute request
	resp, err := hc.client.Do(req)
	if err != nil {
		msg := fmt.Sprintf("%v", err)
		response.ErrorMessage = &msg
		response.ResponseTime = time.Since(t1)
		return response
	}

	// Read body
	response.Body, err = ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		msg := fmt.Sprintf("%v", err)
		response.ErrorMessage = &msg
		response.ResponseTime = time.Since(t1)
		return response
	}

	//Handle response code
	response.Code = resp.StatusCode
	if resp.StatusCode == 403 {
		msg := "Not authorized"
		response.ErrorMessage = &msg
		response.ResponseTime = time.Since(t1)
		return response
	} else 	if resp.StatusCode == 404 {
		msg := "Not authorized"
		response.ErrorMessage = &msg
		response.ResponseTime = time.Since(t1)
		return response
	}

	response.ResponseTime = time.Since(t1)
	return response
}
