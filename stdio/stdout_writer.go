package stdio

import (
	"io"
	"os"

	"github.com/Kilemonn/flow/bidetwriter"
)

func CreateStdOutWriter() (io.WriteCloser, error) {
	return bidetwriter.NewBidetWriter(os.Stdout), nil
}
