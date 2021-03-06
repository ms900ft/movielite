package movielite

import (
	"bufio"
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Port                 int
	Mode                 string
	TMDBImageDir         string
	TMDBImageURL         string
	TMDBApiKey           string
	TargetDir            string
	SQLDebug             bool
	DataBase             string
	WebDav               bool
	Watchdirectory       string
	ScanDirectories      []string
	ServerURL            string
	Player               string
	UseAuthentication    bool
	Secret               string
	InitialAdminPassword string
}

func GetConfig(path string) *Config {
	viper.SetConfigName("movielite.yaml")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	// default values
	// xxx todo
	viper.SetDefault("MovieServerUrl", "http://localhost:8000")
	viper.SetDefault("Rescan.Delay", 10)
	viper.SetDefault("language", "de-DE")
	viper.SetDefault("Port", 8000)
	viper.SetDefault("TMDB.ImageURL", "https://image.tmdb.org/t/p")
	viper.SetDefault("Mode", "prod")
	viper.SetDefault("DataBase.DBname", "./movielite.db")
	viper.SetDefault("TMDB.ImageDir", "./images")
	viper.SetDefault("Player", "vlc")
	viper.SetDefault("UseAuthentication", true)
	viper.SetDefault("InitialAdminPassword", "password")
	if path != "" {
		file, err := os.Open(path)
		if err != nil {
			log.Fatalf("cannot read config file %s: %s", path, err)
		}

		viper.ReadConfig(bufio.NewReader(file))
	} else {
		err := viper.ReadInConfig() // Find and read the config file
		if err != nil {             // Handle errors reading the config file
			log.Fatalf("error reading config file: %s", err)
		}
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
	c.Player = viper.GetString("Player")
	c.Secret = viper.GetString("Secret")
	c.UseAuthentication = viper.GetBool("UseAuthentication")
	c.InitialAdminPassword = viper.GetString("InitialAdminPassword")

	return &c
}
