package main

import (
    "log"
    "net/http"
    "url-shortener/handlers"
)

func main() {
    http.HandleFunc("/shorten", handlers.ShortenURL)
    http.HandleFunc("/redirect/",handlers.RedirectURL)
    log.Println("Server started at http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
