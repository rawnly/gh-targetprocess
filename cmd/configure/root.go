package configure

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/rawnly/gh-targetprocess/internal/config"
	"github.com/spf13/cobra"
)

var defaultsCmd = &cobra.Command{
	Use:   "defaults",
	Short: "Configure default flags (comment, body)",
	RunE: func(cmd *cobra.Command, args []string) error {
		return config.InitDefaults()
	},
}

var Cmd = &cobra.Command{
	Use:   "configure",
	Short: "Configure the gh-targetprocess CLI",
	RunE: func(cmd *cobra.Command, args []string) error {
		var confirmed bool
		confirm := huh.NewConfirm().
			Title("Are you sure you want to wipe your existing configuration?").
			Value(&confirmed)

		if err := confirm.Run(); err != nil {
			return err
		}

		if !confirmed {
			fmt.Println("Operation cancelled.")
			return nil
		}

		if err := config.Reset(); err != nil {
			return err
		}

		if err := config.Init(); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	Cmd.AddCommand(defaultsCmd)
}
