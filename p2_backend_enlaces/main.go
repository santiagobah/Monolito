package main

import (
	"fmt"
	"net/http"

	"enlaces/controllers"
	"enlaces/models"
)

func main() {
	models.Conectar_db()

	http.HandleFunc("/nodos", controllers.Handler_enlaces)

	fmt.Println("backend enlaces prendido en el puerto 8080")
	http.ListenAndServe(":8080", nil)
}
