package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

// Helper to check if a number is prime
func isPrime(n int) bool {
	if n <= 1 {
		return false
	}
	if n <= 3 {
		return true
	}
	if n%2 == 0 || n%3 == 0 {
		return false
	}
	for i := 5; i*i <= n; i += 6 {
		if n%i == 0 || n%(i+2) == 0 {
			return false
		}
	}
	return true
}

// Generate 3 unique random primes with the given digit count
func generatePrimes(digit int) ([]int, error) {
	if digit < 1 || digit > 6 {
		return nil, fmt.Errorf("digit must be between 1 and 6")
	}
	min := 1
	for i := 1; i < digit; i++ {
		min *= 10
	}
	max := min*10 - 1

	primes := []int{}
	rand.Seed(time.Now().UnixNano())
	attempts := 0
	for len(primes) < 5 && attempts < 10000 {
		num := rand.Intn(max-min+1) + min
		if isPrime(num) {
			already := false
			for _, p := range primes {
				if p == num {
					already = true
					break
				}
			}
			if !already {
				primes = append(primes, num)
			}
		}
		attempts++
	}
	if len(primes) < 5 {
		return nil, fmt.Errorf("could not find 5 prime numbers with %d digits", digit)
	}
	return primes, nil
}

type requestPayload struct {
	Digit int `json:"digit"`
}

type responsePayload struct {
	Primes []int  `json:"primes,omitempty"`
	Error  string `json:"error,omitempty"`
}

func generatePrimeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(responsePayload{Error: "method not allowed"})
		return
	}
	var req requestPayload
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(responsePayload{Error: "invalid request body"})
		return
	}
	primes, err := generatePrimes(req.Digit)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(responsePayload{Error: err.Error()})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responsePayload{Primes: primes})
}

func main() {
	http.HandleFunc("/generate_prime", generatePrimeHandler)
	fmt.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
