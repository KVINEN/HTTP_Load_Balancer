package main

import (
	"log"
	"net"
	"net/url"
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
	for _, b := range s.backends {
		status := isBackendAlive(b.URL)
		b.SetAlive(status)

		msg := "ok"
		if !status {
			msg = "dead"
		}

		log.Printf("%s [%s]\n", b.URL, msg)
	}
}
