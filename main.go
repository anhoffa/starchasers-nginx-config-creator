package main

import (
	logging "github.com/ipfs/go-log/v2"
	"os"
)

var log = logging.Logger("config-creator")
var nginxConfigFilePath = os.Getenv("NGINX_CONFIG_FILE_PATH")
var nginxConfigDirPath = os.Getenv("NGINX_CONFIG_DIR_PATH")

func init() {
	if nginxConfigFilePath == "" {
		log.Warnw("NGINX_CONFIG_FILE_PATH is not set. Using the default value...", "default", "/etc/nginx/nginx.conf")
		nginxConfigFilePath = "/etc/nginx/nginx.conf"
	}

	if nginxConfigDirPath == "" {
		nginxConfigDirPath = "/etc/nginx/"
	}
}

func main() {
	logging.SetAllLoggers(logging.LevelInfo)
	log.Info("Starting the services...")
	// todo: add request to directus (if 5xx then:...)
	// todo 2: start nginx by app instead of script.sh

	if err := startNginx(); err != nil {
		log.Fatalw("Failed to start Nginx: %v", "error", err)
		os.Exit(1)
	}

	if err := checkForSavedConfig(nginxConfigFilePath); err != nil {
		log.Warnw("No saved Nginx configuration found. Starting the services without it...", "error", err)
	}

	startServer()
}
