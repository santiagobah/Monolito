package models

// aqui esta la estructura de un nodo de la red (router, switch, server, endpoint)
type Nodo struct {
	ID     int    `json:"id"`
	Nombre string `json:"nombre"`
	Tipo   string `json:"tipo"` // Router, Switch, Server o Endpoint
	IP     string `json:"ip"`
}

// esto inserta un nodo nuevo directo en la bd, sin validar mucho la verdad
func CrearNodo(n Nodo) (int, error) {
	var id int
	query := "INSERT INTO nodos (nombre, tipo, ip) VALUES ($1, $2, $3) RETURNING id"
	err := DB.QueryRow(query, n.Nombre, n.Tipo, n.IP).Scan(&id)
	return id, err
}

// trae todos los nodos, sin paginacion ni nada de eso porque es un proyecto escolar
func ObtenerNodos() ([]Nodo, error) {
	rows, err := DB.Query("SELECT id, nombre, tipo, ip FROM nodos ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lista []Nodo
	for rows.Next() {
		var n Nodo
		rows.Scan(&n.ID, &n.Nombre, &n.Tipo, &n.IP)
		lista = append(lista, n)
	}
	return lista, nil
}
