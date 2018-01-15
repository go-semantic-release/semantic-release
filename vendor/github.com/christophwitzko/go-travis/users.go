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

// UsersService handles communication with the users
// related methods of the Travis CI API.
type UsersService struct {
	client *Client
}

// User represents a Travis CI user.
type User struct {
	Id            uint   `json:"id,omitempty"`
	Name          string `json:"name,omitempty"`
	Login         string `json:"commit_id,omitempty"`
	Email         string `json:"email,omitempty"`
	GravatarId    string `json:"gravatar_id,omitempty"`
	IsSyncing     bool   `json:"is_syncing,omitempty"`
	SyncedAt      string `json:"synced_at,omitempty"`
	CorrectScopes bool   `json:"correct_scopes,omitempty"`
	CreatedAt     string `json:"created_at,omitempty"`
}

type getUserResponse struct {
	User User `json:"user"`
}

// GetAuthenticated fetches the currently authenticated user from Travis CI API.
// This request always needs to be authenticated.
//
// Travis CI API docs: http://docs.travis-ci.com/api/#users
func (us *UsersService) GetAuthenticated() (*User, *http.Response, error) {
	u, err := urlWithOptions("/users/", nil)
	if err != nil {
		return nil, nil, err
	}

	req, err := us.client.NewRequest("GET", u, nil, nil)
	if err != nil {
		return nil, nil, err
	}

	var userResp getUserResponse
	resp, err := us.client.Do(req, &userResp)
	if err != nil {
		return nil, resp, err
	}

	return &userResp.User, resp, err
}

// Get fetches the user with the provided id from the Travis CI API.
//
// Travis CI API docs: http://docs.travis-ci.com/api/#users
func (us *UsersService) Get(userId uint) (*User, *http.Response, error) {
	u, err := urlWithOptions(fmt.Sprintf("/users/%d", userId), nil)
	if err != nil {
		return nil, nil, err
	}

	req, err := us.client.NewRequest("GET", u, nil, nil)
	if err != nil {
		return nil, nil, err
	}

	var userResp getUserResponse
	resp, err := us.client.Do(req, &userResp)
	if err != nil {
		return nil, resp, err
	}

	return &userResp.User, resp, err
}

// Sync triggers a new sync with GitHub.
// Might return status 409 if the user is currently syncing.
// This request always needs to be authenticated.
//
// Travis CI API docs: http://docs.travis-ci.com/api/#users
func (us *UsersService) Sync() (*http.Response, error) {
	u, err := urlWithOptions("/users/sync", nil)
	if err != nil {
		return nil, err
	}

	req, err := us.client.NewRequest("POST", u, nil, nil)
	if err != nil {
		return nil, err
	}

	resp, err := us.client.Do(req, nil)
	if err != nil {
		return nil, err
	}

	return resp, err
}
