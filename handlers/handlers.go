package handlers

import (
    "encoding/json"
    "log"
    "net/http"
    "strings"
    "url-shortener/utils"
)

var urlMap = make(map[string]string)

type URLRequest struct {
    URL string `json:"url"`
}

type URLResponse struct {
    ShortURL string `json:"short_url"`
}

func ShortenURL(w http.ResponseWriter, r *http.Request) {
    log.Println("Received request to shorten URL")
    var req URLRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        log.Printf("Error decoding request body: %v", err)
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    shortURL := utils.GenerateShortURL()
    urlMap[shortURL] = req.URL
    log.Printf("URL shortened: %s -> %s", req.URL, shortURL)

    resp := URLResponse{ShortURL: "http://localhost:8080/redirect/" + shortURL}
    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(resp); err != nil {
        log.Printf("Error encoding response: %v", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func RedirectURL(w http.ResponseWriter, r *http.Request) {
    log.Println("Received request to redirect URL")
    shortURL := strings.TrimPrefix(r.URL.Path, "/redirect/")
    log.Printf("Looking up short code: %s", shortURL)
    originalURL, exists := urlMap[shortURL]

    if !exists {
        log.Printf("URL not found for short code: %s", shortURL)
        http.Error(w, "URL not found", http.StatusNotFound)
        return
    }

    log.Printf("Redirecting short code %s to %s", shortURL, originalURL)
    http.Redirect(w, r, originalURL, http.StatusFound)
}
