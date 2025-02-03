package stdio

import (
	"io"
	"os"
)

func CreateStdInReader() (io.ReadCloser, error) {
	return os.Stdin, nil
}
