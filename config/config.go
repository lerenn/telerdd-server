package config

import (
	"errors"
	config "github.com/lerenn/go-config"
	cst "github.com/lerenn/telerdd-server/constants"
)

func New() (*config.Config, error) {
	conf := config.New()
	if err := conf.Read(cst.CONFIG_FILE); err != nil {
		return conf, errors.New("Error when reading config file.")
	}
	return conf, nil
}
