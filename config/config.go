package config

import (
	"os"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
)

const CONFIG_PATH = "./config.yml"

// Configuration represents the whole user configuration configured in a yaml file.
type Configuration struct {
	Server struct {
		ListenAddr   string
		Port         string   `validate:"required"`
		AllowOrigins []string 

		SSL struct {
			Enabled         bool
			RedirectToHttps bool  
			Port            string `validate:"required_if=Enabled true"`
			CertFile        string `validate:"required_if=Enabled true"`
			CertKey         string `validate:"required_if=Enabled true"`
		}
	}

	Database struct {
		Dialect    string `validate:"required"`
		Connection string `validate:"required"`
	}

	SMTP struct {
		Username string `validate:"required"`
		Password string `validate:"required"`
		Host     string `validate:"required,fqdn"`
		Port     string `validate:"required"`
		From     string `validate:"required,email"`
	}

	Security struct {
		TokenSize uint8 `validate:"required,gte=12"`
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

	if err = validator.New().Struct(config); err != nil {
		panic(err)
	}

	return config
}
