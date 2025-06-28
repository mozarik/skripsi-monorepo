package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGeneratePrimes(t *testing.T) {
	tests := []struct {
		digit       int
		expectsErr  bool
		description string
	}{
		{1, false, "Valid digit 1"},
		{2, false, "Valid digit 2"},
		{6, false, "Valid digit 6"},
		{0, true, "Invalid digit 0"},
		{7, true, "Invalid digit 7"},
		{-1, true, "Negative digit"},
	}

	for _, tt := range tests {
		reqBody, _ := json.Marshal(requestPayload{Digit: tt.digit})
		req := httptest.NewRequest(http.MethodPost, "/generate_prime", bytes.NewReader(reqBody))
		w := httptest.NewRecorder()

		generatePrimeHandler(w, req)

		resp := w.Result()
		var respPayload responsePayload
		json.NewDecoder(resp.Body).Decode(&respPayload)

		if tt.expectsErr {
			if resp.StatusCode == http.StatusOK {
				t.Errorf("%s: expected error but got success", tt.description)
			}
			if respPayload.Error == "" {
				t.Errorf("%s: expected error message but got none", tt.description)
			}
		} else {
			if resp.StatusCode != http.StatusOK {
				t.Errorf("%s: expected success but got status %d", tt.description, resp.StatusCode)
			} else {
				log.Printf("%s: success with primes %v", tt.description, respPayload.Primes)
			}
			if len(respPayload.Primes) != 3 {
				t.Errorf("%s: expected 3 primes but got %d", tt.description, len(respPayload.Primes))
			}
		}
	}
}

func TestGeneratePrimeHandler_MethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/generate_prime", nil)
	w := httptest.NewRecorder()

	generatePrimeHandler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405 Method Not Allowed but got %d", resp.StatusCode)
	} else {
		log.Printf("MethodNotAllowed test: success with status %d", resp.StatusCode)
	}
}

func TestGeneratePrimeHandler_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/generate_prime", bytes.NewReader([]byte("invalid json")))
	w := httptest.NewRecorder()

	generatePrimeHandler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 400 Bad Request but got %d", resp.StatusCode)
	} else {
		log.Printf("InvalidJSON test: success with status %d", resp.StatusCode)
	}
}
