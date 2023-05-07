package internal

import (
	"crypto/sha256"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCacheItemNotFoundError_Error(t *testing.T) {
	sha := sha256.Sum256([]byte("test"))
	err := &CacheItemNotFoundError{Key: sha[:]}
	assert.Error(t, err)
	assert.Regexp(t, "Cache not found: .+", err.Error())
}

func TestCache_Set(t *testing.T) {
	t.Run("Set and get cache", func(t *testing.T) {
		cm := testCacheManager(t)
		c, err := cm.Open()
		assert.NoError(t, err)
		defer c.Close()

		sha := sha256.Sum256([]byte("test-key"))
		// no cache item because cache is not set yet
		item, err := c.Get(sha[:])
		assert.Error(t, err)
		assert.Nil(t, item)

		err = c.Set(sha[:], []byte("test-value"))
		assert.NoError(t, err)

		item, err = c.Get(sha[:])
		assert.NoError(t, err)
		assert.Equal(t, []byte("test-value"), item)
	})

	t.Run("Set with maxLength", func(t *testing.T) {
		cm := testCacheManager(t)
		cm.MaxLength = 10
		c, err := cm.Open()
		assert.NoError(t, err)
		defer c.Close()

		for i := 0; i < 30; i++ {
			sha := sha256.Sum256([]byte(fmt.Sprintf("test-key-%d", i)))
			c.Set(sha[:], []byte(fmt.Sprintf("test-value-%d", i)))
		}

		for i := 0; i < 30; i++ {
			sha := sha256.Sum256([]byte(fmt.Sprintf("test-key-%d", i)))
			item, err := c.Get(sha[:])
			if i < 20 {
				// max length is 10, so first 0 - 19 items is deleted
				assert.Error(t, err)
				assert.Nil(t, item)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, []byte(fmt.Sprintf("test-value-%d", i)), item)
			}
		}

		assert.Equal(t, 10, c.Size())

	})
}

func TestCacheManager_Refresh(t *testing.T) {
	cm := testCacheManager(t)
	c, err := cm.Open()
	assert.NoError(t, err)
	defer c.Close()

	for i := 0; i < 30; i++ {
		sha := sha256.Sum256([]byte(fmt.Sprintf("test-key-%d", i)))
		c.Set(sha[:], []byte(fmt.Sprintf("test-value-%d", i)))
	}
	assert.Equal(t, 30, c.Size())

	// refresh and reopen cache
	err = cm.Refresh()
	assert.NoError(t, err)
	c, err = cm.Open()
	assert.NoError(t, err)
	defer c.Close()

	assert.Equal(t, 0, c.Size())

}
