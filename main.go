package main

import (
    "log"
    "net/http"
    "url-shortener/handlers"
    "url-shortener/storage"
)

func main() {
    // Initialize Redis store
    store := storage.NewRedisStore("localhost:6379")
    handlers.SetStore(store)

    http.HandleFunc("/shorten", handlers.ShortenURL)
    http.HandleFunc("/redirect/", handlers.RedirectURL)
    http.HandleFunc("/topdomains", handlers.GetTopDomains)

    log.Println("Server started at http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
