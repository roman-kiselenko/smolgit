package config

import (
	"bytes"
	_ "embed"
	"errors"
	"html/template"
	"strings"

	"smolgit/pkg/model"
)

type Config struct {
	LogColor            bool                `koanf:"log.color"`
	LogJSON             bool                `koanf:"log.json"`
	LogLevel            string              `koanf:"log.level"`
	ServerDisabled      bool                `koanf:"server.disabled"`
	ServerJWTKey        string              `koanf:"server.jwt_key"`
	ServiceAuthDisabled bool                `koanf:"server.auth_disabled"`
	ServerAddr          string              `koanf:"server.addr"`
	ServerAuthAccounts  []map[string]string `koanf:"server.auth.accounts"`
	ServerBrand         string              `koanf:"server.brand"`
	Users               []model.User        `koanf:"git.users"`
	SSHAddr             string              `koanf:"ssh.addr"`
	GitPath             string              `koanf:"git.path"`
	GitBase             string              `koanf:"git.base"`
	Version             string
}

type Users struct {
	Users map[string]model.User
}

func (c *Config) FindUserByKey(key string) (model.User, error) {
	for _, u := range c.Users {
		for _, k := range u.Keys {
			if strings.HasPrefix(k, key) {
				return u, nil
			}
		}
	}
	return model.User{}, errors.New("user not found")
}

func (c *Config) FindUserByName(name string) (model.User, error) {
	for _, u := range c.Users {
		if u.Name == name {
			return u, nil
		}
	}
	return model.User{}, errors.New("user not found")
}

//go:embed config.TEMPLATE.yaml
var configTemplate string

func GenerateConfig() ([]byte, error) {
	return executeTemplate(configTemplate)
}

func executeTemplate(tmpl string) ([]byte, error) {
	x, err := template.New("").Parse(tmpl)
	if err != nil {
		return nil, err
	}
	var b bytes.Buffer
	if err := x.Execute(&b, map[string]string{}); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}
