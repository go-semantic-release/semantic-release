package update

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
)

const npmrc = "//registry.npmjs.org/:_authToken=${NPM_TOKEN}\n"

func init() {
	Register("package.json", packageJson)
}

func updateJsonFile(newVersion string, file *os.File) error {
	var data map[string]json.RawMessage
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		return err
	}
	data["version"] = json.RawMessage("\"" + newVersion + "\"")
	file.Seek(0, 0)
	file.Truncate(0)
	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")
	if err := enc.Encode(data); err != nil {
		return err
	}
	return nil
}

func packageJson(newVersion string, file *os.File) error {
	if err := updateJsonFile(newVersion, file); err != nil {
		return err
	}

	packageLockPath := path.Join(path.Dir(file.Name()), "package-lock.json")
	plFile, err := os.OpenFile(packageLockPath, os.O_RDWR, 0)
	if err == nil {
		defer plFile.Close()
		if err := updateJsonFile(newVersion, plFile); err != nil {
			return err
		}
	}

	npmrcPath := path.Join(path.Dir(file.Name()), ".npmrc")
	if _, err = os.Stat(npmrcPath); os.IsNotExist(err) {
		return ioutil.WriteFile(npmrcPath, []byte(npmrc), 0644)
	}
	return err
}
