package main

import (
	"fmt"
	"github.com/anhoffa/starchasers-nginx-config-creator/utils"
	"os"
	"os/exec"
	"path/filepath"
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
	dir := filepath.Dir(nginxConfigFilePath)
	tmpPath := filepath.Join(dir, "nginx.tmp")

	err := os.WriteFile(tmpPath, []byte(newConfig), 0644)
	if err != nil {
		return fmt.Errorf("error writing the modified config file: %w", err)
	}

	if err := os.Rename(tmpPath, filepath.Join(dir, "nginx.conf")); err != nil {
		return err
	}

	return nil
}

// setupNginxConfig generates, validates, modifies, and (re)loads the Nginx configuration.
func setupNginxConfig(loadFunc func() error, jsonConfig []byte) error {
	config, err := generateNginxConfig(jsonConfig)
	if err != nil {
		log.Errorw("Error generating NGINX config", "error", err)
		return fmt.Errorf("error generating NGINX config: %w", err)
	}

	if err := validateNginxConfig(config); err != nil {
		log.Errorw("Error validating NGINX config", "error", err)
		return fmt.Errorf("error validating NGINX config: %w", err)
	}

	log.Infow("Config: ", "config", config)

	if err := modifyNginxConfig(config); err != nil {
		log.Errorw("Error modifying NGINX config", "error", err)
	}

	if err := loadFunc(); err != nil {
		log.Fatalw("Failed to start Nginx: %v", "error", err)
	}

	return nil
}
