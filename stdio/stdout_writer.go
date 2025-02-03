package stdio

import (
	"bufio"
	"io"
	"os"

	"github.com/Kilemonn/flow/bidetwriter"
)

func CreateStdOutWriter() (io.WriteCloser, error) {
	return bidetwriter.NewBidetWriter(bufio.NewWriter(os.Stdout)), nil
}
