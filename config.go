package movielight

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Port int
	Mode string
}

func GetConfig() *Config {
	viper.SetConfigName("movielight")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	// default values
	// xxx todo
	viper.SetDefault("MovieServerUrl", "http://localhost:8000")
	viper.SetDefault("Rescan.Delay", 10)
	viper.SetDefault("language", "de-DE")
	viper.SetDefault("Port", 8000)
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("error config file: %s", err))
	}
	viper.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	// read values
	c := Config{}
	c.Port = viper.GetInt("Port")
	c.Mode = viper.GetString("Mode")
	return &c
}
