package controllers

import (
	"encoding/json"
	"net/http"

	"enlaces/models"
	"enlaces/views"
)

func Handler_enlaces(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")

	if r.Method == "OPTIONS" {
		return
	}

	if r.Method == "POST" {
		var e models.Enlace
		json.NewDecoder(r.Body).Decode(&e)

		id, error_ejec := models.Crear_enlace(e)

		if error_ejec != nil {
			views.Env_error(w, "Error al crear el enlace "+error_ejec.Error(), 500)
		}

		e.Id = id
		views.Env_json(w, e)
		return
	}

	if r.Method == "GET" {

		lista, error_ejec := models.Get_enlaces()

		if error_ejec != nil {
			views.Env_error(w, "Error al traer el enlace "+error_ejec.Error(), 500)
			return
		}
		views.Env_json(w, lista)
		return
	}

}
