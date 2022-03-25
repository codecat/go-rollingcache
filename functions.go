package rollingcache

import (
	"time"
)

var rollingCaches map[string]*Cache = make(map[string]*Cache)

// Returns the amount of rolling caches currently running.
func Count() int {
	return len(rollingCaches)
}

// Gets a cache by URL. The interval specifies the duration between requests. The lifetime of the cache is defined by
// the interval multiplied by 10.
func Get(url string, interval time.Duration) ([]byte, error) {
	return GetWithOptions(url, Options{
		Interval: interval,
		Lifetime: interval * 10,
	})
}

// Gets a cache by URL and a timeout duration for the first request response. The interval specifies the duration
// between requests. The lifetime of the cache is defined by the interval multiplied by 10.
func GetWithTimeout(url string, interval time.Duration, timeout time.Duration) ([]byte, error) {
	return GetWithOptionsTimeout(url, Options{
		Interval: interval,
		Lifetime: interval * 10,
	}, timeout)
}

// Gets a cache by URL and the given options. This will panic if the interval is not specified.
func GetWithOptions(url string, options Options) ([]byte, error) {
	return GetWithOptionsTimeout(url, options, 0)
}

// Gets a cache by URL, the given options, and a timeout duration for the first request response. This will panic if
// the interval is not specified.
func GetWithOptionsTimeout(url string, options Options, timeout time.Duration) ([]byte, error) {
	ret, ok := rollingCaches[url]
	if !ok {
		// Start new cache routine
		ret = Start(url, options)
	}

	if timeout == 0 {
		return ret.Get(), nil
	} else {
		return ret.GetWithTimeout(timeout)
	}
}

// Starts a new cache
func Start(url string, options Options) *Cache {
	// Panic if the interval is not set
	if options.Interval == 0 {
		panic("A rolling cache interval must be set!")
	}

	// Create the new cache
	ret := &Cache{
		URL:         url,
		Options:     options,
		LastRequest: time.Now(),
	}
	rollingCaches[url] = ret

	// Begin its update loop
	go ret.updateLoop()

	return ret
}
