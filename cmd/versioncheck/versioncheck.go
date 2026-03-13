package versioncheck

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rawnly/gh-targetprocess/internal/utils"
	"golang.org/x/mod/semver"
)

func CheckAndNotify(ctx context.Context, w io.Writer, currentVersion string) {
	if currentVersion == "dev" || currentVersion == "" {
		return
	}

	cache := &VersionCache{}
	if err := ensureGlobalConfigDir(); err != nil {
		return
	}

	if time.Since(cache.LastCheckTime) < checkInterval {
		return
	}

	latestVersion, err := fetchLatestVersion(ctx)
	if err != nil {
		return
	}

	if !utils.IsPiped() {
		if isOutdated(currentVersion, latestVersion) {
			fmt.Fprintf(w, "\nA newer version of Satispay CLI is available: %s (current: %s)\nRun '%s' to update.\n", latestVersion, currentVersion, "satispay update")

			if confirm("Do you want to update now?") {
				if err := utils.AutoUpdate(ctx); err != nil {
					fmt.Fprintf(w, "\nFailed to auto-update: %s", err.Error())
				}
			}
		}
	}
}

func isOutdated(current, latest string) bool {
	if !strings.HasPrefix(current, "v") {
		current = "v" + current
	}

	if !strings.HasPrefix(latest, "v") {
		latest = "v" + latest
	}

	if strings.Contains(semver.Prerelease(current), "dev") {
		return false
	}

	return semver.Compare(current, latest) < 0
}

func ensureGlobalConfigDir() error {
	configDir, err := globalConfigDirPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(configDir, 0o755); err != nil {
		return fmt.Errorf("creating config directory: %w", err)
	}

	return nil
}

func globalConfigDirPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("getting home directory: %w", err)
	}

	return filepath.Join(home, globalConfigDirName), nil
}

func fetchLatestVersion(ctx context.Context) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, httpTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, githubAPIURL, nil)
	if err != nil {
		return "", fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "satispay-cli")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("fetching release info: %w", err)
	}

	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	// limit response to 1MB
	body, err := io.ReadAll(io.LimitReader(res.Body, 1<<20))
	if err != nil {
		return "", fmt.Errorf("reading response: %w", err)
	}

	version, err := parseGithubRelease(body)
	if err != nil {
		return "", fmt.Errorf("parsing response: %w", err)
	}

	return version, nil
}

func parseGithubRelease(body []byte) (string, error) {
	var release GHRelease
	if err := json.Unmarshal(body, &release); err != nil {
		return "", fmt.Errorf("parsin JSON: %w", err)
	}

	if release.Prerelease {
		return "", errors.New("only prerelease available")
	}

	if release.TagName == "" {
		return "", errors.New("empty tag name")
	}

	return release.TagName, nil
}

func confirm(prompt string) bool {
	for {
		fmt.Printf("%s (Y/n): ", prompt)
		var input string
		fmt.Scanln(&input)
		input = strings.TrimSpace(strings.ToLower(input))
		switch input {
		case "n", "no":
			return false
		case "", "y", "yes":
			return true
		default:
			fmt.Println("Please enter Y or N")
		}
	}
}
