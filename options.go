package rollingcache

import "time"

type Options struct {
	// Interval between rolling requests, will panic if not set
	Interval time.Duration
	// The maximum time this rolling cache can live for without being requested, or 0 for an unlimited time
	Lifetime time.Duration

	// Maximum amount of retries before stopping the request, or -1 for unlimited retries
	MaxRetries int

	// When true, this logs the actions the cache takes
	Debug bool
}
