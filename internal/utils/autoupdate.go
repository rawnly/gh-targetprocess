package utils

import (
	"context"

	"github.com/cli/go-gh/v2"
)

func AutoUpdate(ctx context.Context) error {
	return gh.ExecInteractive(ctx, "extensions", "upgrade", "gh-targetprocess")
}
