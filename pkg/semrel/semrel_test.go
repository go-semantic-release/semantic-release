package semrel

import (
	"fmt"
	"testing"

	"github.com/Masterminds/semver"
	"github.com/go-semantic-release/semantic-release/pkg/config"
)

func TestCalculateChange(t *testing.T) {
	commits := []*Commit{
		{SHA: "a", Change: &Change{Major: true, Minor: false, Patch: false}},
		{SHA: "b", Change: &Change{Major: false, Minor: true, Patch: false}},
		{SHA: "c", Change: &Change{Major: false, Minor: false, Patch: true}},
	}
	change := CalculateChange(commits, &Release{})
	if !change.Major || !change.Minor || !change.Patch {
		t.Fail()
	}
	change = CalculateChange(commits, &Release{SHA: "a"})
	if change.Major || change.Minor || change.Patch {
		t.Fail()
	}
	newVersion := GetNewVersion(&config.Config{}, commits, &Release{SHA: "b", Version: "1.0.0"})
	if newVersion != "2.0.0" {
		t.Fail()
	}
}

func TestApplyChange(t *testing.T) {
	NoChange := &Change{Major: false, Minor: false, Patch: false}
	PatchChange := &Change{Major: false, Minor: false, Patch: true}
	MinorChange := &Change{Major: false, Minor: true, Patch: true}
	MajorChange := &Change{Major: true, Minor: true, Patch: true}

	testCases := []struct {
		currentVersion                  string
		change                          *Change
		expectedVersion                 string
		allowInitialDevelopmentVersions bool
	}{
		// No Previous Releases
		{"0.0.0", NoChange, "1.0.0", false},
		{"0.0.0", PatchChange, "1.0.0", false},
		{"0.0.0", MinorChange, "1.0.0", false},
		{"0.0.0", MajorChange, "1.0.0", false},

		{"0.0.0", NoChange, "0.1.0", true},
		{"0.0.0", PatchChange, "0.1.0", true},
		{"0.0.0", MinorChange, "0.1.0", true},
		{"0.0.0", MajorChange, "1.0.0", true},

		{"1.0.0", NoChange, "", false},
		{"1.0.0", PatchChange, "1.0.1", false},
		{"1.0.0", MinorChange, "1.1.0", false},
		{"1.0.0", MajorChange, "2.0.0", false},
		{"0.1.0", NoChange, "1.0.0", false},
		{"0.1.0", NoChange, "", true},

		{"2.0.0-beta", MajorChange, "2.0.0-beta.1", false},
		{"2.0.0-beta.2", MajorChange, "2.0.0-beta.3", false},
		{"2.0.0-beta.1.1", MajorChange, "2.0.0-beta.2", false},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Version: %s, Change: %v, Expected: %s", tc.currentVersion, tc.change, tc.expectedVersion), func(t *testing.T) {
			current, err := semver.NewVersion(tc.currentVersion)

			if err != nil {
				t.Errorf("failed to create version: %v", err)
			}

			actual := ApplyChange(current.String(), tc.change, tc.allowInitialDevelopmentVersions)

			// Handle no new version case
			if actual != "" && tc.expectedVersion != "" {
				if actual != tc.expectedVersion {
					t.Errorf("expected: %s, got: %s", tc.expectedVersion, actual)
				}
			}
		})
	}
}
