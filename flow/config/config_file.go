package config

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/Kilemonn/flow/flow/bidetwriter"
)

type ConfigFile struct {
	ID   string
	Path string
}

// [ConfigModel.GetID]
func (c ConfigFile) GetID() string {
	return c.ID
}

// [ConfigModel.Validate]
func (c ConfigFile) Validate() error {
	// TODO: Should we fail on input files that don't exist?
	if _, err := os.Stat(c.Path); errors.Is(err, os.ErrNotExist) {
		file, err := os.Create(c.Path)
		if err != nil {
			return fmt.Errorf("failed to create file with ID [%s] and path [%s] with error %s", c.ID, c.Path, err.Error())
		}
		file.Close()
	}
	return nil
}

// [ConfigModel.Reader]
func (c ConfigFile) Reader() (io.ReadCloser, error) {
	return os.OpenFile(c.Path, os.O_RDONLY, os.ModeType)
}

// [ConfigModel.Writer]
func (c ConfigFile) Writer() (io.WriteCloser, error) {
	file, err := os.OpenFile(c.Path, os.O_CREATE, os.ModeAppend)
	if err != nil {
		return nil, err
	}
	return bidetwriter.NewBidetWriter(bufio.NewWriter(file)), nil
}
