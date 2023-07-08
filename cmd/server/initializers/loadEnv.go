// loadEnv.go
//
// PURPOSE:
// Use LoadConfig to load the configuration from environment variables in your application

// IMPORTANT:
// Set the appropriate mapstructure tags on the struct fields to specify the corresponding
// environment variable names.

package initializers

// The github.com/spf13/viper package provides a convenient way to manage
// configuration files and environment variables.
import (
	"github.com/spf13/viper"
)

type Config struct {
	DBHost         string `mapstructure:"POSTGRES_HOST"`
	DBUserName     string `mapstructure:"POSTGRES_USER"`
	DBUserPassword string `mapstructure:"POSTGRES_PASSWORD"`
	DBName         string `mapstructure:"POSTGRES_DB"`
	DBPort         string `mapstructure:"POSTGRES_PORT"`

	ClientOrigin string `mapstructure:"CLIENT_ORIGIN"`
}

func LoadConfig(path string) (config Config, err error) {
	// Set the directory where the configuration file is located
	viper.AddConfigPath(path)

	// Set the config. file type to environment variable format
	viper.SetConfigType("env")

	// Set the config. file name to "app", the name of the environment
	// variable file to be loaded (e.g., .env).
	viper.SetConfigName("app")

	// Add AutomaticEnv() to to automatically map environment variables
	// to the corresponding struct fields based on the mapstructure tags
	viper.AutomaticEnv()

	// Read the configuration file set with viper.AddConfigPath(path)
	// to load the configuration values from the environment variables
	// into the viper instance
	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	// Map the environment variable values to the corresponding fields
	// in the Config struct based on the mapstructure tags
	err = viper.Unmarshal(&config)
	return
}
