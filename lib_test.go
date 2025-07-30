// Copyright 2025 Roma Hlushko
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package flipperirfile

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
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
