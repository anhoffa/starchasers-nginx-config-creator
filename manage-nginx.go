package main

import (
	"fmt"
	"os"
	"os/exec"
)

// reloadNginxConfiguration reloads the Nginx configuration without stopping the service
func reloadNginxConfiguration() error {
	cmd := exec.Command("nginx", "-s", "reload")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to reload Nginx configuration: %w", err)
	}
	fmt.Println("Nginx configuration reloaded successfully")
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
