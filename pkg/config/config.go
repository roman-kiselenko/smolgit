package config

import (
	"bytes"
	_ "embed"
	"errors"
	"html/template"
	"strings"

	"smolgit/pkg/model"
)

func (c *Config) Validate() error {
	if c.ServerAddr == "" {
		return errors.New("ServerAddr is empty")
	}
	if c.ServerJWTKey == "" {
		return errors.New("ServerJWTKey is empty")
	}
	return nil
}

type Config struct {
	LogColor            bool                     `koanf:"log.color"`
	LogJSON             bool                     `koanf:"log.json"`
	LogLevel            string                   `koanf:"log.level"`
	ServerDisabled      bool                     `koanf:"server.disabled"`
	ServerJWTKey        string                   `koanf:"server.jwt_key"`
	ServiceAuthDisabled bool                     `koanf:"server.auth_disabled"`
	ServerAddr          string                   `koanf:"server.addr"`
	ServerAuthAccounts  []map[string]string      `koanf:"server.auth.accounts"`
	ServerBrand         string                   `koanf:"server.brand"`
	SSHAddr             string                   `koanf:"ssh.addr"`
	GitPath             string                   `koanf:"git.path"`
	GitBase             string                   `koanf:"git.base"`
	GitUsers            []map[string]interface{} `koanf:"git.users"`
	Version             string
}

// Open chat with model

func (c *Config) Prepare() error {
	if c.GitPath == "" {
		c.GitPath = "/usr/local/bin/git"
	}
	if c.GitBase == "" {
		c.GitBase = "/Users/roman.kiselenko/Documents/dev/opensource/smolgit"
	}
	return nil
}

func (c *Config) FindUserByKey(key string) (model.User, error) {
	for _, u := range c.GitUsers {
		for _, k := range u["keys"].([]interface{}) {
			p, _ := u["permissions"].(string)
			if strings.HasPrefix(k.(string), key) {
				return model.User{
					Name:        u["name"].(string),
					Permissions: p,
				}, nil
			}
		}
	}
	return model.User{}, errors.New("user not found")
}

func (c *Config) FindUserByName(name string) (model.User, error) {
	for _, u := range c.GitUsers {
		if u["name"].(string) == name {
			p, _ := u["permissions"].(string)
			return model.User{
				Name:        u["name"].(string),
				Permissions: p,
			}, nil
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
