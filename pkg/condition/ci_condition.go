package condition

type CICondition interface {
	Name() string
	RunCondition(map[string]string) error
	GetCurrentBranch() string
	GetCurrentSHA() string
}
