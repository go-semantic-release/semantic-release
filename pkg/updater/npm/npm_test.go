package npm

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNpmUpdater(t *testing.T) {
	require := require.New(t)

	updater := &Updater{}

	nVer := "1.2.3"
	nVerJSON := json.RawMessage("\"" + nVer + "\"")
	npmrcPath := "../../../test/.npmrc"
	pkgJsonPath := "../../../test/package.json"
	os.Remove(npmrcPath)

	err := updater.Apply(pkgJsonPath, nVer)
	require.NoError(err)
	npmfile, err := ioutil.ReadFile(npmrcPath)
	require.NoError(err)
	require.Equal([]byte(npmrc), npmfile, "invalid .npmrc")
	f, err := os.OpenFile(pkgJsonPath, os.O_RDONLY, 0)
	require.NoError(err)
	defer f.Close()
	var data map[string]json.RawMessage
	err = json.NewDecoder(f).Decode(&data)
	require.NoError(err)

	require.Equal(nVerJSON, data["version"], "invalid version")

	plF, err := os.OpenFile("../../../test/package-lock.json", os.O_RDONLY, 0)
	require.NoError(err, "fixture package-lock.json missing")
	var plData map[string]json.RawMessage
	err = json.NewDecoder(plF).Decode(&plData)
	require.NoError(err)
	require.Equal(nVerJSON, plData["version"], "invalid version")
}

func TestNpmrc(t *testing.T) {
	require := require.New(t)

	nVer := "1.2.3"
	npmrcPath := "../../../test/.npmrc"
	pkgJsonPath := "../../../test/package.json"

	err := ioutil.WriteFile(npmrcPath, []byte("TEST"), 0644)
	require.NoError(err)

	updater := &Updater{}
	err = updater.Apply(pkgJsonPath, nVer)
	require.NoError(err)
	npmfile, err := ioutil.ReadFile(npmrcPath)
	require.NoError(err)
	require.Equal([]byte("TEST"), npmfile, "invalid .npmrc")
	f, err := os.OpenFile(pkgJsonPath, os.O_RDONLY, 0)
	require.NoError(err)
	defer f.Close()
	require.NoError(err)
	var data map[string]json.RawMessage
	err = json.NewDecoder(f).Decode(&data)
	require.NoError(err)
	require.Equal(json.RawMessage("\""+nVer+"\""), data["version"], "invalid version")
}
