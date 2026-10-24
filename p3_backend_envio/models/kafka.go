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
	TopicEnviosCreated = "envios.created"
	GroupEnvios        = "envios-service"
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
	if err := EnsureTopic(KafkaBroker, TopicEnviosCreated, 3); err != nil {
		log.Printf("[envios-kafka] warning ensure topic %s: %v", TopicEnviosCreated, err)
	}

	KafkaWriter = &kafka.Writer{
		Addr:         kafka.TCP(KafkaBroker),
		Topic:        TopicEnviosCreated,
		Balancer:     &kafka.Hash{},
		RequiredAcks: kafka.RequireAll,
		BatchTimeout: 10 * time.Millisecond,
	}
}

func PublicarEnvioCreado(p Paquete) {
	if KafkaWriter == nil {
		return
	}
	payload, err := json.Marshal(p)
	if err != nil {
		log.Printf("[envios-kafka] error encoding event: %v", err)
		return
	}

	err = KafkaWriter.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte(strconv.Itoa(p.Id)),
		Value: payload,
		Headers: []kafka.Header{
			{Key: "event-type", Value: []byte("EnvioCreado")},
		},
	})
	if err != nil {
		log.Printf("[envios-kafka] error publishing message: %v", err)
	} else {
		log.Printf("[envios-kafka] evento publicado para envio id=%d (ruta: %s)", p.Id, p.Ruta)
	}
}

func IniciarConsumidorKafka(ctx context.Context) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{KafkaBroker},
		Topic:       TopicEnviosCreated,
		GroupID:     GroupEnvios,
		StartOffset: kafka.FirstOffset,
		MinBytes:    1,
		MaxBytes:    10e6,
	})
	defer r.Close()

	log.Printf("[envios-kafka] consumidor iniciado para el topico %s (grupo: %s)", TopicEnviosCreated, GroupEnvios)
	for {
		msg, err := r.FetchMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return
			}
			log.Printf("[envios-kafka] error fetching message: %v", err)
			select {
			case <-ctx.Done():
				return
			case <-time.After(1 * time.Second):
				continue
			}
		}

		var p Paquete
		if err := json.Unmarshal(msg.Value, &p); err == nil {
			log.Printf("[envios-kafka] evento procesado asincronamente: paquete_id=%d origen=%d destino=%d ruta=%s (partition=%d offset=%d)",
				p.Id, p.Origen, p.Destino, p.Ruta, msg.Partition, msg.Offset)
		}

		if err := r.CommitMessages(ctx, msg); err != nil {
			log.Printf("[envios-kafka] error committing offset: %v", err)
		}
	}
}
