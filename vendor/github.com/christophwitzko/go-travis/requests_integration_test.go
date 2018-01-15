// Copyright (c) 2015 Ableton AG, Berlin. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build integration

package travis

import "testing"

func TestRequestsService_ListFromRepository_without_options(t *testing.T) {
	t.Parallel()

	_, _, _, err := integrationClient.Requests.ListFromRepository(integrationRepo, nil)
	ok(t, err)
}

// This test is commented on purpose as the Travis CI Api returns a 500 error code
// at the moment
// func TestRequestsService_ListFromRepository_with_options(t *testing.T) {
// 	t.Parallel()

// 	opt := &RequestsListOptions{Limit: 1}
// 	requests, _, _, err := integrationClient.Requests.ListFromRepository(integrationRepo, opt)
// 	ok(t, err)

// 	assert(
// 		t,
// 		len(requests) == 1,
// 		"Requests.List returned %d requests; expected %d", len(requests), 1,
// 	)
// }

func TestRequestsService_Get(t *testing.T) {
	t.Parallel()

	requests, _, _, err := integrationClient.Requests.ListFromRepository(integrationRepo, nil)
	if requests == nil || len(requests) == 0 {
		t.Skip("No requests found for the provided integration repo. skipping test")
	}
	requestId := requests[0].Id

	request, _, _, err := integrationClient.Requests.Get(requestId)
	ok(t, err)

	assert(
		t,
		request.Id == requestId,
		"Requests.Get returned request with id %d; expected %d", requestId,
	)
}
