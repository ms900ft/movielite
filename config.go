package movielight

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Port            int
	Mode            string
	TMDBImageDir    string
	TMDBImageURL    string
	TMDBApiKey      string
	TargetDir       string
	SQLDebug        bool
	DataBase        string
	WebDav          bool
	Watchdirectory  string
	ScanDirectories []string
	ServerURL       string
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
	viper.SetDefault("TMDB.ImageURL", "https://image.tmdb.org/t/p")
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
	c.TMDBImageURL = viper.GetString("TMDB.ImageUrl")
	c.TMDBImageDir = viper.GetString("TMDB.ImageDir")
	c.TMDBApiKey = viper.GetString("TMDB.ApiKey")
	c.TargetDir = viper.GetString("TargetDirectory")
	c.DataBase = viper.GetString("DataBase.DBname")
	c.SQLDebug = viper.GetBool("SQLDebug")
	c.WebDav = viper.GetBool("WebDav")
	c.ScanDirectories = viper.GetStringSlice("Directories")
	c.Watchdirectory = viper.GetString("WatchDir")
	c.ServerURL = viper.GetString("MovieServerUrl")
	return &c
}
