package views

import (
	"encoding/json"
	"net/http"
)

func Env_json(w http.ResponseWriter, datos interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(datos)
}

func Env_error(w http.ResponseWriter, mensaje string, codigo int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(codigo)
	json.NewEncoder(w).Encode(map[string]string{"error": mensaje})
}
