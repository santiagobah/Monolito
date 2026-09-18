package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
)

var workers = []string{
	"http://envios_1:8080",
	"http://envios_2:8080",
}

var contador_uso_worker = 0

func recibir_peticion(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.WriteHeader(http.StatusOK)
		return
	}

	worker_destino := workers[contador_uso_worker%len(workers)]
	contador_uso_worker++

	url_destino, _ := url.Parse(worker_destino)
	proxy := httputil.NewSingleHostReverseProxy(url_destino)
	proxy.ServeHTTP(w, r)
}

func main() {
	http.HandleFunc("/", recibir_peticion)
	fmt.Println("load balancer en puerto 9100")
	http.ListenAndServe(":9100", nil)
}
