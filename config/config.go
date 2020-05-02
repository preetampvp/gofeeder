package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
)

type FeedConfig struct {
	Name               string `json:"name"`
	UrlBase            string `json:"urlBase"`
	UrlPath            string `json:"urlPath"`
	LinkSelector       string `json:"linkSelector"`
	NextPageSelector   string `json:"nextPageSelector"`
	NextLinkIsRelative bool   `json:"nextLinkIsRelative"`
}

type Config struct {
	Feeds []FeedConfig `json:"feeds"`
}

func GetConfig() *Config {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic("User home dir not found")
	}

	configFile := path.Join(homeDir, ".gofeeder.json")
	_, err = os.Stat(configFile)
	if os.IsNotExist(err) {
		panic(".gofeeder.json config file required in home dir")
	}

	config := Config{}
	file, err := ioutil.ReadFile(configFile)
	if err != nil {
		panic("Unable to read .gofeeder.json")
	}

	err = json.Unmarshal([]byte(file), &config)
	return &config
}
