package views

import (
	"encoding/json"
	"net/http"
)

// esta funcion nomas manda cualquier cosa como JSON, pa no repetir codigo en cada controller
func EnviarJSON(w http.ResponseWriter, datos interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(datos)
}

func EnviarError(w http.ResponseWriter, mensaje string, codigo int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(codigo)
	json.NewEncoder(w).Encode(map[string]string{"error": mensaje})
}
