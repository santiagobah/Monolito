package controllers

import (
	"encoding/json"
	"net/http"

	"fase1backend/models"
	"fase1backend/views"
)

// este handler atiende GET (listar) y POST (crear) para /nodos
func NodosHandler(w http.ResponseWriter, r *http.Request) {
	// esto es pa que el front en otro puerto le pueda pegar sin bronca de CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		return
	}

	if r.Method == "POST" {
		var n models.Nodo
		json.NewDecoder(r.Body).Decode(&n)

		id, err := models.CrearNodo(n)
		if err != nil {
			views.EnviarError(w, "no se pudo crear el nodo: "+err.Error(), 500)
			return
		}
		n.ID = id
		views.EnviarJSON(w, n)
		return
	}

	if r.Method == "GET" {
		lista, err := models.ObtenerNodos()
		if err != nil {
			views.EnviarError(w, "error trayendo nodos", 500)
			return
		}
		views.EnviarJSON(w, lista)
		return
	}
}
