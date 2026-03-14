package telemetry

import (
	"context"
	"io"
	"os"
	"os/exec"
	"syscall"
)

func spawnDetached(payload string) {
	executable, err := os.Executable()
	if err != nil {
		return
	}

	cmd := exec.CommandContext(context.Background(), executable, "__send_analytics_event", payload)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
	cmd.Dir = "/"
	cmd.Env = os.Environ()
	cmd.Stdout = io.Discard
	cmd.Stderr = io.Discard

	if err := cmd.Start(); err != nil {
		return
	}

	_ = cmd.Process.Release()
}
