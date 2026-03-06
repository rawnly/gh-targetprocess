package logging

import (
	"fmt"
	"io"

	"github.com/rawnly/gh-targetprocess/internal/utils"
)

func GetLogger(w io.Writer) func(data ...any) {
	return func(data ...any) {
		if utils.IsPiped() {
			return
		}

		msg := data[0].(string)
		data = data[1:]

		fmt.Fprintf(w, msg+"\n", data...)
	}
}
