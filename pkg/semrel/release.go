package semrel

import (
	"sort"
	"strings"

	"github.com/Masterminds/semver"
)

type Releases []*Release

func (r Releases) Len() int {
	return len(r)
}

func (r Releases) Less(i, j int) bool {
	return semver.MustParse(r[j].Version).LessThan(semver.MustParse(r[i].Version))
}

func (r Releases) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (releases Releases) GetLatestRelease(vrange string) (*Release, error) {
	if len(releases) == 0 {
		return &Release{SHA: "", Version: "0.0.0"}, nil
	}

	sort.Sort(releases)

	var lastRelease *Release
	for _, r := range releases {
		if semver.MustParse(r.Version).Prerelease() == "" {
			lastRelease = r
			break
		}
	}

	if vrange == "" {
		if lastRelease != nil {
			return lastRelease, nil
		}
		return &Release{SHA: "", Version: "0.0.0"}, nil
	}

	constraint, err := semver.NewConstraint(vrange)
	if err != nil {
		return nil, err
	}
	for _, r := range releases {
		if constraint.Check(semver.MustParse(r.Version)) {
			return r, nil
		}
	}

	nver, err := semver.NewVersion(vrange)
	if err != nil {
		return nil, err
	}

	splitPre := strings.SplitN(vrange, "-", 2)
	if len(splitPre) == 1 {
		return &Release{SHA: lastRelease.SHA, Version: nver.String()}, nil
	}

	npver, err := nver.SetPrerelease(splitPre[1])
	if err != nil {
		return nil, err
	}
	return &Release{SHA: lastRelease.SHA, Version: npver.String()}, nil
}
