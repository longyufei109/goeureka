package goeureka

//File  : request.go
//Author: Simon
//Describe: Defines all request for client request
//Date  : 2020-12-03 11:12:23

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// executeQuery request for eureka server
func executeQuery(requestAction RequestAction) ([]byte, error) {
	request := newHttpRequest(requestAction)

	var DefaultTransport http.RoundTripper = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	resp, err := DefaultTransport.RoundTrip(request)
	if err != nil {
		return []byte(nil), err
	} else {
		defer resp.Body.Close()
		responseBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return []byte(nil), err
		}
		return responseBody, nil
	}
}

// isDoHttpRequest return request eureka server is ok
func isDoHttpRequest(requestAction RequestAction) bool {
	request := newHttpRequest(requestAction)
	var DefaultTransport http.RoundTripper = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	resp, err := DefaultTransport.RoundTrip(request)
	if resp != nil && resp.StatusCode > 299 {
		defer resp.Body.Close()
		log.Printf("HTTP request failed with status (%d)", resp.StatusCode)
		return false
	} else if err != nil {
		log.Printf("HTTP request failed with error (%s)", err.Error())
		return false
	} else {
		return true
		defer resp.Body.Close()
	}
	return false
}

// newHttpRequest build request for eureka
func newHttpRequest(requestAction RequestAction) *http.Request {
	var (
		err     error
		request *http.Request
	)
	//log.Printf("DoHttpRequest URL(%v)",requestAction.Url)
	// load body and template for request
	if requestAction.Body != "" { // add body
		reader := strings.NewReader(requestAction.Body)
		request, err = http.NewRequest(requestAction.Method, requestAction.Url, reader)
	} else if requestAction.Template != "" { // add template
		reader := strings.NewReader(requestAction.Template)
		request, err = http.NewRequest(requestAction.Method, requestAction.Url, reader)
	} else {
		request, err = http.NewRequest(requestAction.Method, requestAction.Url, nil)
	}
	if err != nil {
		log.Fatal(err)
	}

	// Add headers for request
	request.Header = map[string][]string{
		"Accept":       {requestAction.Accept},
		"Content-Type": {requestAction.ContentType},
	}
	if requestAction.Username != "" && requestAction.Password != "" {
		request.Header.Set("Authorization", fmt.Sprintf("Basic %s",
			base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s",
				requestAction.Username, requestAction.Password)))))
	}
	return request
}
