package client

import (
	"fmt"
	"io"
	"net/http"

	"trains/internal/cache"
)

// FetchFromNetwork fetches content from network
func FetchFromNetwork(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch URL %s: %w", url, err)
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body from %s: %w", url, err)
	}
	
	return string(body), nil
}

// FetchWithCache fetches content with caching support
func FetchWithCache(url string) (string, error) {
	// Try to load from cache first
	if content, found := cache.LoadFromCache(url); found {
		return content, nil
	}
	
	// Fetch from network
	fmt.Printf("🌐 Fetching from network: %s\n", url)
	content, err := FetchFromNetwork(url)
	if err != nil {
		return "", err
	}
	
	// Save to cache
	if err := cache.SaveToCache(url, content); err != nil {
		fmt.Printf("⚠️  Warning: failed to save to cache: %v\n", err)
		// Don't fail the entire operation if caching fails
	}
	
	return content, nil
}