package controllers

import (
	"encoding/json"
	"net/http"

	"fase1backend/models"
	"fase1backend/views"
)

type PeticionEnvio struct {
	Origen  int `json:"origen"`
	Destino int `json:"destino"`
}

// este handler simula que se manda un paquete de un nodo a otro y guarda el log
func EnviarPaqueteHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		return
	}

	if r.Method != "POST" {
		views.EnviarError(w, "metodo no permitido", 405)
		return
	}

	var peticion PeticionEnvio
	json.NewDecoder(r.Body).Decode(&peticion)

	resultado, err := models.SimularEnvio(peticion.Origen, peticion.Destino)
	if err != nil {
		views.EnviarError(w, err.Error(), 400)
		return
	}

	views.EnviarJSON(w, resultado)
}
