package main

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/go-redis/redis"
	"github.com/joho/godotenv"
	"gozu/domain"
	"gozu/repo"
	"gozu/usecase"
	"gozu/utils"
	"log"
	"os"
	"os/signal"
	"sort"
	"strings"
	"syscall"
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

func ReadConfig() kafka.ConfigMap {
	m := make(map[string]kafka.ConfigValue)

	m["bootstrap.servers"] = strings.TrimSpace(os.Getenv("KAFKA_BOOTSTRAP_SERVERS"))
	m["security.protocol"] = strings.TrimSpace(os.Getenv("KAFKA_SECURITY_PROTOCOL"))
	m["sasl.mechanisms"] = strings.TrimSpace(os.Getenv("KAFKA_SASL_MECHANISM"))
	m["sasl.username"] = strings.TrimSpace(os.Getenv("KAFKA_API_KEY"))
	m["sasl.password"] = strings.TrimSpace(os.Getenv("KAFKA_API_SECRET"))

	return m
}

func getProducer() (*kafka.Producer, error) {
	kafkaConfig := ReadConfig()

	fmt.Println(kafkaConfig)
	kafkaProducer, err := kafka.NewProducer(&kafkaConfig)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return kafkaProducer, nil
}

func getConsumer(topic string) (*kafka.Consumer, error) {
	kafkaConfig := ReadConfig()
	kafkaConfig["group.id"] = "zulkan"
	//kafkaConfig["group.id"] = "lkc-r56x0k"
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
func readKafka(consumer domain.QueueUseCase, processor domain.QueueProcessor) {
	// Process messages
	go func() {
		for IsRunning {
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

func readFiles(useCase domain.StockUseCase) {
	entries, err := os.ReadDir("./subsetdata")
	if err != nil {
		log.Fatal(err)
	}
	files := []string{}
	for _, e := range entries {
		files = append(files, e.Name())
	}
	sort.Strings(files)

	for _, e := range files {
		content, _ := os.ReadFile("./subsetdata/" + e) //nolint:gosec
		lineData := strings.Split(string(content), "\n")
		for _, data := range lineData {
			err = processRawDataToKafka(data, useCase)

			if err != nil {
				return
			}
		}
	}
}

func processRawDataToKafka(data string, useCase domain.StockUseCase) error {
	dataMap := utils.ReadJSON[map[string]interface{}](data)
	dataType := (*dataMap)["type"]
	record := domain.StockRecord{}
	record.Type = dataType.(string)

	if dataType == "A" {
		if (*dataMap)["quantity"] == "0" {
			record.Quantity = 0
			record.Price = utils.ToInt64((*dataMap)["price"].(string))
			record.StockCode = (*dataMap)["stock_code"].(string)
		}
	}
	if dataType == "E" || dataType == "P" {
		record.Quantity = utils.ToInt64((*dataMap)["executed_quantity"].(string))
		record.Price = utils.ToInt64((*dataMap)["execution_price"].(string))
		record.StockCode = (*dataMap)["stock_code"].(string)
	}

	if record.Price != 0 {
		err := useCase.ProcessFileData(utils.ToJSON(record))
		if err != nil {
			//log.Fatal(err)
			log.Println(err)
			return err
		}
	}
	return nil
}

func main() {
	topicName := "zul-test"
	redisClient := initRedis()
	stockRepo := repo.NewStockRepo(redisClient)

	var err error
	var kafkaProducer *kafka.Producer
	var kafkaConsumer *kafka.Consumer
	kafkaProducer, err = getProducer()
	if err != nil {
		log.Fatal(err)
	}
	kafkaConsumer, err = getConsumer(topicName)
	if err != nil {
		log.Fatal(err)
	}
	queue := usecase.NewQueueUseCase(topicName, kafkaConsumer, kafkaProducer)

	stockUseCase := usecase.NewStockUseCase(queue, stockRepo)

	readFiles(stockUseCase)
	readKafka(queue, stockUseCase)

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
