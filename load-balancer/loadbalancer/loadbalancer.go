package loadbalancer

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync/atomic"
)

// LoadBalancer struct for managing servers and state
type LoadBalancer struct {
	servers []*url.URL
	current uint32
}

// NewLoadBalancer initializes a LoadBalancer with the given server URLs
func NewLoadBalancer(serverURLs []string) *LoadBalancer {
	var servers []*url.URL
	for _, serverURL := range serverURLs {
		parsedURL, err := url.Parse(serverURL)
		if err != nil {
			log.Fatalf("Invalid server URL: %s", err)
		}
		servers = append(servers, parsedURL)
	}

	return &LoadBalancer{
		servers: servers,
		current: 0,
	}
}

// GetNextServer retrieves the next server URL using round-robin
func (lb *LoadBalancer) GetNextServer() *url.URL {
	index := atomic.AddUint32(&lb.current, 1)
	return lb.servers[index%uint32(len(lb.servers))]
}

// ServeHTTP forwards the request to the selected server
func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	target := lb.GetNextServer()
	proxy := httputil.NewSingleHostReverseProxy(target)

	// Log the redirection
	fmt.Printf("Forwarding request to: %s\n", target.String())

	// Rewrite the request's URL to point to the target server
	r.Host = target.Host
	r.URL.Host = target.Host
	r.URL.Scheme = target.Scheme

	// Serve the request using the proxy
	proxy.ServeHTTP(w, r)
}
