package controllers

import (
	"encoding/json"
	"net/http"

	"nodos/models"
	"nodos/views"
)

func Handler_nodos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")

	if r.Method == "OPTIONS" {
		return
	}

	if r.Method == "POST" {
		var n models.Nodo_n
		json.NewDecoder(r.Body).Decode(&n)

		id, error_ejec := models.Crear_nodo(n)

		if error_ejec != nil {
			views.Env_error(w, "Error al crear el nodo "+error_ejec.Error(), 500)
		}

		n.Id = id
		views.Env_json(w, n)
		return
	}

	if r.Method == "GET" {

		lista, error_ejec := models.Get_nodos()

		if error_ejec != nil {
			views.Env_error(w, "Error al traer el nodo "+error_ejec.Error(), 500)
			return
		}
		views.Env_json(w, lista)
		return
	}

}
