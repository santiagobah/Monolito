package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

type consulServiceEntry struct {
	Service struct {
		Service string `json:"Service"`
		Address string `json:"Address"`
		Port    int    `json:"Port"`
	} `json:"Service"`
}

type Backend struct {
	URL   *url.URL
	Proxy *httputil.ReverseProxy
}

type LoadBalancer struct {
	mu         sync.RWMutex
	backends   []*Backend
	counter    uint64
	consulAddr string
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func queryConsul(consulAddr, service string) ([]*Backend, error) {
	reqURL := fmt.Sprintf("%s/v1/health/service/%s?passing=true", consulAddr, service)
	client := http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get(reqURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var entries []consulServiceEntry
	if err := json.NewDecoder(resp.Body).Decode(&entries); err != nil {
		return nil, err
	}

	backends := make([]*Backend, 0, len(entries))
	for _, e := range entries {
		target, err := url.Parse(fmt.Sprintf("http://%s:%d", e.Service.Address, e.Service.Port))
		if err != nil {
			continue
		}
		backends = append(backends, &Backend{
			URL:   target,
			Proxy: httputil.NewSingleHostReverseProxy(target),
		})
	}
	return backends, nil
}

func (lb *LoadBalancer) refreshBackends() {
	backends, err := queryConsul(lb.consulAddr, "envios")
	if err != nil {
		log.Printf("consul query for service %q failed: %v", "envios", err)
		return
	}
	lb.mu.Lock()
	lb.backends = backends
	lb.mu.Unlock()
}

func (lb *LoadBalancer) startDiscoveryLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			lb.refreshBackends()
		}
	}()
}

func (lb *LoadBalancer) next() *Backend {
	lb.mu.RLock()
	defer lb.mu.RUnlock()
	n := len(lb.backends)
	if n == 0 {
		return nil
	}
	idx := atomic.AddUint64(&lb.counter, 1) % uint64(n)
	return lb.backends[idx]
}

func recibir_peticion(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.WriteHeader(http.StatusOK)
		return
	}

	backend := globalLB.next()
	if backend == nil {
		http.Error(w, "no healthy backend available", http.StatusServiceUnavailable)
		return
	}

	backend.Proxy.ServeHTTP(w, r)
}

var globalLB *LoadBalancer

func main() {
	consulAddr := getEnv("CONSUL_HTTP_ADDR", "http://consul:8500")
	globalLB = &LoadBalancer{consulAddr: consulAddr}

	for attempt := 0; attempt < 5; attempt++ {
		globalLB.refreshBackends()
		if len(globalLB.backends) > 0 {
			break
		}
		time.Sleep(2 * time.Second)
	}

	globalLB.startDiscoveryLoop(5 * time.Second)

	http.HandleFunc("/", recibir_peticion)
	fmt.Println("load balancer en puerto 9100")
	http.ListenAndServe(":9100", nil)
}
