package main

import (
	"fmt"
	"net/http"

	"nodos/controllers"
	"nodos/models"
)

func main() {
	models.Conectar_db()

	http.HandleFunc("/nodos", controllers.Handler_nodos)

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("vivo y coleando"))
	})

	fmt.Println("backend nodos prendido en el puerto 8080")
	http.ListenAndServe(":8080", nil)
}
