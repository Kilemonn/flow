package config

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gopkg.in/yaml.v3"
)

// Entry point to read in the provided file, resolve the connections, readers and writers and apply the configuration
func ApplyConfigurationFromFile(filepath string) {
	config, err := readConfig(filepath)
	if err != nil {
		fmt.Printf("Failed to apply configuration from filepath [%s]. Err: [%s].\n", filepath, err.Error())
		return
	}

	err = config.Initialise()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer config.Close()

	signalCtx, signalStopFunc := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer signalStopFunc()
	ctx, cancelFunc := context.WithCancel(signalCtx)

	// TODO: If we have stdin configured, we need to start another go routine that is grabbing content from stdin
	go applyConfig(ctx, cancelFunc, config.Conns, config.Settings)
	<-ctx.Done()
}

// Read and return a Config from the provided filepath
func readConfig(filePath string) (Config, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return Config{}, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)

	return config, err
}

func applyConfig(ctx context.Context, cancelFunc context.CancelFunc, connections []Connection, settings ConfigSettings) {
	// TODO: This needs to be smarter and understand the "flow" of information and call the correct reader and writers in the correct order
	startTime := time.Now()
	for {
		for _, connection := range connections {
			written, err := io.Copy(connection.Writer, connection.Reader)
			if err != nil {
				// TODO: Add a debug flag to enable this
				fmt.Printf("Error occurred when copying content from reader [%s] to writer(s) [%s]. Error: [%s]\n", connection.ReaderId, connection.WriterIds, err.Error())
			}

			if written > 0 {
				fmt.Printf("Wrote [%d] bytes from reader [%s] to writer(s) [%s].\n", written, connection.ReaderId, connection.WriterIds)
				startTime = time.Now()
			}
		}

		select {
		// This will only be detected if a OS signal is received
		case <-ctx.Done():
			return
		default:
			timeDifference := time.Since(startTime)
			if settings.Timeout > 0 && timeDifference.Seconds() >= float64(settings.Timeout) {
				cancelFunc()
				return
			}
		}
	}
}
