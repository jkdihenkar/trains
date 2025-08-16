package cache

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"trains/internal/types"
)

const (
	cacheDir    = "./cache"
	cacheExpiry = 24 * time.Hour // Cache expires after 24 hours
)

// initCache creates cache directory if it doesn't exist
func InitCache() error {
	return os.MkdirAll(cacheDir, 0755)
}

// getCacheFilePath generates MD5-hashed cache file path for URL
func getCacheFilePath(url string) string {
	hash := md5.Sum([]byte(url))
	filename := fmt.Sprintf("%x.json", hash)
	return filepath.Join(cacheDir, filename)
}

// LoadFromCache loads content from cache if available and not expired
func LoadFromCache(url string) (string, bool) {
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
	var entry types.CacheEntry
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

// SaveToCache saves content to cache
func SaveToCache(url, content string) error {
	filePath := getCacheFilePath(url)
	
	entry := types.CacheEntry{
		URL:       url,
		Content:   content,
		Timestamp: time.Now(),
	}
	
	data, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal cache entry for %s: %w", url, err)
	}
	
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write cache file %s: %w", filePath, err)
	}
	
	fmt.Printf("💾 Cached response for %s\n", url)
	return nil
}