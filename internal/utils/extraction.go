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
