package main

import (
	"os"

	"github.com/dracory/envenc"
)

func main() {
	envenc.NewCli().Run(os.Args)
}
