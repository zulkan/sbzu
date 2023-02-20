package domain

//go:generate mockery --name QueueUseCase
//go:generate mockery --name QueueProcessor

type QueueUseCase interface {
	ReadMessage() (string, string, error) // return key, value, and error
	GetTopicName() string
	PublishMessage(key, value string) error
}
type QueueProcessor interface {
	ProcessQueueMessage(rawData string) error
}
