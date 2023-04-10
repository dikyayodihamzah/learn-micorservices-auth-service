package helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	_ "github.com/joho/godotenv/autoload"
)

var (
	KafkaHost              = os.Getenv("KAFKA_HOST")
	KafkaPort              = os.Getenv("KAFKA_PORT")
	KafkaTopic             = os.Getenv("KAFKA_TOPIC")
	KafkaTopicLog          = os.Getenv("KAFKA_TOPIC_LOG")
	KafkaTopicAuth         = os.Getenv("KAFKA_TOPIC_AUTH")
	KafkaTopicProfile      = os.Getenv("KAFKA_TOPIC_PROFILE")
	KafkaTopicRole         = os.Getenv("KAFKA_TOPIC_ROLE")
	KafkaConsumerGroup     = os.Getenv("KAFKA_CONSUMER_GROUP")
	KafkaAddressFamily     = os.Getenv("KAFKA_ADDRESS_FAMILY")
	KafkaSessionTimeout, _ = strconv.Atoi(os.Getenv("KAFKA_SESSION_TIMEOUT"))
	KafkaAutoOffsetReset   = os.Getenv("KAFKA_AUTO_OFFSET_RESET")
)

func NewKafkaProducer() *kafka.Producer {
	broker := fmt.Sprintf("%s:%s", KafkaHost, KafkaPort)

	producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": broker})
	if err != nil {
		log.Fatalf("Error on creating kafka producer: %s\n", err.Error())
	}

	return producer
}

func NewKafkaConsumer() *kafka.Consumer {
	broker := fmt.Sprintf("%s:%s", KafkaHost, KafkaPort)

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":     broker,
		"broker.address.family": KafkaAddressFamily,
		"group.id":              KafkaConsumerGroup,
		"session.timeout.ms":    KafkaSessionTimeout,
		"auto.offset.reset":     KafkaAutoOffsetReset,
	})
	if err != nil {
		log.Fatalf("Error on creating kafka consumer: %s\n", err.Error())
	}

	log.Println("Kafka consumer created")

	return consumer
}

func ProduceToKafka(data interface{}, action, kafkaTopi string) {
	broker := fmt.Sprintf("%s:%s", KafkaHost, KafkaPort)

	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.server": broker})

	if err != nil {
		log.Printf("KAFKA | Failed to create producer: %s\n", err)
	}

	value := new(bytes.Buffer)
	json.NewEncoder(value).Encode(data)

	p.ProduceChannel() <- &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &KafkaTopic, Partition: kafka.PartitionAny},
		Value:          value.Bytes(),
		Headers:        []kafka.Header{{Key: "method", Value: []byte(action)}},
	}

	for e := range p.Events() {
		switch ev := e.(type) {
		case *kafka.Message:
			m := ev
			if m.TopicPartition.Error != nil {
				log.Printf("KAFKA | Delivery failed: %v\n", m.TopicPartition.Error)
				p.Close()
			} else {
				log.Printf("KAFKA | Delivered message to topic %s [%d] at offset %v\n",
					*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
				p.Close()
			}
			return

		default:
			if strings.Contains(ev.String(), "failed:") {
				p.Close()
				PanicIfError("message broker error")
			} else {
				log.Printf("KAFKA | Ignored event: %s\n", ev)
			}
		}
	}

	p.Close()
}
