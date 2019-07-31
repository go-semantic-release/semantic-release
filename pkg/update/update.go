package update

import (
	"errors"
	"os"
	"path"
	"sync"
)

type Updater func(string, *os.File) error

var updatersMu sync.RWMutex
var updaters = make(map[string]Updater)

var ErrNoUpdater = errors.New("no updater registered")

func Register(name string, u Updater) {
	updatersMu.Lock()
	defer updatersMu.Unlock()
	updaters[name] = u
}

func Apply(file, newVersion string) error {
	name := path.Base(file)
	ufn, ok := updaters[name]
	if !ok {
		return ErrNoUpdater
	}
	f, err := os.OpenFile(file, os.O_RDWR, 0)
	if err != nil {
		return err
	}
	defer f.Close()
	return ufn(newVersion, f)
}
