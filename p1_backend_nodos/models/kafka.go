package models

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/segmentio/kafka-go"
)

var (
	KafkaWriter *kafka.Writer
	KafkaBroker string
)

const (
	TopicNodosCreated = "nodos.created"
	GroupNodos        = "nodos-service"
)

func getenvKafka(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func EnsureTopic(broker, topic string, partitions int) error {
	conn, err := kafka.Dial("tcp", broker)
	if err != nil {
		return err
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		return err
	}
	ctrl, err := kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		return err
	}
	defer ctrl.Close()

	return ctrl.CreateTopics(kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     partitions,
		ReplicationFactor: 1,
	})
}

func InitKafka() {
	KafkaBroker = getenvKafka("KAFKA_BROKER", "localhost:9092")
	if err := EnsureTopic(KafkaBroker, TopicNodosCreated, 3); err != nil {
		log.Printf("[nodos-kafka] warning ensure topic %s: %v", TopicNodosCreated, err)
	}

	KafkaWriter = &kafka.Writer{
		Addr:         kafka.TCP(KafkaBroker),
		Topic:        TopicNodosCreated,
		Balancer:     &kafka.Hash{},
		RequiredAcks: kafka.RequireAll,
		BatchTimeout: 10 * time.Millisecond,
	}
}

func PublicarNodoCreado(n Nodo_n) {
	if KafkaWriter == nil {
		return
	}
	payload, err := json.Marshal(n)
	if err != nil {
		log.Printf("[nodos-kafka] error encoding event: %v", err)
		return
	}

	err = KafkaWriter.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte(strconv.Itoa(n.Id)),
		Value: payload,
		Headers: []kafka.Header{
			{Key: "event-type", Value: []byte("NodoCreado")},
		},
	})
	if err != nil {
		log.Printf("[nodos-kafka] error publishing message: %v", err)
	} else {
		log.Printf("[nodos-kafka] evento publicado para nodo id=%d (%s)", n.Id, n.Name)
	}
}

func IniciarConsumidorKafka(ctx context.Context) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{KafkaBroker},
		Topic:       TopicNodosCreated,
		GroupID:     GroupNodos,
		StartOffset: kafka.FirstOffset,
		MinBytes:    1,
		MaxBytes:    10e6,
	})
	defer r.Close()

	log.Printf("[nodos-kafka] consumidor iniciado para el topico %s (grupo: %s)", TopicNodosCreated, GroupNodos)
	for {
		msg, err := r.FetchMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return
			}
			log.Printf("[nodos-kafka] error fetching message: %v", err)
			select {
			case <-ctx.Done():
				return
			case <-time.After(1 * time.Second):
				continue
			}
		}

		var n Nodo_n
		if err := json.Unmarshal(msg.Value, &n); err == nil {
			log.Printf("[nodos-kafka] evento procesado asincronamente: nodo_id=%d name=%s type=%s ip=%s (partition=%d offset=%d)",
				n.Id, n.Name, n.Type, n.Ip, msg.Partition, msg.Offset)
		}

		if err := r.CommitMessages(ctx, msg); err != nil {
			log.Printf("[nodos-kafka] error committing offset: %v", err)
		}
	}
}
