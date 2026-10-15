package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

type ConsulService struct {
	Node           string `json:"Node"`
	Address        string `json:"Address"`
	ServiceAddress string `json:"ServiceAddress"`
	ServicePort    int    `json:"ServicePort"`
}

var routes map[string]string

func loadRoutes() {
	file, err := os.ReadFile("routes.json")
	if err != nil {
		panic(err)
	}
	json.Unmarshal(file, &routes)
}

func getBackendFromConsul(serviceName string) (string, error) {
	consulHost := os.Getenv("CONSUL_HOST")
	if consulHost == "" {
		consulHost = "consul"
	}
	urlStr := fmt.Sprintf("http://%s:8500/v1/health/service/%s?passing=true", consulHost, serviceName)
	resp, err := http.Get(urlStr)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var services []ConsulService
	if err := json.Unmarshal(body, &services); err != nil {
		return "", err
	}

	if len(services) == 0 {
		return "", fmt.Errorf("no hay instancias saludables para el servicio %s", serviceName)
	}

	s := services[0]
	addr := s.ServiceAddress
	if addr == "" {
		addr = s.Address
	}
	return fmt.Sprintf("http://%s:%d", addr, s.ServicePort), nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	var serviceName string
	var matchedPrefix string

	for prefix, sName := range routes {
		if strings.HasPrefix(r.URL.Path, prefix) {
			serviceName = sName
			matchedPrefix = prefix
			break
		}
	}

	if serviceName == "" {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	backendURLStr, err := getBackendFromConsul(serviceName)
	if err != nil {
		http.Error(w, "Service Unavailable: "+err.Error(), http.StatusServiceUnavailable)
		return
	}

	targetURL, _ := url.Parse(backendURLStr)
	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	r.URL.Path = strings.TrimPrefix(r.URL.Path, matchedPrefix)
	if r.URL.Path == "" {
		r.URL.Path = "/"
	}

	proxy.ServeHTTP(w, r)
}

func main() {
	loadRoutes()
	http.HandleFunc("/", handler)
	fmt.Println("Gateway corriendo en el puerto 8080...")
	http.ListenAndServe(":8080", nil)
}
