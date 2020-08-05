package semrel

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/go-semantic-release/semantic-release/pkg/config"
)

func CalculateChange(commits []*Commit, latestRelease *Release) *Change {
	change := &Change{}
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

func ApplyChange(rawVersion string, rawChange *Change, allowInitialDevelopmentVersions bool) string {
	version := semver.MustParse(rawVersion)
	change := &Change{
		Major: rawChange.Major,
		Minor: rawChange.Minor,
		Patch: rawChange.Patch,
	}
	if !allowInitialDevelopmentVersions && version.Major() == 0 {
		change.Major = true
	}

	if allowInitialDevelopmentVersions && version.Major() == 0 && version.Minor() == 0 {
		change.Minor = true
	}
	if !change.Major && !change.Minor && !change.Patch {
		return ""
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
		return newVersion.String()
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
	return newVersion.String()
}

func GetNewVersion(conf *config.Config, commits []*Commit, latestRelease *Release) string {
	return ApplyChange(latestRelease.Version, CalculateChange(commits, latestRelease), conf.AllowInitialDevelopmentVersions)
}
