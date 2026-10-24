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
	TopicEnlacesCreated = "enlaces.created"
	GroupEnlaces        = "enlaces-service"
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
	if err := EnsureTopic(KafkaBroker, TopicEnlacesCreated, 3); err != nil {
		log.Printf("[enlaces-kafka] warning ensure topic %s: %v", TopicEnlacesCreated, err)
	}

	KafkaWriter = &kafka.Writer{
		Addr:         kafka.TCP(KafkaBroker),
		Topic:        TopicEnlacesCreated,
		Balancer:     &kafka.Hash{},
		RequiredAcks: kafka.RequireAll,
		BatchTimeout: 10 * time.Millisecond,
	}
}

func PublicarEnlaceCreado(e Enlace) {
	if KafkaWriter == nil {
		return
	}
	payload, err := json.Marshal(e)
	if err != nil {
		log.Printf("[enlaces-kafka] error encoding event: %v", err)
		return
	}

	err = KafkaWriter.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte(strconv.Itoa(e.Id)),
		Value: payload,
		Headers: []kafka.Header{
			{Key: "event-type", Value: []byte("EnlaceCreado")},
		},
	})
	if err != nil {
		log.Printf("[enlaces-kafka] error publishing message: %v", err)
	} else {
		log.Printf("[enlaces-kafka] evento publicado para enlace id=%d (%d -> %d)", e.Id, e.Nodo_origen_id, e.Nodo_destino_id)
	}
}

func IniciarConsumidorKafka(ctx context.Context) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{KafkaBroker},
		Topic:       TopicEnlacesCreated,
		GroupID:     GroupEnlaces,
		StartOffset: kafka.FirstOffset,
		MinBytes:    1,
		MaxBytes:    10e6,
	})
	defer r.Close()

	log.Printf("[enlaces-kafka] consumidor iniciado para el topico %s (grupo: %s)", TopicEnlacesCreated, GroupEnlaces)
	for {
		msg, err := r.FetchMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return
			}
			log.Printf("[enlaces-kafka] error fetching message: %v", err)
			select {
			case <-ctx.Done():
				return
			case <-time.After(1 * time.Second):
				continue
			}
		}

		var e Enlace
		if err := json.Unmarshal(msg.Value, &e); err == nil {
			log.Printf("[enlaces-kafka] evento procesado asincronamente: enlace_id=%d origen=%d destino=%d (partition=%d offset=%d)",
				e.Id, e.Nodo_origen_id, e.Nodo_destino_id, msg.Partition, msg.Offset)
		}

		if err := r.CommitMessages(ctx, msg); err != nil {
			log.Printf("[enlaces-kafka] error committing offset: %v", err)
		}
	}
}
