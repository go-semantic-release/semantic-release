package condition

type CI interface {
	Name() string
	RunCondition(map[string]interface{}) error
	GetCurrentBranch() string
	GetCurrentSHA() string
}
