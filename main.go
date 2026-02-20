package main

// imported libraries/packages
import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"        // this is a library package
	"sync/atomic" // this is a library sub-package
)

// struct to hold our backends
type Backend struct {
	URL          *url.URL
	Alive        bool
	mux          sync.RWMutex
	ReverseProxy *httputil.ReverseProxy
}

// we need a way to track all backends
type ServerPool struct {
	backends []*Backend
	current  uint64
}

// global instance of our pool
var serverPool ServerPool

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

// SetAlive for the backend
func (b *Backend) SetAlive(alive bool) {
	b.mux.Lock()
	b.Alive = alive
	b.mux.Unlock()
}

// IsAlive returns true when backend is alive
func (b *Backend) IsAlive() (alive bool) {
	b.mux.RLock()
	alive = b.Alive
	b.mux.RUnlock()
	return
}

// load balance balances the incoming request
func lb(w http.ResponseWriter, r *http.Request) {
	peer := serverPool.GetNextPeer()
	if peer != nil {
		peer.ReverseProxy.ServeHTTP(w, r)
		return
	}
	http.Error(w, "Service not available", http.StatusServiceUnavailable)
}

/*
server := http.Server {
	Addr: fmt.Sprintf(":%d", port)
	Handler: http.HandlerFunc(lb)
}


// atcively check for healty backends
proxy.ErrorHandler = func(writer http.ResponseWriter, request *http.Request, e error) {
	log.Printf("[%s] %s\n", serverUrl.Host, e.Error())
	retries := GetRetryFromContext(request)

	// tries the backend 3 times
	if retries < 3 {
		select {
		case <- time.After(10 * time.Millisecond):
			ctx := context.WithValue(request.Context(), Retry, retries + 1)
			proxy.ServeHTTP(write, request.WithContext(ctx))
		}
		return
	}

	// after 3 tries mark backend as down
	ServerPool.MarkBackendStatus(serverUrl, false)

	// if the same request routing for few attempts with different backends, increase the count
	attemps := GetAttampetsForContext(request)
	log.Printf("%s(%s) Attampting retry %d\n", request.RemoteAddr, request.URL.Path, attemps)
	ctx := context.WithValue(request.Context(), Attempts, attemps + 1)
	lb(writer, request.WithContext(ctx))
}
*/

func main() {

	// define backend URLs
	tokens := []string{"http://localhost:8081", "http://localhost:8082"}

	for _, tok := range tokens {
		serverUrl, _ := url.Parse(tok)

		// initialize the ReverseProxy for this backend
		proxy := httputil.NewSingleHostReverseProxy(serverUrl)

		// add to the pool
		serverPool.backends = append(serverPool.backends, &Backend{
			URL:          serverUrl,
			Alive:        true,
			ReverseProxy: proxy,
		})
	}
}
