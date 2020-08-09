package condition

type CICondition interface {
	Name() string
	RunCondition(map[string]interface{}) error
	GetCurrentBranch() string
	GetCurrentSHA() string
}
