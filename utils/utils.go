package utils

import (
    "math/rand"
    "net/url"
    "strings"
    "time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func init() {
    rand.Seed(time.Now().UnixNano())
}

func GenerateShortURL() string {
    b := make([]byte, 6)
    for i := range b {
        b[i] = letterBytes[rand.Intn(len(letterBytes))]
    }
    return string(b)
}

func GetDomain(inputURL string) string {
    u, err := url.Parse(inputURL)
    if err != nil {
        return ""
    }
    domain := strings.Split(u.Hostname(), ".")
    if len(domain) < 2 {
        return u.Hostname()
    }
    return domain[len(domain)-2]
}
