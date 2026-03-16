package versioncheck

import "time"

type VersionCache struct {
	LastCheckTime time.Time `json:"last_check_time"`
}

type GHRelease struct {
	TagName    string `json:"tag_name"`
	Prerelease bool   `json:"prerelease"`
}

var githubAPIURL = "https://api.github.com/repos/rawnly/gh-targetprocess/releases/latest"

const (
	checkInterval = 24 * time.Hour
	httpTimeout   = 2 * time.Second
)
