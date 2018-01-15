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
)

// BranchesService handles communication with the branches
// related methods of the Travis CI API.
type BranchesService struct {
	client *Client
}

// Branch represents a Travis CI build
type Branch struct {
	Id           uint   `json:"id,omitempty"`
	RepositoryId uint   `json:"repository_id,omitempty"`
	CommitId     uint   `json:"commit_id,omitempty"`
	Number       string `json:"number,omitempty"`
	// Config       Config `json:"config,omitempty"`
	State       string `json:"state,omitempty"`
	StartedAt   string `json:"started_at,omitempty"`
	FinishedAt  string `json:"finished_at,omitempty"`
	Duration    uint   `json:"duration,omitempty"`
	JobIds      []uint `json:"job_ids,omitempty"`
	PullRequest bool   `json:"pull_request,omitempty"`
}

// listBranchesResponse represents the response of a call
// to the Travis CI list branches endpoint.
type listBranchesResponse struct {
	Branches []Branch `json:"branches"`
}

// getBranchResponse represents the response of a call
// to the Travis CI get branch endpoint.
type getBranchResponse struct {
	Branch *Branch `json:"branch"`
}

// List the branches of a given repository.
//
// Travis CI API docs: http://docs.travis-ci.com/api/#builds
func (bs *BranchesService) ListFromRepository(slug string) ([]Branch, *http.Response, error) {
	u, err := urlWithOptions(fmt.Sprintf("/repos/%v/branches", slug), nil)
	if err != nil {
		return nil, nil, err
	}

	req, err := bs.client.NewRequest("GET", u, nil, nil)
	if err != nil {
		return nil, nil, err
	}

	var branchesResp listBranchesResponse
	resp, err := bs.client.Do(req, &branchesResp)
	if err != nil {
		return nil, resp, err
	}

	return branchesResp.Branches, resp, err
}

// Get fetches a branch based on the provided repository slug
// and it's id.
//
// Travis CI API docs: http://docs.travis-ci.com/api/#builds
func (bs *BranchesService) Get(repoSlug string, branchId uint) (*Branch, *http.Response, error) {
	u, err := urlWithOptions(fmt.Sprintf("/repos/%v/branches/%d", repoSlug, branchId), nil)
	if err != nil {
		return nil, nil, err
	}

	req, err := bs.client.NewRequest("GET", u, nil, nil)
	if err != nil {
		return nil, nil, err
	}

	var branchResp getBranchResponse
	resp, err := bs.client.Do(req, &branchResp)
	if err != nil {
		return nil, resp, err
	}

	return branchResp.Branch, resp, err
}

// Get fetches a branch based on the provided repository slug
// and its name.
//
// Travis CI API docs: http://docs.travis-ci.com/api/#builds
func (bs *BranchesService) GetFromSlug(repoSlug string, branchSlug string) (*Branch, *http.Response, error) {
	u, err := urlWithOptions(fmt.Sprintf("/repos/%v/branches/%v", repoSlug, branchSlug), nil)
	if err != nil {
		return nil, nil, err
	}

	req, err := bs.client.NewRequest("GET", u, nil, nil)
	if err != nil {
		return nil, nil, err
	}

	var branchResp getBranchResponse
	resp, err := bs.client.Do(req, &branchResp)
	if err != nil {
		return nil, resp, err
	}

	return branchResp.Branch, resp, err
}
