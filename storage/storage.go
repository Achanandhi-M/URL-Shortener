package storage

import (
    "log"

    "github.com/go-redis/redis"
)

type Store interface {
    Save(shortURL, originalURL string) error
    Get(shortURL string) (string, error)
    GetShortURL(originalURL string) (string, error)
    GetAllOriginalURLs() ([]string, error)
    FlushDB() error
}

type RedisStore struct {
    client *redis.Client
}

func NewRedisStore(addr string) *RedisStore {
    client := redis.NewClient(&redis.Options{
        Addr: addr,
    })

    return &RedisStore{client: client}
}

func (r *RedisStore) Save(shortURL, originalURL string) error {
    log.Printf("Saving short URL %s with original URL %s", shortURL, originalURL)
    if err := r.client.Set("url:"+shortURL, originalURL, 0).Err(); err != nil {
        log.Printf("Error saving short URL %s: %v", shortURL, err)
        return err
    }
    if err := r.client.Set("original:"+originalURL, shortURL, 0).Err(); err != nil {
        log.Printf("Error saving original URL %s: %v", originalURL, err)
        return err
    }
    return nil
}

func (r *RedisStore) Get(shortURL string) (string, error) {
    log.Printf("Getting original URL for short URL %s", shortURL)
    originalURL, err := r.client.Get("url:" + shortURL).Result()
    if err != nil {
        log.Printf("Error getting original URL for short URL %s: %v", shortURL, err)
        return "", err
    }
    log.Printf("Found original URL %s for short URL %s", originalURL, shortURL)
    return originalURL, nil
}

func (r *RedisStore) GetShortURL(originalURL string) (string, error) {
    log.Printf("Getting short URL for original URL %s", originalURL)
    shortURL, err := r.client.Get("original:" + originalURL).Result()
    if err != nil {
        log.Printf("Error getting short URL for original URL %s: %v", originalURL, err)
        return "", err
    }
    log.Printf("Found short URL %s for original URL %s", shortURL, originalURL)
    return shortURL, nil
}

func (r *RedisStore) GetAllOriginalURLs() ([]string, error) {
    log.Println("Getting all original URLs")
    keys, err := r.client.Keys("original:*").Result()
    if err != nil {
        log.Printf("Error getting original URL keys: %v", err)
        return nil, err
    }

    var urls []string
    for _, key := range keys {
        shortURL, err := r.client.Get(key).Result()
        if err != nil {
            log.Printf("Error getting short URL for key %s: %v", key, err)
            continue
        }
        originalURL, err := r.client.Get("url:" + shortURL).Result()
        if err != nil {
            log.Printf("Error getting original URL for short URL %s: %v", shortURL, err)
            continue
        }
        urls = append(urls, originalURL)
    }
    return urls, nil
}


func (r *RedisStore) FlushDB() error {
    log.Println("Flushing the database")
    return r.client.FlushDB().Err()
}
