package handlers

import (
    "encoding/json"
    "errors"
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"
)


var ErrURLNotFound = errors.New("url not found")

type mockStore struct {
    data map[string]string
}

func (m *mockStore) Save(shortURL, originalURL string) error {
    m.data[shortURL] = originalURL
    return nil
}

func (m *mockStore) Get(shortURL string) (string, error) {
    if url, ok := m.data[shortURL]; ok {
        return url, nil
    }
    return "", ErrURLNotFound
}

func (m *mockStore) GetShortURL(originalURL string) (string, error) {
    for shortURL, url := range m.data {
        if url == originalURL {
            return shortURL, nil
        }
    }
    return "", ErrURLNotFound
}

func (m *mockStore) GetAllOriginalURLs() ([]string, error) {
    urls := make([]string, 0, len(m.data))
    for _, url := range m.data {
        urls = append(urls, url)
    }
    return urls, nil
}


func (m *mockStore) FlushDB() error {
    m.data = make(map[string]string)
    return nil
}

func TestShortenURL(t *testing.T) {
    store := &mockStore{data: make(map[string]string)}
    SetStore(store)

    req := httptest.NewRequest("POST", "/shorten", strings.NewReader(`{"url":"https://www.example.com"}`))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()
    ShortenURL(w, req)

    resp := w.Result()
    if resp.StatusCode != http.StatusOK {
        t.Errorf("Expected status OK, got %v", resp.Status)
    }

    var response URLResponse
    if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
        t.Fatalf("Failed to decode response: %v", err)
    }

    if !strings.HasPrefix(response.ShortURL, "http://localhost:8080/redirect/") {
        t.Errorf("Unexpected short URL: %s", response.ShortURL)
    }
}

func TestRedirectURL(t *testing.T) {
    store := &mockStore{data: make(map[string]string)}
    SetStore(store)

    shortURL := "test123"
    originalURL := "https://www.example.com"
    store.Save(shortURL, originalURL)

    req := httptest.NewRequest("GET", "/redirect/"+shortURL, nil)
    w := httptest.NewRecorder()
    RedirectURL(w, req)

    resp := w.Result()
    if resp.StatusCode != http.StatusFound {
        t.Errorf("Expected status Found, got %v", resp.Status)
    }

    location, err := resp.Location()
    if err != nil {
        t.Fatal(err)
    }

    if location.String() != originalURL {
        t.Errorf("Expected redirect to %s, got %s", originalURL, location)
    }
}

func TestGetTopDomains(t *testing.T) {
    store := &mockStore{data: make(map[string]string)}
    store.Save("1", "https://www.example.com")
    store.Save("2", "https://www.example.com")
    store.Save("3", "https://sub.domain.com")
    SetStore(store)

    req := httptest.NewRequest("GET", "/metrics", nil)
    w := httptest.NewRecorder()
    GetTopDomains(w, req)

    resp := w.Result()
    if resp.StatusCode != http.StatusOK {
        t.Errorf("Expected status OK, got %v", resp.Status)
    }

    var domainCounts []DomainCount
    if err := json.NewDecoder(resp.Body).Decode(&domainCounts); err != nil {
        t.Fatalf("Failed to decode response: %v", err)
    }

    if len(domainCounts) == 0 {
        t.Fatalf("Expected non-empty domain counts")
    }

    for _, dc := range domainCounts {
        t.Logf("Domain: %s, Count: %d", dc.Domain, dc.Count)
    }
}
