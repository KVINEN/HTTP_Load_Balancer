package main

import("fmt")

// struct to hold our backends
type Backend struct {
	URL *url.URL
	Alive bool
	mux sync.RWMutex
	ReverseProxy *httputil.ReverseProxy
}

// we need a way to track all backends
type ServerPool struct {
	backends []*Backend
	current uint64
}

u, _ := url.Parse("http://localhost:8080")
rb := httputil.NewSingleHostReverseProxy(u)

// initialize server and add this as handler
http.HandlerFunc(rp.ServeHTTP)

/* 
	we need a count, this is because we want to skip dead backends.
	automaticly increas value by one and return the index by modding
	with the length of the slice. 
	this means the value always will be
	between 0 and length of the slice. Interested in partilcular index, 
	not total value.
*/ 
func (s *ServerPool) NextIndex() int {
	return int(atomic.AddUint64(&s.current, uint64(1)) % uint64(len(s.backends)))
}

// GetNextPeer returns next active peer to take a connection
func (s *ServerPool) GetNextPeer() *Backend {
	// loop entire backends to find an Alive backend
	next := s.NextIndex()
	l := len(s.backends) + next // start from nex and move full cycle 
	for i := next; i < l; i++ {
		idx := i % len(s.backends) // take an index by modding with length 
		// if we have an alive backend, use it and store if its not the original one
		if s.backends[idx].IsAlive() {
			if i != next {
				atomic.StoreUint64(&s.current, uint64(idx)) // mark the current one
			}
			return s.backend[idx]
		}
	}
	return nil // nil = zero value in golang
}