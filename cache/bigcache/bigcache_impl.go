package bigcache

import (
	"time"

	"github.com/allegro/bigcache/v2"
	"github.com/daheige/bigcache-demo/cache"
)

var _ cache.Cache = (*cacheImpl)(nil)

// New create a cache entry
func New(lifeWindow time.Duration, opts ...CacheOption) (cache.Cache, error) {
	config := bigcache.DefaultConfig(lifeWindow)
	for _, o := range opts {
		o(&config)
	}

	cacheEntry, err := bigcache.NewBigCache(config)
	if err != nil {
		return nil, err
	}

	c := &cacheImpl{cacheEntry}

	return c, nil
}

type cacheImpl struct {
	*bigcache.BigCache
}

// Stats returns cache's statistics
func (c *cacheImpl) Stats() cache.Stats {
	status := c.BigCache.Stats()
	s := cache.Stats{
		Hits:       status.Hits,
		Misses:     status.Misses,
		DelHits:    status.DelHits,
		DelMisses:  status.DelMisses,
		Collisions: status.Collisions,
	}

	return s
}

// KeyMetadata returns number of times a cached resource was requested.
func (c *cacheImpl) KeyMetadata(key string) cache.Metadata {
	cacheMeta := c.BigCache.KeyMetadata(key)
	m := cache.Metadata{
		RequestCount: cacheMeta.RequestCount,
	}

	return m
}

// ErrEntryNotFound is an error type struct which is returned when entry was not found for provided key
var ErrEntryNotFound = bigcache.ErrEntryNotFound

// GetWithInfo reads entry for the key with Response info.
// It returns an ErrEntryNotFound when
// no entry exists for the given key.
func (c *cacheImpl) GetWithInfo(key string) ([]byte, cache.Response, error) {
	b, res, err := c.BigCache.GetWithInfo(key)
	response := cache.Response{
		EntryStatus: cache.RemoveReason(res.EntryStatus),
	}

	return b, response, err
}

// CacheOption bigcache config option
type CacheOption func(c *bigcache.Config)

// WithShards 设置shards
func WithShards(num int) CacheOption {
	return func(c *bigcache.Config) {
		c.Shards = num
	}
}

// Time after which entry can be evicted
// LifeWindow time.Duration
func WithLifeWindow(t time.Duration) CacheOption {
	return func(c *bigcache.Config) {
		c.LifeWindow = t
	}
}

// Interval between removing expired entries (clean up).
// If set to <= 0 then no action is performed.
// Setting to < 1 second is counterproductive — bigcache has a one second resolution.
// CleanWindow time.Duration
func WithCleanWindow(t time.Duration) CacheOption {
	return func(c *bigcache.Config) {
		c.CleanWindow = t
	}
}

// WithMaxEntriesInWindow Max number of entries in life window.
// Used only to calculate initial size for cache shards.
// When proper value is set then additional memory allocation does not occur.
func WithMaxEntriesInWindow(n int) CacheOption {
	return func(c *bigcache.Config) {
		c.MaxEntriesInWindow = n
	}
}

// WithMaxEntrySize Max size of entry in bytes. Used only to calculate initial size for cache shards.
func WithMaxEntrySize(size int) CacheOption {
	return func(c *bigcache.Config) {
		c.MaxEntrySize = size
	}
}

// WithStatsEnabled StatsEnabled if true calculate the number of times a cached resource was requested.
func WithStatsEnabled() CacheOption {
	return func(c *bigcache.Config) {
		c.StatsEnabled = true
	}
}

// WithVerbose Verbose mode prints information about new memory allocation
func WithVerbose(v bool) CacheOption {
	return func(c *bigcache.Config) {
		c.Verbose = v
	}
}

// WithHasher Hasher used to map between string keys and unsigned 64bit integers, by default fnv64 hashing is used.
func WithHasher(h bigcache.Hasher) CacheOption {
	return func(c *bigcache.Config) {
		c.Hasher = h
	}
}

// WithHardMaxCacheSize HardMaxCacheSize is a limit for cache size in MB.
// Cache will not allocate more memory than this limit.
// It can protect application from consuming all available memory on machine, therefore from running OOM Killer.
// Default value is 0 which means unlimited size. When the limit is higher than 0 and reached then
// the oldest entries are overridden for the new ones.
func WithHardMaxCacheSize(hardMaxCacheSize int) CacheOption {
	return func(c *bigcache.Config) {
		c.HardMaxCacheSize = hardMaxCacheSize
	}
}

// WithOnRemove OnRemove is a callback fired when the oldest entry is
// removed because of its expiration time or no space left
// for the new entry, or because delete was called.
// Default value is nil which means no callback and it prevents from unwrapping the oldest entry.
// ignored if OnRemoveWithMetadata is specified.
func WithOnRemove(fn func(key string, entry []byte)) CacheOption {
	return func(c *bigcache.Config) {
		c.OnRemove = fn
	}
}

// WithOnRemoveWithMetadata OnRemoveWithMetadata is a callback fired
// when the oldest entry is removed because of its expiration time or no space left
// for the new entry, or because delete was called.
// A structure representing details about that specific entry.
// Default value is nil which means no callback and it prevents from unwrapping the oldest entry.
func WithOnRemoveWithMetadata(fn func(key string, entry []byte, keyMetadata bigcache.Metadata)) CacheOption {
	return func(c *bigcache.Config) {
		c.OnRemoveWithMetadata = fn
	}
}

// WithOnRemoveWithReason OnRemoveWithReason is a callback fired
// when the oldest entry is removed because of its expiration time or no space left
// for the new entry, or because delete was called.
// A constant representing the reason will be passed through.
// Default value is nil which means no callback and it prevents from unwrapping the oldest entry.
// Ignored if OnRemove is specified.
func WithOnRemoveWithReason(fn func(key string, entry []byte, reason bigcache.RemoveReason)) CacheOption {
	return func(c *bigcache.Config) {
		c.OnRemoveWithReason = fn
	}
}

// WithLogger Logger is a logging interface and used in combination with `Verbose`
// Defaults to bigcache.`DefaultLogger()`
func WithLogger(logger bigcache.Logger) CacheOption {
	return func(c *bigcache.Config) {
		c.Logger = logger
	}
}

// WithOnRemoveFilterSet OnRemoveFilterSet sets which remove reasons will trigger a call to OnRemoveWithReason.
// Filtering out reasons prevents bigcache from unwrapping them, which saves cpu.
func WithOnRemoveFilterSet(reasons ...bigcache.RemoveReason) CacheOption {
	return func(c *bigcache.Config) {
		c.OnRemoveFilterSet(reasons...)
	}
}
