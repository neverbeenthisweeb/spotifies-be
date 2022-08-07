package config

import (
	"sync"

	"github.com/spf13/viper"
)

type Config struct {
	Spotify       Spotify `mapstructure:",squash"`
	Gin           Gin     `mapstructure:",squash"`
	ServerAddress string  `mapstructure:"SERVER_ADDRESS"`
	SessionKey    string  `mapstructure:"SESSION_KEY"`
}

type Spotify struct {
	RedirectURI string `mapstructure:"SPOTIFY_REDIRECT_URI"`
	ClientID    string `mapstructure:"SPOTIFY_CLIENT_ID"`
}

type Gin struct {
	Release bool `mapstructure:"GIN_RELEASE"`
}

var (
	config     *Config
	configLock = &sync.Mutex{}
)

func LoadConfig() Config {
	configLock.Lock()
	defer configLock.Unlock()

	if config != nil {
		return *config
	}

	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	// Override values that it has read from config file with
	// the values of the corresponding environment variables if they exist.
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		panic(err)
	}

	return *config
}
