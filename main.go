package main

import (
	"context"
	"fmt"

	"github.com/rawnly/gh-targetprocess/cmd"
	"github.com/rawnly/gh-targetprocess/internal"
	"github.com/rawnly/gh-targetprocess/internal/config"
	"github.com/rawnly/gh-targetprocess/pkg/targetprocess"
	"github.com/spf13/cobra"
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
	cfg, err := config.Load()
	cobra.CheckErr(err)

	tp := targetprocess.New(cfg.URL, cfg.Token)

	ctx := context.Background()
	ctx = internal.InitContext(ctx, cfg, tp)

	cmd.Execute(ctx)
}
