package config

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func getConfigFile() string {
	homeDirectory, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	return homeDirectory + "/.gatorconfig.json"
}

func Read() (Config, error) {
	var c Config

	b, err := os.ReadFile(getConfigFile())
	if err != nil {
		return c, err
	}

	err = json.Unmarshal(b, &c)
	if err != nil {
		return c, err
	}
	return c, nil
}

func (c *Config) SetUser(user string) error {
	c.CurrentUserName = user
	data, err := json.Marshal(c)
	if err != nil {
		return err
	}

	return os.WriteFile(getConfigFile(), data, 0666)
}
