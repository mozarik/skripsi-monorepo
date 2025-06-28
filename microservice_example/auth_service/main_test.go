package main

import (
	"encoding/base64"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func init() {
	// Disable log timestamps
	log.SetFlags(0)
}

func TestAuthHandler(t *testing.T) {
	tests := []struct {
		name           string
		authHeader     string
		expectedStatus int
		expectedBody   string
	}{
		{"Valid Credentials", "Basic " + base64.StdEncoding.EncodeToString([]byte("user:pass")), http.StatusOK, "OK"},
		{"Invalid Credentials", "Basic " + base64.StdEncoding.EncodeToString([]byte("wrong:pass")), http.StatusUnauthorized, "Unauthorized\n"},
		{"Missing Header", "", http.StatusUnauthorized, "Unauthorized\n"},
		{"Malformed Header", "Bearer token", http.StatusUnauthorized, "Unauthorized\n"},
		{"Invalid Base64", "Basic invalid_base64", http.StatusUnauthorized, "Unauthorized\n"},
		{"Missing Password", "Basic " + base64.StdEncoding.EncodeToString([]byte("user")), http.StatusUnauthorized, "Unauthorized\n"},
		{"Empty Credentials", "Basic ", http.StatusUnauthorized, "Unauthorized\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log.Printf("[START] Test case: %s", tt.name)

			req, err := http.NewRequest("GET", "/auth", nil)
			if err != nil {
				t.Fatal(err)
			}

			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(authHandler)
			handler.ServeHTTP(rr, req)

			status := rr.Code
			body := rr.Body.String()

			if status != tt.expectedStatus {
				t.Errorf("[FAIL] %s: expected status %d, got %d", tt.name, tt.expectedStatus, status)
			} else {
				log.Printf("[PASS] %s: status %d as expected", tt.name, status)
			}

			if body != tt.expectedBody {
				t.Errorf("[FAIL] %s: expected body %q, got %q", tt.name, tt.expectedBody, body)
			} else {
				log.Printf("[PASS] %s: body %q as expected", tt.name, body)
			}

			log.Printf("[END] Test case: %s\n", tt.name)
		})
	}
}

// authHandler is extracted for testing
func authHandler(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	if !strings.HasPrefix(auth, "Basic ") {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	payload, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(auth, "Basic "))
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	pair := strings.SplitN(string(payload), ":", 2)
	if len(pair) != 2 || pair[0] != validUser || pair[1] != validPass {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
