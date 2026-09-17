package models

type Enlace struct {
	Id              int `json:"id"`
	Nodo_origen_id  int `json:"nodo_origen_id"`
	Nodo_destino_id int `json:"nodo_destino_id"`
}

func Crear_enlace(e Enlace) (int, error) {
	var id_v int
	query := "INSERT INTO enlaces (nodo_origen_id, nodo_destino_id) VALUES ($1, $2) RETURNING id"
	err_crecion_enlace := DB.QueryRow(query, e.Nodo_origen_id, e.Nodo_destino_id).Scan(&id_v)
	return id_v, err_crecion_enlace
}

func Get_enlaces() ([]Enlace, error) {
	rows, err_get_enlaces := DB.Query("SELECT id, nodo_origen_id, nodo_destino_id FROM enlaces")
	if err_get_enlaces != nil {
		return nil, err_get_enlaces
	}

	var lista_enlaces []Enlace
	for rows.Next() {
		var e_n Enlace
		rows.Scan(&e_n.Id, &e_n.Nodo_origen_id, &e_n.Nodo_destino_id)
		lista_enlaces = append(lista_enlaces, e_n)
	}
	return lista_enlaces, nil
}
