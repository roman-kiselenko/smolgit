package config

type Config struct {
	Database string `envconfig:"DATABASE" required:"true" default:"./database.sqlite"`
	Debug    bool   `envconfig:"DEBUG" default:"true"`
	Address  string `envconfig:"HTTP_ADDRESS" default:":8080"`
}
