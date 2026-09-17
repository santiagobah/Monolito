package main

import (
	"fmt"
	"net/http"

	"envios/controllers"
	"envios/models"
)

func main() {
	models.Conectar_db()

	http.HandleFunc("/envios", controllers.Handler_envios)

	fmt.Println("backend envios prendido en el puerto 8080")
	http.ListenAndServe(":8080", nil)
}
