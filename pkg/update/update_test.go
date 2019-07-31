package update

import (
	"os"
	"strings"
	"testing"
)

func TestRegisterApply(t *testing.T) {
	nVer := "1.2.3"
	Register("package.json", func(newVersion string, file *os.File) error {
		if newVersion != nVer {
			t.Fatal("invalid version")
		}
		return nil
	})
	if err := Apply("../../test/package.json", nVer); err != nil {
		t.Fatal(err)
	}
}

func TestApply(t *testing.T) {
	if err := Apply("invalidFile", ""); err != ErrNoUpdater {
		t.Fatal(err)
	}
	if err := Apply("not/existing/package.json", ""); !strings.Contains(err.Error(), "no such file or directory") {
		t.Fatal(err)
	}
}
