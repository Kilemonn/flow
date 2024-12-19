package file

import (
	"io"
	"io/fs"
	"os"
	"sync"
)

type SyncFileReadWriter struct {
	file  *os.File
	mutex sync.Mutex
}

// NewSynchronisedFileReadWriter create a new SyncFileReadWriter.
// Do not pass os.O_APPEND as a flag, it is not required.
func NewSynchronisedFileReadWriter(filepath string, flags int) (SyncFileReadWriter, error) {
	file, err := os.OpenFile(filepath, flags, fs.ModeType)
	if err != nil {
		return SyncFileReadWriter{}, err
	}

	// Move to the start of the file to make sure the reading pointer is ready
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return SyncFileReadWriter{}, err
	}

	return SyncFileReadWriter{
		file: file,
	}, nil
}

// [io.Reader]
func (rw *SyncFileReadWriter) Read(b []byte) (int, error) {
	rw.mutex.Lock()
	defer rw.mutex.Unlock()

	return rw.file.Read(b)
}

// [io.Writer]
// Moves the file pointer to the end, performs the read, then returns it back to its prior position ready to Read
func (rw *SyncFileReadWriter) Write(b []byte) (int, error) {
	rw.mutex.Lock()
	defer rw.mutex.Unlock()

	// Get current position
	currentPos, err := rw.file.Seek(0, io.SeekCurrent)
	if err != nil {
		return 0, err
	}

	// Move to the end to perform the write
	_, err = rw.file.Seek(0, io.SeekEnd)
	if err != nil {
		return 0, err
	}

	// Return back to previous position after write
	defer rw.file.Seek(currentPos, io.SeekStart)

	return rw.file.Write(b)
}

// [io.Closer]
func (rw *SyncFileReadWriter) Close() error {
	rw.mutex.Lock()
	defer rw.mutex.Unlock()

	return rw.file.Close()
}
