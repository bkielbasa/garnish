package garnish

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseCacheControl(t *testing.T) {
	headers := map[string]struct {
		givenHeader    string
		shouldBeCached bool
		cacheTime      time.Duration
	}{
		"empty header": {
			givenHeader:    "",
			shouldBeCached: false,
		},
		"no-cache": {
			givenHeader:    "no-cache",
			shouldBeCached: false,
		},
		"private": {
			givenHeader:    "private",
			shouldBeCached: false,
		},
		"max-age": {
			givenHeader:    "max-age=123",
			shouldBeCached: true,
			cacheTime:      time.Second * 123,
		},
		"max-age uppercase": {
			givenHeader:    "MAX-AGE=123",
			shouldBeCached: true,
			cacheTime:      time.Second * 123,
		},
		"s-max-age": {
			givenHeader:    "s-max-age=123",
			shouldBeCached: true,
			cacheTime:      time.Second * 123,
		},
		"s-max-age overrides max-age": {
			givenHeader:    "s-max-age=123,s-max-age=321",
			shouldBeCached: true,
			cacheTime:      time.Second * 321,
		},
	}

	for name, tcase := range headers {
		t.Run(name, func(t *testing.T) {
			shouldBeCached, cacheTime := parseCacheControl(tcase.givenHeader)
			assert.Equal(t, tcase.shouldBeCached, shouldBeCached)
			assert.Equal(t, tcase.cacheTime, cacheTime)
		})
	}
}
