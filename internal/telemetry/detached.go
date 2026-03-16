package telemetry

import (
	"encoding/json"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/denisbrodbeck/machineid"
	"github.com/posthog/posthog-go"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var (
	PostHogAPIKey   = "phc_dev_key"
	PostHogEndpoint = "https://eu.i.posthog.com"
)

type EventPayload struct {
	Event      string         `json:"event"`
	DistinctID string         `json:"distinct_id"`
	Properties map[string]any `json:"properties"`
	Timestamp  time.Time      `json:"timestamp"`
}

type silentLogger struct{}

func (silentLogger) Logf(_ string, _ ...interface{})   {}
func (silentLogger) Debugf(_ string, _ ...interface{}) {}
func (silentLogger) Warnf(_ string, _ ...interface{})  {}
func (silentLogger) Errorf(_ string, _ ...interface{}) {}

func BuildEventPayload(cmd *cobra.Command, version string) *EventPayload {
	if cmd == nil {
		return nil
	}

	machineID, err := machineid.ProtectedID("gh-targetprocess")
	if err != nil {
		return nil
	}

	var flags []string
	cmd.Flags().Visit(func(f *pflag.Flag) {
		flags = append(flags, f.Name)
	})

	properties := map[string]any{
		"command":     cmd.CommandPath(),
		"cli_version": version,
		"os":          runtime.GOOS,
		"arch":        runtime.GOARCH,
	}

	if len(flags) > 0 {
		properties["flags"] = strings.Join(flags, ",")
	}

	return &EventPayload{
		Event:      "command_executed",
		DistinctID: machineID,
		Properties: properties,
		Timestamp:  time.Now(),
	}
}

func TrackCommandDetached(cmd *cobra.Command, version string) {
	if os.Getenv("GH_TARGETPROCESS_TELEMETRY_DISABLE") != "" {
		return
	}

	if cmd == nil || cmd.Hidden {
		return
	}

	payload := BuildEventPayload(cmd, version)
	if payload == nil {
		return
	}

	if json, err := json.Marshal(payload); err == nil {
		spawnDetached(string(json))
	}
}

func SendEvent(payloadJson string) {
	var payload EventPayload

	if err := json.Unmarshal([]byte(payloadJson), &payload); err != nil {
		return
	}

	// allow building without analytics
	if PostHogAPIKey == "phc_dev_key" {
		return
	}

	client, err := posthog.NewWithConfig(PostHogAPIKey, posthog.Config{
		Endpoint:     PostHogEndpoint,
		Logger:       silentLogger{},
		DisableGeoIP: posthog.Ptr(true),
	})

	if err != nil {
		return
	}

	defer func() {
		_ = client.Close()
	}()

	props := posthog.NewProperties()
	for k, v := range payload.Properties {
		props.Set(k, v)
	}

	// best effort, fire-n-forget
	_ = client.Enqueue(posthog.Capture{
		DistinctId: payload.DistinctID,
		Event:      payload.Event,
		Properties: props,
		Timestamp:  payload.Timestamp,
	})
}
