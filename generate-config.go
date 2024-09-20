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

	tmplBytes, err := os.ReadFile("nginxTemplate.tmpl")
	if err != nil {
		return "", fmt.Errorf("failed to read the template file: %w", err)
	}

	tmpl, err := template.New("nginxTemplate.tmpl").Parse(string(tmplBytes))
	if err != nil {
		return "", fmt.Errorf("failed to parse the template: %w", err)
	}

	var result bytes.Buffer
	err = tmpl.Execute(&result, prepareTemplateParams(&config))
	if err != nil {
		return "", fmt.Errorf("failed to execute the template: %w", err)
	}

	return result.String(), nil
}

func prepareTemplateParams(config *Config) TemplateParams {
	var containers []Container
	var existingNames = map[string]bool{}

	for _, domain := range config.Domains {
		_, exists := existingNames[domain.ContainerName]
		if !exists {
			existingNames[domain.ContainerName] = true
			containers = append(containers, Container{ContainerName: domain.ContainerName, Ip: domain.Ip})
		}
	}
	return TemplateParams{
		Containers: containers,
		Domains:    config.Domains,
	}
}

func validateNginxConfig(config string) error {
	f, err := os.CreateTemp("", "nginx-conf-test-")
	if err != nil {
		return err
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			log.Warnw("Error removing temporary file", "error", err)
		}
	}(f.Name())

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
	return nil
}
