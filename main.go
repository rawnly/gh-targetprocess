package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/rawnly/gh-targetprocess/cmd"
	"github.com/rawnly/gh-targetprocess/internal"
	"github.com/rawnly/gh-targetprocess/internal/config"
	"github.com/rawnly/gh-targetprocess/pkg/targetprocess"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init_viper(ctx context.Context) error {
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

			if err = config.Init(ctx); err != nil {
				return err
			}

			return nil
		}

		return err
	}

	return nil
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM)

	defer func() {
		<-sigChan
		cancel()
	}()

	migrated, err := config.MigrateConfig()
	cobra.CheckErr(err)

	if migrated {
		fmt.Println()
		fmt.Println("Found an old configuration file. Migrated to the new config.")
		fmt.Println()
	}

	if err := init_viper(ctx); err != nil {
		cobra.CheckErr(err)
	}

	cfg, err := config.Load(ctx)
	cobra.CheckErr(err)

	tp := targetprocess.New(cfg.URL, cfg.Token)

	ctx = internal.InitContext(ctx, cfg, tp)

	root := cmd.NewRootCMD()
	err = root.ExecuteContext(ctx)

	if err != nil {
		switch {
		case strings.Contains(err.Error(), "unknown command") || strings.Contains(err.Error(), "unknown flag"):
			showSuggestion(root, err)
		default:
			fmt.Fprintln(root.OutOrStderr(), err)
		}
	}

	cancel()
	os.Exit(1)
}

func showSuggestion(cmd *cobra.Command, err error) {
	fmt.Fprintf(cmd.OutOrStderr(), cmd.UsageString())
	fmt.Fprintf(cmd.OutOrStderr(), "\nError: Invalid ussage: %v\n", err)
}
