package config

import (
	"github.com/pelletier/go-toml"
	"os"
)

var tomlConfig *toml.TomlTree

// LoadConfigFile loads and parses the config file
func LoadConfigFile(filepath string) error {
	var err error
	tomlConfig, err = toml.LoadFile(filepath)
	return err
}

// Env returns the config value form within the current application environment
func Env() *toml.TomlTree {
	// get the current env
	env := os.Getenv("FOLKPOCKET_ENV")
	if env == "" {
		env = "dev"
	}

	// return the config tree for that env
	return tomlConfig.Get(env).(*toml.TomlTree)
}
