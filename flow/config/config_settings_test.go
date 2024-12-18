package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfigSettings_Default(t *testing.T) {
	settings := ConfigSettings{}
	require.Equal(t, 0, settings.Timeout)
}

func TestMapToArrayExistsOrNot(t *testing.T) {
	m := make(map[string][]string)
	v := m["key"]
	require.Equal(t, []string(nil), v)
	require.Equal(t, 0, len(v))
}
