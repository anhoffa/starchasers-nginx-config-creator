package main

import (
	"fmt"
	"github.com/annahoffa/starchasers-nginx-config-creator/utils"
	"os"
	"os/exec"
)

// startNginx starts the Nginx service in the foreground as the child process.
func startNginx() error {
	if err := exec.Command("nginx", "-g", "daemon off;").Start(); err != nil {
		return err
	}
	log.Info("Nginx started successfully")
	return nil
}

// reloadNginxConfig reloads the Nginx configuration without stopping the service.
func reloadNginxConfig() error {
	if err := utils.RedirectCmdOutput("nginx", "-s", "reload"); err != nil {
		return fmt.Errorf("failed to reload Nginx configuration: %w", err)
	}
	log.Info("Nginx configuration reloaded successfully")
	return nil
}

// modifyNginxConfig modifies the Nginx configuration file by overwriting it with a new configuration.
// It should be used only after the configuration has been validated with validateNginxConfig.
func modifyNginxConfig(newConfig string) error {
	dir := os.Getenv("NGINX_CONFIG_DIR_PATH")

	err := os.WriteFile(dir+"nginx.tmp", []byte(newConfig), 0644)
	if err != nil {
		return fmt.Errorf("error writing the modified config file: %w", err)
	}

	if err := os.Rename(dir+"nginx.tmp", dir+"nginx.conf"); err != nil {
		return err
	}

	return nil
}
