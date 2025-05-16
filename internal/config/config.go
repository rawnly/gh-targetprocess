package config

import (
	"errors"
	"os"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/rawnly/gh-targetprocess/pkg/targetprocess"
	"github.com/spf13/viper"
	"github.com/zalando/go-keyring"
)

type Me struct {
	ID string `json:"id"`
}

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
		return nil, Init()
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

func Reset() error {
	if err := keyring.Delete(serviceName, getUserName()); err != nil {
		return err
	}

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	viper.Set("url", "")
	if err := viper.WriteConfig(); err != nil {
		return err
	}

	return nil
}

func Init() error {
	var baseURL string
	var token string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Base URL").
				Placeholder("https://<my-company>.tpondemand.com").
				Validate(func(s string) error {
					if s == "" {
						return errors.New("base URL cannot be empty")
					}

					if !strings.HasPrefix(s, "http:") && !strings.HasPrefix(s, "https:") {
						return errors.New("URL must start with http:// or https://")
					}

					return nil
				}).
				Value(&baseURL),
			huh.NewInput().
				Title("Access Token").
				Description("https://www.ibm.com/docs/en/app-connect/12.0.x?topic=t-obtaining-connection-values-targetprocess").
				Validate(func(token string) error {
					tp := targetprocess.New(baseURL, token)

					if err := tp.Test("/v1/Users/loggeduser"); err != nil {
						return err
					}

					return nil
				}).
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
