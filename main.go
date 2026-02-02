package main

import {
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
	"time"
}

//simple struct to hold out backends
type Backend struct {
	URL				*url url
	Alive			bool
	mux				sync.RWMutex
	ReverseProxy 	*httpputil.ReverseProxy
}

//struct to track all backends in out load balancer
type ServerPool struct {
	backends 	[]*Backend
	current 	uint64
}