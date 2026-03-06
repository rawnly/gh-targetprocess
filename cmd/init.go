/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/rawnly/gh-targetprocess/internal"
	"github.com/rawnly/gh-targetprocess/internal/utils"
	"github.com/spf13/cobra"
)

func NewInitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "init <tp-id>",
		SilenceErrors: true,
		SilenceUsage:  true,
		Example:       "init 210045",
		Args:          cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := context.WithTimeout(cmd.Context(), time.Second*5)
			defer cancel()

			dryRun, err := cmd.Flags().GetBool("dry-run")
			if err != nil {
				return err
			}

			// 1 fetch the tp-id
			idOrUrl := args[0]
			id := utils.ExtractTicketID(&idOrUrl)
			tp := internal.GetTargetProcess(ctx)

			assignable, err := tp.GetAssignable(ctx, *id)
			if err != nil {
				return err
			}

			branchName := strings.Join([]string{
				"feature",
				*id,
				strings.ToLower(strings.ReplaceAll(assignable.Name, " ", "_")),
			}, "/")

			if dryRun {
				fmt.Fprintf(cmd.OutOrStdout(), "dry-run: %s\n", branchName)
				return nil
			}

			checkoutCmd := exec.Command("git", "checkout", "-b", branchName)
			checkoutCmd.Stdout = cmd.OutOrStdout()
			checkoutCmd.Stderr = cmd.OutOrStderr()

			return checkoutCmd.Start()
		},
	}

	cmd.Flags().Bool("dry-run", false, "run without creating the branch")

	return cmd
}
