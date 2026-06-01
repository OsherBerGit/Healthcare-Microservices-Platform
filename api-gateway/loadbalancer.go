package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type Backend struct {
	URL          *url.URL
	ReverseProxy *httputil.ReverseProxy
	Alive        bool
	HealthPath   string
	mux          sync.RWMutex
}

func (b *Backend) SetAlive(alive bool) {
	b.mux.Lock()
	b.Alive = alive
	b.mux.Unlock()
}

func (b *Backend) IsAlive() bool {
	b.mux.RLock()
	alive := b.Alive
	b.mux.RUnlock()
	return alive
}

type LoadBalancer struct {
	backends []*Backend
	counter  uint32
}

func NewLoadBalancer(targetUrls []string) *LoadBalancer {
	var backends []*Backend

	for _, target := range targetUrls {
		parsedUrl, err := url.Parse(target)
		if err != nil {
			log.Fatalf("Error parsing target URL %s: %v", target, err)
		}

		healthPath := "/health"
		if strings.Contains(target, "patients") {
			healthPath = "/api/patients"
		} else if strings.Contains(target, "appointments") {
			healthPath = "/api/appointments"
		}

		backend := &Backend{
			URL:          parsedUrl,
			ReverseProxy: httputil.NewSingleHostReverseProxy(parsedUrl),
			Alive:        true,
			HealthPath:   healthPath,
		}
		backends = append(backends, backend)
	}

	lb := &LoadBalancer{
		backends: backends,
		counter:  0,
	}

	go lb.healthCheckLoop()

	return lb
}

func (lb *LoadBalancer) healthCheckLoop() {
	ticker := time.NewTicker(10 * time.Second)

	for {
		<-ticker.C

		var wg sync.WaitGroup

		for _, b := range lb.backends {
			wg.Add(1)

			go func(backend *Backend) {
				defer wg.Done()

				alive := isBackendAlive(backend)
				backend.SetAlive(alive)

				if !alive {
					log.Printf("HealthCheck: Backend %s is DOWN", backend.URL)
				}
			}(b)
		}
		wg.Wait()
	}
}

func isBackendAlive(b *Backend) bool {
	client := http.Client{Timeout: 2 * time.Second}

	targetURL := b.URL.String() + b.HealthPath

	resp, err := client.Get(targetURL)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode != http.StatusBadGateway && resp.StatusCode != http.StatusServiceUnavailable
}

func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	backendsCount := uint32(len(lb.backends))

	for i := uint32(0); i < backendsCount; i++ {
		count := atomic.AddUint32(&lb.counter, 1)
		index := count % backendsCount
		backend := lb.backends[index]

		if backend.IsAlive() {
			backend.ReverseProxy.ServeHTTP(w, r)
			return
		}
	}

	log.Printf("All backends are down!")
	http.Error(w, "Service Unavailable: All upstream nodes are down", http.StatusServiceUnavailable)
}
