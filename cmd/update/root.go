package update

import (
	"errors"

	"github.com/cli/go-gh/v2"
	"github.com/rawnly/gh-targetprocess/internal"
	"github.com/rawnly/gh-targetprocess/internal/utils"
	"github.com/spf13/cobra"
)

func init() {
	Cmd.Flags().BoolP("title", "t", false, "updates the title of the PR")
	Cmd.Flags().BoolP("body", "b", false, "updates the body of the PR")
}

var Cmd = &cobra.Command{
	Use:   "update",
	Short: "Update current PR with TargetProcess data",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		tp := internal.GetTargetProcess(ctx)
		config := internal.GetConfig(ctx)

		flags := cmd.Flags()

		shouldUpdateTitle, err := flags.GetBool("title")
		if err != nil {
			return err
		}

		shouldUpdateBody, err := flags.GetBool("body")
		if err != nil {
			return err
		}

		var rawId *string
		if len(args) > 0 {
			rawId = &args[0]
		}

		id := utils.ExtractTicketID(rawId)
		if id == nil {
			return errors.New("invalid ticket ID")
		}

		assignable, err := tp.GetAssignable(*id)
		if err != nil {
			return err
		}

		arguments := []string{"pr", "edit"}

		if shouldUpdateTitle {
			t := assignable.GetPRTitle()
			arguments = append(arguments, "--title", t)
		}

		if shouldUpdateBody {
			b := assignable.GetPRBody(config.URL)
			arguments = append(arguments, "--body", b)
		}

		if err := gh.ExecInteractive(ctx, arguments...); err != nil {
			// error is show by GH
			return nil
		}

		return nil
	},
}
