package helpers

import "os"

func IsCliMode() bool {
	return len(os.Args) > 1
}
