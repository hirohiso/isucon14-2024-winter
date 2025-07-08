package main

import (
	"context"
	"sync"
	"time"
)

// Simple in-memory cache with TTL
type Cache struct {
	data map[string]cacheItem
	mu   sync.RWMutex
}

type cacheItem struct {
	value     interface{}
	expiresAt time.Time
}

func NewCache() *Cache {
	return &Cache{
		data: make(map[string]cacheItem),
	}
}

func (c *Cache) Set(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = cacheItem{
		value:     value,
		expiresAt: time.Now().Add(ttl),
	}
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, exists := c.data[key]
	if !exists {
		return nil, false
	}

	if time.Now().After(item.expiresAt) {
		// Item expired, remove it
		delete(c.data, key)
		return nil, false
	}

	return item.value, true
}

func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.data, key)
}

func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data = make(map[string]cacheItem)
}

// Global cache instance
var cache *Cache

func init() {
	cache = NewCache()
}

// Cache chair models (static data)
func getCachedChairModels(ctx context.Context) ([]ChairModel, error) {
	const cacheKey = "chair_models"

	if cached, ok := cache.Get(cacheKey); ok {
		return cached.([]ChairModel), nil
	}

	models := []ChairModel{}
	if err := db.SelectContext(ctx, &models, "SELECT * FROM chair_models"); err != nil {
		return nil, err
	}

	// Cache for 1 hour (chair models rarely change)
	cache.Set(cacheKey, models, time.Hour)

	return models, nil
}

// Cache latest ride status
func getCachedLatestRideStatus(ctx context.Context, rideID string) (string, error) {
	cacheKey := "ride_status_" + rideID

	if cached, ok := cache.Get(cacheKey); ok {
		return cached.(string), nil
	}

	status := ""
	if err := db.GetContext(ctx, &status, `SELECT status FROM ride_statuses WHERE ride_id = ? ORDER BY created_at DESC LIMIT 1`, rideID); err != nil {
		return "", err
	}

	// Cache for 5 seconds (statuses change frequently)
	cache.Set(cacheKey, status, 5*time.Second)

	return status, nil
}

// Invalidate ride status cache when status changes
func invalidateRideStatusCache(rideID string) {
	cacheKey := "ride_status_" + rideID
	cache.Delete(cacheKey)
}

// Cache user by access token
func getCachedUserByToken(ctx context.Context, accessToken string) (*User, error) {
	cacheKey := "user_token_" + accessToken

	if cached, ok := cache.Get(cacheKey); ok {
		return cached.(*User), nil
	}

	user := &User{}
	if err := db.GetContext(ctx, user, "SELECT * FROM users WHERE access_token = ?", accessToken); err != nil {
		return nil, err
	}

	// Cache for 10 minutes (user data changes infrequently)
	cache.Set(cacheKey, user, 10*time.Minute)

	return user, nil
}

// Cache chair by access token
func getCachedChairByToken(ctx context.Context, accessToken string) (*Chair, error) {
	cacheKey := "chair_token_" + accessToken

	if cached, ok := cache.Get(cacheKey); ok {
		return cached.(*Chair), nil
	}

	chair := &Chair{}
	if err := db.GetContext(ctx, chair, "SELECT * FROM chairs WHERE access_token = ?", accessToken); err != nil {
		return nil, err
	}

	// Cache for 10 minutes (chair data changes infrequently)
	cache.Set(cacheKey, chair, 10*time.Minute)

	return chair, nil
}

// Cache owner by access token
func getCachedOwnerByToken(ctx context.Context, accessToken string) (*Owner, error) {
	cacheKey := "owner_token_" + accessToken

	if cached, ok := cache.Get(cacheKey); ok {
		return cached.(*Owner), nil
	}

	owner := &Owner{}
	if err := db.GetContext(ctx, owner, "SELECT * FROM owners WHERE access_token = ?", accessToken); err != nil {
		return nil, err
	}

	// Cache for 10 minutes (owner data changes infrequently)
	cache.Set(cacheKey, owner, 10*time.Minute)

	return owner, nil
}
