package updater

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type testFileUpdater struct {
	req  *require.Assertions
	nVer string
}

func (tfu *testFileUpdater) ForFiles() string {
	return "package\\.json"
}

func (tfu *testFileUpdater) Apply(file, newVersion string) error {
	tfu.req.Equal(newVersion, tfu.nVer, "invalid version")
	return nil
}

func TestChainedUpdater(t *testing.T) {
	require := require.New(t)
	nVer := "1.2.3"
	tfu := &testFileUpdater{require, nVer}
	updaters := &ChainedUpdater{Updaters: []FilesUpdater{tfu}}
	if err := updaters.Apply("../../test/package.json", nVer); err != nil {
		require.NoError(err)
	}
}
