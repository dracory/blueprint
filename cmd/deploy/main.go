package main

import (
	"log"
)

// main is the entry point of the application.
func main() {
	config := InitConfig()

	err := BuildApp(config)
	if err != nil {
		log.Fatal(err)
		return
	}

	err = UploadFiles(config)
	if err != nil {
		log.Fatal(err)
		return
	}

	err = UploadExecutable(config)
	if err != nil {
		log.Fatal(err)
		return
	}

	err = ReplaceExecutable(config)
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Printf("✅ Deployed! ")
}
