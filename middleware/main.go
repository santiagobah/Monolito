package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

type consulServiceEntry struct {
	Service struct {
		Service string `json:"Service"`
		Address string `json:"Address"`
		Port    int    `json:"Port"`
	} `json:"Service"`
}

var (
	serviciosMapa map[string]string
	consulAddr    string
	mutexVivos    sync.RWMutex
	backendsMap   = make(map[string][]*url.URL)
)

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func queryConsulService(service string) ([]*url.URL, error) {
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

	var targets []*url.URL
	for _, e := range entries {
		u, err := url.Parse(fmt.Sprintf("http://%s:%d", e.Service.Address, e.Service.Port))
		if err == nil {
			targets = append(targets, u)
		}
	}
	return targets, nil
}

func RefreshConsulDiscovery() {
	serviciosConsul := []string{"nodos", "enlaces", "envios"}
	for _, srv := range serviciosConsul {
		urls, err := queryConsulService(srv)
		if err == nil && len(urls) > 0 {
			mutexVivos.Lock()
			backendsMap[srv] = urls
			mutexVivos.Unlock()
		} else {
			if fallback, ok := serviciosMapa[srv]; ok {
				u, err := url.Parse(fallback)
				if err == nil {
					mutexVivos.Lock()
					backendsMap[srv] = []*url.URL{u}
					mutexVivos.Unlock()
				}
			}
		}
	}
}

func Carga_servicios() {
	serv_disp, err_lec_serv_disp := os.ReadFile("servicios.json")
	if err_lec_serv_disp != nil {
		panic("No se detectaron los servicios" + err_lec_serv_disp.Error())
	}

	serviciosMapa = make(map[string]string)
	err_lec_serv_disp = json.Unmarshal(serv_disp, &serviciosMapa)

	if err_lec_serv_disp != nil {
		panic("Formato erróneo de los servicios declarados" + err_lec_serv_disp.Error())
	}

	fmt.Println("Se trajeron los servivios: ", serviciosMapa)
}

func Peticiones_mid(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.WriteHeader(http.StatusOK)
		return
	}

	key := strings.TrimPrefix(r.URL.Path, "/")
	parts := strings.Split(key, "/")
	if len(parts) > 0 {
		key = parts[0]
	}

	mutexVivos.RLock()
	targets, ok := backendsMap[key]
	mutexVivos.RUnlock()

	if !ok || len(targets) == 0 {
		w.Write([]byte(`Servicio no disponible`))
		return
	}

	urlDestino := targets[0]
	proxy := httputil.NewSingleHostReverseProxy(urlDestino)
	proxy.ServeHTTP(w, r)
}

func main() {
	consulAddr = getEnv("CONSUL_HTTP_ADDR", "http://consul:8500")
	Carga_servicios()
	RefreshConsulDiscovery()

	go func() {
		for {
			time.Sleep(5 * time.Second)
			RefreshConsulDiscovery()
		}
	}()

	http.HandleFunc("/", Peticiones_mid)
	fmt.Println("middleware en el puerto 8100")
	http.ListenAndServe(":8100", nil)
}
