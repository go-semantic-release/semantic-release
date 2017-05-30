package update

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
)

func TestPackageJson(t *testing.T) {
	f, err := os.OpenFile("../test-fixtures/package.json", os.O_RDWR, 0)
	if err != nil {
		t.Fatal("fixture missing")
	}
	defer f.Close()
	nVer := "1.2.3"
	packageJson(nVer, f)
	npmfile, err := ioutil.ReadFile("../test-fixtures/.npmrc")
	if err != nil || bytes.Compare(npmfile, []byte(npmrc)) != 0 {
		t.Fatal("invalid .npmrc")
	}
	f.Seek(0, 0)
	var data map[string]json.RawMessage
	json.NewDecoder(f).Decode(&data)
	if bytes.Compare(data["version"], []byte("\""+nVer+"\"")) != 0 {
		t.Fatal("invalid version")
	}
}
