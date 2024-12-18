package config

import (
	"io"
	"os"
	"testing"

	"github.com/Kilemonn/flow/flow/testutil"
	"github.com/stretchr/testify/require"
)

func TestFilesValid(t *testing.T) {
	fileThatDoesntExist := "fileThatDoesntExist.txt"
	testutil.WithTempFile(t, func(fileName string) {
		fileConfigs := []ConfigFile{
			{
				ID:   "FileExists",
				Path: fileName,
			},
			{
				ID:   "FileDoesNotExist",
				Path: fileThatDoesntExist,
			},
		}

		for _, fCon := range fileConfigs {
			err := fCon.Validate()
			require.Nil(t, err)
		}

		require.Nil(t, os.Remove(fileThatDoesntExist))
	})
}

// TestFileWriterAndReader ensure that using File will write to the provided file, and the reader can read from it
func TestFileWriterAndReader(t *testing.T) {
	testutil.WithTempFile(t, func(filename string) {
		fileConfig := ConfigFile{
			Path: filename,
		}
		writer, err := fileConfig.Writer()
		require.Nil(t, err)

		data := "some file writer data"
		_, err = writer.Write([]byte(data))
		require.Nil(t, err)

		reader, err := fileConfig.Reader()
		require.Nil(t, err)
		read, err := io.ReadAll(reader)
		require.Nil(t, err)
		require.Equal(t, data, string(read))

		// Ensure that if more data is written that is it is picked up during the next read
		data = "more data to write to file"
		_, err = writer.Write([]byte(data))
		require.Nil(t, err)
		read, err = io.ReadAll(reader)
		require.Nil(t, err)
		require.Equal(t, data, string(read))
	})
}
