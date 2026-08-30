package controllers

import (
	"encoding/json"
	"net/http"

	"fase1backend/models"
	"fase1backend/views"
)

func EnlacesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		return
	}

	if r.Method == "POST" {
		var e models.Enlace
		json.NewDecoder(r.Body).Decode(&e)

		id, err := models.CrearEnlace(e)
		if err != nil {
			views.EnviarError(w, "no se pudo crear el enlace", 500)
			return
		}
		e.ID = id
		views.EnviarJSON(w, e)
		return
	}

	if r.Method == "GET" {
		lista, err := models.ObtenerEnlaces()
		if err != nil {
			views.EnviarError(w, "error trayendo enlaces", 500)
			return
		}
		views.EnviarJSON(w, lista)
		return
	}
}
