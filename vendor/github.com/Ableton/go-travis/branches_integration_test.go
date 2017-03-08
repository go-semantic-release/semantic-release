// Copyright (c) 2015 Ableton AG, Berlin. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build integration

package travis

import "testing"

func TestBranchesService_ListFromRepository(t *testing.T) {
	t.Parallel()

	_, _, err := integrationClient.Branches.ListFromRepository(integrationRepo)
	ok(t, err)
}

// func TestBranchesService_Get(t *testing.T) {
// 	t.Parallel()

// 	branches, _, err := integrationClient.Branches.ListFromRepository(integrationRepo)
// 	branchId := branches[0].Id

// 	branch, _, err := integrationClient.Branches.Get(integrationRepo, branchId)
// 	ok(t, err)

// 	assert(
// 		t,
// 		branch.Id == branchId,
// 		"Branches.Get return branch with Id %d; expected %d", branch.Id, branchId,
// 	)
// }
