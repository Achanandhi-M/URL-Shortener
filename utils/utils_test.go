package utils

import (
    "testing"
    "strings"
)


func TestGenerateShortURL(t *testing.T) {
    shortURL := GenerateShortURL()

    if len(shortURL) != 6 {
        t.Errorf("Expected length 6, but got %d", len(shortURL))
    }

    for _, c := range shortURL {
        if !strings.ContainsRune(letterBytes, c) {
            t.Errorf("Generated short URL contains invalid characters: %s", shortURL)
        }
    }
}


func TestGetDomain(t *testing.T) {
    cases := []struct {
        input    string
        expected string
    }{
        {"https://www.example.com/path", "example"},
        {"https://example.com", "example"},
        {"http://example.com:8080", "example"},
        {"", ""},
    }

    for _, c := range cases {
        result := GetDomain(c.input)
        if result != c.expected {
            t.Errorf("GetDomain(%q) == %q, expected %q", c.input, result, c.expected)
        }
    }
}
