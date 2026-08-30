package models

// un enlace es como el cable que conecta dos nodos
type Enlace struct {
	ID            int `json:"id"`
	NodoOrigenID  int `json:"nodo_origen_id"`
	NodoDestinoID int `json:"nodo_destino_id"`
}

func CrearEnlace(e Enlace) (int, error) {
	var id int
	query := "INSERT INTO enlaces (nodo_origen_id, nodo_destino_id) VALUES ($1, $2) RETURNING id"
	err := DB.QueryRow(query, e.NodoOrigenID, e.NodoDestinoID).Scan(&id)
	return id, err
}

func ObtenerEnlaces() ([]Enlace, error) {
	rows, err := DB.Query("SELECT id, nodo_origen_id, nodo_destino_id FROM enlaces")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lista []Enlace
	for rows.Next() {
		var e Enlace
		rows.Scan(&e.ID, &e.NodoOrigenID, &e.NodoDestinoID)
		lista = append(lista, e)
	}
	return lista, nil
}
