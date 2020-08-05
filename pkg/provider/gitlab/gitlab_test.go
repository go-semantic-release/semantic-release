package gitlab

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strconv"
	"testing"

	"github.com/Masterminds/semver"
	"github.com/go-semantic-release/semantic-release/pkg/provider"
	"github.com/stretchr/testify/require"
	"github.com/xanzy/go-gitlab"
)

func TestNewGitlabRepository(t *testing.T) {
	require := require.New(t)

	repo, err := NewRepository(context.TODO(), "", "", "", "")
	require.Nil(repo)
	require.EqualError(err, "project id is required")

	repo, err = NewRepository(context.TODO(), "", "token", "", "1")
	require.NotNil(repo)
	require.NoError(err)

	repo, err = NewRepository(context.TODO(), "https://mygitlab.com", "token", "", "1")
	require.NotNil(repo)
	require.NoError(err)
	require.Equal("https://mygitlab.com/api/v4/", repo.client.BaseURL().String(), "invalid custom instance initialization")
}

func createGitlabCommit(sha, message string) *gitlab.Commit {
	return &gitlab.Commit{ID: sha, Message: message}
}

func createGitlabTag(name, sha string) *gitlab.Tag {
	return &gitlab.Tag{Name: name, Commit: &gitlab.Commit{
		ID: sha,
	}}
}

var (
	GITLAB_PROJECT_ID    = 12324322
	GITLAB_DEFAULTBRANCH = "master"
	GITLAB_PROJECT       = gitlab.Project{DefaultBranch: GITLAB_DEFAULTBRANCH, Visibility: gitlab.PrivateVisibility, ID: GITLAB_PROJECT_ID}
	GITLAB_COMMITS       = []*gitlab.Commit{
		createGitlabCommit("abcd", "feat(app): new feature"),
		createGitlabCommit("dcba", "Fix: bug"),
		createGitlabCommit("cdba", "Initial commit"),
		createGitlabCommit("efcd", "chore: break\nBREAKING CHANGE: breaks everything"),
	}
	GITLAB_TAGS = []*gitlab.Tag{
		createGitlabTag("test-tag", "deadbeef"),
		createGitlabTag("v1.0.0", "deadbeef"),
		createGitlabTag("v2.0.0", "deadbeef"),
		createGitlabTag("v2.1.0-beta", "deadbeef"),
		createGitlabTag("v3.0.0-beta.2", "deadbeef"),
		createGitlabTag("v3.0.0-beta.1", "deadbeef"),
		createGitlabTag("2020.04.19", "deadbeef"),
	}
)

//nolint:errcheck
func GitlabHandler(w http.ResponseWriter, r *http.Request) {
	// Rate Limit headers
	if r.Method == "GET" && r.URL.Path == "/api/v4/" {
		json.NewEncoder(w).Encode(struct{}{})
		return
	}

	if r.Header.Get("PRIVATE-TOKEN") == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if r.Method == "GET" && r.URL.Path == fmt.Sprintf("/api/v4/projects/%d", GITLAB_PROJECT_ID) {
		json.NewEncoder(w).Encode(GITLAB_PROJECT)
		return
	}

	if r.Method == "GET" && r.URL.Path == fmt.Sprintf("/api/v4/projects/%d/repository/commits", GITLAB_PROJECT_ID) {
		json.NewEncoder(w).Encode(GITLAB_COMMITS)
		return
	}

	if r.Method == "GET" && r.URL.Path == fmt.Sprintf("/api/v4/projects/%d/repository/tags", GITLAB_PROJECT_ID) {
		json.NewEncoder(w).Encode(GITLAB_TAGS)
		return
	}

	if r.Method == "POST" && r.URL.Path == fmt.Sprintf("/api/v4/projects/%d/releases", GITLAB_PROJECT_ID) {
		var data map[string]string
		json.NewDecoder(r.Body).Decode(&data)
		r.Body.Close()
		if data["tag_name"] != "v2.0.0" {
			http.Error(w, "invalid tag name", http.StatusBadRequest)
			return
		}
		fmt.Fprint(w, "{}")
		return
	}

	http.Error(w, "invalid route", http.StatusNotImplemented)
}

func getNewGitlabTestRepo(t *testing.T) (*GitLabRepository, *httptest.Server) {
	ts := httptest.NewServer(http.HandlerFunc(GitlabHandler))
	repo, err := NewRepository(context.TODO(), ts.URL, "gitlab-examples-ci", "", strconv.Itoa(GITLAB_PROJECT_ID))
	require.NoError(t, err)

	return repo, ts
}

func TestGitlabGetInfo(t *testing.T) {
	repo, ts := getNewGitlabTestRepo(t)
	defer ts.Close()
	repoInfo, err := repo.GetInfo()
	require.NoError(t, err)
	require.Equal(t, GITLAB_DEFAULTBRANCH, repoInfo.DefaultBranch)
	require.True(t, repoInfo.Private)
}

func TestGitlabGetCommits(t *testing.T) {
	repo, ts := getNewGitlabTestRepo(t)
	defer ts.Close()
	commits, err := repo.GetCommits("")
	require.NoError(t, err)
	require.Len(t, commits, 4)

	for i, c := range commits {
		require.Equal(t, c.SHA, GITLAB_COMMITS[i].ID)
		require.Equal(t, c.RawMessage, GITLAB_COMMITS[i].Message)
	}
}

func TestGitlabGetReleases(t *testing.T) {
	repo, ts := getNewGitlabTestRepo(t)
	defer ts.Close()

	testCases := []struct {
		vrange          string
		re              *regexp.Regexp
		expectedSHA     string
		expectedVersion string
	}{
		{"", nil, "deadbeef", "2020.4.19"},
		{"", regexp.MustCompile("^v[0-9]*"), "deadbeef", "2.0.0"},
		{"2-beta", nil, "deadbeef", "2.1.0-beta"},
		{"3-beta", nil, "deadbeef", "3.0.0-beta.2"},
		{"4-beta", nil, "deadbeef", "4.0.0-beta"},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("VersionRange: %s, RE: %s", tc.vrange, tc.re), func(t *testing.T) {
			releases, err := repo.GetReleases(tc.re)
			require.NoError(t, err)
			release, err := releases.GetLatestRelease(tc.vrange)
			require.NoError(t, err)
			require.Equal(t, tc.expectedSHA, release.SHA)
			require.Equal(t, tc.expectedVersion, release.Version)
		})
	}
}

func TestGitlabCreateRelease(t *testing.T) {
	repo, ts := getNewGitlabTestRepo(t)
	defer ts.Close()
	newVersion := semver.MustParse("2.0.0")
	err := repo.CreateRelease(&provider.RepositoryRelease{NewVersion: newVersion, SHA: "deadbeef"})
	require.NoError(t, err)
}
