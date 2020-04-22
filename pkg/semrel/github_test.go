package semrel

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"testing"

	"github.com/Masterminds/semver"
	"github.com/google/go-github/v30/github"
)

func TestNewGithubRepository(t *testing.T) {
	repo, err := NewGithubRepository(context.TODO(), "", "", "")
	if repo != nil || err == nil {
		t.Fatal("invalid initialization")
	}
	repo, err = NewGithubRepository(context.TODO(), "", "owner/test-repo", "token")
	if repo == nil || err != nil {
		t.Fatal("invalid initialization")
	}
	repo, err = NewGithubRepository(context.TODO(), "github.enterprise", "owner/test-repo", "token")
	if repo.Client.BaseURL.Host != "github.enterprise" || err != nil {
		t.Fatal("invalid enterprise initialization")
	}
}

func createCommit(sha, message string) *github.RepositoryCommit {
	return &github.RepositoryCommit{SHA: &sha, Commit: &github.Commit{Message: &message}}
}

func createRef(ref, sha string) *github.Reference {
	return &github.Reference{Ref: &ref, Object: &github.GitObject{SHA: &sha}}
}

var (
	GITHUB_REPO_PRIVATE  = true
	GITHUB_DEFAULTBRANCH = "master"
	GITHUB_REPO          = github.Repository{DefaultBranch: &GITHUB_DEFAULTBRANCH, Private: &GITHUB_REPO_PRIVATE}
	GITHUB_COMMITS       = []*github.RepositoryCommit{
		createCommit("abcd", "feat(app): new feature"),
		createCommit("dcba", "Fix: bug"),
		createCommit("cdba", "Initial commit"),
		createCommit("efcd", "chore: break\nBREAKING CHANGE: breaks everything"),
	}
	GITHUB_TAGS = []*github.Reference{
		createRef("refs/tags/test-tag", "deadbeef"),
		createRef("refs/tags/v1.0.0", "deadbeef"),
		createRef("refs/tags/v2.0.0", "deadbeef"),
		createRef("refs/tags/v2.1.0-beta", "deadbeef"),
		createRef("refs/tags/v3.0.0-beta.2", "deadbeef"),
		createRef("refs/tags/v3.0.0-beta.1", "deadbeef"),
		createRef("refs/tags/2020.04.19", "deadbeef"),
	}
)

func githubHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Authorization") != "Bearer token" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	if r.Method == "GET" && r.URL.Path == "/repos/owner/test-repo" {
		json.NewEncoder(w).Encode(GITHUB_REPO)
		return
	}
	if r.Method == "GET" && r.URL.Path == "/repos/owner/test-repo/commits" {
		json.NewEncoder(w).Encode(GITHUB_COMMITS)
		return
	}
	if r.Method == "GET" && r.URL.Path == "/repos/owner/test-repo/git/refs/tags" {
		json.NewEncoder(w).Encode(GITHUB_TAGS)
		return
	}
	if r.Method == "POST" && r.URL.Path == "/repos/owner/test-repo/git/refs" {
		var data map[string]string
		json.NewDecoder(r.Body).Decode(&data)
		r.Body.Close()
		if data["sha"] != "deadbeef" || data["ref"] != "refs/tags/v2.0.0" {
			http.Error(w, "invalid sha or ref", http.StatusBadRequest)
			return
		}
		fmt.Fprint(w, "{}")
		return
	}
	if r.Method == "POST" && r.URL.Path == "/repos/owner/test-repo/releases" {
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

func getNewTestRepo(t *testing.T) (*GithubRepository, *httptest.Server) {
	repo, err := NewGithubRepository(context.TODO(), "", "owner/test-repo", "token")
	if err != nil {
		t.Fatal(err)
		return nil, nil
	}
	ts := httptest.NewServer(http.HandlerFunc(githubHandler))
	repo.Client.BaseURL, _ = url.Parse(ts.URL + "/")
	return repo, ts
}

func TestGithubGetInfo(t *testing.T) {
	repo, ts := getNewTestRepo(t)
	defer ts.Close()
	defaultBranch, isPrivate, err := repo.GetInfo()
	if err != nil {
		t.Fatal(err)
	}
	if defaultBranch != GITHUB_DEFAULTBRANCH || !isPrivate {
		t.Fatal("invalid response")
	}
}

func compareCommit(c *Commit, t, s string, change Change) bool {
	if c.Type != t || c.Scope != s {
		return false
	}
	if c.Change.Major != change.Major ||
		c.Change.Minor != change.Minor ||
		c.Change.Patch != change.Patch {
		return false
	}
	return true
}

func TestGithubGetCommits(t *testing.T) {
	repo, ts := getNewTestRepo(t)
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

func TestGithubGetLatestRelease(t *testing.T) {
	repo, ts := getNewTestRepo(t)
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

func TestGithubCreateRelease(t *testing.T) {
	repo, ts := getNewTestRepo(t)
	defer ts.Close()
	newVersion, _ := semver.NewVersion("2.0.0")
	err := repo.CreateRelease("", newVersion, false, "", "deadbeef")
	if err != nil {
		t.Fatal(err)
	}
}
