// Copyright (c) 2015 Ableton AG, Berlin. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build integration

package travis

import "testing"

func TestBuildsService_List_without_options(t *testing.T) {
	t.Parallel()

	builds, _, _, _, err := integrationClient.Builds.List(nil)
	ok(t, err)

	assert(
		t,
		len(builds) > 0,
		"Builds.List returned no builds",
	)
}

func TestBuildsService_List_with_options(t *testing.T) {
	t.Parallel()

	slug := integrationRepo
	number := "1"
	opt := &BuildListOptions{Slug: slug, Number: number}

	builds, _, _, _, err := integrationClient.Builds.List(opt)
	ok(t, err)

	assert(
		t,
		len(builds) == 1,
		"Builds.List returned no builds",
	)

	// Weirdly, the returned payload does not contained the slug
	// you are filtering on...
	// assert(
	// 	t,
	// 	builds[0].Slug == slug,
	// 	"Builds.List first returned build instance Slug = %s; expected %s", builds[0].Slug, slug,
	// )

	assert(
		t,
		builds[0].Number == "1",
		"Builds.List first returned build instance Number = %s; expected %s", builds[0].Number, number,
	)
}

func TestBuildsService_ListFromRepository_without_options(t *testing.T) {
	t.Parallel()

	_, _, _, _, err := integrationClient.Builds.ListFromRepository(integrationRepo, nil)
	ok(t, err)
}

func TestBuildsService_ListFromRepository_with_options(t *testing.T) {
	t.Parallel()

	opt := &BuildListOptions{EventType: "push"}

	builds, _, _, _, err := integrationClient.Builds.ListFromRepository(integrationRepo, opt)
	ok(t, err)

	if builds != nil {
		for _, b := range builds {
			assert(
				t,
				b.EventType == "push",
				"Builds.ListFromRepository returned builds with EventType != push",
			)
		}
	}
}

func TestBuildsService_Get(t *testing.T) {
	t.Parallel()

	// Fetch the reference repository first build in order
	// to have an existing build id to test against
	builds, _, _, _, err := integrationClient.Builds.ListFromRepository(integrationRepo, &BuildListOptions{Number: "1"})
	if builds == nil || len(builds) == 0 {
		t.Skip("No builds found for the provided integration repo. skipping test")
	}
	buildId := builds[0].Id

	build, _, _, _, err := integrationClient.Builds.Get(buildId)
	ok(t, err)

	assert(
		t,
		build.Id == buildId,
		"Builds.Get return build with Id %d; expected %d", build.Id, buildId,
	)

	assert(
		t,
		build.Number == "1",
		"Builds.Get return build with Number %s; expected %s", build.Number, "1",
	)
}
