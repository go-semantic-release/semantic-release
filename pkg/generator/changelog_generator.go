package generator

type ChangelogGenerator interface {
	Generate(*ChangelogGeneratorConfig) string
}
