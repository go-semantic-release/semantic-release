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

// RepositoriesService handles communication with the builds
// related methods of the Travis CI API.
type RepositoriesService struct {
	client *Client
}

// Repository represents a Travis CI repository
type Repository struct {
	Id                  uint   `json:"id"`
	Slug                string `json:"slug"`
	Description         string `json:"description"`
	LastBuildId         uint   `json:"last_build_id"`
	LastBuildNumber     string `json:"last_build_number"`
	LastBuildState      string `json:"last_build_state"`
	LastBuildDuration   uint   `json:"last_build_duration"`
	LastBuildStartedAt  string `json:"last_build_started_at"`
	LastBuildFinishedAt string `json:"last_build_finished_at"`
	GithubLanguage      string `json:"github_language"`
}

// RepositoryListOptions specifies the optional parameters to the
// RepositoriesService.Findmethod.
type RepositoryListOptions struct {
	// list of repository ids to fetch, cannot be combined with other parameters
	Ids []uint `url:"ids,omitempty"`

	// filter by user that has access to it (github login)
	Member string `url:"member,omitempty"`

	// filter by owner name (first segment of slug)
	OwnerName string `url:"owner_name,omitempty"`

	// filter by slug
	Slug string `url:"slug,omitempty"`

	// filter by search term
	Search string `url:"search,omitempty"`

	// if true, will only return repositories that are enabled
	Active bool `url:"active,omitempty"`
}

// listRepositoriesResponse represents the response of a call
// to the Travis CI list builds endpoint.
type listRepositoriesResponse struct {
	Repositories []Repository `json:"repos"`
}

// getRepositoryResponse represents the response of a call
// to the Travis CI get repository endpoint.
type getRepositoryResponse struct {
	Repository Repository `json:"repo"`
}

// Find repositories using the provided options.
// If no options are provided, a list of repositories with recent
// activity for the authenticated user is returned.
//
// Travis CI API docs: http://docs.travis-ci.com/api/#repositories
func (rs *RepositoriesService) Find(opt *RepositoryListOptions) ([]Repository, *http.Response, error) {
	u, err := urlWithOptions("/repos", opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := rs.client.NewRequest("GET", u, nil, nil)
	if err != nil {
		return nil, nil, err
	}

	var reposResp listRepositoriesResponse
	resp, err := rs.client.Do(req, &reposResp)
	if err != nil {
		return nil, resp, err
	}

	return reposResp.Repositories, resp, err
}

// GetBySlug fetches a repository based on the provided slug.
//
// Travis CI API docs: http://docs.travis-ci.com/api/#repositories
func (rs *RepositoriesService) GetFromSlug(slug string) (*Repository, *http.Response, error) {
	u, err := urlWithOptions(fmt.Sprintf("/repos/%s", slug), nil)
	if err != nil {
		return nil, nil, err
	}

	req, err := rs.client.NewRequest("GET", u, nil, nil)
	if err != nil {
		return nil, nil, err
	}

	var repoResp getRepositoryResponse
	resp, err := rs.client.Do(req, &repoResp)
	if err != nil {
		return nil, resp, err
	}

	return &repoResp.Repository, resp, err
}

// Get fetches a repository based on the provided id.
//
// Travis CI API docs: http://docs.travis-ci.com/api/#repositories
func (rs *RepositoriesService) Get(id uint) (*Repository, *http.Response, error) {
	u, err := urlWithOptions(fmt.Sprintf("/repos/%d", id), nil)
	if err != nil {
		return nil, nil, err
	}

	req, err := rs.client.NewRequest("GET", u, nil, nil)
	if err != nil {
		return nil, nil, err
	}

	var repoResp getRepositoryResponse
	resp, err := rs.client.Do(req, &repoResp)
	if err != nil {
		return nil, resp, err
	}

	return &repoResp.Repository, resp, err
}
