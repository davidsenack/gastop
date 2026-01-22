package adapter

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"sync"
	"time"
)

// Adapter executes Gas Town CLI commands and caches results.
type Adapter struct {
	gtPath   string
	bdPath   string
	townRoot string
	timeout  time.Duration

	mu    sync.RWMutex
	cache map[string]*cacheEntry
}

type cacheEntry struct {
	data      interface{}
	fetchedAt time.Time
	stale     bool
}

// New creates a new adapter with the given configuration.
func New(gtPath, bdPath, townRoot string) *Adapter {
	if gtPath == "" {
		gtPath = "gt"
	}
	if bdPath == "" {
		bdPath = "bd"
	}
	return &Adapter{
		gtPath:   gtPath,
		bdPath:   bdPath,
		townRoot: townRoot,
		timeout:  2 * time.Second,
		cache:    make(map[string]*cacheEntry),
	}
}

// SetTimeout sets the command execution timeout.
func (a *Adapter) SetTimeout(d time.Duration) {
	a.timeout = d
}

// execGT runs a gt command and returns the output.
func (a *Adapter) execGT(ctx context.Context, args ...string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, a.timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, a.gtPath, args...)
	if a.townRoot != "" {
		cmd.Dir = a.townRoot
	}

	out, err := cmd.Output()
	if ctx.Err() == context.DeadlineExceeded {
		return nil, fmt.Errorf("command timed out: gt %v", args)
	}
	return out, err
}

// execBD runs a bd command and returns the output.
func (a *Adapter) execBD(ctx context.Context, args ...string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, a.timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, a.bdPath, args...)
	if a.townRoot != "" {
		cmd.Dir = a.townRoot
	}

	out, err := cmd.Output()
	if ctx.Err() == context.DeadlineExceeded {
		return nil, fmt.Errorf("command timed out: bd %v", args)
	}
	return out, err
}

// getCache returns cached data if available and not too old.
func (a *Adapter) getCache(key string, maxAge time.Duration) (interface{}, bool) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	entry, ok := a.cache[key]
	if !ok {
		return nil, false
	}
	if time.Since(entry.fetchedAt) > maxAge {
		entry.stale = true
	}
	return entry.data, true
}

// setCache stores data in the cache.
func (a *Adapter) setCache(key string, data interface{}) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.cache[key] = &cacheEntry{
		data:      data,
		fetchedAt: time.Now(),
		stale:     false,
	}
}

// ClearCache clears all cached data.
func (a *Adapter) ClearCache() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.cache = make(map[string]*cacheEntry)
}

// IsCacheStale returns true if the cached data for the key is stale.
func (a *Adapter) IsCacheStale(key string) bool {
	a.mu.RLock()
	defer a.mu.RUnlock()

	entry, ok := a.cache[key]
	if !ok {
		return false
	}
	return entry.stale
}

// parseJSON unmarshals JSON output into the target.
func parseJSON(data []byte, target interface{}) error {
	if len(data) == 0 {
		return nil
	}
	return json.Unmarshal(data, target)
}
