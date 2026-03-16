/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/rawnly/gh-targetprocess/internal"
	"github.com/rawnly/gh-targetprocess/internal/logging"
	"github.com/rawnly/gh-targetprocess/internal/utils"
	"github.com/rawnly/gh-targetprocess/pkg/targetprocess"
	"github.com/spf13/cobra"
)

func NewInitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "init <tp-id>",
		SilenceErrors: true,
		SilenceUsage:  true,
		Example:       "init 210045",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			logf := logging.GetLogger(cmd.OutOrStdout())

			dryRun, err := cmd.Flags().GetBool("dry-run")
			if err != nil {
				return err
			}

			if len(args) == 0 {
				return errors.New("invalid assignable ID")
			}

			// 1 fetch the tp-id
			idOrUrl := args[0]
			id := utils.ExtractTicketID(&idOrUrl)
			tp := internal.GetTargetProcess(ctx)

			assignable, err := tp.GetAssignable(ctx, *id)
			if err != nil {
				return err
			}

			sanitized := strings.ToLower(strings.ReplaceAll(assignable.Name, " ", "_"))
			re := regexp.MustCompile(`[^a-z0-9_\-]`)
			sanitized = re.ReplaceAllString(sanitized, "")
			sanitized = regexp.MustCompile(`[_\-]{2,}`).ReplaceAllString(sanitized, "_")
			sanitized = strings.Trim(sanitized, "_-")
			if len(sanitized) > 60 {
				sanitized = strings.TrimRight(sanitized[:60], "_-")
			}

			branchName := strings.Join([]string{
				"feature",
				fmt.Sprintf("%s_%s", *id, sanitized),
			}, "/")

			if dryRun {
				logf("dry-run: %s\n", branchName)
				return nil
			}

			checkoutCmd := exec.Command("git", "checkout", "-b", branchName)
			checkoutCmd.Stdout = cmd.OutOrStdout()
			checkoutCmd.Stderr = cmd.OutOrStderr()

			if err := tp.UpdateState(ctx, assignable.ID, targetprocess.EntityStateInProgress); err != nil {
				return fmt.Errorf("updating US state: %w", err)
			}

			return checkoutCmd.Start()
		},
	}

	cmd.Flags().Bool("dry-run", false, "run without creating the branch")

	return cmd
}
