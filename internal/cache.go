package internal

import (
	"encoding/hex"
	"fmt"
	bolt "go.etcd.io/bbolt"
	"os"
	"sync"
	"time"
)

const (
	CacheBucketCaches = "caches"

	CacheBucketOrder = "order"
)

// CacheItemNotFoundError is an error that is returned when the cache is not found.
type CacheItemNotFoundError struct {
	// Key is sha256 hash
	Key []byte
}

// Error returns the error message.
func (e *CacheItemNotFoundError) Error() string {
	return fmt.Sprintf("Cache not found: %s", hex.EncodeToString(e.Key[:]))
}

type CacheManager struct {
	DBPath    string
	MaxLength int
	cache     *Cache
	lock      sync.RWMutex
}

func (m *CacheManager) Open() (*Cache, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	if m.cache != nil {
		// already opened
		return m.cache, nil
	}

	opt := *bolt.DefaultOptions
	opt.Timeout = 1 * time.Second
	db, err := bolt.Open(m.DBPath, os.FileMode(0600), &opt)
	if err != nil {
		return nil, err
	}

	m.cache = &Cache{
		m:         m,
		db:        db,
		maxLength: m.MaxLength,
	}

	return m.cache, nil
}

func (m *CacheManager) Close() error {
	m.lock.Lock()
	defer m.lock.Unlock()

	if m.cache == nil {
		// already closed
		return nil
	}
	err := m.cache.Close()
	if err != nil {
		return err
	}
	m.cache = nil
	return nil
}

func (m *CacheManager) Refresh() error {
	if err := m.Close(); err != nil {
		return nil
	}
	if err := os.RemoveAll(m.DBPath); err != nil {
		return err
	}
	cache, err := m.Open()
	if err != nil {
		return err
	}
	return cache.Init()
}

type Cache struct {
	m         *CacheManager
	db        *bolt.DB
	maxLength int
}

func (c *Cache) Close() error {
	if c.db == nil {
		// already closed
		return nil
	}
	err := c.db.Close()
	if err != nil {
		return err
	}
	c.db = nil
	c.m.cache = nil
	return nil
}

func (c *Cache) Init() error {
	err := c.db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists([]byte(CacheBucketCaches)); err != nil {
			return err
		}
		if _, err := tx.CreateBucketIfNotExists([]byte(CacheBucketOrder)); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *Cache) Set(key []byte, value []byte) error {
	return c.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(CacheBucketCaches))

		if b.Get(key) == nil {
			// if the key is not found, the item is new.
			// So we need to add it to the order bucket
			// and remove the oldest item if the cache is full.

			// add key to order bucket
			orderBucket := tx.Bucket([]byte(CacheBucketOrder))
			orderNo, _ := orderBucket.NextSequence()
			if err := orderBucket.Put(uint64tob(orderNo), key); err != nil {
				return err
			}

			// plus 1 means that the new item is added
			keyN := orderBucket.Stats().KeyN + 1
			// remove the oldest item if the cache is full
			for c.maxLength != 0 && keyN > c.maxLength {
				// get the oldest key
				oldestOrderNo, oldestKey := orderBucket.Cursor().First()
				if oldestOrderNo != nil {
					// remove the oldest orderNo from order bucket
					if err := orderBucket.Delete(oldestOrderNo); err != nil {
						return err
					}
					// remove the oldest key from cache bucket
					if err := b.Delete(oldestKey); err != nil {
						return err
					}
					keyN--
				}
			}
		}

		err := b.Put(key, value)
		return err
	})
}

func (c *Cache) Get(key []byte) ([]byte, error) {
	var value []byte
	err := c.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(CacheBucketCaches))
		value = b.Get(key)
		if value == nil {
			return &CacheItemNotFoundError{Key: key}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return value, nil
}

func (c *Cache) Size() (size int) {
	err := c.db.View(func(tx *bolt.Tx) error {
		orderBucket := tx.Bucket([]byte(CacheBucketOrder))
		size = orderBucket.Stats().KeyN
		return nil
	})
	if err != nil {
		panic(err)
	}
	return size
}
