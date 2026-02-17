package utils

import (
	"os/user"
	"path"
)

func ExpandPath(p string) string {
	if p[:2] == "~/" {
		usr, err := user.Current()
		if err != nil {
			panic(err.Error())
		}

		return path.Join(usr.HomeDir, p[1:])
	}

	return p
}
