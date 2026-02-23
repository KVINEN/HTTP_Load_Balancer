# HTTP_Load_Balancer

A functional, Layer 7 HTTP Load Balancer prototype built in Go as part of my "1 New Project Every Month of 2026" challange.

## 📅 The 2026 Challange

This project is Part 2 of 12 in my year-long journey to build different projects in 2026. The focus is on creating functional prototypes or MVPs (Minimum Viable Products) rather than 100% feature-complete applications.

## 🚀 Overview

HTTP_Load_Balancer is a concurrent networking tool designed to distribute incomming HTTP trafic across multiple  backend servers. It ensures high availablility by automatically skipping unhealthy backends and retrying failed requests.

### Features
* **Round-Robin Scheduling:** Evenly distribute requests across the server pool using atomic counters for tread-safe indexing.
* **Active Health Checks:** A background goroutine monitors backend health every 2 minutes by attempting TCP connections.
* **Dynamic Retries & Failover:** Automatically retries failed requests up to 3 times before marking a backend as down and routing to a new peer.
* **Standard Library Design:** Built entirely using Go's standard library, leveraging ```net/http/httputil``` for reverse proxying and ```sync``` for concurrency management.

## 🛠️ Technical Details
* **Language:** Go 1.25.6
* **Port:** Defaults to 8080 (configurable via flags).
* **Concurrency:** Uses ```sync.RWMutex``` for safe state sharing and ```sync.WaitGroup``` for parallel health checks.
* **Proxying:** Implements a custom ```ErrorHandler``` to manage backend failures and context-based retry tracking.

## 🏗️ Project Structure
```
.
├── main.go          # Entry point, reverse proxy logic, and server pool management
├── healthcheck.go   # TCP-based health checking and status updates
├── go.mod           # Project dependencies and Go version
├── LICENSE          # MIT License
└── README.md        # Project documentation
```

##🚦 Getting Started

### Prerequisites

* Go compiler (version 1.25.6 or higher recommended).

### Compilation & Running

To run the load balancer, provide a comma-separated list of your backend server URLs:
```
go run main.go healthcheck.go -backends="http://localhost:8081,http://localhost:8082" -port=8080
```

### Testing

You can verify the load balancer is working by sending requests to the configured port:
```
# In a separate terminal
curl http://localhost:8080
```
The server will log which backend the request is being proxied to, or return a "Service not available" error if no backends are alive.

## 📝 License

This project is licensed under the MIT License.