package main

import (
	"io"
	"net/http"
	"os"
	"strings"
)

func startServer() {
	http.HandleFunc("/generate-config", authMiddleware(handleNginxConfigRequest))
	log.Infof("Server is running on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		token = strings.TrimPrefix(token, "Bearer ")
		if token != os.Getenv("API_KEY") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}
}

func handleNginxConfigRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	log.Info("Received request")

	jsonConfig, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	config, err := generateNginxConfig(jsonConfig)
	if err != nil {
		log.Errorw("Error generating NGINX config", "error", err)
		http.Error(w, "Error generating NGINX config", http.StatusBadRequest)
		return
	}

	if _, err := validateNginxConfig(config); err != nil {
		log.Errorw("Error validating NGINX config", "error", err)
		http.Error(w, "Error validating NGINX config", http.StatusBadRequest)
		return
	}

	log.Infow("Config: ", "config", config)

	if err := modifyNginxConfig(config); err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		log.Errorw("Error modifying NGINX config", "error", err)
		os.Exit(1)
		return
	}
	if err := reloadNginxConfig(); err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		log.Errorw("Error reloading NGINX config", "error", err)
		os.Exit(1)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
