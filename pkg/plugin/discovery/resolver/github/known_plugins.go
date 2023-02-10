package github

var knownPlugins = map[string]string{
	"provider-github":              "go-semantic-release/provider-github",
	"provider-gitlab":              "go-semantic-release/provider-gitlab",
	"changelog-generator-default":  "go-semantic-release/changelog-generator-default",
	"commit-analyzer-default":      "go-semantic-release/commit-analyzer-cz",
	"condition-default":            "go-semantic-release/condition-default",
	"condition-github":             "go-semantic-release/condition-github",
	"condition-gitlab":             "go-semantic-release/condition-gitlab",
	"files-updater-npm":            "go-semantic-release/files-updater-npm",
	"provider-git":                 "go-semantic-release/provider-git",
	"condition-bitbucket":          "go-semantic-release/condition-bitbucket",
	"files-updater-helm":           "go-semantic-release/files-updater-helm",
	"hooks-goreleaser":             "go-semantic-release/hooks-goreleaser",
	"hooks-npm-binary-releaser":    "go-semantic-release/hooks-npm-binary-releaser",
	"hooks-plugin-registry-update": "go-semantic-release/hooks-plugin-registry-update",
	"hooks-exec":                   "go-semantic-release/hooks-exec",
}
