package storage

import (
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
)

func setupTestRedis() (*RedisStore, *miniredis.Miniredis) {
	
	mr, err := miniredis.Run()
	if err != nil {
		panic(err)
	}

	
	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	return &RedisStore{client: client}, mr
}

func TestSaveAndGet(t *testing.T) {
	store, mr := setupTestRedis()
	defer mr.Close()

	shortURL := "short"
	originalURL := "http://example.com"

	
	err := store.Save(shortURL, originalURL)
	assert.NoError(t, err)


	retrievedURL, err := store.Get(shortURL)
	assert.NoError(t, err)
	assert.Equal(t, originalURL, retrievedURL)
}

func TestGetShortURL(t *testing.T) {
	store, mr := setupTestRedis()
	defer mr.Close()

	shortURL := "short"
	originalURL := "http://example.com"

	
	err := store.Save(shortURL, originalURL)
	assert.NoError(t, err)

	
	retrievedShortURL, err := store.GetShortURL(originalURL)
	assert.NoError(t, err)
	assert.Equal(t, shortURL, retrievedShortURL)
}

func TestGetAllOriginalURLs(t *testing.T) {
	store, mr := setupTestRedis()
	defer mr.Close()

	shortURL1 := "short1"
	originalURL1 := "http://example1.com"
	shortURL2 := "short2"
	originalURL2 := "http://example2.com"

	
	err := store.Save(shortURL1, originalURL1)
	assert.NoError(t, err)
	err = store.Save(shortURL2, originalURL2)
	assert.NoError(t, err)


	allURLs, err := store.GetAllOriginalURLs()
	assert.NoError(t, err)
	assert.Contains(t, allURLs, originalURL1)
	assert.Contains(t, allURLs, originalURL2)
}
