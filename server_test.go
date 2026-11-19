package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAuthMiddleware(t *testing.T) {
	t.Setenv("API_KEY", "secret-key")

	cases := []struct {
		name       string
		header     string
		wantStatus int
		wantCalled bool
	}{
		{"no header", "", http.StatusUnauthorized, false},
		{"wrong key", "Bearer wrong", http.StatusUnauthorized, false},
		{"bearer key", "Bearer secret-key", http.StatusOK, true},
		{"bare key", "secret-key", http.StatusOK, true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			called := false
			handler := authMiddleware(func(w http.ResponseWriter, r *http.Request) {
				called = true
				w.WriteHeader(http.StatusOK)
			})

			req := httptest.NewRequest(http.MethodPost, "/generate-config", strings.NewReader("{}"))
			if tc.header != "" {
				req.Header.Set("Authorization", tc.header)
			}
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)

			if rec.Code != tc.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tc.wantStatus)
			}
			if called != tc.wantCalled {
				t.Errorf("next handler called = %v, want %v", called, tc.wantCalled)
			}
		})
	}
}

func TestHandleNginxConfigRequestRejectsNonPost(t *testing.T) {
	for _, method := range []string{http.MethodGet, http.MethodPut, http.MethodDelete} {
		req := httptest.NewRequest(method, "/generate-config", nil)
		rec := httptest.NewRecorder()
		handleNginxConfigRequest(rec, req)

		if rec.Code != http.StatusMethodNotAllowed {
			t.Errorf("%s: status = %d, want %d", method, rec.Code, http.StatusMethodNotAllowed)
		}
	}
}
