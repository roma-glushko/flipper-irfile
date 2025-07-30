package flipperirlib

import (
	"github.com/stretchr/testify/require"
	"io"
	"os"
	"testing"
)

func TestMarshalParsedSignals(t *testing.T) {
	c := libFile(t, "./testdata/parsed.ir")

	lib, err := Unmarshal(c)

	require.NoError(t, err)
	require.NotEmpty(t, lib.Signals)
	require.Equal(t, lib.Filetype, FiletypeSignalLib)
	require.Equal(t, lib.Version, "1")
	require.Len(t, lib.Signals, 3)

	rawLib, err := Marshal(lib)
	require.NoError(t, err)
	require.NotEmpty(t, rawLib)
	require.Equal(t, string(c), string(rawLib))
}

func TestMarshalRawSignals(t *testing.T) {
	c := libFile(t, "./testdata/raw.ir")

	lib, err := Unmarshal(c)

	require.NoError(t, err)
	require.NotEmpty(t, lib.Signals)
	require.Equal(t, lib.Filetype, FiletypeSignalLib)
	require.Equal(t, lib.Version, "1")
	require.Len(t, lib.Signals, 3)

	rawLib, err := Marshal(lib)
	require.NoError(t, err)
	require.NotEmpty(t, rawLib)
	require.Equal(t, string(c), string(rawLib))
}

func libFile(t *testing.T, path string) []byte {
	file, err := os.Open(path)
	require.NoError(t, err)

	defer func() {
		err = file.Close()
		require.NoError(t, err)
	}()

	data, err := io.ReadAll(file)
	require.NoError(t, err)

	return data
}
