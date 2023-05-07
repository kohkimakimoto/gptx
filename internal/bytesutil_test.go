package internal

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUint64tobAndBtouint64(t *testing.T) {
	for _, v := range []uint64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 123, 1234, 12345, 123456, 1234567, 12345678} {
		assert.Equal(t, v, btouint64(uint64tob(v)))
	}
}

type example struct {
	Str string
	Num int
	B   bool
	M   map[string]string
}

func TestSerializeAndDeserialize(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		b, err := serialize("test")
		assert.NoError(t, err)
		var ret string
		err = deserialize(b, &ret)
		assert.NoError(t, err)
		assert.Equal(t, "test", ret)
	})

	t.Run("int", func(t *testing.T) {
		b, err := serialize(123)
		assert.NoError(t, err)
		var ret int
		err = deserialize(b, &ret)
		assert.NoError(t, err)
		assert.Equal(t, 123, ret)
	})

	t.Run("bool", func(t *testing.T) {
		b, err := serialize(true)
		assert.NoError(t, err)
		var ret bool
		err = deserialize(b, &ret)
		assert.NoError(t, err)
		assert.Equal(t, true, ret)
	})

	t.Run("map", func(t *testing.T) {
		b, err := serialize(map[string]string{"test": "test"})
		assert.NoError(t, err)
		var ret map[string]string
		err = deserialize(b, &ret)
		assert.NoError(t, err)
		assert.Equal(t, map[string]string{"test": "test"}, ret)
	})

	t.Run("struct", func(t *testing.T) {
		b, err := serialize(example{
			Str: "test",
			Num: 123,
			B:   true,
			M:   map[string]string{"test": "test"},
		})
		assert.NoError(t, err)
		var ret example
		err = deserialize(b, &ret)
		assert.NoError(t, err)
		assert.Equal(t, example{
			Str: "test",
			Num: 123,
			B:   true,
			M:   map[string]string{"test": "test"},
		}, ret)
	})
}
