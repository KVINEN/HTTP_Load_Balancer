package main

import (
	"log"
	"net"
	"net/url"
	"sync"
	"time"
)

// check if a backend by establising TCP connection
func isBackendAlive(u *url.URL) bool {
	timeout := 2 * time.Second

	conn, err := net.DialTimeout("tcp", u.Host, timeout)
	if err != nil {
		log.Println("Site unreachable, error: ", err)
	}

	// always close connection to avoid leaking file descriptors
	_ = conn.Close()
	return true
}

func (s *ServerPool) HealthCheck() {
	var wg sync.WaitGroup // create wate group
	for _, b := range s.backends {
		wg.Add(1) // tell group we have one more task
		go func(b *Backend) {
			defer wg.Done()
			alive := isBackendAlive(b.URL)
			b.SetAlive(alive)
		}(b)

		wg.Wait()
	}
}
