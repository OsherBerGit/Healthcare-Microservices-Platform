package main

import (
	"log"
	"net/http"
)

func main() {
	patientServiceURLs := []string{
		"http://localhost:8081",
		"http://localhost:8083",
	}

	appointmentServiceURLs := []string{
		"http://localhost:8082",
	}

	mux := http.NewServeMux()

	patientLB := NewLoadBalancer(patientServiceURLs)
	appointmentLB := NewLoadBalancer(appointmentServiceURLs)

	mux.Handle("/api/patients/", verifyJWTMiddleware(patientLB))
	mux.Handle("/api/appointments/", idempotencyMiddleware(appointmentLB))

	mux.HandleFunc("/health", healthCheckHandler)

	log.Println("Starting Healthcare API Gateway with Load Balancing on port 8080...")
	if err := http.ListenAndServe(":8080", securityHeadersMiddleware(loggingMiddleware(mux))); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("API Gateway is Up and Running\n"))
}
