package config

import (
	"context"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Kilemonn/flow/flow/testutil"
	"github.com/stretchr/testify/require"
	goSerial "go.bug.st/serial"
)

func getTestStruct() Config {
	return Config{
		Connections: []ConfigConnection{
			{
				ReaderID: "InputFile1",
				WriterID: "Serial1",
			},
			{
				ReaderID: "Serial1",
				WriterID: "OutputFile1",
			},
		},
		Nodes: ConfigNodes{
			Ports: []ConfigPort{
				{
					ID:      "Serial1",
					Channel: "COM4",
					Mode: goSerial.Mode{
						BaudRate: 9600,
						DataBits: 8,
					},
				},
			},
			Files: []ConfigFile{
				{
					ID:   "InputFile1",
					Path: "inputfile1.txt",
				},
				{
					ID:   "OutputFile1",
					Path: "outputfile1.txt",
				},
			},
			Sockets: []ConfigSocket{
				{
					ID:       "socketID",
					Protocol: "UDP",
					Port:     54321,
					Address:  "127.0.0.1",
				},
			},
			Ipcs: []ConfigIPC{
				{
					ID:      "IPC1",
					Channel: "Channel1",
				},
			},
		},
		Conns: []Connection{},
	}
}

func TestDuplicateIDs_NoDuplicates(t *testing.T) {
	nodes := ConfigNodes{
		Ports: []ConfigPort{
			{
				ID: "myPortID",
			},
			{
				ID: "myPortID2",
			},
		},
		Files: []ConfigFile{
			{
				ID: "myFileID",
			},
			{
				ID: "myFileID2",
			},
		},
	}

	c := Config{Nodes: nodes}
	err := c.componentIDsAreUnique()
	require.NoError(t, err)
}

func TestDuplicateIDs_WithDuplicates(t *testing.T) {
	ports := []ConfigPort{
		{
			ID: "myPortID",
		},
		{
			ID: "myID",
		},
	}

	files := []ConfigFile{
		{
			ID: "myFileID",
		},
		{
			ID: "myID",
		},
	}

	c := &Config{Nodes: ConfigNodes{Ports: ports, Files: files}}
	err := c.componentIDsAreUnique()
	require.Error(t, err)

	c = &Config{Nodes: ConfigNodes{Files: files}}
	err = c.componentIDsAreUnique()
	require.NoError(t, err)

	c = &Config{Nodes: ConfigNodes{Ports: ports}}
	err = c.componentIDsAreUnique()
	require.NoError(t, err)

	ports = []ConfigPort{
		{
			ID: "PortID",
		},
		{
			ID: "PortID",
		},
	}

	c = &Config{Nodes: ConfigNodes{Ports: ports}}
	err = c.componentIDsAreUnique()
	require.Error(t, err)

	files = []ConfigFile{
		{
			ID: "FileID",
		},
		{
			ID: "FileID",
		},
	}

	c = &Config{Nodes: ConfigNodes{Files: files}}
	err = c.componentIDsAreUnique()
	require.Error(t, err)
}

func TestYamlReadWrite(t *testing.T) {
	testutil.WithTempFile(t, func(fileName string) {
		initialConfig := getTestStruct()
		var finalConfig Config

		err := initialConfig.writeConfig(fileName)
		require.NoError(t, err)

		finalConfig, err = readConfig(fileName)
		require.NoError(t, err)

		require.Equal(t, initialConfig, finalConfig)
	})
}

func TestCombineToReadersAndWriters(t *testing.T) {
	testutil.WithTempFile(t, func(inputFile string) {
		testutil.WithTempFile(t, func(outputFile string) {
			connections := []ConfigConnection{
				{
					ReaderID: "InputFile",
					WriterID: "OutputFile",
				},
				{
					ReaderID: "OutputFile",
					WriterID: "stdout",
				},
			}

			fileConfig := []ConfigFile{
				{
					ID:   "InputFile",
					Path: inputFile,
				},
				{
					ID:   "OutputFile",
					Path: outputFile,
				},
			}
			config := &Config{
				Connections: connections,
				Nodes: ConfigNodes{
					Files: fileConfig,
				},
			}
			err := config.Initialise()
			require.NoError(t, err)

			require.NotNil(t, config.readers[StdIn])
			require.NotNil(t, config.readers["InputFile"])
			require.NotNil(t, config.readers["OutputFile"])
			require.Nil(t, config.writers["InputFile"])

			require.NotNil(t, config.writers[StdOut])
			require.NotNil(t, config.writers["OutputFile"])
		})
	})
}

func TestApplyConfig_FileToFile(t *testing.T) {
	testutil.WithTempFile(t, func(inputFile string) {
		testutil.WithTempFile(t, func(outputFile string) {

			connections := []ConfigConnection{
				{
					ReaderID: "InputFile",
					WriterID: "OutputFile",
				},
			}

			fileConfig := []ConfigFile{
				{
					ID:   "InputFile",
					Path: inputFile,
				},
				{
					ID:   "OutputFile",
					Path: outputFile,
				},
			}

			config := Config{
				Connections: connections,
				Nodes: ConfigNodes{
					Files: fileConfig,
				},
			}
			err := config.Initialise()
			require.NoError(t, err)

			content := "Wow some great content to write"
			err = os.WriteFile(inputFile, []byte(content), os.ModeAppend)
			require.NoError(t, err)

			settings := ConfigSettings{Timeout: 1}
			testutil.TakesAtleast(t, time.Duration(settings.Timeout*int(time.Second)), func() {
				ctx, cancelFunc := context.WithCancel(context.Background())
				defer config.Close()
				go applyConfig(ctx, cancelFunc, config.Conns, settings)
				<-ctx.Done()
			})

			read, err := os.ReadFile(outputFile)
			require.NoError(t, err)
			require.Equal(t, content, string(read))
		})
	})
}

func TestApplyConfig_StdInToFileToStdOut(t *testing.T) {
	content := "TestApplyConfig_StdInToFileToStdOut"
	testutil.WithBytesInStdIn(t, []byte(content), func() {
		testutil.WithTempFile(t, func(file string) {
			stdout := testutil.CaptureStdout(t, func() {
				fileConfig := []ConfigFile{
					{
						ID:   "FileID",
						Path: file,
					},
				}

				connections := []ConfigConnection{
					{
						ReaderID: "stdin",
						WriterID: "FileID",
					},
					{
						ReaderID: "FileID",
						WriterID: "stdout",
					},
				}

				config := Config{
					Connections: connections,
					Nodes: ConfigNodes{
						Files: fileConfig,
					},
				}
				err := config.Initialise()
				require.NoError(t, err)

				settings := ConfigSettings{Timeout: 1}
				testutil.TakesAtleast(t, time.Duration(settings.Timeout*int(time.Second)), func() {
					ctx, cancelFunc := context.WithCancel(context.Background())
					defer config.Close()
					go applyConfig(ctx, cancelFunc, config.Conns, settings)
					<-ctx.Done()
				})
			})

			writtenToStdOut, err := io.ReadAll(stdout)
			require.NoError(t, err)

			writtenToFile, err := os.ReadFile(file)
			require.NoError(t, err)

			// Since there is other things being printed to std out, we want to make sure that the stdout contains the content
			require.True(t, strings.Contains(string(writtenToStdOut), content))
			require.Equal(t, content, string(writtenToFile))
		})
	})
}

// A scenario where a writer is defined in multiple connections
func TestApplyConfig_MultipleWriters(t *testing.T) {
	content := "TestApplyConfig_MultipleWriters"
	testutil.WithBytesInStdIn(t, []byte(content), func() {
		testutil.WithTempFile(t, func(file string) {
			testutil.WithTempFile(t, func(file2 string) {
				fileConfig := []ConfigFile{
					{
						ID:   "File1",
						Path: file,
					},
					{
						ID:   "File2",
						Path: file2,
					},
				}

				connections := []ConfigConnection{
					{
						ReaderID: "stdin",
						WriterID: "File1",
					},
					{
						ReaderID: "stdin",
						WriterID: "File2",
					},
					{
						ReaderID: "File1",
						WriterID: "File2",
					},
				}

				config := Config{
					Connections: connections,
					Nodes: ConfigNodes{
						Files: fileConfig,
					},
				}

				err := config.Initialise()
				require.NoError(t, err)

				settings := ConfigSettings{Timeout: 1}
				testutil.TakesAtleast(t, time.Duration(settings.Timeout*int(time.Second)), func() {
					ctx, cancelFunc := context.WithCancel(context.Background())
					defer config.Close()
					go applyConfig(ctx, cancelFunc, config.Conns, settings)
					<-ctx.Done()
				})

				writtenToFile, err := os.ReadFile(file)
				require.NoError(t, err)

				writtenToFile2, err := os.ReadFile(file2)
				require.NoError(t, err)

				require.Equal(t, content, string(writtenToFile))
				require.Equal(t, content+content, string(writtenToFile2))
			})
		})
	})
}

func TestApplyConfig_WithUDPSockets(t *testing.T) {
	content := "TestApplyConfig_WithUDPSockets"
	testutil.WithBytesInStdIn(t, []byte(content), func() {
		testutil.WithTempFile(t, func(outputFile string) {
			fileConfig := []ConfigFile{
				{
					ID:   "outputfile",
					Path: outputFile,
				},
			}

			socketConfig := []ConfigSocket{
				{
					ID:       "sender-socket",
					Protocol: "udp",
					Port:     45231,
					Address:  "127.0.0.1",
				},
				{
					ID:       "recv-socket",
					Protocol: "udp",
					Port:     45231,
					Address:  "127.0.0.1",
				},
			}

			connections := []ConfigConnection{
				{
					ReaderID: "stdin",
					WriterID: "sender-socket",
				},
				{
					ReaderID: "recv-socket",
					WriterID: "outputfile",
				},
			}

			config := Config{
				Connections: connections,
				Nodes: ConfigNodes{
					Files:   fileConfig,
					Sockets: socketConfig,
				},
			}

			err := config.Initialise()
			require.NoError(t, err)

			settings := ConfigSettings{Timeout: 2}
			testutil.TakesAtleast(t, time.Duration(settings.Timeout*int(time.Second)), func() {
				ctx, cancelFunc := context.WithCancel(context.Background())
				defer config.Close()
				go applyConfig(ctx, cancelFunc, config.Conns, settings)
				<-ctx.Done()
			})

			writtenToFile, err := os.ReadFile(outputFile)
			require.NoError(t, err)

			require.Equal(t, content, string(writtenToFile))

		})
	})
}

// Ensure that all the data from multiple TCP connections are read in from the TCP reader.
func TestApplyConfig_WithMultipleTCPSockets(t *testing.T) {
	content := "TestApplyConfig_WithTCPSockets"
	socketPort := uint16(64621)
	testutil.WithBytesInStdIn(t, []byte(content), func() {
		testutil.WithTempFile(t, func(outputFile string) {
			fileConfig := []ConfigFile{
				{
					ID:   "outputfile",
					Path: outputFile,
				},
			}

			socketConfig := []ConfigSocket{
				{
					ID:       "sender-socket",
					Protocol: "tcp",
					Port:     socketPort,
					Address:  "127.0.0.1",
				},
				{
					ID:       "recv-socket",
					Protocol: "tcp",
					Port:     socketPort,
					Address:  "127.0.0.1",
				},
				{
					ID:       "sender-socket2",
					Protocol: "tcp",
					Port:     socketPort,
					Address:  "127.0.0.1",
				},
			}

			connections := []ConfigConnection{
				{
					ReaderID: "stdin",
					WriterID: "sender-socket",
				},
				{
					ReaderID: "stdin",
					WriterID: "sender-socket2",
				},
				{
					ReaderID: "recv-socket",
					WriterID: "outputfile",
				},
			}

			config := Config{
				Connections: connections,
				Nodes: ConfigNodes{
					Files:   fileConfig,
					Sockets: socketConfig,
				},
			}

			err := config.Initialise()
			require.NoError(t, err)

			settings := ConfigSettings{Timeout: 1}
			testutil.TakesAtleast(t, time.Duration(settings.Timeout*int(time.Second)), func() {
				ctx, cancelFunc := context.WithCancel(context.Background())
				defer config.Close()
				go applyConfig(ctx, cancelFunc, config.Conns, settings)
				<-ctx.Done()
			})

			writtenToFile, err := os.ReadFile(outputFile)
			require.NoError(t, err)

			require.Equal(t, content+content, string(writtenToFile))
		})
	})
}

// Ensure that all data from multiple IPC connections is received and forwarded.
func TestApplyConfig_WithMultipleIPC(t *testing.T) {
	content := "TestApplyConfig_WithIPC"
	testutil.WithBytesInStdIn(t, []byte(content), func() {
		testutil.WithTempFile(t, func(outputFile string) {
			fileConfig := []ConfigFile{
				{
					ID:   "outputfile",
					Path: outputFile,
				},
			}

			ipcConfig := []ConfigIPC{
				{
					ID:      "sender-ipc",
					Channel: "config-test-channel",
				},
				{
					ID:      "sender-ipc2",
					Channel: "config-test-channel",
				},
				{
					ID:      "recv-ipc",
					Channel: "config-test-channel",
				},
			}

			connections := []ConfigConnection{
				{
					ReaderID: "stdin",
					WriterID: "sender-ipc",
				},
				{
					ReaderID: "stdin",
					WriterID: "sender-ipc2",
				},
				{
					ReaderID: "recv-ipc",
					WriterID: "outputfile",
				},
			}

			config := Config{
				Connections: connections,
				Nodes: ConfigNodes{
					Files: fileConfig,
					Ipcs:  ipcConfig,
				},
			}

			err := config.Initialise()
			require.NoError(t, err)

			settings := ConfigSettings{Timeout: 1}
			testutil.TakesAtleast(t, time.Duration(settings.Timeout*int(time.Second)), func() {
				ctx, cancelFunc := context.WithCancel(context.Background())
				defer config.Close()
				go applyConfig(ctx, cancelFunc, config.Conns, settings)
				<-ctx.Done()
			})

			writtenToFile, err := os.ReadFile(outputFile)
			require.NoError(t, err)

			require.Equal(t, content+content, string(writtenToFile))
		})
	})
}
