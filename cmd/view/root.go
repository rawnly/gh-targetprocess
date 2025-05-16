package view

import (
	"errors"

	"github.com/charmbracelet/glamour"
	"github.com/cli/browser"
	"github.com/rawnly/gh-targetprocess/internal"
	"github.com/rawnly/gh-targetprocess/internal/utils"
	"github.com/spf13/cobra"
)

func init() {
	Cmd.Flags().BoolP("web", "w", false, "Open the ticket in a web browser")
}

var Cmd = &cobra.Command{
	Use:   "view",
	Short: "View the current ticket",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		tp := internal.GetTargetProcess(ctx)
		config := internal.GetConfig(ctx)

		flags := cmd.Flags()

		web, err := flags.GetBool("web")
		cobra.CheckErr(err)

		var arg *string
		if len(args) > 0 {
			arg = &args[0]
		}

		id := utils.ExtractTicketID(arg)
		if id == nil {
			return errors.New("invalid ticket ID")
		}

		assignable, err := tp.GetAssignable(*id)
		cobra.CheckErr(err)

		if web {
			url := assignable.URL(config.URL)

			return browser.OpenURL(url)
		}

		glamour.Render(assignable.GetPRBody(config.URL), "")

		return nil
	},
}
