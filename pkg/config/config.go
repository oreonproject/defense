// oreon/defense · watchthelight <wtl>

package config

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Config struct {
	General       General       `toml:"general"`
	Firewall      Firewall      `toml:"firewall"`
	Notifications Notifications `toml:"notifications"`
	Scanning      Scanning      `toml:"scanning"`
	ClamAV        ClamAV        `toml:"clamav"`
}

type General struct {
	RealTimeProtection bool   `toml:"real_time_protection"`
	LogLevel           string `toml:"log_level"`
}

type Firewall struct {
	Enabled bool `toml:"enabled"`
}

type Notifications struct {
	Level string `toml:"level"`
}

type Scanning struct {
	Exclusions     []string `toml:"exclusions"`
	QuickScanPaths []string `toml:"quick_scan_paths"`
}

type ClamAV struct {
	SocketPath string `toml:"socket_path"`
}

func Default() *Config {
	return &Config{
		General: General{
			RealTimeProtection: true,
			LogLevel:           "info",
		},
		Firewall: Firewall{
			Enabled: true,
		},
		Notifications: Notifications{
			Level: "all",
		},
		Scanning: Scanning{
			Exclusions: []string{},
			QuickScanPaths: []string{
				"/tmp",
				"/var/tmp",
			},
		},
		ClamAV: ClamAV{
			SocketPath: "/var/run/clamav/clamd.sock",
		},
	}
}

func Load(path string) (*Config, error) {
	cfg := Default()

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, err
	}

	if _, err := toml.Decode(string(data), cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) Save(path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return toml.NewEncoder(f).Encode(c)
}
