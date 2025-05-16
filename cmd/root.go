package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/cli/go-gh"
	"github.com/rawnly/gh-targetprocess/cmd/configure"
	"github.com/rawnly/gh-targetprocess/cmd/view"
	"github.com/rawnly/gh-targetprocess/internal"
	"github.com/rawnly/gh-targetprocess/internal/utils"
	"github.com/rawnly/gh-targetprocess/pkg/targetprocess"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(view.Cmd)
	rootCmd.AddCommand(configure.Cmd)
}

var rootCmd = &cobra.Command{
	Use:        "gh-targetprocess",
	Short:      "gh-targetprocess is a tool to create PRs starting from a Targetprocess ID or branch",
	Example:    "gh targetprocess 12345",
	ArgAliases: []string{"id", "url"},
	// DisableFlagParsing: true,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		cfg := internal.GetConfig(ctx)
		tp := internal.GetTargetProcess(ctx)

		// we ignore the error the directory may not be a git repo
		branch, _ := utils.CurrentBranch()

		var id *string
		if branch != "" {
			id = utils.GetTicketIDFromBranch(branch)
		}

		if len(args) > 0 {
			if id == nil {
				idOrURL := args[0]
				id = utils.ExtractIDFromURL(idOrURL)

				if id == nil {
					id = &idOrURL
				}
			}
		}

		if id == nil {
			return errors.New("invalid ticket id or url")
		}

		assignable := targetprocess.Assignable{}
		if err := tp.Get(fmt.Sprintf("/v1/Assignables/%s", *id), &assignable); err != nil {
			return err
		}

		if _, _, err := gh.Exec("pr", "create", "--title", assignable.GetPRTitle(), "--body", assignable.GetPRBody(cfg.URL), "-w"); err != nil {
			return err
		}

		return nil
	},
}

func Execute(ctx context.Context) {
	err := rootCmd.ExecuteContext(ctx)
	if err != nil {
		os.Exit(1)
	}
}
