package config

import (
	"errors"
	"strings"

	"smolgit/pkg/model"
)

type Config struct {
	LogColor           bool                     `koanf:"log.color"`
	LogJSON            bool                     `koanf:"log.json"`
	LogLevel           string                   `koanf:"log.level"`
	ServerAddr         string                   `koanf:"server.addr"`
	ServerAuthEnabled  bool                     `koanf:"server.auth.enabled"`
	ServerAuthAccounts []map[string]string      `koanf:"server.auth.accounts"`
	ServerBrand        string                   `koanf:"server.brand"`
	SSHAddr            string                   `koanf:"ssh.addr"`
	GitPath            string                   `koanf:"git.path"`
	GitBase            string                   `koanf:"git.base"`
	Users              []map[string]interface{} `koanf:"git.users"`
	Version            string
}

func (c *Config) Validate() error {
	// TODO
	return nil
}

func (c *Config) FindUserByKey(key string) (model.User, error) {
	for _, u := range c.Users {
		for _, k := range u["keys"].([]interface{}) {
			if strings.HasPrefix(k.(string), key) {
				return model.User{
					User: u["login"].(string),
				}, nil
			}
		}
	}
	return model.User{}, errors.New("user not found")
}
