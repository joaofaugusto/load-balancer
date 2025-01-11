package main

import (
	"log"
	"net/http"

	"github.com/joaofaugusto/load-balancer/loadbalancer"
)

func main() {
	// List of backend servers
	servers := []string{
		"http://localhost:8080",
		"http://localhost:8081",
		"http://localhost:8082",
	}

	// Create a LoadBalancer instance
	lb := loadbalancer.NewLoadBalancer(servers)

	// Start the load balancer server
	log.Println("Starting load balancer on http://localhost:8000...")
	if err := http.ListenAndServe(":8000", lb); err != nil {
		log.Fatalf("Failed to start load balancer: %v", err)
	}
}
