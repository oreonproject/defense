// oreon/defense · watchthelight <wtl>

package config

import (
	"os"
	"path/filepath"
)

const (
	SystemConfigPath = "/etc/oreon/defense.toml"
	SocketPath       = "/run/oreon/defense.sock"
	LogPath          = "/var/log/oreon/defense.log"
	DataPath         = "/var/lib/oreon/defense"
	QuarantinePath   = "/var/lib/oreon/defense/quarantine"
	DatabasePath     = "/var/lib/oreon/defense/defense.db"
)

func UserConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".config", "oreon", "defense.toml")
}

// EnsureDirectories creates all required directories with appropriate permissions.
// Should be called by the daemon on startup.
func EnsureDirectories() error {
	dirs := []struct {
		path string
		perm os.FileMode
	}{
		{"/etc/oreon", 0755},
		{"/run/oreon", 0755},
		{"/var/log/oreon", 0755},
		{DataPath, 0700},
		{QuarantinePath, 0700},
	}

	for _, d := range dirs {
		if err := os.MkdirAll(d.path, d.perm); err != nil {
			return err
		}
	}
	return nil
}

// ExpandPath expands ~ to the user's home directory.
func ExpandPath(path string) string {
	if len(path) == 0 || path[0] != '~' {
		return path
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return path
	}
	return filepath.Join(home, path[1:])
}
