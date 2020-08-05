package semrel

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReleases(t *testing.T) {
	testCases := []struct {
		Releases      Releases
		VRange        string
		LatestVersion string
	}{
		{
			Releases{},
			"",
			"0.0.0",
		},
		{
			Releases{
				{SHA: "a", Version: "0.1.0"},
			},
			"",
			"0.1.0",
		},
		{
			Releases{
				{SHA: "a", Version: "1.0.0"},
				{SHA: "b", Version: "1.1.0"},
				{SHA: "c", Version: "2.0.0-beta"},
				{SHA: "d", Version: "0.1.0"},
			},
			"",
			"1.1.0",
		},
		{
			Releases{
				{SHA: "a", Version: "1.0.0"},
				{SHA: "b", Version: "1.1.0"},
				{SHA: "c", Version: "2.0.0"},
				{SHA: "c", Version: "2.1.0-beta"},
				{SHA: "d", Version: "0.1.0"},
			},
			"2-beta",
			"2.1.0-beta",
		},
		{
			Releases{
				{SHA: "a", Version: "1.0.0"},
				{SHA: "b", Version: "1.1.0"},
				{SHA: "c", Version: "3.0.0-rc.1"},
				{SHA: "c", Version: "3.0.0-rc.2"},
				{SHA: "d", Version: "0.1.0"},
			},
			"3-rc",
			"3.0.0-rc.2",
		},
		{
			Releases{
				{SHA: "a", Version: "1.0.0"},
				{SHA: "b", Version: "1.1.0"},
				{SHA: "c", Version: "3.0.0-rc.1"},
				{SHA: "c", Version: "3.0.0-rc.2"},
				{SHA: "d", Version: "0.1.0"},
			},
			"4-alpha",
			"4.0.0-alpha",
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TestReleases: %d, LV: %s", i, tc.LatestVersion), func(t *testing.T) {
			lr, err := tc.Releases.GetLatestRelease(tc.VRange)
			require.NoError(t, err)
			require.Equal(t, tc.LatestVersion, lr.Version)
		})
	}
}
