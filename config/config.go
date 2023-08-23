package config

import (
	"os"

	"gopkg.in/validator.v2"
	"gopkg.in/yaml.v3"
)

const CONFIG_PATH = "./config.yml"

// Configuration represents the whole user configuration configured in a yaml file.
type Configuration struct {
	Server struct {
		ListenAddr   string
		Port         string `validate:"nonzero"`
		AllowOrigins []string
	}

	Database struct {
		Dialect    string `validate:"nonzero"`
		Connection string `validate:"nonzero"`
	}

	SMTP struct {
		Username string `validate:"nonzero"`
		Password string `validate:"nonzero"`
		Host     string `validate:"nonzero"`
		Port     string `validate:"nonzero"`
		From     string `validate:"nonzero"`
	}

	Security struct {
		TokenSize int `validate:"nonzero"`
	}
}

// Get reads the configuration file and returns its content as a Configuration.
func Get() *Configuration {
	content, err := os.ReadFile(CONFIG_PATH)
	if err != nil {
		panic(err)
	}

	config := new(Configuration)
	if err := yaml.Unmarshal(content, &config); err != nil {
		panic(err)
	}

	if err = validator.Validate(config); err != nil {
		panic(err)
	}

	return config
}
