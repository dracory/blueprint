package main

import "os"

func isCliMode() bool {
	return len(os.Args) > 1
}
