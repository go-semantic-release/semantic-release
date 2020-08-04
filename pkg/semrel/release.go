package semrel

import (
	"sort"
	"strings"

	"github.com/Masterminds/semver"
)

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
