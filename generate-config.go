package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"text/template"
)

func generateNginxConfig(jsonConfig []byte) (string, error) {
	var config Config
	err := json.Unmarshal(jsonConfig, &config)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshall the supplied json: %w", err)
	}

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

func validateNginxConfig(config string) (bool, error) {
	f, err := os.CreateTemp("", "nginx-conf-test-")
	if err != nil {
		return false, err
	}

	if _, err := f.WriteString(config); err != nil {
		return false, err
	}
	if err := f.Close(); err != nil {
		return false, err
	}

	cmd := exec.Command("nginx", "-t", "-c", f.Name())

	if err := cmd.Run(); err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			return false, nil
		}
		return false, err
	}

	log.Info("Nginx configuration is valid")

	// todo: what if it'll break halfway the saving process?
	//  also move the path to env after discussing it
	if err := saveInPersistentVolume(f.Name(), "/app/persistent/nginxBackup.conf"); err != nil {
		return false, err
	}

	return true, nil
}
