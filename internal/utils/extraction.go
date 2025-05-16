package utils

import (
	"regexp"
)

func GetTicketIDFromBranch(branch string) *string {
	re := regexp.MustCompile(`\w+/(\d+)_.*`)
	matches := re.FindStringSubmatch(branch)
	if len(matches) > 1 {
		return &matches[1]
	}
	return nil
}

func ExtractIDFromURL(url string) *string {
	re := regexp.MustCompile(`https?://\w+\.tpondemand\.com/entity/(\d+)([\w+-]+)`)
	matches := re.FindStringSubmatch(url)
	if len(matches) > 1 {
		return &matches[1]
	}
	return nil
}

// ExtractTicketID extracts the ticket ID from the current branch or a given URL.
func ExtractTicketID(idOrURL *string) *string {
	// we ignore the error the directory may not be a git repo
	branch, _ := CurrentBranch()

	var id *string
	if branch != "" {
		id = GetTicketIDFromBranch(branch)
	}

	if id == nil && idOrURL != nil {
		id = ExtractIDFromURL(*idOrURL)

		if id == nil {
			id = idOrURL
		}
	}

	return id
}
