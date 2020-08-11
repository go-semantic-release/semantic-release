package npm

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"

	"github.com/go-semantic-release/semantic-release/v2/pkg/plugin"
	"github.com/go-semantic-release/semantic-release/v2/pkg/updater"
	"github.com/urfave/cli/v2"
)

const npmrc = "//registry.npmjs.org/:_authToken=${NPM_TOKEN}\n"

type Updater struct {
}

func (u *Updater) ForFiles() string {
	return "package\\.json"
}

func updateJsonFile(fName, newVersion string) error {
	file, err := os.OpenFile(fName, os.O_RDWR, 0)
	if err != nil {
		return err
	}
	defer file.Close()
	var data map[string]json.RawMessage
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		return err
	}
	data["version"] = json.RawMessage("\"" + newVersion + "\"")
	if _, err := file.Seek(0, 0); err != nil {
		return err
	}
	if err := file.Truncate(0); err != nil {
		return err
	}
	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")
	if err := enc.Encode(data); err != nil {
		return err
	}
	return nil
}

func (u *Updater) Apply(file, newVersion string) error {
	if err := updateJsonFile(file, newVersion); err != nil {
		return err
	}

	packageLockPath := path.Join(path.Dir(file), "package-lock.json")
	if _, err := os.Stat(packageLockPath); err == nil {
		if err := updateJsonFile(packageLockPath, newVersion); err != nil {
			return err
		}
	}

	if os.Getenv("NPM_CONFIG_USERCONFIG") != "" {
		return nil
	}

	var err error
	npmrcPath := path.Join(path.Dir(file), ".npmrc")
	if _, err = os.Stat(npmrcPath); os.IsNotExist(err) {
		return ioutil.WriteFile(npmrcPath, []byte(npmrc), 0644)
	}

	return err
}

func Main(c *cli.Context) error {
	plugin.Serve(&plugin.ServeOpts{
		FilesUpdater: func() updater.FilesUpdater {
			return &Updater{}
		},
	})
	return nil
}
