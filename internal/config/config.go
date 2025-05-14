package config

import (
	"github.com/charmbracelet/huh"
	"github.com/spf13/viper"
)

type Config struct {
	URL   string `json:"url"`
	Token string `json:"token"`
}

func Load() (*Config, error) {
	var C *Config
	err := viper.Unmarshal(&C)
	return C, err
}

func (c *Config) Save() error {
	viper.Set("url", c.URL)
	viper.Set("token", c.Token)
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
