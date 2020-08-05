package updater

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRegisterApply(t *testing.T) {
	require := require.New(t)
	nVer := "1.2.3"
	Register("package.json", func(newVersion string, file *os.File) error {
		require.Equal(newVersion, nVer, "invalid version")
		return nil
	})
	if err := Apply("../../test/package.json", nVer); err != nil {
		require.NoError(err)
	}
}

func TestApply(t *testing.T) {
	require := require.New(t)
	err := Apply("invalidFile", "")
	require.Equal(err, ErrNoUpdater)

	err = Apply("not/existing/package.json", "")
	var pathErr *os.PathError
	ok := errors.As(err, &pathErr)

	if ok {
		require.Equal("open", pathErr.Op)
		require.Equal("not/existing/package.json", pathErr.Path)
	} else {
		require.NoError(err)
	}
}
