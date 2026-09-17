package controllers

import (
	"encoding/json"
	"net/http"

	"envios/models"
	"envios/views"
)

type Peticion_envio struct {
	Origen  int `json:"origen"`
	Destino int `json:"destino"`
}

func Handler_envios(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")

	if r.Method == "OPTIONS" {
		return
	}

	if r.Method != "POST" {
		views.Env_error(w, "método no permitido", 405)
		return
	}

	var pet_env Peticion_envio
	json.NewDecoder(r.Body).Decode(&pet_env)

	estado_envio, err_envio := models.Envio_paquetes(pet_env.Origen, pet_env.Destino)

	if err_envio != nil {
		views.Env_error(w, err_envio.Error(), 400)
		return
	}

	views.Env_json(w, estado_envio)
}
