package plugin

import (
	"testing"

	"github.com/Masterminds/semver/v3"

	"github.com/stretchr/testify/require"
)

func parseConstraint(c string) *semver.Constraints {
	constraint, _ := semver.NewConstraint(c)
	return constraint
}

func TestGetPluginInfo(t *testing.T) {
	testCases := []struct {
		t       string
		input   string
		results *PluginInfo
	}{
		{"provider", "git", &PluginInfo{
			Type: "provider", Name: "git", NormalizedName: "provider-git", Resolver: "default",
		}},
		{"provider", "github:myorg/myplugin", &PluginInfo{
			Type: "provider", Name: "myplugin", NormalizedName: "provider-github-myorg-myplugin", Resolver: "github",
		}},
		{"provider", "github:myorg/myplugin", &PluginInfo{
			Type: "provider", Name: "myplugin", NormalizedName: "provider-github-myorg-myplugin", Resolver: "github",
		}},
		{"provider", "github:myorg/myplugin@^1.0.0", &PluginInfo{
			Type: "provider", Name: "myplugin", NormalizedName: "provider-github-myorg-myplugin", Resolver: "github",
			Constraint: parseConstraint("^1.0.0"),
		}},
		{"provider", "git@2.0.0", &PluginInfo{
			Type: "provider", Name: "git", NormalizedName: "provider-git", Resolver: "default",
			Constraint: parseConstraint("2.0.0"),
		}},
	}
	for _, testCase := range testCases {
		results, err := GetPluginInfo(testCase.t, testCase.input)
		require.NoError(t, err)
		require.Equal(t, testCase.results, results)
	}
}
