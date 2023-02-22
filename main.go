package main

import (
	"context"
	"fmt"
	"gozu/constant"
	"gozu/delivery/grpc"
	gozuQueue "gozu/delivery/queue"
	"gozu/domain"
	"gozu/repo"
	"gozu/usecase"
	"log"
	"os"
	"os/signal"
	"sort"
	"strings"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/go-redis/redis"
	"github.com/joho/godotenv"
)

var (
	RedisHost     string
	RedisPort     string
	RedisPassword string

	IsRunning = true
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	RedisHost = os.Getenv("REDIS_HOST")
	RedisPort = os.Getenv("REDIS_PORT")
	RedisPassword = os.Getenv("REDIS_PASSWORD")
}
func initRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     RedisHost + ":" + RedisPort,
		Password: RedisPassword,
		DB:       0,
	})

	_, err := client.Ping().Result()
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func readFiles(useCase domain.StockUseCase) {
	entries, err := os.ReadDir("./subsetdata")
	if err != nil {
		log.Fatal(err)
	}
	var files []string
	for _, e := range entries {
		files = append(files, e.Name())
	}
	sort.Strings(files)

	for _, e := range files {
		content, _ := os.ReadFile("./subsetdata/" + e) //nolint:gosec //because just looping our own files
		lineData := strings.Split(string(content), "\n")
		for _, data := range lineData {
			err = useCase.ProcessFileData(context.Background(), data)

			if err != nil {
				return
			}
		}
	}
}

func main() {
	redisClient := initRedis()
	stockRepo := repo.NewStockRepo(redisClient)

	var err error
	var kafkaProducer *kafka.Producer
	var kafkaConsumer *kafka.Consumer
	kafkaProducer, err = gozuQueue.GetProducer()
	if err != nil {
		log.Fatal(err)
	}
	kafkaConsumer, err = gozuQueue.GetConsumer(constant.KafkaStockTopic)
	if err != nil {
		log.Fatal(err)
	}
	queue := usecase.NewQueueUseCase(constant.KafkaStockTopic, kafkaConsumer, kafkaProducer)

	stockUseCase := usecase.NewStockUseCase(queue, stockRepo)

	go func() {
		readFiles(stockUseCase)
	}()
	gozuQueue.ReadAndProcessKafka(&IsRunning, queue, stockUseCase)

	// init grpc server
	go func() {
		grpc.NewGrpcServer(stockUseCase)
	}()

	// Set up a channel for handling Ctrl-C, etc
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	for IsRunning {
		if sig := <-sigchan; true {
			fmt.Printf("Caught signal %v: terminating\n", sig)
			IsRunning = false
		}
	}
}
