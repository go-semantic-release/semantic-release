package semrel

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/go-semantic-release/semantic-release/pkg/config"
)

var commitPattern = regexp.MustCompile(`^(\w*)(?:\((.*)\))?\: (.*)$`)
var breakingPattern = regexp.MustCompile("BREAKING CHANGES?")

type Change struct {
	Major, Minor, Patch bool
}

type Commit struct {
	SHA     string
	Raw     []string
	Type    string
	Scope   string
	Message string
	Change  Change
}

func NewCommit(sha, msg string) *Commit {
	c := new(Commit)
	c.SHA = sha
	c.Raw = strings.Split(msg, "\n")
	found := commitPattern.FindAllStringSubmatch(c.Raw[0], -1)
	if len(found) < 1 {
		return c
	}
	c.Type = strings.ToLower(found[0][1])
	c.Scope = found[0][2]
	c.Message = found[0][3]
	c.Change = Change{
		Major: breakingPattern.MatchString(msg),
		Minor: c.Type == "feat",
		Patch: c.Type == "fix",
	}
	return c
}

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

func (releases Releases) GetLatestRelease(vrange string) (*Release, error) {
	if len(releases) == 0 {
		return &Release{"", &semver.Version{}}, nil
	}

	sort.Sort(releases)

	var lastRelease *Release
	for _, r := range releases {
		if r.Version.Prerelease() == "" {
			lastRelease = r
			break
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
	for _, r := range releases {
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

func CalculateChange(commits []*Commit, latestRelease *Release) Change {
	var change Change
	for _, commit := range commits {
		if latestRelease.SHA == commit.SHA {
			break
		}
		change.Major = change.Major || commit.Change.Major
		change.Minor = change.Minor || commit.Change.Minor
		change.Patch = change.Patch || commit.Change.Patch
	}
	return change
}

func ApplyChange(version *semver.Version, change Change, allowInitialDevelopmentVersions bool) *semver.Version {
	if !allowInitialDevelopmentVersions && version.Major() == 0 {
		change.Major = true
	}

	if allowInitialDevelopmentVersions && version.Major() == 0 && version.Minor() == 0 {
		change.Minor = true
	}

	if !change.Major && !change.Minor && !change.Patch {
		return nil
	}
	var newVersion semver.Version
	preRel := version.Prerelease()
	if preRel == "" {
		switch {
		case change.Major:
			newVersion = version.IncMajor()
		case change.Minor:
			newVersion = version.IncMinor()
		case change.Patch:
			newVersion = version.IncPatch()
		}
		return &newVersion
	}
	preRelVer := strings.Split(preRel, ".")
	if len(preRelVer) > 1 {
		idx, err := strconv.ParseInt(preRelVer[1], 10, 32)
		if err != nil {
			idx = 0
		}
		preRel = fmt.Sprintf("%s.%d", preRelVer[0], idx+1)
	} else {
		preRel += ".1"
	}
	newVersion, _ = version.SetPrerelease(preRel)
	return &newVersion
}

func GetNewVersion(conf *config.Config, commits []*Commit, latestRelease *Release) *semver.Version {
	return ApplyChange(latestRelease.Version, CalculateChange(commits, latestRelease), conf.AllowInitialDevelopmentVersions)
}
