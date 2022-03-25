package rollingcache

import (
	"net/http"
)

var HttpClient http.Client = http.Client{}
var HttpHeaders map[string]string = make(map[string]string)
