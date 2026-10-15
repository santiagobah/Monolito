package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"envios/controllers"
	"envios/models"
)

type Registration struct {
	ID      string `json:"ID"`
	Name    string `json:"Name"`
	Address string `json:"Address"`
	Port    int    `json:"Port"`
	Check   *Check `json:"Check,omitempty"`
}

type Check struct {
	HTTP                           string `json:"HTTP"`
	Interval                       string `json:"Interval"`
	DeregisterCriticalServiceAfter string `json:"DeregisterCriticalServiceAfter"`
}

func registerWithConsul() {
	consulHost := os.Getenv("CONSUL_HOST")
	if consulHost == "" {
		consulHost = "consul"
	}
	serviceID := os.Getenv("SERVICE_ID")
	if serviceID == "" {
		serviceID = "nodos-1" // (Ajustar según corresponda en enlaces o envios)
	}
	serviceName := os.Getenv("SERVICE_NAME")
	if serviceName == "" {
		serviceName = "nodos" // (Ajustar según corresponda en enlaces o envios)
	}

	reg := Registration{
		ID:      serviceID,
		Name:    serviceName,
		Address: serviceID,
		Port:    8080,
		Check: &Check{
			HTTP:                           fmt.Sprintf("http://%s:8080/health", serviceID),
			Interval:                       "5s",
			DeregisterCriticalServiceAfter: "10s",
		},
	}

	body, _ := json.Marshal(reg)
	url := fmt.Sprintf("http://%s:8500/v1/agent/service/register", consulHost)

	client := &http.Client{}

	for {
		// CAMBIO CLAVE: Usar PUT en lugar de POST para Consul
		req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(body))
		if err == nil {
			req.Header.Set("Content-Type", "application/json")
			resp, err := client.Do(req)
			if err == nil && resp.StatusCode == http.StatusOK {
				resp.Body.Close()
				fmt.Println("Registrado en Consul exitosamente")
				break
			}
			if resp != nil {
				resp.Body.Close()
			}
		}
		fmt.Println("Esperando a que Consul esté listo...")
		time.Sleep(2 * time.Second)
	}
}

func deregisterFromConsul() {
	consulHost := os.Getenv("CONSUL_HOST")
	if consulHost == "" {
		consulHost = "consul"
	}
	serviceID := os.Getenv("SERVICE_ID")
	if serviceID == "" {
		serviceID = "envios-1"
	}
	url := fmt.Sprintf("http://%s:8500/v1/agent/service/deregister/%s", consulHost, serviceID)
	req, _ := http.NewRequest(http.MethodPut, url, nil)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err == nil {
		resp.Body.Close()
		fmt.Println("Desregistrado de Consul")
	}
}

func main() {
	models.Conectar_db()

	go registerWithConsul()

	http.HandleFunc("/envios", controllers.Handler_envios)

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("vivo y coleando"))
	})

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		deregisterFromConsul()
		os.Exit(0)
	}()

	fmt.Println("backend envios prendido en el puerto 8080")
	http.ListenAndServe(":8080", nil)
}
