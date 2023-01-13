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
		results *Info
	}{
		{"provider", "git", &Info{
			Type:                "provider",
			Name:                "git",
			NormalizedName:      "provider-git",
			ShortNormalizedName: "provider-git",
			Resolver:            "default",
		}},
		{"provider", "github:myorg/myplugin", &Info{
			Type:                "provider",
			Name:                "myplugin",
			NormalizedName:      "provider-github-myorg-myplugin",
			ShortNormalizedName: "provider-myplugin",
			Resolver:            "github",
			RepoSlug:            "myorg/myplugin",
		}},
		{"ci_condition", "github:myorg/myplugin", &Info{
			Type:                "ci_condition",
			Name:                "myplugin",
			NormalizedName:      "condition-github-myorg-myplugin",
			ShortNormalizedName: "condition-myplugin",
			Resolver:            "github",
			RepoSlug:            "myorg/myplugin",
		}},
		{"provider", "github:myorg/myplugin@^1.0.0", &Info{
			Type:                "provider",
			Name:                "myplugin",
			NormalizedName:      "provider-github-myorg-myplugin",
			ShortNormalizedName: "provider-myplugin",
			Resolver:            "github",
			RepoSlug:            "myorg/myplugin",
			Constraint:          parseConstraint("^1.0.0"),
		}},
		{"provider", "git@2.0.0", &Info{
			Type:                "provider",
			Name:                "git",
			NormalizedName:      "provider-git",
			ShortNormalizedName: "provider-git",
			Resolver:            "default",
			Constraint:          parseConstraint("2.0.0"),
		}},
		{"hooks", "registry:logger", &Info{
			Type:                "hooks",
			Name:                "logger",
			NormalizedName:      "hooks-registry-logger",
			ShortNormalizedName: "hooks-logger",
			Resolver:            "registry",
		}},
		{"hooks", "myresolver:@myorg/myplugin", &Info{
			Type:                "hooks",
			Name:                "myplugin",
			NormalizedName:      "hooks-myresolver-@myorg-myplugin",
			ShortNormalizedName: "hooks-myplugin",
			Resolver:            "myresolver",
			RepoSlug:            "@myorg/myplugin",
		}},
		{"hooks", "myresolver:@myorg/myplugin@1.2.3", &Info{
			Type:                "hooks",
			Name:                "myplugin",
			NormalizedName:      "hooks-myresolver-@myorg-myplugin",
			ShortNormalizedName: "hooks-myplugin",
			Resolver:            "myresolver",
			RepoSlug:            "@myorg/myplugin",
			Constraint:          parseConstraint("1.2.3"),
		}},
	}
	for _, testCase := range testCases {
		results, err := GetPluginInfo(testCase.t, testCase.input)
		require.NoError(t, err)
		require.Equal(t, testCase.results, results)
	}
}
