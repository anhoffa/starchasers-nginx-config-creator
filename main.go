package main

import (
	logging "github.com/ipfs/go-log/v2"
)

var log = logging.Logger("config-creator")

func main() {
	log.Info("Starting the services...")
	// todo: move path to env and add request to directus
	//  if 5xx then:
	if err := checkForSavedConfig("/app/persistent/nginxBackup.conf"); err != nil {
		log.Warnw("No saved Nginx configuration found. Starting the services without it...", "error", err)
	}

	startServer()
}
