package garnish

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

const cacheControl = "Cache-Control"
const ccNoCache = "no-cache"
const ccNoStore = "no-store"
const ccPrivate = "private"

var maxAgeReg = regexp.MustCompile(`max-age=(\d+)`)

// MustCompile parses a regular expression and returns,
// if successful, a Regexp object that can be used to match against text.
var sharedMaxAgeReg = regexp.MustCompile(`s-maxage=(\d+)`)

/**
Parse a string to the header Cache-Control
*/
func parseCacheControl(cc string) (cache bool, duration time.Duration) {
	if cc == ccPrivate || cc == ccNoCache || cc == ccNoStore || cc == "" {
		return false, 0
	}

	directives := strings.Split(cc, ",")
	for _, directive := range directives {
		directive = strings.ToLower(directive)
		age := maxAgeReg.FindStringSubmatch(directive)
		if len(age) > 0 {
			d, err := strconv.Atoi(age[1])
			if err != nil {
				return false, 0
			}
			cache = true
			duration = time.Duration(d) * time.Second
		}

		age = sharedMaxAgeReg.FindStringSubmatch(directive)
		if len(age) > 0 {
			d, err := strconv.Atoi(age[1])
			if err != nil {
				return false, 0
			}
			cache = true
			duration = time.Duration(d) * time.Second
		}
	}

	return //return cache, duration
}
