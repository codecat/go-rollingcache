# Rolling cache
This is a simple "rolling cache" module for Go. It works as a continuous request loop, where the data is only updated in a specific interval.

## Example usage
On-demand request API:
```go
// Start a cache loop with:
// - An interval of 30 seconds
// - An inactivity lifetime of 5 minutes (interval multiplied by 10)
// And return the byte data into `res`.
res, err := rollingcache.Get("https://httpbin.org/get", 30 * time.Second)
```

Continuous request API:
```go
// Start a cache loop with some options. The returned object is a rollingcache.Cache pointer.
// You can use this pointer to get the cache data.
cache := rollingcache.Start("https://httpbin.org/get", rollingcache.Options{
	Interval: 10 * time.Second,
})
```
