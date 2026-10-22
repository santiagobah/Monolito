package main

import (
	"fmt"
	"net/http"

	"enlaces/controllers"
	"enlaces/models"
)

func main() {
	models.Conectar_db()

	http.HandleFunc("/enlaces", controllers.Handler_enlaces)

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("vivo y coleando"))
	})

	fmt.Println("backend enlaces prendido en el puerto 8080")
	http.ListenAndServe(":8080", nil)
}
