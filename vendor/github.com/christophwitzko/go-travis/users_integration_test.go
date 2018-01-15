// Copyright (c) 2015 Ableton AG, Berlin. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build integration

package travis

import (
	"testing"
	"time"
)

func TestUsersService_GetAuthenticated(t *testing.T) {
	t.Parallel()

	if auth := integrationClient.IsAuthenticated(); !auth {
		t.Skip("test client is unauthenticated. skipping.")
	}

	user, _, err := integrationClient.Users.GetAuthenticated()
	ok(t, err)

	assert(
		t,
		user != nil,
		"UsersService.GetAuthenticated returned nil user",
	)
}

func TestUsersService_Get(t *testing.T) {
	t.Parallel()

	if auth := integrationClient.IsAuthenticated(); !auth {
		t.Skip("test client is unauthenticated. skipping.")
	}

	authenticatedUser, _, err := integrationClient.Users.GetAuthenticated()
	userId := authenticatedUser.Id

	user, _, err := integrationClient.Users.Get(userId)
	ok(t, err)

	assert(
		t,
		authenticatedUser != nil,
		"UsersService.Get returned nil user",
	)

	assert(
		t,
		user.Id == userId,
		"UsersService.Get returned user with id %d; expected %d", user.Id, userId,
	)
}

func TestUsersService_Sync(t *testing.T) {
	t.Parallel()

	if auth := integrationClient.IsAuthenticated(); !auth {
		t.Skip("test client is unauthenticated. skipping.")
	}

	userNow, _, err := integrationClient.Users.GetAuthenticated()
	_, err = integrationClient.Users.Sync()
	ok(t, err)
	time.Sleep(5 * time.Second)

	userThen, _, err := integrationClient.Users.GetAuthenticated()
	if userThen != nil {
		// User might have been destroyed since last sync
		assert(
			t,
			userThen.SyncedAt > userNow.SyncedAt,
			"UsersService.Sync does not have updated the user synced_at marker",
		)
	}
}
