package helper

import (
	// "fmt"
	"os"
	"strconv"

	// "github.com/confluentinc/confluent-kafka-go/kafka"
)

var (
	KafkaHost              = os.Getenv("KAFKA_HOST")
	KafkaPort              = os.Getenv("KAFKA_PORT")
	KafkaTopic             = os.Getenv("KAFKA_TOPIC")
	KafkaLogTopic          = os.Getenv("KAFKA_LOG_TOPIC")
	KafkaTopicAuth         = os.Getenv("KAFKA_TOPIC_AUTH")
	KafkaTopicProfile      = os.Getenv("KAFKA_TOPIC_PROFILE")
	KafkaTopicRoleSystem   = os.Getenv("KAFKA_TOPIC_ROLE_SYSTEM")
	KafkaTopicDepartement  = os.Getenv("KAFKA_TOPIC_DEPARTEMENT")
	KafkaTopicSiteLocation = os.Getenv("KAFKA_TOPIC_SITE_LOCATION")
	KafkaConsumerGroup     = os.Getenv("KAFKA_CONSUMER_GROUP")
	KafkaAddressFamily     = os.Getenv("KAFKA_ADDRESS_FAMILY")
	KafkaSessionTimeout, _ = strconv.Atoi(os.Getenv("KAFKA_SESSION_TIMEOUT"))
	KafkaAutoOffsetReset   = os.Getenv("KAFKA_AUTO_OFFSET_RESET")
)

func ProduceKafka(data interface{}, action, kafkaTopik string) {
	// broker := fmt.Sprintf("%s:%s", )

	// p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.server": broker})
}