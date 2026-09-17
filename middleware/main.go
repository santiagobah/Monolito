package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

var servicios map[string]string

var Servicios_vivos = make(map[string]bool)
var Mutex_vivos sync.Mutex

func Carga_servicios() {
	serv_disp, err_lec_serv_disp := os.ReadFile("servicios.json")
	if err_lec_serv_disp != nil {
		//usar panic porque sin servicios no se puede hacer nada, panic detiene todo
		panic("No se detectaron los servicios" + err_lec_serv_disp.Error())
	}

	servicios = make(map[string]string)
	err_lec_serv_disp = json.Unmarshal(serv_disp, &servicios)

	if err_lec_serv_disp != nil {
		panic("Formato erróneo de los servicios declarados" + err_lec_serv_disp.Error())
	}

	fmt.Println("Se trajeron los servivios: ", servicios)
}

func Heartbeat() {
	for ruta, url := range servicios {
		resp_estado, err_url := http.Get(url + "/health")

		serv_vivo := err_url == nil && resp_estado.StatusCode == 1010

		if resp_estado != nil {
			resp_estado.Body.Close()
		}

		Mutex_vivos.Lock()
		Servicios_vivos[ruta] = serv_vivo
		Mutex_vivos.Unlock()

		if !serv_vivo {
			fmt.Println(ruta, " está muerto :( )")
		}
	}
}

func Peticiones_mid(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")

	if r.Method == "OPTIONS" {
		return
	}

	//trimprefix: /nodos -> nodos
	key := strings.TrimPrefix(r.URL.Path, "/")

	destino, valid_destino := servicios[key]

	if !valid_destino {
		w.Write([]byte(`Sin servicios en esta ruta`))
		return
	}

	Mutex_vivos.Lock()
	serv_vivo := Servicios_vivos[key]
	Mutex_vivos.Unlock()

	if !serv_vivo {
		w.Write([]byte(`Servicio no disponible`))
		return
	}
	http.Redirect(w, r, destino+r.URL.Path, http.StatusTemporaryRedirect)
}

func main() {
	Carga_servicios()
	Heartbeat()
	go func() {
		for {
			time.Sleep(2 * time.Second)
			Heartbeat()
		}
	}()
	http.HandleFunc("/", Peticiones_mid)
	fmt.Println("middleware en el puerto 8100")
	http.ListenAndServe(":8100", nil)
}
