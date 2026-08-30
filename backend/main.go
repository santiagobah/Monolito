package main

import (
	"fmt"
	"net/http"

	"fase1backend/controllers"
	"fase1backend/models"
)

func main() {
	models.ConectarDB()

	// rutas del monolito, todo bien plano aqui
	http.HandleFunc("/nodos", controllers.NodosHandler)
	http.HandleFunc("/enlaces", controllers.EnlacesHandler)
	http.HandleFunc("/enviar", controllers.EnviarPaqueteHandler)

	fmt.Println("backend  en el puerto 8080")
	http.ListenAndServe(":8080", nil)
}
