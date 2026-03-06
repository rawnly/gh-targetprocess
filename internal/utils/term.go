package utils

import "os"

func IsPiped() bool {
	fi, err := os.Stdout.Stat()
	if err != nil {
		return false
	}

	return (fi.Mode() & os.ModeCharDevice) == 0
}
