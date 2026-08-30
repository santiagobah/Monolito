package models

import (
	"errors"
	"strconv"
)

// esta lista temporal guarda las rutas que se van procesando ahorita mismo
// OJO PROFE: aqui esta el problema de concurrencia a proposito, no tiene mutex
var bufferPaquetes []string
var contadorPaquetes int // cuenta cuantos paquetes van, sin protegerlo (a proposito, para la fase 1)

type LogPaquete struct {
	ID      int    `json:"id"`
	Origen  int    `json:"origen"`
	Destino int    `json:"destino"`
	Ruta    string `json:"ruta"`
}

// busca un camino de un nodo a otro usando los enlaces (bfs bien basico, sin nada elegante)
func buscarRuta(origen, destino int) ([]int, error) {
	enlaces, err := ObtenerEnlaces()
	if err != nil {
		return nil, err
	}

	// armamos un mapa de adyacencia bien facil
	adyacencia := make(map[int][]int)
	for _, e := range enlaces {
		adyacencia[e.NodoOrigenID] = append(adyacencia[e.NodoOrigenID], e.NodoDestinoID)
		adyacencia[e.NodoDestinoID] = append(adyacencia[e.NodoDestinoID], e.NodoOrigenID) // el cable jala pa los dos lados
	}

	visitados := make(map[int]bool)
	padres := make(map[int]int)
	cola := []int{origen}
	visitados[origen] = true
	encontrado := false

	for len(cola) > 0 {
		actual := cola[0]
		cola = cola[1:]

		if actual == destino {
			encontrado = true
			break
		}

		for _, vecino := range adyacencia[actual] {
			if !visitados[vecino] {
				visitados[vecino] = true
				padres[vecino] = actual
				cola = append(cola, vecino)
			}
		}
	}

	if !encontrado {
		return nil, errors.New("no hay ruta entre esos nodos")
	}

	// reconstruimos la ruta de atras pa adelante
	var ruta []int
	nodoActual := destino
	for nodoActual != origen {
		ruta = append([]int{nodoActual}, ruta...)
		nodoActual = padres[nodoActual]
	}
	ruta = append([]int{origen}, ruta...)

	return ruta, nil
}

// esta funcion simula el envio del paquete y guarda el log en la bd
// FASE 1: aqui NO hay mutex, si mandas 2 paquetes bien pegaditos se puede
// pisar el contador y el buffer (eso es lo que hay que evidenciar en esta fase)
func SimularEnvio(origen, destino int) (LogPaquete, error) {
	ruta, err := buscarRuta(origen, destino)
	if err != nil {
		return LogPaquete{}, err
	}

	// convertimos la ruta a texto tipo "1 -> 2 -> 3"
	rutaTexto := ""
	for i, nodoID := range ruta {
		if i > 0 {
			rutaTexto += " -> "
		}
		rutaTexto += strconv.Itoa(nodoID)
	}

	// race condition aqui: dos gorutinas (dos requests al mismo tiempo) pueden
	// leer y escribir estas variables globales sin ningun control
	contadorPaquetes = contadorPaquetes + 1
	bufferPaquetes = append(bufferPaquetes, rutaTexto)

	var id int
	query := "INSERT INTO logs_paquetes (origen, destino, ruta) VALUES ($1, $2, $3) RETURNING id"
	errDB := DB.QueryRow(query, origen, destino, rutaTexto).Scan(&id)
	if errDB != nil {
		return LogPaquete{}, errDB
	}

	resultado := LogPaquete{ID: id, Origen: origen, Destino: destino, Ruta: rutaTexto}
	return resultado, nil
}
