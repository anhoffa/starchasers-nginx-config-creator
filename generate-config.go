package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/annahoffa/starchasers-nginx-config-creator/utils"
	"os"
	"text/template"
)

func generateNginxConfig(jsonConfig []byte) (string, error) {
	var config Config
	err := json.Unmarshal(jsonConfig, &config)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshall the supplied json: %w", err)
	}

	allowedIp := os.Getenv("HEALTHCHECK_ALLOWED_IP")
	if allowedIp == "" {
		return "", fmt.Errorf("HEALTHCHECK_ALLOWED_IP is not set")
	}
	config.HealthcheckAllowedIp = allowedIp

	tmplBytes, err := os.ReadFile("nginxTemplate.tmpl")
	if err != nil {
		return "", fmt.Errorf("failed to read the template file: %w", err)
	}

	tmpl, err := template.New("nginxTemplate.tmpl").Parse(string(tmplBytes))
	if err != nil {
		return "", fmt.Errorf("failed to parse the template: %w", err)
	}

	var result bytes.Buffer
	err = tmpl.Execute(&result, config)
	if err != nil {
		return "", fmt.Errorf("failed to execute the template: %w", err)
	}

	return result.String(), nil
}

func validateNginxConfig(config string) error {
	f, err := os.CreateTemp("", "nginx-conf-test-")
	if err != nil {
		return err
	}

	if _, err := f.WriteString(config); err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}

	if err := utils.RedirectCmdOutput("nginx", "-t", "-c", f.Name()); err != nil {
		return err
	}

	log.Info("Nginx configuration is valid")

	// todo: what if it'll break halfway the saving process?
	//  also move the path to env after discussing it
	if err := saveInPersistentVolume(f.Name(), "/app/persistent/nginxBackup.conf"); err != nil {
		return err
	}

	return nil
}
