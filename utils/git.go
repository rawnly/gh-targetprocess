package utils

import (
	"bytes"
	"os/exec"
	"strings"
)

func CurrentBranch() (string, error) {
	cmd := exec.Command("git", "branch", "--show-current")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	branch := strings.TrimSpace(out.String())
	return branch, nil
}
