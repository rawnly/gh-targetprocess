package config

import (
	"os"

	"github.com/charmbracelet/huh"
	"github.com/spf13/viper"
	"github.com/zalando/go-keyring"
)

type Config struct {
	URL   string `json:"url"`
	Token string `json:"token"`
}

func getUserName() string {
	return os.Getenv("USER")
}

const serviceName = "gh-targetprocess.access_token"

func Load() (*Config, error) {
	url := viper.GetString("url")
	token, err := keyring.Get(serviceName, getUserName())
	if err != nil {
		return nil, err
	}

	return &Config{URL: url, Token: token}, err
}

func (c *Config) Save() error {
	if err := keyring.Set(serviceName, getUserName(), c.Token); err != nil {
		return err
	}

	viper.Set("url", c.URL)
	return viper.WriteConfig()
}

func Init() error {
	var baseURL string
	var token string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Base URL").
				Value(&baseURL),
			huh.NewInput().
				Title("TP Access Token").
				Value(&token),
		),
	)

	if err := form.Run(); err != nil {
		return err
	}

	config := Config{
		URL:   baseURL,
		Token: token,
	}

	return config.Save()
}
