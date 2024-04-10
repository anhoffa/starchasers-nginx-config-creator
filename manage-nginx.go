package main

import (
	"fmt"
	"github.com/annahoffa/starchasers-nginx-config-creator/utils"
	"os"
)

// reloadNginxConfig reloads the Nginx configuration without stopping the service
func reloadNginxConfig() error {
	if err := utils.RedirectCmdOutput("nginx", "-s", "reload"); err != nil {
		return fmt.Errorf("failed to reload Nginx configuration: %w", err)
	}
	log.Info("Nginx configuration reloaded successfully")
	return nil
}

func modifyNginxConfig(newConfig string) error {
	nginxConfigPath := os.Getenv("NGINX_CONFIG_PATH")

	err := os.WriteFile(nginxConfigPath+"nginx.tmp", []byte(newConfig), 0644)
	if err != nil {
		return fmt.Errorf("error writing the modified config file: %w", err)
	}

	if err := os.Rename(nginxConfigPath+"nginx.tmp", nginxConfigPath+"nginx.conf"); err != nil {
		return err
	}

	return nil
}
