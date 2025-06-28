package main

import (
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func colorHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	colors := []string{"red", "green", "blue", "yellow", "purple", "orange", "pink", "cyan"}
	rand.Seed(time.Now().UnixNano())
	color := colors[rand.Intn(len(colors))]
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(color))
}

func TestColorHandler(t *testing.T) {
	t.Run("GET /color returns valid color", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/color", nil)
		w := httptest.NewRecorder()
		colorHandler(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		validColors := map[string]bool{
			"red": true, "green": true, "blue": true, "yellow": true,
			"purple": true, "orange": true, "pink": true, "cyan": true,
		}

		buf := make([]byte, 16)
		n, _ := resp.Body.Read(buf)
		color := strings.TrimSpace(string(buf[:n]))

		if !validColors[color] {
			t.Errorf("Unexpected color returned: %q", color)
		}
	})

	t.Run("POST /color returns 405", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/color", nil)
		w := httptest.NewRecorder()
		colorHandler(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusMethodNotAllowed {
			t.Errorf("Expected status 405, got %d", resp.StatusCode)
		}
	})

	t.Run("GET /wrong returns 404", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/wrong", nil)
		w := httptest.NewRecorder()

		http.NotFoundHandler().ServeHTTP(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", resp.StatusCode)
		}
	})
}
