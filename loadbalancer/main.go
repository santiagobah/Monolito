package main

import (
	"fmt"
	"net/http"
)

var workers = []string{
	"http://localhost:8091",
	"http://localhost:8092",
}

var contador_uso_worker = 0

func recibir_peticion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")

	if r.Method == "OPTIONS" {
		return
	}

	worker_destino := workers[contador_uso_worker%len(workers)]
	contador_uso_worker++

	http.Redirect(w, r, worker_destino+r.URL.Path, http.StatusTemporaryRedirect)

}

func main() {
	http.HandleFunc("/", recibir_peticion)
	fmt.Println("load balancer en puerto 9100")
	http.ListenAndServe(":9100", nil)
}
