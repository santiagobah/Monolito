package models

type Nodo_n struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
	Ip   string `json:"ip"`
}

func Crear_nodo(n Nodo_n) (int, error) {
	var Id_v int
	query := "INSERT INTO nodos (nombre, tipo, ip) VALUES ($1, $2, $3) RETURNING id"
	error_db_crear_nodo := DB.QueryRow(query, n.Name, n.Type, n.Ip).Scan(&Id_v)
	return Id_v, error_db_crear_nodo
}

func Get_nodos() ([]Nodo_n, error) {
	rows, error_get_nodos := DB.Query("SELECT id, nombre, tipo, ip FROM nodos ORDER BY id")
	if error_get_nodos != nil {
		return nil, error_get_nodos
	}
	defer rows.Close()

	var lista_nodos []Nodo_n
	for rows.Next() {
		var n_ag Nodo_n
		rows.Scan(&n_ag.Id, &n_ag.Name, &n_ag.Type, &n_ag.Ip)
		lista_nodos = append(lista_nodos, n_ag)
	}
	return lista_nodos, nil
}
