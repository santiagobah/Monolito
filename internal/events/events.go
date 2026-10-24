// Package events defines the event contracts shared between services.
package events

import "time"

const (
	TopicOrderCreated   = "orders.created"
	TopicNodosCreated   = "nodos.created"
	TopicEnlacesCreated = "enlaces.created"
	TopicEnviosCreated  = "envios.created"

	GroupInventory    = "inventory-service"
	GroupNotification = "notification-service"
	GroupNodos        = "nodos-service"
	GroupEnlaces      = "enlaces-service"
	GroupEnvios       = "envios-service"
)

type Item struct {
	SKU      string `json:"sku"`
	Quantity int32  `json:"quantity"`
}

type OrderCreated struct {
	EventID       string    `json:"event_id"`
	OrderID       string    `json:"order_id"`
	CustomerID    string    `json:"customer_id"`
	CustomerEmail string    `json:"customer_email"`
	Items         []Item    `json:"items"`
	CreatedAt     time.Time `json:"created_at"`
}

type NodoCreated struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	IP        string    `json:"ip"`
	CreatedAt time.Time `json:"created_at"`
}

type EnlaceCreated struct {
	ID            int       `json:"id"`
	NodoOrigenID  int       `json:"nodo_origen_id"`
	NodoDestinoID int       `json:"nodo_destino_id"`
	CreatedAt     time.Time `json:"created_at"`
}

type EnvioCreated struct {
	ID        int       `json:"id"`
	Origen    int       `json:"origen"`
	Destino   int       `json:"destino"`
	Ruta      string    `json:"ruta"`
	CreatedAt time.Time `json:"created_at"`
}
