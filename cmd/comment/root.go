package comment

import (
	"bytes"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"text/template"

	"github.com/cli/go-gh/v2"
	"github.com/rawnly/gh-targetprocess/internal"
	"github.com/rawnly/gh-targetprocess/internal/logging"
	"github.com/rawnly/gh-targetprocess/internal/utils"
	"github.com/spf13/cobra"
)

type PullRequestInfoAuthor struct {
	Login string `json:"login"`
}

type PullRequestInfo struct {
	Number int                   `json:"number"`
	Url    string                `json:"url"`
	Author PullRequestInfoAuthor `json:"author"`
}

func NewCommentCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "comment",
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			logf := logging.GetLogger(cmd.OutOrStdout())
			tp := internal.GetTargetProcess(cmd.Context())

			templateStr, err := cmd.Flags().GetString("template")
			if err != nil {
				return err
			}

			isDryRun, err := cmd.Flags().GetBool("dry-run")
			if err != nil {
				return err
			}

			stdout, _, err := gh.Exec("pr", "view", "--json", "url,author,number")
			if err != nil {
				return err
			}

			var pr PullRequestInfo
			if err := json.Unmarshal(stdout.Bytes(), &pr); err != nil {
				return err
			}

			if strings.TrimSpace(templateStr) == "" {
				templateStr = "PR: {{.Url}}"
			}

			comment, err := T(templateStr, pr)
			if err != nil {
				return err
			}

			var arg0 *string
			if len(args) > 0 {
				arg0 = &args[0]
			}

			id := utils.ExtractTicketID(arg0)
			if id == nil {
				return errors.New("invalid ticket ID or URL")
			}

			assignableID, err := strconv.Atoi(*id)
			if err != nil {
				return err
			}

			if !isDryRun {
				if err := tp.PostComment(cmd.Context(), comment, assignableID); err != nil {
					return err
				}
			}

			logf("Ticket [%d] commented with:\n> %s", assignableID, comment)

			return nil
		},
	}

	cmd.Flags().Bool("dry-run", false, "run without executing changes")
	cmd.Flags().StringP("template", "t", "", "Teplate for the comment available variables: (.Url, .ID, .Number)")

	return cmd
}

func T(tmpl string, data any) (string, error) {
	var b bytes.Buffer
	if err := template.Must(template.New("").Parse(tmpl)).Execute(&b, data); err != nil {
		return "", err
	}

	return b.String(), nil
}
