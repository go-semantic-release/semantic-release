package npm

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPackageJson(t *testing.T) {
	require := require.New(t)

	f, err := os.OpenFile("../../../test/package.json", os.O_RDWR, 0)
	require.NoError(err, "fixture package.json missing")
	defer f.Close()
	nVer := "1.2.3"
	nVerJSON := json.RawMessage("\"" + nVer + "\"")
	npmrcPath := "../../../test/.npmrc"
	os.Remove(npmrcPath)
	err = packageJson(nVer, f)
	require.NoError(err)
	npmfile, err := ioutil.ReadFile(npmrcPath)
	require.NoError(err)
	require.Equal([]byte(npmrc), npmfile, "invalid .npmrc")
	_, err = f.Seek(0, 0)
	require.NoError(err)
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
	f, err := os.OpenFile("../../../test/package.json", os.O_RDWR, 0)
	require.NoError(err, "fixture missing")
	defer f.Close()
	nVer := "1.2.3"
	npmrcPath := "../../../test/.npmrc"
	err = ioutil.WriteFile(npmrcPath, []byte("TEST"), 0644)
	require.NoError(err)
	err = packageJson(nVer, f)
	require.NoError(err)
	npmfile, err := ioutil.ReadFile(npmrcPath)
	require.NoError(err)
	require.Equal([]byte("TEST"), npmfile, "invalid .npmrc")
	_, err = f.Seek(0, 0)
	require.NoError(err)
	var data map[string]json.RawMessage
	err = json.NewDecoder(f).Decode(&data)
	require.NoError(err)
	require.Equal(json.RawMessage("\""+nVer+"\""), data["version"], "invalid version")
}
