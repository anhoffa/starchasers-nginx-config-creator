package main

import (
	logging "github.com/ipfs/go-log/v2"
	"os"
)

var log = logging.Logger("config-creator")
var nginxConfigFilePath = os.Getenv("NGINX_CONFIG_FILE_PATH")

func init() {
	if nginxConfigFilePath == "" {
		log.Warnw("NGINX_CONFIG_FILE_PATH is not set. Using the default value...", "default", "/etc/nginx/nginx.conf")
		nginxConfigFilePath = "/etc/nginx/nginx.conf"
	}
}

func main() {
	logging.SetAllLoggers(logging.LevelInfo)
	log.Info("Starting the services...")
	// todo: add request to directus (if 5xx then:...)

	if err := checkForSavedConfig(nginxConfigFilePath); err != nil {
		log.Warnw("No saved Nginx configuration found. Starting the services without it...", "error", err)
		if err := setupNginxConfig(func() error {
			return nil
		}, []byte("{}")); err != nil {
			log.Fatalw("Failed to setup Nginx: %v", "error", err)
		}
	}

	err := startNginx()
	if err != nil {
		if err := setupNginxConfig(startNginx, []byte("{}")); err != nil {
			log.Fatalw("Failed to setup Nginx: %v", "error", err)
		}
		log.Fatalw("Failed to start Nginx: %v", "error", err)
	}

	startServer()
}
