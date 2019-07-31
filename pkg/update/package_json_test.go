package update

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
)

func TestPackageJson(t *testing.T) {
	f, err := os.OpenFile("../../test/package.json", os.O_RDWR, 0)
	if err != nil {
		t.Fatal("fixture missing")
	}
	defer f.Close()
	nVer := "1.2.3"
	npmrcPath := "../../test/.npmrc"
	os.Remove(npmrcPath)
	packageJson(nVer, f)
	npmfile, err := ioutil.ReadFile(npmrcPath)
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

func TestNpmrc(t *testing.T) {
	f, err := os.OpenFile("../../test/package.json", os.O_RDWR, 0)
	if err != nil {
		t.Fatal("fixture missing")
	}
	defer f.Close()
	nVer := "1.2.3"
	npmrcPath := "../../test/.npmrc"
	ioutil.WriteFile(npmrcPath, []byte("TEST"), 0644)
	packageJson(nVer, f)
	npmfile, err := ioutil.ReadFile(npmrcPath)
	if err != nil || bytes.Compare(npmfile, []byte("TEST")) != 0 {
		t.Fatal("invalid .npmrc")
	}
	f.Seek(0, 0)
	var data map[string]json.RawMessage
	json.NewDecoder(f).Decode(&data)
	if bytes.Compare(data["version"], []byte("\""+nVer+"\"")) != 0 {
		t.Fatal("invalid version")
	}
}
