package updater

import (
	"path"
	"regexp"
)

type Updater interface {
	Init(map[string]string) error
	Name() string
	Version() string
	Apply(file, newVersion string) error
}

type FilesUpdater interface {
	ForFiles() string
	Updater
}

type ChainedUpdater struct {
	Updaters []FilesUpdater
}

func (u *ChainedUpdater) Apply(file, newVersion string) error {
	for _, fu := range u.Updaters {
		re, err := regexp.Compile(fu.ForFiles())
		if err != nil {
			return err
		}
		if !re.MatchString(path.Base(file)) {
			continue
		}
		if err := fu.Apply(file, newVersion); err != nil {
			return err
		}
	}
	return nil
}
