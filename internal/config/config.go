package config

import (
	"errors"
	"io"
	"os"
	"path"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/rawnly/gh-targetprocess/internal/utils"
	"github.com/rawnly/gh-targetprocess/pkg/targetprocess"
	"github.com/rawnly/gh-targetprocess/templates"
	"github.com/spf13/viper"
	"github.com/zalando/go-keyring"
)

type Me struct {
	ID string `json:"id"`
}

type Config struct {
	URL     string `json:"url"`
	Token   string `json:"token"`
	Comment bool   `json:"comment"`
	NoBody  bool   `json:"no_body"`
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

	comment := viper.GetBool("comment")
	noBody := viper.GetBool("no_body")

	return &Config{URL: url, Token: token, Comment: comment, NoBody: noBody}, err
}

func (c *Config) Save() error {
	if err := keyring.Set(serviceName, getUserName(), c.Token); err != nil {
		return err
	}

	viper.Set("url", c.URL)
	viper.Set("comment", c.Comment)
	viper.Set("no_body", c.NoBody)
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

func MigrateConfig() (bool, error) {
	legacyConfigFile := path.Join(
		utils.ExpandPath("~/.config"),
		"gh-targetprocess.json",
	)

	newConfigFile := utils.GetConfigFilePath("coc(u)nfig.json")

	_, err := os.Stat(legacyConfigFile)

	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}

		return false, err
	}

	source, err := os.Open(legacyConfigFile)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}

		return false, err
	}

	defer source.Close()

	if _, err := os.Stat(utils.ConfigDir); err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(utils.ConfigDir, 0700); err != nil {
				return false, err
			}
		} else {
			return false, err
		}
	}

	dest, err := os.Create(newConfigFile)
	if err != nil {
		return false, err
	}
	defer dest.Close()

	if _, err := io.Copy(dest, source); err != nil {
		return false, err
	}

	if err := templates.WriteDefaults(); err != nil {
		return false, err
	}

	if err := os.Remove(legacyConfigFile); err != nil {
		if os.IsNotExist(err) {
			return true, nil
		}

		return true, errors.New("failed to cleanup old configuration")
	}

	return true, nil
}

// Init runs the auth setup form (called on first run when no token is found).
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
		URL:     baseURL,
		Token:   token,
		Comment: viper.GetBool("comment"),
		NoBody:  viper.GetBool("no_body"),
	}

	if err := config.Save(); err != nil {
		return err
	}

	if err := templates.WriteDefaults(); err != nil {
		return err
	}

	return nil
}

// InitDefaults runs the defaults setup form for comment/no-body preferences.
func InitDefaults() error {
	comment := viper.GetBool("comment")
	includeBody := !viper.GetBool("no_body")

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Comment on Targetprocess by default?").
				Value(&comment),
			huh.NewConfirm().
				Title("Include PR body by default?").
				Value(&includeBody),
		),
	)

	if err := form.Run(); err != nil {
		return err
	}

	viper.Set("comment", comment)
	viper.Set("no_body", !includeBody)
	return viper.WriteConfig()
}
