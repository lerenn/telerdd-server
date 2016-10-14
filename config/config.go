package config

import (
	"errors"

	config "github.com/lerenn/go-config"
)

// New configuration for application
func New() (*config.Config, error) {
	conf := config.New()
	if err := conf.Read(ConfigFile); err != nil {
		return conf, errors.New("Error when reading config file.")
	}
	return conf, nil
}
