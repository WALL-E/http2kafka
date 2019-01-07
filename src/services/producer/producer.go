package main

import (
	"fmt"
	"log"
	"services"

	"github.com/Shopify/sarama"
	"github.com/gin-gonic/gin"
)

var cfg *configs.MqConfig
var producer sarama.SyncProducer

func init() {

	fmt.Print("init kafka producer, it may take a few seconds to init the connection\n")

	var err error

	cfg = &configs.MqConfig{}
	configs.LoadJsonConfig(cfg, "kafka.json")

	mqConfig := sarama.NewConfig()
	mqConfig.Net.SASL.Enable = false

	mqConfig.Net.TLS.Enable = false
	mqConfig.Producer.Return.Successes = true

	if err = mqConfig.Validate(); err != nil {
		msg := fmt.Sprintf("Kafka producer config invalidate. config: %v. err: %v", *cfg, err)
		fmt.Println(msg)
		panic(msg)
	}

	producer, err = sarama.NewSyncProducer(cfg.Servers, mqConfig)
	if err != nil {
		msg := fmt.Sprintf("Kafak producer create fail. err: %v", err)
		fmt.Println(msg)
		panic(msg)
	}

}

func produce(topic string, key string, content string) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.StringEncoder(content),
	}

	_, _, err := producer.SendMessage(msg)
	if err != nil {
		msg := fmt.Sprintf("Send Error topic: %v. key: %v. content: %v", topic, key, content)
		fmt.Println(msg)
		return err
	}
	fmt.Printf("Send OK topic:%s key:%s value:%s\n", topic, key, content)

	return nil
}

func main() {
	//the key of the kafka messages
	//do not set the same the key for all messages, it may cause partition im-balance
	//key := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
	//value := "this is a kafka message!"
	//produce(cfg.Topics[0], key, value)

	gin.DisableConsoleColor()

	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  0,
			"message": "ok",
		})
	})

	r.POST("/kafka/:topic/upload", func(c *gin.Context) {
		topic := c.Param("topic")

		file, _ := c.FormFile("file")
		log.Println(file.Filename)

		c.JSON(200, gin.H{
			"status":  0,
			"message": "ok",
			"info":    fmt.Sprintf("topic: %v", topic),
		})
	})

	r.Run(":8080") // listen and serve on 0.0.0.0:8080
}
