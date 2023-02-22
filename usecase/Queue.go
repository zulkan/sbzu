package usecase

import (
	"context"
	"gozu/domain"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type kafkaQueue struct {
	KafkaConsumer *kafka.Consumer
	KafkaProducer *kafka.Producer
	TopicName     string
}

func (k *kafkaQueue) ReadMessage(context.Context) (key string, value string, err error) {
	ev, err := k.KafkaConsumer.ReadMessage(100 * time.Millisecond)

	if err != nil {
		return "", "", err
	}

	return string(ev.Key), string(ev.Value), nil
}

func (k *kafkaQueue) GetTopicName(context.Context) string {
	return k.TopicName
}

func (k *kafkaQueue) PublishMessage(_ context.Context, key, value string) error {
	// fmt.Println("PublishMessage ", k.TopicName, key, value)
	return k.KafkaProducer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &k.TopicName, Partition: kafka.PartitionAny},
		Key:            []byte(key),
		Value:          []byte(value),
	}, nil)
}

func NewQueueUseCase(topicName string, kafkaConsumer *kafka.Consumer, kafkaProducer *kafka.Producer) domain.QueueUseCase {
	// To debug producer
	// go func() {
	//	for e := range kafkaProducer.Events() {
	//		if ev, ok := e.(*kafka.Message); ok {
	//			if ev.TopicPartition.Error != nil {
	//				log.Printf("Failed to deliver message: %v\n", ev.TopicPartition)
	//			} else {
	//				log.Printf("Produced event to topic %s: key = %-10s value = %s\n",
	//					*ev.TopicPartition.Topic, string(ev.Key), string(ev.Value))
	//			}
	//		}
	//	}
	// }()

	return &kafkaQueue{KafkaConsumer: kafkaConsumer, KafkaProducer: kafkaProducer, TopicName: topicName}
}
