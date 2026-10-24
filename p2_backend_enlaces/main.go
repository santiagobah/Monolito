package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"enlaces/controllers"
	"enlaces/models"
)

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func registerWithConsul(consulAddr string) error {
	instanceID := getEnv("INSTANCE_ID", "enlaces")
	payload := map[string]any{
		"ID":      instanceID,
		"Name":    "enlaces",
		"Address": instanceID,
		"Port":    8080,
		"Check": map[string]any{
			"HTTP":                           fmt.Sprintf("http://%s:8080/health", instanceID),
			"Interval":                       "5s",
			"Timeout":                        "2s",
			"DeregisterCriticalServiceAfter": "1m",
		},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	var lastErr error
	for attempt := 1; attempt <= 10; attempt++ {
		req, err := http.NewRequest(http.MethodPut, consulAddr+"/v1/agent/service/register", bytes.NewReader(body))
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				log.Printf("registered %s with Consul at %s", instanceID, consulAddr)
				return nil
			}
			lastErr = fmt.Errorf("consul returned status %d", resp.StatusCode)
		} else {
			lastErr = err
		}
		time.Sleep(time.Duration(attempt) * time.Second)
	}
	return fmt.Errorf("giving up registering with consul: %w", lastErr)
}

func deregisterFromConsul(consulAddr string) {
	instanceID := getEnv("INSTANCE_ID", "enlaces")
	req, err := http.NewRequest(http.MethodPut, consulAddr+"/v1/agent/service/deregister/"+instanceID, nil)
	if err != nil {
		log.Printf("failed to build deregister request: %v", err)
		return
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("failed to deregister %s from consul: %v", instanceID, err)
		return
	}
	resp.Body.Close()
	log.Printf("deregistered %s from consul", instanceID)
}

func main() {
	models.Conectar_db()
	models.InitKafka()

	consulAddr := getEnv("CONSUL_HTTP_ADDR", "http://consul:8500")
	if err := registerWithConsul(consulAddr); err != nil {
		log.Printf("WARNING: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/enlaces", controllers.Handler_enlaces)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("vivo y coleando"))
	})

	srv := &http.Server{Addr: ":8080", Handler: mux}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go models.IniciarConsumidorKafka(ctx)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		instanceID := getEnv("INSTANCE_ID", "enlaces")
		log.Printf("shutting down %s", instanceID)
		deregisterFromConsul(consulAddr)
		cancel()
		if models.KafkaWriter != nil {
			models.KafkaWriter.Close()
		}
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()
		srv.Shutdown(shutdownCtx)
	}()

	fmt.Println("backend enlaces prendido en el puerto 8080")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
