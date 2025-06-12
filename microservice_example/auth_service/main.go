package main

import (
    "encoding/base64"
    "net/http"
    "strings"
)

const (
    validUser = "user"
    validPass = "pass"
)

func main() {
    http.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
        auth := r.Header.Get("Authorization")
        if !strings.HasPrefix(auth, "Basic ") {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        payload, _ := base64.StdEncoding.DecodeString(strings.TrimPrefix(auth, "Basic "))
        pair := strings.SplitN(string(payload), ":", 2)
        if len(pair) != 2 || pair[0] != validUser || pair[1] != validPass {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("OK"))
    })
    http.ListenAndServe(":8080", nil)
}
