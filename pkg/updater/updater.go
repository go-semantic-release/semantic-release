package updater

import (
	"fmt"
	"path"
	"regexp"
)

type FilesUpdater interface {
	Init(map[string]string) error
	Name() string
	Version() string
	ForFiles() string
	Apply(file, newVersion string) error
}

type ChainedUpdater struct {
	Updaters []FilesUpdater
}

func (u *ChainedUpdater) Init(conf map[string]string) error {
	for _, fu := range u.Updaters {
		err := fu.Init(conf)
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *ChainedUpdater) GetNameVersionPairs() []string {
	ret := make([]string, len(u.Updaters))
	for i, fu := range u.Updaters {
		ret[i] = fmt.Sprintf("%s@%s", fu.Name(), fu.Version())
	}
	return ret
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
