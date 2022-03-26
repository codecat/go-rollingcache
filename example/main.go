package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	rollingcache "github.com/codecat/go-rollingcache"
)

func testserver() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Duration(rand.Int()%5) * time.Second)
		if rand.Int()%10 < 7 {
			w.WriteHeader(500)
			w.Write([]byte("Internal server error"))
			return
		}
		w.WriteHeader(200)
		w.Write([]byte("It works: " + time.Now().Format(time.RFC3339)))
	})
	http.ListenAndServe(":8900", nil)
}

func main() {
	go testserver()

	rollingcache.HttpHeaders["User-Agent"] = "github.com/codecat/go-rollingcache test"

	cache := rollingcache.Start("http://127.0.0.1:8900/", rollingcache.Options{
		Interval:     10 * time.Second,
		FailInterval: 3 * time.Second,
		MaxRetries:   3,
		Debug:        true,
	})

	for i := 0; i < 100; i++ {
		fmt.Printf("Request %d ... ", i+1)
		res := cache.Get()
		fmt.Printf("got: %s\n", string(res))
		time.Sleep(1337 * time.Millisecond)
	}

	fmt.Printf("Done!\n")
}
