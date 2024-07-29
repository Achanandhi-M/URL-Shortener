package handlers

import (
    "encoding/json"
    "log"
    "net/http"
    "sort"
    "strings"
    "url-shortener/storage"
    "url-shortener/utils"
)

var store storage.Store

func SetStore(s storage.Store) {
    store = s
}

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

    // Check if the URL already exists
    shortURL, err := store.GetShortURL(req.URL)
    if err == nil {
        // URL already exists, return existing short URL
        resp := URLResponse{ShortURL: "http://localhost:8080/redirect/" + shortURL}
        w.Header().Set("Content-Type", "application/json")
        if err := json.NewEncoder(w).Encode(resp); err != nil {
            log.Printf("Error encoding response: %v", err)
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
        return
    }

    // Generate new short URL
    shortURL = utils.GenerateShortURL()
    if err := store.Save(shortURL, req.URL); err != nil {
        log.Printf("Error saving to storage: %v", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
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
    originalURL, err := store.Get(shortURL)
    if err != nil {
        log.Printf("URL not found for short code: %s", shortURL)
        http.Error(w, "URL not found", http.StatusNotFound)
        return
    }

    log.Printf("Redirecting short code %s to %s", shortURL, originalURL)
    http.Redirect(w, r, originalURL, http.StatusFound)
}

type DomainCount struct {
    Domain string `json:"domain"`
    Count  int    `json:"count"`
}

func GetTopDomains(w http.ResponseWriter, r *http.Request) {
    log.Println("Received request for top domains")

    domainCounts := make(map[string]int)
    keys, err := store.GetAllOriginalURLs()
    if err != nil {
        log.Printf("Error retrieving all original URLs: %v", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    log.Printf("Original URLs: %v", keys)
    for _, originalURL := range keys {
        domain := utils.GetDomain(originalURL)
        if domain != "" {
            domainCounts[domain]++
        }
    }

    log.Printf("Domain counts: %v", domainCounts)

    var sortedDomains []DomainCount
    for domain, count := range domainCounts {
        sortedDomains = append(sortedDomains, DomainCount{Domain: domain, Count: count})
    }

    sort.Slice(sortedDomains, func(i, j int) bool {
        return sortedDomains[i].Count > sortedDomains[j].Count
    })

    if len(sortedDomains) > 3 {
        sortedDomains = sortedDomains[:3]
    }

    log.Printf("Top domains: %v", sortedDomains)

    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(sortedDomains); err != nil {
        log.Printf("Error encoding response: %v", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}
