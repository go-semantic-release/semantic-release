// Copyright (c) 2015 Ableton AG, Berlin. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// Fragments of this file have been copied from the go-github (https://github.com/google/go-github)
// project, and is therefore licensed under the following copyright:
// Copyright 2013 The go-github AUTHORS. All rights reserved.

package travis

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"reflect"
	"strconv"

	"github.com/google/go-querystring/query"
	"github.com/oleiade/reflections"
)

const (
	TRAVIS_USER_AGENT string = "go-travis/1.0.0"

	TRAVIS_REQUEST_ACCEPT_HEADER string = "application/vnd.travis-ci.2+json"
	TRAVIS_REQUEST_CONTENT_TYPE  string = "application/json"

	TRAVIS_API_DEFAULT_URL string = "https://api.travis-ci.org/"
	TRAVIS_API_PRO_URL     string = "https://api.travis-ci.com/"

	TRAVIS_RESPONSE_PER_PAGE uint64 = 25
)

// A Client manages communication with the Travis CI API.
type Client struct {
	// HTTP client used to communicate with the API
	client *http.Client

	// Headers to attach to every requests made with the client.
	// As a default, Headers will be used to provide Travis API authentication
	// token and other necessary headers.
	// However these could be updated per-request through a parameters.
	Headers map[string]string

	// Base URL for api requests. Defaults to the public Travis API, but
	// can be set to an alternative endpoint to use with Travis Pro or Enterprise.
	// BaseURL should alway be terminated by a slash.
	BaseURL *url.URL

	// User agent used when communicating with the Travis API
	UserAgent string

	// Services used to manipulate API entities
	Authentication *AuthenticationService
	Repositories   *RepositoriesService
	Builds         *BuildsService
	Jobs           *JobsService
	Branches       *BranchesService
	Logs           *LogsService
	Commits        *CommitsService
	Requests       *RequestsService
	Users          *UsersService
}

// NewClient returns a new Travis API client.
// If travisToken is not provided, the client can be authenticated at any time,
// using it's Authentication exposed service.
func NewClient(baseUrl string, travisToken string) *Client {
	bu, _ := url.Parse(baseUrl)
	bh := map[string]string{
		"Content-Type": TRAVIS_REQUEST_CONTENT_TYPE,
		"User-Agent":   TRAVIS_USER_AGENT,
		"Accept":       TRAVIS_REQUEST_ACCEPT_HEADER,
		"Host":         bu.Host,
	}

	c := &Client{
		client:    http.DefaultClient,
		Headers:   bh,
		BaseURL:   bu,
		UserAgent: TRAVIS_USER_AGENT,
	}

	c.Authentication = &AuthenticationService{client: c}
	c.Repositories = &RepositoriesService{client: c}
	c.Builds = &BuildsService{client: c}
	c.Jobs = &JobsService{client: c}
	c.Branches = &BranchesService{client: c}
	c.Logs = &LogsService{client: c}
	c.Commits = &CommitsService{client: c}
	c.Requests = &RequestsService{client: c}
	c.Users = &UsersService{client: c}

	if travisToken != "" {
		c.Authentication.UsingTravisToken(travisToken)
	}

	return c
}

// NewDefaultClient returns a new Travis API client bound to the public travis API.
// If travisToken is not provided, the client can be authenticated at any time,
// using it's Authentication exposed service.
func NewDefaultClient(travisToken string) *Client {
	return NewClient(TRAVIS_API_DEFAULT_URL, travisToken)
}

// NewRequest creates an API request. A relative URL can be provided in urlStr,
// in which case it is resolved relative to the BaseURL of the Client.
// Relative URLs should always be specified without a preceding slash. If
// specified, the value pointed to by body is JSON encoded and included as the
// request body. If specified, the map provided by headers will be used to udate
// request headers.
func (c *Client) NewRequest(method, urlStr string, body interface{}, headers map[string]string) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	u := c.BaseURL.ResolveReference(rel)

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	var h map[string]string = c.Headers
	if headers != nil {
		for k, v := range headers {
			h[k] = v
		}
	}

	for k, v := range h {
		req.Header.Set(k, v)
	}

	return req, nil
}

// Do sends an API request and returns the API response.  The API response is
// JSON decoded and stored in the value pointed to by v, or returned as an
// error if an API error has occurred.  If v implements the io.Writer
// interface, the raw response body will be written to v, without attempting to
// first decode it.
func (c *Client) Do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	err = checkResponse(resp)
	if err != nil {
		return resp, err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			io.Copy(w, resp.Body)
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
		}
	}

	return resp, err
}

// IsAuthenticated indicates if Authorization headers were
// found in Client.Headers mapping.
func (c *Client) IsAuthenticated() bool {
	authHeader, ok := c.Headers["Authorization"]

	if !ok || (ok && authHeader == "token ") {
		return false
	}

	return true
}

// ErrorResponse reports an error caused by an API request.
// ErrorResponse implemented the Error interface.
type ErrorResponse struct {
	// HTTP response that caused this error
	Response *http.Response

	// Error message produced by Travis API
	Message string `json:"error"`
}

func (er *ErrorResponse) Error() string {
	return fmt.Sprintf(
		"%v %v: %d %v",
		er.Response.Request.Method,
		er.Response.Request.URL.String(),
		er.Response.StatusCode,
		er.Message,
	)
}

// checkResponse checks the API response for errors; and returns them
// if present.
// A Response is considered an error if it has a status code outside the 2XX range.
func checkResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}

	errorResponse := &ErrorResponse{Response: r}

	data, err := ioutil.ReadAll(r.Body)
	if err == nil && data != nil {
		json.Unmarshal(data, errorResponse)
	}

	return errorResponse
}

func urlWithOptions(s string, opt interface{}) (string, error) {
	rv := reflect.ValueOf(opt)
	if rv.Kind() == reflect.Ptr && rv.IsNil() {
		return s, nil
	}

	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	qs, err := query.Values(opt)
	u.RawQuery = qs.Encode()

	return u.String(), nil
}

type Paginator interface {
	GetNextPage(interface{}) error
}

type ListOptions struct {
	AfterNumber uint `url:"after_number,omitempty"`
}

// GetNextPage provided a collection of resources (Builds or Jobs),
// will update the ListOptions to fetch the next resource page on next call.
func (into *ListOptions) GetNextPage(from interface{}) error {
	if reflect.TypeOf(from).Kind() != reflect.Slice {
		return fmt.Errorf("provided interface{} does not represent a slice")
	}

	slice := reflect.ValueOf(from)
	if slice.Len() == 0 {
		return fmt.Errorf("provided interface{} is a zero sized slice")
	}

	lastElem := slice.Index(slice.Len() - 1).Interface()
	has, _ := reflections.HasField(lastElem, "Number")
	if !has {
		return fmt.Errorf("last element of the provided slice does not have a Number attribute")
	}

	value, err := reflections.GetField(lastElem, "Number")
	if err != nil {
		return err
	}

	// We rely on travis sending us numbers representations here
	// so no real need to check for errors
	number, _ := strconv.ParseUint(value.(string), 10, 64)
	into.AfterNumber = uint(math.Max(float64(number), 0))

	return nil
}
