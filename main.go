package main

import (
	"log"
)

func main() {
	// todo: move path to env and add request to directus
	//  if 5xx then:
	if err := checkForSavedConfig("/app/persistent/nginxBackup.conf"); err != nil {
		log.Println("No saved Nginx configuration found. Starting the services without it...")
	}

	startServer()
}
