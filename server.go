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

	if err := setupNginxConfig(reloadNginxConfig, jsonConfig); err != nil {
		http.Error(w, "Failed to setup Nginx config", http.StatusInternalServerError)
		log.Errorw("Failed to setup Nginx config", "error", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
