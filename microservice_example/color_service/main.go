package main

import (
    "fmt"
    "math/rand"
    "net/http"
    "time"
)

func main() {
    http.HandleFunc("/color", func(w http.ResponseWriter, r *http.Request) {
        colors := []string{"red", "green", "blue", "yellow", "purple", "orange", "pink", "cyan"}
        rand.Seed(time.Now().UnixNano())
        color := colors[rand.Intn(len(colors))]
        fmt.Fprint(w, color)
    })
    http.ListenAndServe(":8080", nil)
}
