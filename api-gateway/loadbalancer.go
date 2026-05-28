package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync/atomic"
)

type LoadBalancer struct {
	proxies []*httputil.ReverseProxy
	counter uint32
}

func NewLoadBalancer(targetUrls []string) *LoadBalancer {
	var proxies []*httputil.ReverseProxy

	for _, target := range targetUrls {
		parsedUrl, err := url.Parse(target)
		if err != nil {
			log.Fatalf("Error parsing target URL %s: %v", target, err)
		}
		proxies = append(proxies, httputil.NewSingleHostReverseProxy(parsedUrl))
	}

	return &LoadBalancer{
		proxies: proxies,
		counter: 0,
	}
}

// ServeHTTP makes the LoadBalancer implement the http.Handler interface
func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	count := atomic.AddUint32(&lb.counter, 1)

	index := count % uint32(len(lb.proxies))

	log.Printf("LoadBalancer: Routing request to instance %d", index)

	lb.proxies[index].ServeHTTP(w, r)
}
