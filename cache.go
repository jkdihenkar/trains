package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const (
	cacheDir    = "./cache"
	cacheExpiry = 24 * time.Hour // Cache expires after 24 hours
)

// initCache creates cache directory if it doesn't exist
func initCache() error {
	return os.MkdirAll(cacheDir, 0755)
}

// getCacheFilePath generates MD5-hashed cache file path for URL
func getCacheFilePath(url string) string {
	hash := md5.Sum([]byte(url))
	filename := fmt.Sprintf("%x.json", hash)
	return filepath.Join(cacheDir, filename)
}

// loadFromCache loads content from cache if available and not expired
func loadFromCache(url string) (string, bool) {
	filePath := getCacheFilePath(url)
	
	// Check if cache file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return "", false
	}
	
	// Read cache file
	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading cache file: %v\n", err)
		return "", false
	}
	
	// Parse cache entry
	var entry CacheEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		fmt.Printf("Error parsing cache entry: %v\n", err)
		return "", false
	}
	
	// Check if cache is expired
	if time.Since(entry.Timestamp) > cacheExpiry {
		fmt.Printf("Cache expired for %s\n", url)
		return "", false
	}
	
	fmt.Printf("💾 Cache hit for %s (cached %v ago)\n", url, time.Since(entry.Timestamp).Round(time.Minute))
	return entry.Content, true
}

// saveToCache saves content to cache
func saveToCache(url, content string) error {
	filePath := getCacheFilePath(url)
	
	entry := CacheEntry{
		URL:       url,
		Content:   content,
		Timestamp: time.Now(),
	}
	
	data, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling cache entry: %v", err)
	}
	
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("error writing cache file: %v", err)
	}
	
	fmt.Printf("💾 Cached response for %s\n", url)
	return nil
}

// fetchFromNetwork fetches content from network
func fetchFromNetwork(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("error fetching URL: %v", err)
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}
	
	return string(body), nil
}

// fetchWithCache fetches content with caching support
func fetchWithCache(url string) (string, error) {
	// Try to load from cache first
	if content, found := loadFromCache(url); found {
		return content, nil
	}
	
	// Fetch from network
	fmt.Printf("🌐 Fetching from network: %s\n", url)
	content, err := fetchFromNetwork(url)
	if err != nil {
		return "", err
	}
	
	// Save to cache
	if err := saveToCache(url, content); err != nil {
		fmt.Printf("⚠️  Warning: failed to save to cache: %v\n", err)
		// Don't fail the entire operation if caching fails
	}
	
	return content, nil
}