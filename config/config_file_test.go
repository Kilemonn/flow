package config

import (
	"io"
	"os"
	"testing"

	"github.com/Kilemonn/flow/testutil"
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
			require.NoError(t, err)
		}

		require.NoError(t, os.Remove(fileThatDoesntExist))
	})
}

// TestFileWriterAndReader ensure that using File will write
// to the provided file, and the reader can read from it
func TestFileWriterAndReader(t *testing.T) {
	testutil.WithTempFile(t, func(filename string) {
		fileConfig := ConfigFile{
			Path: filename,
		}
		writer, err := fileConfig.Writer()
		require.NoError(t, err)

		data := "some file writer data"
		_, err = writer.Write([]byte(data))
		require.NoError(t, err)

		reader, err := fileConfig.Reader()
		require.NoError(t, err)
		read, err := io.ReadAll(reader)
		require.NoError(t, err)
		require.Equal(t, data, string(read))

		// Ensure that if more data is written that is it is picked up during the next read
		data = "more data to write to file"
		_, err = writer.Write([]byte(data))
		require.NoError(t, err)
		read, err = io.ReadAll(reader)
		require.NoError(t, err)
		require.Equal(t, data, string(read))
	})
}

// TestFileWithNoTruncateFlag Ensure that without the truncate flag, the file
// is not truncated and is appended to
func TestFileWithNoTruncateFlag(t *testing.T) {
	testutil.WithTempFile(t, func(filePath string) {
		fileConfig := ConfigFile{
			Path: filePath,
		}
		require.False(t, fileConfig.Trunc)

		// Create file and write content to it
		data := "TestFileWithNoFlag"
		require.NoError(t, os.WriteFile(filePath, []byte(data), os.ModeType))

		// With no flag provided the writer should "create" and not truncate the file
		require.NoError(t, fileConfig.Validate())
		w, err := fileConfig.Writer()
		require.NoError(t, err)
		defer w.Close()

		n, err := w.Write([]byte(data))
		require.NoError(t, err)
		require.Equal(t, len(data), n)

		read, err := os.ReadFile(filePath)
		require.NoError(t, err)
		require.Equal(t, data+data, string(read))
	})
}

// TestFileWithTruncFlag Ensure that with the truncate flag that the file is truncated
func TestFileWithTruncFlag(t *testing.T) {
	testutil.WithTempFile(t, func(filePath string) {
		fileConfig := ConfigFile{
			Path:  filePath,
			Trunc: true,
		}
		require.True(t, fileConfig.Trunc)

		// Create file and write content to it
		data := "TestFileWithTruncFlag"
		require.NoError(t, os.WriteFile(filePath, []byte(data), os.ModeType))

		// With the trunc flag provided, after getting the writer, the file should be truncated
		require.NoError(t, fileConfig.Validate())
		w, err := fileConfig.Writer()
		require.NoError(t, err)
		defer w.Close()

		// Now we can run the content again and expect only 1 copy of it to exist in the file
		n, err := w.Write([]byte(data))
		require.NoError(t, err)
		require.Equal(t, len(data), n)

		read, err := os.ReadFile(filePath)
		require.NoError(t, err)
		require.Equal(t, data, string(read))
	})
}

// TestFileAsReaderAndWriter_NoTruncate_ReadAllFirst check behaviour when truncate is disabled along with
// reading all of the file content before writing more then reading the remaining
func TestFileAsReaderAndWriter_NoTruncate_ReadAllFirst(t *testing.T) {
	testutil.WithTempFile(t, func(filePath string) {
		fileConfig := ConfigFile{
			Path: filePath,
		}
		require.False(t, fileConfig.Trunc)

		initialContent := "initial"
		require.NoError(t, os.WriteFile(filePath, []byte(initialContent), os.ModeAppend))

		r, err := fileConfig.Reader()
		require.NoError(t, err)
		defer r.Close()

		w, err := fileConfig.Writer()
		require.NoError(t, err)
		defer w.Close()

		// Read until EOF
		read, err := io.ReadAll(r)
		require.NoError(t, err)
		require.Equal(t, initialContent, string(read))

		// Write some content
		content := "TestFileAsReaderAndWriter_NoTruncate_ReadAllFirst"
		n, err := w.Write([]byte(content))
		require.NoError(t, err)
		require.Equal(t, len(content), n)

		// Now read again
		read, err = io.ReadAll(r)
		require.NoError(t, err)
		require.Equal(t, content, string(read))
	})
}

// TestFileAsReaderAndWriter_NoTruncate_WriteThenRead check behaviour of truncate disabled along with
// writing additional content before reading all content in the file
func TestFileAsReaderAndWriter_NoTruncate_WriteThenRead(t *testing.T) {
	testutil.WithTempFile(t, func(filePath string) {
		fileConfig := ConfigFile{
			Path: filePath,
		}
		require.False(t, fileConfig.Trunc)

		initialContent := "initial"
		require.NoError(t, os.WriteFile(filePath, []byte(initialContent), os.ModeAppend))

		r, err := fileConfig.Reader()
		require.NoError(t, err)
		defer r.Close()

		w, err := fileConfig.Writer()
		require.NoError(t, err)
		defer w.Close()

		// Write some content
		content := "TestFileAsReaderAndWriter_NoTruncate_WriteThenRead"
		n, err := w.Write([]byte(content))
		require.NoError(t, err)
		require.Equal(t, len(content), n)

		// Now read again
		read, err := io.ReadAll(r)
		require.NoError(t, err)
		require.Equal(t, initialContent+content, string(read))
	})
}
