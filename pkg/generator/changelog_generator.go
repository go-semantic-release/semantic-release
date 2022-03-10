package generator

type ChangelogGenerator interface {
	Init(map[string]string) error
	Name() string
	Version() string
	Generate(*ChangelogGeneratorConfig) string
}
