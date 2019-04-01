package config

import (
	"github.com/BurntSushi/toml"
	"path"
	"runtime"
	"sync"
)

type Config struct {
	Prod *baseConfig  // Production Config
}

type baseConfig struct {
	Host          string
	Port          uint16
	ServerCrt     string `toml:"server_crt"`
	ServerKey     string `toml:"server_key"`
	ClientCrt     string `toml:"client_crt"`
	ClientKey     string `toml:"client_key"`
	CertAuth      string `toml:"cert_auth"`
}

var once sync.Once
var conf *Config
var err error

func GetConfig() (*Config, error) {
	once.Do(func() {
		_, filename, _, _ := runtime.Caller(0)
		directory := path.Dir(filename) // Current Package path
		conf = &Config{}
		_, err = toml.DecodeFile(path.Join(directory, "./config.toml"), conf)
	})

	return conf, err
}
