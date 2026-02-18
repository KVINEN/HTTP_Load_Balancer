package main

import("fmt")

//struct to hold our backends
type Backend struct {
	URL *url.URL
	Alive bool
	mux sync.RWMutex
	ReverseProxy *httputil.ReverseProxy
}

//we need a way to track all backends
type ServerPool struct {
	backends []*Backend
	current uint64
}

u, _ := url.Parse("http://localhost:8080")
rb := httputil.NewSingleHostReverseProxy(u)

//initialize server and add this as handler
http.HandlerFunc(rp.ServeHTTP)

//we need a count, this is because we want to skip dead backends 
func (s *ServerPool) NextIndex() int {
	return int(atomic.AddUint64(&s.current, uint64(1)) % uint64(len(s.backends)))
}