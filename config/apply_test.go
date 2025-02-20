package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

// Ensure the yaml unmarshalling is correct
func TestReadConfig(t *testing.T) {
	filepath := "connection.yaml"
	config, err := readConfig(filepath)
	require.NoError(t, err)

	require.NotNil(t, config)
	require.Equal(t, 4, len(config.Connections))

	require.Equal(t, 2, len(config.Nodes.Files))
	require.Equal(t, 2, len(config.Nodes.Sockets))
	require.Equal(t, 1, len(config.Nodes.Ipcs))
	require.Equal(t, 0, len(config.Nodes.Ports))

	require.NotNil(t, config.Settings)
	require.Equal(t, 3, config.Settings.Timeout)
}

func TestReadConfig_invalidFile(t *testing.T) {
	filepath := "somefile.yaml"
	_, err := os.Stat(filepath)
	require.Error(t, err)

	_, err = readConfig(filepath)
	require.Error(t, err)
}

func TestApplyConfigurationFromFile_invalidFile(t *testing.T) {
	filepath := "somefile.yaml"
	_, err := os.Stat(filepath)
	require.Error(t, err)

	ApplyConfigurationFromFile(filepath)
}

func TestApplyConfigurationFromFile(t *testing.T) {
	filepath := "connection.yaml"
	data := "TestApplyConfigurationFromFile"

	inputFile := "input.txt"
	require.NoError(t, os.WriteFile(inputFile, []byte(data), 0666))
	outputFile := "output.txt"
	defer os.Remove(inputFile)
	defer os.Remove(outputFile)

	ApplyConfigurationFromFile(filepath)

	read, err := os.ReadFile(outputFile)
	require.NoError(t, err)
	require.Equal(t, data, string(read))
}
