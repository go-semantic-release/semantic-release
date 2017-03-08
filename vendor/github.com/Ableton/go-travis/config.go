package travis

// Config represents a Travis CI build config
//
// As the endpoints responses fields can change types
// according to the request, it is not possible to represent
// them easily in a statically language like go. Therefore
// Its implementation is purposedly left in an incomplete
// and very limited state.
type Config struct {
	Os            string              `json:"os,omitempty"`
	Language      string              `json:"language,omitempty"`
	Branches      map[string][]string `json:"branches,omitemtpy"`
	BeforeScript  []string            `json:"before_script,omitempty"`
	Script        []string            `json:"script,omitempty"`
	AfterScript   []string            `json:"after_script,omitempty"`
	BeforeInstall []string            `json:"before_install,omitempty"`
	Install       []string            `json:"install,omitempty"`
	AfterSuccess  []string            `json:"after_success,omitempty"`
	AfterFailure  []string            `json:"after_failure,omitempty"`
	Addons        map[string]string   `json:"addons,omitempty"`
	Notifications NotificationMap     `json:"notifications,omitempty"`
}

type NotificationMap map[string]map[string]string
