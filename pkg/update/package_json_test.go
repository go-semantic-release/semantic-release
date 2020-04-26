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
		t.Fatal("fixture package.json missing")
	}
	defer f.Close()
	nVer := "1.2.3"
	nVerJson := []byte("\"" + nVer + "\"")
	npmrcPath := "../../test/.npmrc"
	os.Remove(npmrcPath)
	if err := packageJson(nVer, f); err != nil {
		t.Fatal(err)
	}
	npmfile, err := ioutil.ReadFile(npmrcPath)
	if err != nil || !bytes.Equal(npmfile, []byte(npmrc)) {
		t.Fatal("invalid .npmrc")
	}

	if _, err := f.Seek(0, 0); err != nil {
		t.Fatal(err)
	}
	var data map[string]json.RawMessage
	if err := json.NewDecoder(f).Decode(&data); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(data["version"], nVerJson) {
		t.Fatal("invalid version")
	}

	plF, err := os.OpenFile("../../test/package-lock.json", os.O_RDONLY, 0)
	if err != nil {
		t.Fatal("fixture package-lock.json missing")
	}
	var plData map[string]json.RawMessage
	if err := json.NewDecoder(plF).Decode(&plData); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(plData["version"], nVerJson) {
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
	if err := ioutil.WriteFile(npmrcPath, []byte("TEST"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := packageJson(nVer, f); err != nil {
		t.Fatal(err)
	}
	npmfile, err := ioutil.ReadFile(npmrcPath)
	if err != nil || !bytes.Equal(npmfile, []byte("TEST")) {
		t.Fatal("invalid .npmrc")
	}
	if _, err := f.Seek(0, 0); err != nil {
		t.Fatal(err)
	}
	var data map[string]json.RawMessage
	if err := json.NewDecoder(f).Decode(&data); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(data["version"], []byte("\""+nVer+"\"")) {
		t.Fatal("invalid version")
	}
}
