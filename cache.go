package rollingcache

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Cache contains information about a currently running rolling cache.
type Cache struct {
	// The URL that this cache has to request
	URL string

	// The options for this cache
	Options Options

	// Latest fetched data, or nil when it hasn't fetched the first result yet
	Data []byte
	// Whether the data has been successfully fetched the last time it was requested
	Success bool

	// Time the cache was last requested
	LastRequest time.Time
	// Time the cache was last updated
	LastUpdate time.Time
}

// Returns true if there is available data for the cache.
func (cache *Cache) Available() bool {
	return cache.Data != nil
}

// Blocks forever until there is available data for the cache and then returns it.
func (cache *Cache) Get() []byte {
	// Update our last request time
	cache.LastRequest = time.Now()

	// Wait for rolling cache result to exist
	for {
		//TODO: We can probably do this a bit safer with channels

		// If it's available, stop waiting
		if cache.Available() {
			break
		}

		// Wait
		time.Sleep(1 * time.Millisecond)
	}

	// Return data
	return cache.Data
}

// Blocks until the timeout is reached, or there is available data for the cache and then returns it.
func (cache *Cache) GetWithTimeout(timeout time.Duration) ([]byte, error) {
	// Update our last request time
	cache.LastRequest = time.Now()

	// Determine timeout time
	timeoutTime := time.Now().Add(timeout)

	// Wait for rolling cache result to exist
	for time.Now().Before(timeoutTime) {
		//TODO: We can probably do this a bit safer with channels

		// If it exists, return the data
		if cache.Available() {
			return cache.Data, nil
		}

		// Wait
		time.Sleep(1 * time.Millisecond)
	}

	// When we get here, we've timed out
	return nil, errors.New("timeout while getting rolling cache")
}

func (cache *Cache) updateLoop() {
	for {
		// Stop if we've reached the cache lifetime
		if cache.Options.Lifetime != 0 && time.Now().After(cache.LastRequest.Add(cache.Options.Lifetime)) {
			if cache.Options.Debug {
				fmt.Printf("Reached rolling cache end of life: %s\n", cache.URL)
			}
			break
		}

		// When we perform the request, we might need to retry if it fails
		retry := 0
		for {
			if cache.Options.Debug {
				fmt.Printf("Making rolling cache request: %s\n", cache.URL)
			}

			// Perform the request
			req, _ := http.NewRequest("GET", cache.URL, nil)
			for k, v := range HttpHeaders {
				req.Header.Add(k, v)
			}
			resp, err := HttpClient.Do(req)

			// Check if there was an error
			if err != nil {
				fmt.Printf("Error while requesting rolling cache: %s", err.Error())

				if cache.Options.MaxRetries != -1 && retry >= cache.Options.MaxRetries {
					// We couldn't recover by retrying, so we stop here
					fmt.Printf("Too many server errors when requesting rolling cache: %s\n", resp.Status)
					cache.Success = false
					break
				}
				retry++
				continue
			}

			// We are expecting 2xx status codes, otherwise we consider it a fail
			if resp.StatusCode < 200 || resp.StatusCode > 299 {
				if cache.Options.MaxRetries != -1 && retry >= cache.Options.MaxRetries {
					// We couldn't recover by retrying, so we stop here
					fmt.Printf("Too many server errors when requesting rolling cache: %s\n", resp.Status)
					cache.Success = false
					break
				}
				retry++
				continue
			}

			// If there was an error, see if we can recover by retrying
			if err != nil {
				if cache.Options.MaxRetries != -1 && retry >= cache.Options.MaxRetries {
					// We couldn't recover by retrying, so we stop here
					fmt.Printf("Too many errors when requesting rolling cache: %s\n", err.Error())
					cache.Success = false
					break
				}
				retry++
				continue
			}

			// Get the response data
			data, err := io.ReadAll(resp.Body)

			// If there was an error, see if we can recover by retrying
			if err != nil {
				if cache.Options.MaxRetries != -1 && retry >= cache.Options.MaxRetries {
					// We couldn't recover by retrying, so we stop here
					fmt.Printf("Too many errors when reading rolling cache response body: %s\n", err.Error())
					cache.Success = false
					break
				}
				retry++
				continue
			}

			// Update the data, and when we last updated it
			cache.Data = data
			cache.Success = true
			cache.LastUpdate = time.Now()

			if cache.Options.Debug {
				fmt.Printf("Rolling cache data updated with %d bytes\n", len(data))
			}

			// Break out of the retry loop since we're done
			break
		}

		// Wait the specified amount of time until the next request
		if cache.Success {
			time.Sleep(cache.Options.Interval)
		} else {
			time.Sleep(cache.Options.FailInterval)
		}
	}

	// Remove ourselves from the cache
	delete(rollingCaches, cache.URL)
}
