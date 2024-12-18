package config

import (
	"fmt"
	"io"
	"os"

	"github.com/Kilemonn/flow/flow"
	"gopkg.in/yaml.v3"
)

const (
	StdIn  string = "stdin"
	StdOut string = "stdout"
)

type Config struct {
	ConfConns []ConfigConnection
	Ports     []ConfigPort
	Files     []ConfigFile
	Sockets   []ConfigSocket
	Ipcs      []ConfigIPC
	Settings  ConfigSettings

	models      map[string]ConfigModel    `json:"-"`
	readers     map[string]io.ReadCloser  `json:"-"`
	writers     map[string]io.WriteCloser `json:"-"`
	Connections []Connection              `json:"-"`
}

type Connection struct {
	Reader    io.Reader
	ReaderId  string
	Writer    io.Writer
	WriterIds []string
}

type ConfigConnection struct {
	ReaderID string
	WriterID string
}

// An interface that all Config* objects will implement.
type ConfigModel interface {
	// GetID returns the ID of the model
	GetID() string
	// Validate will perform any validation on behalf of the object, or any pre-setup that is required
	Validate() error
	// Reader will get a reader for the underlying configuration
	Reader() (io.ReadCloser, error)
	// Writer will get a writer for the underlying configuration
	Writer() (io.WriteCloser, error)
}

func (c *Config) Initialise() error {
	err := c.validate()
	if err != nil {
		return err
	}

	err = c.combineToReadersAndWriters()
	if err != nil {
		return err
	}
	c.createConnections()

	return nil
}

// Write the provided Config to the provided filepath
func (c Config) writeConfig(filePath string) error {
	yamlData, err := yaml.Marshal(&c)
	if err != nil {
		return err
	}

	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Writer.Write(f, yamlData)
	return err
}

func (c *Config) validate() error {
	// TODO: Add checks to make sure there is no loops, OR any nodes that are not connected (have no connection)?
	err := c.componentIDsAreUnique()
	if err != nil {
		return err
	}

	for _, model := range c.models {
		err = model.Validate()
		if err != nil {
			return err
		}
	}

	return nil
}

func isInvalidID(id string) bool {
	return id == StdIn || id == StdOut
}

// Check that the IDs of files and ports are unique and also do not clash with stdin or stdout
func (c *Config) componentIDsAreUnique() error {
	c.models = make(map[string]ConfigModel)

	for _, port := range c.Ports {
		if _, exists := c.models[port.ID]; isInvalidID(port.ID) || exists {
			return fmt.Errorf("found port with a duplicate ID [%s] defined or is overriding \"%s\" or \"%s\"", port.ID, StdIn, StdOut)
		} else {
			c.models[port.ID] = &port
		}
	}

	for _, file := range c.Files {
		if _, exists := c.models[file.ID]; isInvalidID(file.ID) || exists {
			return fmt.Errorf("found file with a duplicate ID [%s] defined or is overriding \"%s\" or \"%s\"", file.ID, StdIn, StdOut)
		} else {
			c.models[file.ID] = file
		}
	}

	for _, socket := range c.Sockets {
		if _, exists := c.models[socket.ID]; isInvalidID(socket.ID) || exists {
			return fmt.Errorf("found socket with a duplicate ID [%s] defined or is overriding \"%s\" or \"%s\"", socket.ID, StdIn, StdOut)
		} else {
			c.models[socket.ID] = socket
		}
	}

	for _, ipc := range c.Ipcs {
		if _, exists := c.models[ipc.ID]; isInvalidID(ipc.ID) || exists {
			return fmt.Errorf("found ipc with a duplicate ID [%s] defined or is overriding \"%s\" or \"%s\"", ipc.ID, StdIn, StdOut)
		} else {
			c.models[ipc.ID] = ipc
		}
	}

	return nil
}

// Load all configured readers and writers and load them into the returned map with "id" -> [io.ReadCloser] / [io.WriteCloser] as appropriate.
// StdIn and StdOut are also initialised and returned in these maps.
func (c *Config) combineToReadersAndWriters() error {
	c.readers = make(map[string]io.ReadCloser)
	stdIn, _ := flow.CreateStdInReader()
	c.readers[StdIn] = stdIn

	c.writers = make(map[string]io.WriteCloser)
	stdOut, _ := flow.CreateStdOutWriter()
	c.writers[StdOut] = stdOut

	// Firstly iterate over and ONLY initialise the READER (listening) sockets, since if we connect to ourself we need to make sure
	// the reader is listening first before the writer connects to us (for TCP). See below for the second loop.
	// This is the same for IPC channels.
	for _, connection := range c.ConfConns {
		rID := connection.ReaderID

		if _, exists := c.readers[rID]; !exists {
			if model, ok := c.models[rID]; ok {
				r, err := model.Reader()
				if err != nil {
					return err
				}
				c.readers[rID] = r
			}
		}
	}

	// Move the socket writer init and IPC init to a second loop:
	for _, connection := range c.ConfConns {
		wID := connection.WriterID

		if _, exists := c.writers[wID]; !exists {
			if model, ok := c.models[wID]; ok {
				w, err := model.Writer()
				if err != nil {
					return err
				}
				c.writers[wID] = w
			}
		}
	}

	// -1 from reader and writer count since stdin and stdout are always registered
	fmt.Printf("Configured [%d] writers, [%d] readers and [%d] connections.\n", len(c.writers)-1, len(c.readers)-1, len(c.ConfConns))
	return nil
}

// Create the connection objects which contains the [io.ReadCloser] and its [io.WriteCloser].
// This will look up and resolve multiple writers per reader, and bundle them in a [io.MultiWriter].
func (c *Config) createConnections() {
	convertedReaders := make(map[string]bool)
	c.Connections = make([]Connection, 0)
	for _, conf := range c.ConfConns {
		if _, exists := convertedReaders[conf.ReaderID]; !exists {
			writer, writerIds := c.getWritersForReaderId(conf.ReaderID)
			convertedReaders[conf.ReaderID] = true

			if writer != nil {
				c.Connections = append(c.Connections, Connection{
					Reader:    c.readers[conf.ReaderID],
					ReaderId:  conf.ReaderID,
					Writer:    writer,
					WriterIds: writerIds,
				})
			} else {
				fmt.Printf("Resolved no matching writers for reader with id [%s]", conf.ReaderID)
			}
		}
	}
}

// Get all the [io.WriteCloser] that has the provided [string] as its registered [io.ReadCloser]. If only a single [io.WriteCloser] is resolved it will
// be returned, otherwise if there are multiple they will be wrapped in an [io.MultiWriter].
func (c Config) getWritersForReaderId(readerId string) (io.Writer, []string) {
	w := []io.Writer{}
	writerNames := []string{}
	for _, conf := range c.ConfConns {
		if conf.ReaderID == readerId {
			w = append(w, c.writers[conf.WriterID])
			writerNames = append(writerNames, conf.WriterID)
		}
	}

	if len(w) == 0 {
		return nil, writerNames
	} else if len(w) == 1 {
		return w[0], writerNames
	} else {
		return io.MultiWriter(w...), writerNames
	}
}

// Close all provided reader and writers
// Only the "first" occurring error will be returned, writers are closed first.
func (c Config) Close() error {
	var err error
	for _, w := range c.writers {
		e := w.Close()
		if e != nil && err == nil {
			err = e
		}
	}

	for _, r := range c.readers {
		e := r.Close()
		if e != nil && err == nil {
			err = e
		}
	}

	return err
}
