package file

import (
	"io"
	"os"
	"testing"

	"github.com/Kilemonn/flow/flow/testutil"
	"github.com/stretchr/testify/require"
)

func TestSyncFileReadWriter(t *testing.T) {
	testutil.WithTempFile(t, func(filepath string) {
		rw, err := NewSynchronisedFileReadWriter(filepath, os.O_RDWR)
		require.NoError(t, err)
		defer rw.Close()

		// Check position is at start
		pos, err := rw.file.Seek(0, io.SeekCurrent)
		require.NoError(t, err)
		require.Equal(t, int64(0), pos)

		content := "tw"
		n, err := rw.Write([]byte(content))
		require.NoError(t, err)
		require.Equal(t, len(content), n)

		read := make([]byte, 1)
		c, err := rw.Read(read)
		require.NoError(t, err)
		require.Equal(t, len(read), c)
		require.Equal(t, content[0], string(read)[0])

		// Check position is at position 1
		pos, err = rw.file.Seek(0, io.SeekCurrent)
		require.NoError(t, err)
		require.Equal(t, int64(1), pos)

		// Write another 2
		n, err = rw.Write([]byte(content))
		require.NoError(t, err)
		require.Equal(t, len(content), n)

		// Read the remaining 3
		all, err := io.ReadAll(&rw)
		require.NoError(t, err)
		require.Equal(t, 3, len(all))
		require.Equal(t, string(content[1])+content, string(all))

		// Check position is at position 4
		pos, err = rw.file.Seek(0, io.SeekCurrent)
		require.NoError(t, err)
		require.Equal(t, int64(4), pos)
	})
}
