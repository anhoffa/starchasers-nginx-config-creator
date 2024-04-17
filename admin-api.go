package main

import (
	"net/http"
	"os"
)

func notifyAdmin() {
	var url = os.Getenv("ADMIN_NOTIFY_URL")
	var token = os.Getenv("ADMIN_NOTIFY_TOKEN")
	log.Info("Notifying the admin panel at " + url)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		log.Errorw("Failed to notify the admin panel",
			"error", err)
		os.Exit(1)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	defer func() {
		if resp != nil {
			resp.Body.Close()
		}
	}()
	if err != nil {
		log.Errorw("Failed to notify the admin panel", "error", err)
		if resp != nil {
			log.Errorw("Admin panel response status: ", "status", resp.Status)
		}
		os.Exit(1)
	}

	if resp != nil && resp.StatusCode != 200 {
		log.Errorw("Failed to notify the admin panel", "status", resp.Status)
		os.Exit(1)
	}
	log.Info("Admin panel has been notified")
}
