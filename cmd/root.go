package cmd

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/charmbracelet/glamour"
	"github.com/cli/go-gh/v2"
	"github.com/rawnly/gh-targetprocess/cmd/comment"
	"github.com/rawnly/gh-targetprocess/cmd/configure"
	"github.com/rawnly/gh-targetprocess/cmd/update"
	"github.com/rawnly/gh-targetprocess/cmd/view"
	"github.com/rawnly/gh-targetprocess/internal"
	"github.com/rawnly/gh-targetprocess/internal/logging"
	"github.com/rawnly/gh-targetprocess/internal/utils"
	"github.com/rawnly/gh-targetprocess/pkg/targetprocess"
	"github.com/spf13/cobra"
)

func NewRootCMD() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "gh-targetprocess",
		Short: "gh-targetprocess is a tool to create PRs starting from a Targetprocess ID or branch",
		Example: `
  gh targetprocess --assignee @me
  gh targetprocess --base feature/stacked-base
  gh targetprocess --reviewer monalisa,hubot --reviewer myorg/team-name,
		`,
		Aliases:       []string{},
		Args:          cobra.MaximumNArgs(1),
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			stdout := cmd.OutOrStdout()

			logf := logging.GetLogger(stdout)

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

			httpCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()

			if err := tp.Get(httpCtx, fmt.Sprintf("/v1/Assignables/%s", *id), &assignable); err != nil {
				return err
			}

			title := assignable.GetPRTitle()
			body := assignable.GetPRBody(cfg.URL)
			flags := cmd.Flags()

			assignee, err := flags.GetString("assignee")
			if err != nil {
				return err
			}

			noBody := cfg.NoBody
			if flags.Changed("no-body") {
				noBody, err = flags.GetBool("no-body")
				if err != nil {
					return err
				}
			}

			dryRun, err := flags.GetBool("dry-run")
			if err != nil {
				return err
			}

			draft, err := flags.GetBool("draft")
			if err != nil {
				return err
			}

			label, err := flags.GetString("label")
			if err != nil {
				return err
			}

			reviewer, err := flags.GetString("reviewer")
			if err != nil {
				return err
			}

			milestone, err := flags.GetString("milestone")
			if err != nil {
				return err
			}

			base, err := flags.GetString("base")
			if err != nil {
				return err
			}

			web := cfg.Web
			if flags.Changed("web") {
				web, err = flags.GetBool("web")
				if err != nil {
					return err
				}
			}

			shouldComment := cfg.Comment
			if flags.Changed("comment") {
				shouldComment, err = flags.GetBool("comment")
				if err != nil {
					return err
				}
			}

			titleArg := title
			if dryRun {
				titleArg = "<title>"
			}

			prArgs := []string{"pr", "create", "--title", titleArg}

			if !noBody {
				b := body

				if dryRun {
					b = "<body>"
				}

				prArgs = append(prArgs, "--body", b)
			}

			if reviewer != "" {
				prArgs = append(prArgs, "--reviewer", reviewer)
			}

			if milestone != "" {
				prArgs = append(prArgs, "--milestone", milestone)
			}

			if base != "" {
				prArgs = append(prArgs, "--base", base)
			}

			if draft {
				prArgs = append(prArgs, "--draft")
			}

			if web && !draft {
				prArgs = append(prArgs, "--web")
			}

			if label != "" {
				prArgs = append(prArgs, "--label", label)
			}

			if assignee == "" {
				prArgs = append(prArgs, "-a", "@me")
			} else {
				prArgs = append(prArgs, "-a", assignee)
			}

			if dryRun {
				re, err := regexp.Compile(`\s+`)
				if err != nil {
					return err
				}

				args := strings.TrimSpace(re.ReplaceAllString(strings.Join(prArgs, " "), " "))

				logf("Running in dry-run")
				logf("Executing: `gh %s`", args)
				logf()
				logf()
				logf(title)

				if !noBody {
					r, err := glamour.NewTermRenderer(
						glamour.WithAutoStyle(),
					)
					cobra.CheckErr(err)

					rendered, err := r.Render(body)
					if err != nil {
						return fmt.Errorf("rendering body: %w", err)
					}

					fmt.Print(rendered)
				}
			} else {
				if err := gh.ExecInteractive(cmd.Context(), prArgs...); err != nil {
					return err
				}
			}

			if err := tp.UpdateState(httpCtx, assignable.ID, targetprocess.EntityStateInTest); err != nil {
				return fmt.Errorf("updating entity state: %w", err)
			}

			if !shouldComment {
				return nil
			}

			logf("Posting comment on Targetprocess...")

			if !dryRun {
				stdout, _, err := gh.Exec("pr", "view", "--json", "url", "-q", ".url")
				if err != nil {
					return err
				}

				url := stdout.String()

				commentBody := fmt.Sprintf("PR: %s", url)

				logf("Commented: ", commentBody)
				if err := tp.PostComment(ctx, commentBody, assignable.ID); err != nil {
					return err
				}
			}

			logf("Comment posted successfully")

			return nil
		},
	}

	cmd.AddCommand(view.Cmd)
	cmd.AddCommand(update.Cmd)
	cmd.AddCommand(configure.Cmd)
	cmd.AddCommand(NewInitCmd())
	cmd.AddCommand(comment.NewCommentCommand())

	cmd.Flags().BoolP("draft", "d", false, "mark pr as draft")
	cmd.Flags().BoolP("no-body", "", false, "skip body")
	cmd.Flags().BoolP("web", "w", false, "open pr in web browser")
	cmd.Flags().BoolP("dry-run", "", false, "dry-run pr creation")

	// mimic gh pr create [flags]
	cmd.Flags().StringP("assignee", "a", "", "Assign people by their login.")
	cmd.Flags().StringP("reviewer", "r", "", "Request reviews from people or teams by their handle")
	cmd.Flags().StringP("milestone", "m", "", "Add the pull request to a milestone by name")
	cmd.Flags().StringP("label", "l", "", "Add labels by name")
	cmd.Flags().StringP("base", "B", "", "The branch into which you want your code merged")

	cmd.Flags().BoolP("comment", "c", false, "comment on targetprocess US with the pull-request link")

	return cmd
}
