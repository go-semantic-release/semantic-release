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
	"fmt"
	"net/http"
)

// LogssService handles communication with the logs
// related methods of the Travis CI API.
type LogsService struct {
	client *Client
}

// Log represents a Travis CI job log
type Log struct {
	Id    uint   `json:"id,omitempty"`
	JobId uint   `json:"job_id,omitempty"`
	Type  string `json:"type,omitempty"`
	Body  string `json:"body,omitempty"`
}

// getLogResponse represents the response of a call
// to the Travis CI get log endpoint.
type getLogResponse struct {
	Log Log `json:"log,omitempty"`
}

// Get fetches a log based on the provided id.
//
// Travis CI API docs: http://docs.travis-ci.com/api/#logs
func (ls *LogsService) Get(logId uint) (*Log, *http.Response, error) {
	u, err := urlWithOptions(fmt.Sprintf("/logs/%d", logId), nil)
	if err != nil {
		return nil, nil, err
	}

	req, err := ls.client.NewRequest("GET", u, nil, nil)
	if err != nil {
		return nil, nil, err
	}

	var logResp getLogResponse
	resp, err := ls.client.Do(req, &logResp)
	if err != nil {
		return nil, resp, err
	}

	return &logResp.Log, resp, err
}

// Get a job's log based on it's provided id.
//
// Travis CI API docs: http://docs.travis-ci.com/api/#logs
func (ls *LogsService) GetByJob(jobId uint) (*Log, *http.Response, error) {
	u, err := urlWithOptions(fmt.Sprintf("/jobs/%d/log", jobId), nil)
	if err != nil {
		return nil, nil, err
	}

	req, err := ls.client.NewRequest("GET", u, nil, nil)
	if err != nil {
		return nil, nil, err
	}

	var plainText bytes.Buffer
	resp, err := ls.client.Do(req, &plainText)
	if err != nil {
		return nil, resp, err
	}

	return &Log{JobId: jobId, Body: plainText.String()}, resp, err
}
