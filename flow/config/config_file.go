package config

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/Kilemonn/flow/flow/bidetwriter"
	"github.com/Kilemonn/flow/flow/file"
)

type ConfigFile struct {
	ID   string
	Path string
	// Refer to syscall properties for your OS, specifically: syscall.O_APPEND, syscall.O_TRUNC.
	// This value is XOR'd with O_CREATE so the file is always created.
	// If no flag is provided, by default the file is created and truncated.
	Flag int

	file *file.SyncFileReadWriter
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
	err := c.initialiseFile()
	return c.file, err
}

// [ConfigModel.Writer]
func (c ConfigFile) Writer() (io.WriteCloser, error) {
	err := c.initialiseFile()
	if err != nil {
		return nil, err
	}
	return bidetwriter.NewBidetWriter(bufio.NewWriter(c.file)), nil
}

func (c *ConfigFile) initialiseFile() error {
	if c.file == nil {
		mode := os.O_CREATE | os.O_RDWR
		temp, err := file.NewSynchronisedFileReadWriter(c.Path, mode|c.Flag)
		c.file = &temp
		return err
	}
	return nil
}
