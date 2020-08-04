package semrel

type RawCommit struct {
	SHA        string
	RawMessage string
}

type Change struct {
	Major, Minor, Patch bool
}

type Commit struct {
	SHA     string
	Raw     []string
	Type    string
	Scope   string
	Message string
	Change  Change
}
