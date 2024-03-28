package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func startServer() {
	http.HandleFunc("/generate-config", authMiddleware(handleNginxConfigRequest))
	fmt.Println("Server is running on port 8080...")
	if err := http.ListenAndServe("localhost:8080", nil); err != nil {
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

	jsonConfig, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	config, err := generateNginxConfig(jsonConfig)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing NGINX config: %s", err), http.StatusBadRequest)
		return
	}

	if _, err := validateNginxConfig(config); err != nil {
		http.Error(w, fmt.Sprintf("Error validating NGINX config: %s", err), http.StatusBadRequest)
		return
	}

	fmt.Println("Config: ", config) // todo: delete

	// todo: possibly restart the whole app in these cases
	if err := modifyNginxConfig(config); err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	if err := reloadNginxConfiguration(); err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(config); err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
}
