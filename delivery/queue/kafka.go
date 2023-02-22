package queue

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"gozu/constant"
	"gozu/domain"
	"log"
	"os"
	"strings"
)

func ReadAndProcessKafka(isRunning *bool, consumer domain.QueueUseCase, processor domain.QueueProcessor) {
	// Process messages
	go func() {
		for *isRunning {
			_, value, err := consumer.ReadMessage()
			if err != nil {
				// Errors are informational and automatically handled by the consumer
				continue
			}
			// fmt.Printf("Consumed event from topic %s: key = %-10s value = %s\n",
			// consumer.GetTopicName(), key, value)
			err = processor.ProcessQueueMessage(value)

			if err != nil {
				log.Println(err)
			}
		}
	}()
}

func GetConsumer(topic string) (*kafka.Consumer, error) {
	kafkaConfig := ReadConfig()
	kafkaConfig["group.id"] = constant.KAFKA_GROUP_ID
	kafkaConfig["auto.offset.reset"] = "earliest"
	kafkaConfig["session.timeout.ms"] = os.Getenv("KAFKA_SESSION_TIMEOUT")

	kafkaConsumer, err := kafka.NewConsumer(&kafkaConfig)
	fmt.Println(kafkaConfig)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	err = kafkaConsumer.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return kafkaConsumer, nil
}

func ReadConfig() kafka.ConfigMap {
	m := make(map[string]kafka.ConfigValue)

	m["bootstrap.servers"] = strings.TrimSpace(os.Getenv("KAFKA_BOOTSTRAP_SERVERS"))
	m["security.protocol"] = strings.TrimSpace(os.Getenv("KAFKA_SECURITY_PROTOCOL"))
	m["sasl.mechanisms"] = strings.TrimSpace(os.Getenv("KAFKA_SASL_MECHANISM"))
	m["sasl.username"] = strings.TrimSpace(os.Getenv("KAFKA_API_KEY"))
	m["sasl.password"] = strings.TrimSpace(os.Getenv("KAFKA_API_SECRET"))

	return m
}

func GetProducer() (*kafka.Producer, error) {
	kafkaConfig := ReadConfig()

	fmt.Println(kafkaConfig)
	kafkaProducer, err := kafka.NewProducer(&kafkaConfig)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return kafkaProducer, nil
}
