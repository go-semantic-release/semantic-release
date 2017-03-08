// Copyright (c) 2015 Ableton AG, Berlin. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package travis

import (
	"fmt"
	"testing"
)

func TestAuthenticate_UsingGithubToken_with_empty_token(t *testing.T) {
	as := &AuthenticationService{client: NewClient(TRAVIS_API_DEFAULT_URL, "")}

	_, _, err := as.UsingGithubToken("")
	notOk(t, err)
}

func TestAuthenticate_UsingTravisToken_with_empty_token(t *testing.T) {
	as := &AuthenticationService{client: NewClient(TRAVIS_API_DEFAULT_URL, "")}

	err := as.UsingTravisToken("")
	notOk(t, err)
}

func TestAuthenticate_UsingTravisToken_with_string_token(t *testing.T) {
	token := "abc123easyasdoremi"
	as := &AuthenticationService{client: NewClient(TRAVIS_API_DEFAULT_URL, token)}

	err := as.UsingTravisToken(token)
	ok(t, err)

	authHeader := as.client.Headers["Authorization"]
	assert(
		t,
		authHeader == fmt.Sprintf("token %s", token),
		fmt.Sprintf("travis token found in AuthenticationService.Headers: %s; expected %s", authHeader, token),
	)
}
