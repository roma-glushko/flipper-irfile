package flipperirlib

import (
	"fmt"
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
	require.Equal(t, lib.Filetype, "IR library file")
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
	require.Equal(t, lib.Filetype, "IR library file")
	require.Equal(t, lib.Version, "1")
	require.Len(t, lib.Signals, 3)

	rawLib, err := Marshal(lib)
	require.NoError(t, err)
	require.NotEmpty(t, rawLib)
	require.Equal(t, string(c), string(rawLib))
}

func TestProcess(t *testing.T) {
	c := libFile(t, "./testdata/tv.ir")

	lib, err := Unmarshal(c)
	require.NoError(t, err)

	powerSignals := make([]Signal, 0, 50)

	for _, s := range lib.Signals {
		if s.Name == "Power" && s.Type == SignalTypeParsed {
			powerSignals = append(powerSignals, s)
		}
	}

	fmt.Println(fmt.Sprintf("power signals found: %d", len(powerSignals)))

	lib.Signals = powerSignals
	rawLib, err := Marshal(lib)

	err = os.WriteFile("./testdata/power.ir", rawLib, 0644)
	require.NoError(t, err)
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
