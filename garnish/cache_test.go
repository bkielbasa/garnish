package garnish

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStoringAndRetrieving(t *testing.T) {
	c := newCache()
	data := []byte("data to store")
	c.store("key", data, 0)

	assert.Equal(t, data, c.get("key"))
}

func TestNotReachedTimeout(t *testing.T) {
	c := newCache()
	data := []byte("data to store")
	c.store("key", data, time.Millisecond*100)
	time.Sleep(time.Millisecond * 80)

	assert.Equal(t, data, c.get("key"))
}

func TestTimeout(t *testing.T) {
	c := newCache()
	data := []byte("data to store")
	c.store("key", data, time.Millisecond*100)
	time.Sleep(time.Millisecond * 100)

	assert.Equal(t, []byte(nil), c.get("key"))
}
