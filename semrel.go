package semrel

import (
	"context"
	"errors"
	"fmt"
	"log"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/Masterminds/semver"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

var (
	commitPattern   = regexp.MustCompile(`^(\w*)(?:\((.*)\))?\: (.*)$`)
	breakingPattern = regexp.MustCompile("BREAKING CHANGES?")
)

type Change struct {
	Major, Minor, Patch, NoChange bool
}

type Commit struct {
	SHA     string
	Raw     []string
	Type    string
	Scope   string
	Message string
	Change  Change
}
type Commits []*Commit

type Release struct {
	SHA     string
	Version *semver.Version
}

type Releases []*Release

func (r Releases) Len() int {
	return len(r)
}

func (r Releases) Less(i, j int) bool {
	return r[j].Version.LessThan(r[i].Version)
}

func (r Releases) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

type Repository struct {
	Owner  string
	Repo   string
	Ctx    context.Context
	Client *github.Client
}

func NewRepository(ctx context.Context, slug, token string) (*Repository, error) {
	if !strings.Contains(slug, "/") {
		return nil, errors.New("invalid slug")
	}
	repo := new(Repository)
	splited := strings.Split(slug, "/")
	repo.Owner = splited[0]
	repo.Repo = splited[1]
	repo.Ctx = ctx
	repo.Client = github.NewClient(oauth2.NewClient(ctx, oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)))
	return repo, nil
}

func (repo *Repository) GetInfo() (string, bool, error) {
	r, _, err := repo.Client.Repositories.Get(repo.Ctx, repo.Owner, repo.Repo)
	if err != nil {
		return "", false, err
	}
	return r.GetDefaultBranch(), r.GetPrivate(), nil
}

func parseCommit(commit *github.RepositoryCommit) *Commit {
	c := new(Commit)
	c.SHA = commit.GetSHA()
	c.Raw = strings.Split(commit.Commit.GetMessage(), "\n")
	found := commitPattern.FindAllStringSubmatch(c.Raw[0], -1)
	if len(found) < 1 {
		return c
	}
	c.Type = strings.ToLower(found[0][1])
	c.Scope = found[0][2]
	c.Message = found[0][3]
	c.Change = Change{
		Major: breakingPattern.MatchString(commit.Commit.GetMessage()),
		Minor: c.Type == "feat",
		Patch: c.Type == "fix",
	}
	return c
}

func (repo *Repository) GetCommits(branch string) ([]*Commit, error) {
	opts := &github.CommitsListOptions{
		SHA:         branch,
		ListOptions: github.ListOptions{PerPage: 100},
	}

	commits, _, err := repo.Client.Repositories.ListCommits(repo.Ctx, repo.Owner, repo.Repo, opts)

	if err != nil {
		return nil, err
	}

	ret := make([]*Commit, len(commits))

	for i, commit := range commits {
		ret[i] = parseCommit(commit)
	}

	return ret, nil
}

func (repo *Repository) GetLatestRelease(vrange string, prerelease string) (*Release, error) {
	allReleases := make(Releases, 0)
	opts := &github.ReferenceListOptions{"tags", github.ListOptions{PerPage: 100}}
	for {
		refs, resp, err := repo.Client.Git.ListRefs(repo.Ctx, repo.Owner, repo.Repo, opts)
		if resp != nil && resp.StatusCode == 404 {
			return &Release{"", &semver.Version{}}, nil
		}
		if err != nil {
			return nil, err
		}
		for _, r := range refs {
			version, err := semver.NewVersion(strings.TrimPrefix(r.GetRef(), "refs/tags/"))
			if err != nil {
				continue
			}
			allReleases = append(allReleases, &Release{r.Object.GetSHA(), version})
		}
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}
	sort.Sort(allReleases)

	var lastRelease *Release
	for _, r := range allReleases {

		log.Println("Checking version: ", r.Version.String())

		if r.Version.Prerelease() == "" && lastRelease == nil {
			lastRelease = r

			// If there is no prerelease requested, its safe to stop here.
			if prerelease == "" {
				break
			}
		}

		prereleaseParts := strings.Split(r.Version.Prerelease(), ".")

		if prereleaseParts[0] == prerelease {
			// If the last production release is newer just stop here and go with the last production release version.
			if lastRelease != nil && r.Version.LessThan(lastRelease.Version) {
				break
			}

			if prerelease != "" {
				lastRelease = r
				break
			}
		}
	}

	if vrange == "" {
		if lastRelease != nil {
			return lastRelease, nil
		}
		return &Release{"", &semver.Version{}}, nil
	}

	constraint, err := semver.NewConstraint(vrange)
	if err != nil {
		return nil, err
	}
	for _, r := range allReleases {
		if constraint.Check(r.Version) {
			return r, nil
		}
	}

	nver, err := semver.NewVersion(vrange)
	if err != nil {
		return nil, err
	}

	splitPre := strings.SplitN(vrange, "-", 2)
	if len(splitPre) == 1 {
		return &Release{lastRelease.SHA, nver}, nil
	}

	npver, err := nver.SetPrerelease(splitPre[1])
	if err != nil {
		return nil, err
	}
	return &Release{lastRelease.SHA, &npver}, nil
}

func (repo *Repository) CreateRelease(commits []*Commit, latestRelease *Release, newVersion *semver.Version, branch string) error {
	tag := fmt.Sprintf("v%s", newVersion.String())
	changelog := GetChangelog(commits, latestRelease, newVersion)
	hasPrerelease := newVersion.Prerelease() != ""
	opts := &github.RepositoryRelease{
		TagName:         &tag,
		TargetCommitish: &branch,
		Body:            &changelog,
		Prerelease:      &hasPrerelease,
	}
	_, _, err := repo.Client.Repositories.CreateRelease(repo.Ctx, repo.Owner, repo.Repo, opts)
	if err != nil {
		return err
	}
	return nil
}

func CalculateChange(commits []*Commit, latestRelease *Release) Change {
	var change Change
	for _, commit := range commits {
		log.Println("Examining commit", commit.SHA)

		if latestRelease.SHA == commit.SHA {
			change.NoChange = true
			break
		}
		change.Major = change.Major || commit.Change.Major
		change.Minor = change.Minor || commit.Change.Minor
		change.Patch = change.Patch || commit.Change.Patch
	}
	return change
}

func ApplyChange(latestVersion *semver.Version, prerelease string, change Change) *semver.Version {
	if latestVersion.Major() == 0 {
		change.Major = true
	}
	if !change.Major && !change.Minor && !change.Patch {
		if change.NoChange {
			return latestVersion
		}
		return nil
	}

	preRel := latestVersion.Prerelease()
	preRelVer := strings.Split(preRel, ".")
	preRelLabel := preRelVer[0]

	newVersion := *latestVersion
	if preRelLabel == "" {
		switch {
		case change.Major:
			newVersion = latestVersion.IncMajor()
		case change.Minor:
			newVersion = latestVersion.IncMinor()
		case change.Patch:
			newVersion = latestVersion.IncPatch()
		}
	}

	if prerelease != "" && preRelVer[0] != prerelease {
		preRel = prerelease + ".1"
	} else {
		if len(preRelVer) > 1 {
			idx, err := strconv.ParseInt(preRelVer[1], 10, 32)
			if err != nil {
				idx = 0
			}
			preRel = fmt.Sprintf("%s.%d", preRelVer[0], idx+1)
		} else {
			preRel += ".1"
		}
	}

	newVersion, _ = newVersion.SetPrerelease(preRel)
	return &newVersion
}

func GetNewVersion(commits []*Commit, latestRelease *Release, prerelease string) *semver.Version {

	newVersion := ApplyChange(latestRelease.Version, prerelease, CalculateChange(commits, latestRelease))

	return newVersion
}

func trimSHA(sha string) string {
	if len(sha) < 9 {
		return sha
	}
	return sha[:8]
}

func formatCommit(c *Commit) string {
	ret := "* "
	if c.Scope != "" {
		ret += fmt.Sprintf("**%s:** ", c.Scope)
	}
	ret += fmt.Sprintf("%s (%s)\n", c.Message, trimSHA(c.SHA))
	return ret
}

var typeToText = map[string]string{
	"feat":     "Feature",
	"fix":      "Bug Fixes",
	"perf":     "Performance Improvements",
	"revert":   "Reverts",
	"docs":     "Documentation",
	"style":    "Styles",
	"refactor": "Code Refactoring",
	"test":     "Tests",
	"chore":    "Chores",
	"%%bc%%":   "Breaking Changes",
}

func getSortedKeys(m *map[string]string) []string {
	keys := make([]string, len(*m))
	i := 0
	for k := range *m {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	return keys
}

func GetChangelog(commits []*Commit, latestRelease *Release, newVersion *semver.Version) string {
	ret := fmt.Sprintf("## %s (%s)\n\n", newVersion.String(), time.Now().UTC().Format("2006-01-02"))
	typeScopeMap := make(map[string]string)
	for _, commit := range commits {
		if latestRelease.SHA == commit.SHA {
			break
		}
		if commit.Change.Major {
			typeScopeMap["%%bc%%"] += fmt.Sprintf("%s\n```%s\n```\n", formatCommit(commit), strings.Join(commit.Raw[1:], "\n"))
			continue
		}
		if commit.Type == "" {
			continue
		}
		typeScopeMap[commit.Type] += formatCommit(commit)
	}
	for _, t := range getSortedKeys(&typeScopeMap) {
		msg := typeScopeMap[t]
		typeName, found := typeToText[t]
		if !found {
			typeName = t
		}
		ret += fmt.Sprintf("#### %s\n\n%s\n", typeName, msg)
	}
	return ret
}
