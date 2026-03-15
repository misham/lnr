package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Store provides access to persistent configuration.
type Store interface {
	Load() error
	Save() error
	TeamID() string
	SetTeamID(id string)
}

// ViperStore implements Store using Viper backed by a YAML config file.
type ViperStore struct {
	v   *viper.Viper
	dir string
}

// NewViperStore creates a Store that reads/writes config at dir/config.yaml.
func NewViperStore(dir string) *ViperStore {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(dir)
	v.SetEnvPrefix("LNR")
	v.AutomaticEnv()

	return &ViperStore{v: v, dir: dir}
}

func (s *ViperStore) Load() error {
	if err := s.v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil
		}
		return fmt.Errorf("read config: %w", err)
	}
	return nil
}

func (s *ViperStore) Save() error {
	if err := os.MkdirAll(s.dir, 0o700); err != nil { //nolint:gosec // directory needs 0700 for user access
		return fmt.Errorf("create config dir: %w", err)
	}
	if err := os.Chmod(s.dir, 0o700); err != nil { //nolint:gosec // directory needs 0700 for user access
		return fmt.Errorf("chmod config dir: %w", err)
	}

	configPath := filepath.Join(s.dir, "config.yaml")
	if err := s.v.WriteConfigAs(configPath); err != nil {
		return fmt.Errorf("write config: %w", err)
	}
	if err := os.Chmod(configPath, 0o600); err != nil {
		return fmt.Errorf("chmod config: %w", err)
	}
	return nil
}

func (s *ViperStore) TeamID() string {
	return s.v.GetString("team_id")
}

func (s *ViperStore) SetTeamID(id string) {
	s.v.Set("team_id", id)
}
