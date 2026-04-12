package configs

import (
	"github.com/spf13/viper"
)

type Config struct {
	DBDriver        string `mapstructure:"DB_DRIVER"`
	DBHost          string `mapstructure:"DB_HOST"`
	DBPort          string `mapstructure:"DB_PORT"`
	DBUser          string `mapstructure:"DB_USER"`
	DBPassword      string `mapstructure:"DB_PASSWORD"`
	DBName          string `mapstructure:"DB_NAME"`
	WebServerPort   string `mapstructure:"WEB_SERVER_PORT"`
	AppClientID     string `mapstructure:"APP_CLIENT_ID"`
	AppClientSecret string `mapstructure:"APP_CLIENT_SECRET"`
	Dns             string `mapstructure:"DNS"`
	HashKey         string `mapstructure:"HASH_KEY"`
	BlockKey        string `mapstructure:"BLOCK_KEY"`
	EncryptionKey   string `mapstructure:"ENCRYPTION_KEY"`
}

func LoadConfig(path string) *Config {
	var cfg *Config

	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}

	return cfg
}
