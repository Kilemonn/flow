package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Ensure the yaml unmarshalling is correct
func TestReadConfig(t *testing.T) {
	filepath := "connection.yaml"
	config, err := readConfig(filepath)
	require.NoError(t, err)

	require.NotNil(t, config)
	require.Equal(t, 2, len(config.Connections))

	require.Equal(t, 2, len(config.Nodes.Files))
	require.Equal(t, 2, len(config.Nodes.Sockets))
	require.Equal(t, 1, len(config.Nodes.Ipcs))
	require.Equal(t, 0, len(config.Nodes.Ports))

	require.NotNil(t, config.Settings)
	require.Equal(t, 5, config.Settings.Timeout)
}
