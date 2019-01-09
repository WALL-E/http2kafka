package main

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net/http"
	"services"
	"strconv"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	"github.com/gin-gonic/gin"
)

var cfg *configs.MqConfig
var producer sarama.SyncProducer
var nonce int64

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
	gin.DisableConsoleColor()

	r := gin.Default()
	r.MaxMultipartMemory = 1 << 20 // 1 MiB

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  0,
			"message": "ok",
		})
	})

	// 只接受压缩文件，解压缩后的文件采用行格式
	// 详细：http://jira.gmugmu.com:8090/pages/viewpage.action?pageId=3047483
	r.POST("/http2kafka/v1/:topic/upload", func(c *gin.Context) {
		topic := c.Param("topic")

		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  StatusReadFail,
				"message": StatusText(StatusReadFail),
				"info":    fmt.Sprintf("file: %v", file),
			})

			return
		}
		log.Println(file.Filename)

		f, err := file.Open()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  StatusReadFail,
				"message": StatusText(StatusReadFail),
				"info":    fmt.Sprintf("filename: %v", file.Filename),
			})

			return
		}
		defer f.Close()

		// 创建gzip文件读取对象
		gr, err := gzip.NewReader(f)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  StatusGzipReadFail,
				"message": StatusText(StatusGzipReadFail),
				"info":    fmt.Sprintf("filename: %v", file.Filename),
			})

			return
		}
		defer gr.Close()

		// 按行读取解压缩后的文件
		rd := bufio.NewReader(gr)
		for {
			line, err := rd.ReadString('\n')

			if err != nil || io.EOF == err {
				break
			}
			//fmt.Println(line)

			line = strings.Replace(line, " ", "", -1)
			line = strings.Replace(line, "\n", "", -1)

			if len(line) == 0 {
				continue
			}

			if nonce == 0 {
				nonce = time.Now().UTC().UnixNano()
			}
			nonce++

			key := strconv.FormatInt(nonce, 10)
			value := line
			err = produce(cfg.Topics[0], key, value)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"status":  StatusWriteKafkaFail,
					"message": StatusText(StatusWriteKafkaFail),
					"info":    fmt.Sprintf("line: %v", line),
				})

				return
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  0,
			"message": "ok",
			"info":    fmt.Sprintf("topic: %v", topic),
		})
	})

	// 只接收文本数据
	r.POST("/http2kafka/v1/:topic/send", func(c *gin.Context) {
		topic := c.Param("topic")

		data, err := c.GetRawData()

		if nonce == 0 {
			nonce = time.Now().UTC().UnixNano()
		}
		nonce++

		key := strconv.FormatInt(nonce, 10)
		value := string(data[:])
		err = produce(cfg.Topics[0], key, value)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  StatusWriteKafkaFail,
				"message": StatusText(StatusWriteKafkaFail),
				"info":    fmt.Sprintf("data: %v", data),
			})

			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  0,
			"message": "ok",
			"info":    fmt.Sprintf("topic: %v", topic),
		})
	})

	r.Run(":8080") // listen and serve on 0.0.0.0:8080
}
