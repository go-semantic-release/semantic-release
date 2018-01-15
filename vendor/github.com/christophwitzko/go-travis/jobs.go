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
	"fmt"
	"net/http"

	"github.com/fatih/structs"
)

// JobsService handles communication with the jobs
// related methods of the Travis CI API.
type JobsService struct {
	client *Client
}

// Job represents a Travis CI job
type Job struct {
	Id           uint   `json:"id,omitempty"`
	BuildId      uint   `json:"build_id,omitempty"`
	RepositoryId uint   `json:"repository_id,omitempty"`
	CommitId     uint   `json:"commit_id,omitempty"`
	LogId        uint   `json:"log_id,omitempty"`
	Number       string `json:"number,omitempty"`
	// Config        Config `json:"config,omitempty"`
	State         string `json:"state,omitempty"`
	StartedAt     string `json:"started_at,omitempty"`
	FinishedAt    string `json:"finished_at,omitempty"`
	Duration      uint   `json:"duration,omitempty"`
	Queue         string `json:"queue,omitempty"`
	AllowFailure  bool   `json:"allow_failure,omitempty"`
	AnnotationIds []uint `json:"annotation_ids,omitempty"`
}

type findJobsResponse struct {
	Jobs []Job `json:"jobs"`
}

// getJobResponse represents the response of a call
// to the Travis CI get build endpoint.
type getJobResponse struct {
	Job Job `json:"job"`
}

// JobListOptions specifies the optional parameters to the
// JobsService.List method. You need to provide exactly one
// of the below attribute. If you provide State or Queue, a
// maximum of 250 jobs will be returned.
type JobFindOptions struct {
	ListOptions

	// List of job ids
	Ids []uint `url:"ids,omitempty"`

	// Job state to filter by
	State string `url:"state,omitempty"`

	// Job queue to filter by
	Queue string `url:"queue,omitempty"`
}

// IsValid asserts the JobFindOptions instance has one
// and only one value set to a non-zero value.
//
// This method is particularly useful to check a JobFindOptions
// instance before passing it to JobsService.Find method.
func (jfo *JobFindOptions) IsValid() bool {
	s := structs.New(jfo)
	f := s.Fields()

	nonZeroValues := 0

	for _, field := range f {
		if !field.IsZero() {
			nonZeroValues += 1
		}
	}

	return nonZeroValues == 0 || nonZeroValues == 1
}

// Get fetches job with the provided id.
//
// Travis CI API docs: http://docs.travis-ci.com/api/#jobs
func (js *JobsService) Get(id uint) (*Job, *http.Response, error) {
	u, err := urlWithOptions(fmt.Sprintf("/jobs/%d", id), nil)
	if err != nil {
		return nil, nil, err
	}

	req, err := js.client.NewRequest("GET", u, nil, nil)
	if err != nil {
		return nil, nil, err
	}

	var jobResp getJobResponse
	resp, err := js.client.Do(req, &jobResp)
	if err != nil {
		return nil, resp, err
	}

	return &jobResp.Job, resp, err
}

// ListByBuild retrieve a build jobs from its provided id.
//
// Travis CI API docs: http://docs.travis-ci.com/api/#jobs
func (js *JobsService) ListFromBuild(buildId uint) ([]Job, *http.Response, error) {
	u, err := urlWithOptions(fmt.Sprintf("/builds/%d", buildId), nil)
	if err != nil {
		return nil, nil, err
	}

	req, err := js.client.NewRequest("GET", u, nil, nil)
	if err != nil {
		return nil, nil, err
	}

	var buildResp getBuildResponse
	resp, err := js.client.Do(req, &buildResp)
	if err != nil {
		return nil, resp, err
	}

	return buildResp.Jobs, resp, err
}

// Find jobs using the provided options.
// You need to provide exactly one of the opt fields value.
// If you provide State or Queue, a maximum of 250 jobs will be returned.
//
// Travis CI API docs: http://docs.travis-ci.com/api/#jobs
func (js *JobsService) Find(opt *JobFindOptions) ([]Job, *http.Response, error) {
	if opt != nil && !opt.IsValid() {
		return nil, nil, fmt.Errorf(
			"More than one value set in provided JobFindOptions instance",
		)
	}

	u, err := urlWithOptions("/jobs", opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := js.client.NewRequest("GET", u, nil, nil)
	if err != nil {
		return nil, nil, err
	}

	var jobsResp findJobsResponse
	resp, err := js.client.Do(req, &jobsResp)
	if err != nil {
		return nil, resp, err
	}

	return jobsResp.Jobs, resp, err
}

// Cancel job with the provided id.
//
// Travis CI API docs: http://docs.travis-ci.com/api/#jobs
func (js *JobsService) Cancel(id uint) (*http.Response, error) {
	u, err := urlWithOptions(fmt.Sprintf("/jobs/%d/cancel", id), nil)
	if err != nil {
		return nil, err
	}

	req, err := js.client.NewRequest("POST", u, nil, nil)
	if err != nil {
		return nil, err
	}

	resp, err := js.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, err
}

// Restart job with the provided id.
//
// Travis CI API docs: http://docs.travis-ci.com/api/#jobs
func (js *JobsService) Restart(id uint) (*http.Response, error) {
	u, err := urlWithOptions(fmt.Sprintf("/jobs/%d/restart", id), nil)
	if err != nil {
		return nil, err
	}

	req, err := js.client.NewRequest("POST", u, nil, nil)
	if err != nil {
		return nil, err
	}

	resp, err := js.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, err
}
