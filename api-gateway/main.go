package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func main() {
	patientServiceURL := "http://localhost:8081"
	appointmentServiceURL := "http://localhost:8082"

	mux := http.NewServeMux()

	mux.Handle("/api/patients/", newProxy(patientServiceURL))

	mux.Handle("/api/appointments/", idempotencyMiddleware(newProxy(appointmentServiceURL)))

	mux.HandleFunc("/health", healthCheckHandler)

	log.Println("Starting Healthcare API Gateway on port 8080...")
	if err := http.ListenAndServe(":8080", securityHeadersMiddleware(loggingMiddleware(mux))); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Incoming Request: [Method: %s] [Path: %s]", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func securityHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

		next.ServeHTTP(w, r)
	})
}

func idempotencyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idempotencyKey := r.Header.Get("Idempotency-Key")

		if idempotencyKey == "" {
			http.Error(w, "Missing required header: Idempotency-Key", http.StatusBadRequest)
			return
		}

		log.Printf("Received request with Idempotency-Key: %s", idempotencyKey)

		next.ServeHTTP(w, r)
	})
}

func newProxy(targetHost string) *httputil.ReverseProxy {
	targetUrl, err := url.Parse(targetHost)
	if err != nil {
		log.Fatalf("Error parsing target URL: %v", err)
	}
	return httputil.NewSingleHostReverseProxy(targetUrl)
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("API Gateway is Up and Running\n"))
}
