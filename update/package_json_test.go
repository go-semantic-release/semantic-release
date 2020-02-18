package update

import (
	"bytes"
	"encoding/json"
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
	var data map[string]json.RawMessage
	json.NewDecoder(f).Decode(&data)
	if bytes.Equal(data["version"], []byte("\""+nVer+"\"")) {
		t.Fatal("invalid version")
	}
}
