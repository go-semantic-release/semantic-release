package condition

type CICondition interface {
	Name() string
	Version() string
	RunCondition(map[string]string) error
	GetCurrentBranch() string
	GetCurrentSHA() string
}
