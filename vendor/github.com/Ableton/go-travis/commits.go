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

// CommitsService handles communication with the commits
// related parts of the Travis CI API. As commits are not directly
// exposed as an endpoint, most of this service methods will fetch
// commits through the builds or jobs endpoint.
type CommitsService struct {
	client *Client
}

// Commit represents a VCS commit as seen by Travis CI
type Commit struct {
	Id             uint   `json:"id,omitempty"`
	Sha            string `json:"sha,omitempty"`
	Branch         string `json:"branch,omitempty"`
	Message        string `json:"message,omitempty"`
	CommittedAt    string `json:"committed_at,omitempty"`
	AuthorName     string `json:"author_name,omitempty"`
	AuthorEmail    string `json:"author_email,omitempty"`
	CommitterName  string `json:"committer_name,omitempty"`
	CommitterEmail string `json:"committer_email,omitempty"`
	CompareUrl     string `json:"compare_url,omitempty"`
}

// Get fetches the commit that triggered a build based on the build id.
//
// Travis CI API docs: http://docs.travis-ci.com/api/#builds
func (cs *CommitsService) GetFromBuild(buildId uint) (*Commit, *http.Response, error) {
	u, err := urlWithOptions(fmt.Sprintf("/builds/%d", buildId), nil)
	if err != nil {
		return nil, nil, err
	}

	req, err := cs.client.NewRequest("GET", u, nil, nil)
	if err != nil {
		return nil, nil, err
	}

	var buildResp getBuildResponse
	resp, err := cs.client.Do(req, &buildResp)
	if err != nil {
		return nil, resp, err
	}

	return &buildResp.Commit, resp, nil
}

// List last commits attached to a repository builds.
//
// Travis CI API docs: http://docs.travis-ci.com/api/#builds
func (cs *CommitsService) ListFromRepository(repositorySlug string) ([]Commit, *http.Response, error) {
	u, err := urlWithOptions(fmt.Sprintf("/repos/%s/builds", repositorySlug), nil)
	if err != nil {
		return nil, nil, err
	}

	req, err := cs.client.NewRequest("GET", u, nil, nil)
	if err != nil {
		return nil, nil, err
	}

	var buildsResp listBuildsResponse
	resp, err := cs.client.Do(req, &buildsResp)
	if err != nil {
		return nil, resp, err
	}

	return buildsResp.Commits, resp, nil
}
