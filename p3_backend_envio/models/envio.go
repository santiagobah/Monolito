package models

import (
	"errors"
	"strconv"
	"sync"
)

var mutex sync.Mutex

var mem_paquetes []string //es un buffer para los pquetes que se envíen a través de la red

var cont_paquetes int

type Paquete struct {
	Id      int    `json:"id"`
	Origen  int    `json:"origen"`
	Destino int    `json:"destino"`
	Ruta    string `json:"ruta"`
}

type Enlace_paquetes struct {
	Origen  int `json:"origen"`
	Destino int `json:"destino"`
}

func get_enlaces_env() ([]Enlace_paquetes, error) {
	rows, err_enl_env := DB.Query("SELECT nodo_origen_id, nodo_destino_id FROM enlaces")
	if err_enl_env != nil {
		return nil, err_enl_env
	}
	defer rows.Close()
	var lista_enlaces []Enlace_paquetes
	for rows.Next() {
		var e_ag Enlace_paquetes
		rows.Scan(&e_ag.Origen, &e_ag.Destino)
		lista_enlaces = append(lista_enlaces, e_ag)
	}
	return lista_enlaces, nil
}

func Buscar_ruta_nodos(n_origen, n_destino int) ([]int, error) {
	enlaces_nodos, err_get_enl := get_enlaces_env()
	if err_get_enl != nil {
		return nil, err_get_enl
	}
	nod_ady := make(map[int][]int) //nodos adyacentes

	for _, e_ag := range enlaces_nodos {
		nod_ady[e_ag.Origen] = append(nod_ady[e_ag.Origen], e_ag.Destino)
		nod_ady[e_ag.Destino] = append(nod_ady[e_ag.Destino], e_ag.Origen)
	}

	nodos_visitados := make(map[int]bool)
	nodos_padres := make(map[int]int)
	cola_nodos := []int{n_origen}
	nodos_visitados[n_origen] = true
	nodo_encontrado := false

	for len(cola_nodos) > 0 {
		nodo_actual := cola_nodos[0]
		cola_nodos = cola_nodos[1:]
		if nodo_actual == n_destino {
			nodo_encontrado = true
			break
		}

		for _, n_vecino := range nod_ady[nodo_actual] {
			if !nodos_visitados[n_vecino] {
				nodos_visitados[n_vecino] = true
				nodos_padres[n_vecino] = nodo_actual
				cola_nodos = append(cola_nodos, n_vecino)
			}
		}

	}

	if !nodo_encontrado {
		return nil, errors.New("No existe una ruta definida entre estos nodos.")
	}

	var ruta_n []int

	nodo_actual := n_destino

	for nodo_actual != n_origen {
		ruta_n = append([]int{nodo_actual}, ruta_n...)
		nodo_actual = nodos_padres[nodo_actual]
	}

	ruta_n = append([]int{n_origen}, ruta_n...)

	return ruta_n, nil
}

func Envio_paquetes(n_origen, n_destino int) (Paquete, error) {
	ruta_nodos, err_env_paq := Buscar_ruta_nodos(n_origen, n_destino)

	if err_env_paq != nil {
		return Paquete{}, err_env_paq
	}

	ruta_string := ""

	for i, nodo_a_ID := range ruta_nodos {
		if i > 0 {
			ruta_string += "->"
		}
		ruta_string += strconv.Itoa(nodo_a_ID)
	}

	mutex.Lock()

	cont_paquetes += 1
	mem_paquetes = append(mem_paquetes, ruta_string)

	var id_e int
	query := "INSERT INTO logs_paquetes (origen, destino, ruta) VALUES ($1, $2, $3) RETURNING id"
	err_db := DB.QueryRow(query, n_origen, n_destino, ruta_string).Scan(&id_e)

	mutex.Unlock()

	if err_db != nil {
		return Paquete{}, err_db
	}

	paq_final := Paquete{Id: id_e, Origen: n_origen, Destino: n_destino, Ruta: ruta_string}

	return paq_final, nil

}
