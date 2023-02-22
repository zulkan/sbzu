package domain

import "context"

//go:generate mockery --name QueueUseCase
//go:generate mockery --name QueueProcessor

type QueueUseCase interface {
	ReadMessage(ctx context.Context) (string, string, error) // return key, value, and error
	GetTopicName(ctx context.Context) string
	PublishMessage(ctx context.Context, key, value string) error
}
type QueueProcessor interface {
	ProcessQueueMessage(ctx context.Context, rawData string) error
}
