package semrel

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
	gitlab "github.com/xanzy/go-gitlab"
)

func TestNewGitlabRepository(t *testing.T) {
	repo, err := NewGitLabRepository(context.TODO(), "", "", "", "", "")
	if repo != nil || err == nil {
		t.Fatal("invalid initialization")
	}
	repo, err = NewGitLabRepository(context.TODO(), "", "owner/test-repo", "token", "", "1")
	if repo == nil || err != nil || repo.Owner() != "owner" || repo.Repo() != "test-repo" {
		t.Fatal("invalid initialization")
	}
	repo, err = NewGitLabRepository(context.TODO(), "https://mygitlab.com", "owner/test-repo", "token", "", "1")
	if repo.client.BaseURL().String() != "https://mygitlab.com/api/v4/" || err != nil {
		t.Fatal("invalid custom instance initialization")
	}
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
	repo, err := NewGitLabRepository(context.TODO(), ts.URL, "gitlab-examples-ci", "token", "", strconv.Itoa(GITLAB_PROJECT_ID))
	if err != nil {
		t.Fatal(err)
		return nil, nil
	}

	return repo, ts
}

func TestGitlabGetInfo(t *testing.T) {
	repo, ts := getNewGitlabTestRepo(t)
	defer ts.Close()
	defaultBranch, isPrivate, err := repo.GetInfo()
	if err != nil {
		t.Fatal(err)
	}
	if defaultBranch != GITLAB_DEFAULTBRANCH || !isPrivate {
		t.Fatal("invalid response")
	}
}

func TestGitlabGetCommits(t *testing.T) {
	repo, ts := getNewGitlabTestRepo(t)
	defer ts.Close()
	commits, err := repo.GetCommits("")
	if err != nil {
		t.Fatal(err)
	}
	if len(commits) != 4 {
		t.Fatal("invalid response")
	}

	if !compareCommit(commits[0], "feat", "app", Change{false, true, false}) ||
		!compareCommit(commits[1], "fix", "", Change{false, false, true}) ||
		!compareCommit(commits[2], "", "", Change{false, false, false}) ||
		!compareCommit(commits[3], "chore", "", Change{true, false, false}) {
		t.Fatal("invalid commits")
	}
}

func TestGitlabGetLatestRelease(t *testing.T) {
	repo, ts := getNewGitlabTestRepo(t)
	defer ts.Close()
	release, err := repo.GetLatestRelease("", nil)
	if err != nil {
		t.Fatal(err)
	}
	if release.SHA != "deadbeef" || release.Version.String() != "2020.4.19" {
		t.Fatal("invalid tag")
	}

	re := regexp.MustCompile("^v[0-9]*")
	release, err = repo.GetLatestRelease("", re)
	if err != nil {
		t.Fatal(err)
	}
	if release.SHA != "deadbeef" || release.Version.String() != "2.0.0" {
		t.Fatal("invalid tag")
	}

	release, err = repo.GetLatestRelease("2-beta", nil)
	if err != nil {
		t.Fatal(err)
	}
	if release.SHA != "deadbeef" || release.Version.String() != "2.1.0-beta" {
		t.Fatal("invalid tag")
	}

	release, err = repo.GetLatestRelease("3-beta", nil)
	if err != nil {
		t.Fatal(err)
	}
	if release.SHA != "deadbeef" || release.Version.String() != "3.0.0-beta.2" {
		t.Fatal("invalid tag")
	}

	release, err = repo.GetLatestRelease("4-beta", nil)
	if err != nil {
		t.Fatal(err)
	}
	if release.SHA != "deadbeef" || release.Version.String() != "4.0.0-beta" {
		t.Fatal("invalid tag")
	}
}

func TestGitlabCreateRelease(t *testing.T) {
	repo, ts := getNewGitlabTestRepo(t)
	defer ts.Close()
	newVersion, _ := semver.NewVersion("2.0.0")
	err := repo.CreateRelease("", newVersion, false, "", "deadbeef")
	if err != nil {
		t.Fatal(err)
	}
}
