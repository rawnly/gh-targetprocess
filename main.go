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

func init_viper() error {
	viper.SetConfigType("json")
	viper.SetConfigName("config")
	viper.AddConfigPath("$HOME/.config/gh-targetprocess")
	viper.SetEnvPrefix("GH_TARGETPROCESS")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			if err = viper.SafeWriteConfig(); err != nil {
				return err
			}

			if err = config.Init(); err != nil {
				return err
			}

			return nil
		}

		return err
	}

	return nil
}

func main() {
	migrated, err := config.MigrateConfig()
	cobra.CheckErr(err)

	if migrated {
		fmt.Println()
		fmt.Println("Found an old configuration file. Migrated to the new config.")
		fmt.Println()
	}

	if err := init_viper(); err != nil {
		cobra.CheckErr(err)
	}

	cfg, err := config.Load()
	cobra.CheckErr(err)

	tp := targetprocess.New(cfg.URL, cfg.Token)

	ctx := context.Background()
	ctx = internal.InitContext(ctx, cfg, tp)

	cmd.Execute(ctx)
}
