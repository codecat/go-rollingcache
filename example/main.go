package main

import (
	"fmt"
	"time"

	rollingcache "github.com/codecat/go-rollingcache"
)

func main() {
	rollingcache.HttpHeaders["User-Agent"] = "github.com/codecat/go-rollingcache test"

	cache := rollingcache.Start("https://httpbin.org/get", rollingcache.Options{
		Interval: 10 * time.Second,
		Debug:    true,
	})

	for i := 0; i < 100; i++ {
		res := cache.Get()
		fmt.Printf("Request %d data length: %d bytes\n", i+1, len(res))
		time.Sleep(1337 * time.Millisecond)
	}

	fmt.Printf("Done!\n")
}
