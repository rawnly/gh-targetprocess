package main

import (
	"fmt"
	"os"

	"github.com/cli/go-gh"
	"github.com/federicovitale-satispay/gh-targetprocess/config"
	"github.com/federicovitale-satispay/gh-targetprocess/targetprocess"
	"github.com/federicovitale-satispay/gh-targetprocess/utils"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigType("json")
	viper.SetConfigName("gh-targetprocess")
	viper.AddConfigPath("$HOME/.config")
	viper.SetEnvPrefix("GH_TARGETPROCESS")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			if err = viper.SafeWriteConfig(); err != nil {
				fmt.Println(err)
				return
			}

			if err = config.Init(); err != nil {
				fmt.Println(err)
				return
			}

			return
		}

		fmt.Println(err)
		return
	}
}

func main() {
	branch, _ := utils.CurrentBranch()

	cfg, err := config.Load()
	if err != nil {
		fmt.Println(err)
		return
	}

	var id *string
	if branch != "" {
		id = utils.GetTicketIDFromBranch(branch)
	}

	if id == nil {
		idOrURL := os.Args[1:][0]
		id = utils.ExtractIDFromURL(idOrURL)

		if id == nil {
			id = &idOrURL
		}
	}

	if id == nil {
		fmt.Println("Usage: gh targetprocess [task-id]")
		return
	}

	tp := targetprocess.New(cfg.URL, cfg.Token)
	assignable := targetprocess.Assignable{}
	if err = tp.Get(fmt.Sprintf("/v1/Assignables/%s", *id), &assignable); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(assignable.GetPRBody(cfg.URL))

	if _, _, err := gh.Exec("pr", "create", "--title", assignable.GetPRTitle(), "--body", assignable.GetPRBody(cfg.URL), "-w"); err != nil {
		fmt.Println(err)
		return
	}
}
