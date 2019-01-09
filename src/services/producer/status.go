package main

const (
	StatusOK = 0

	// 客户端引起的错误
	StatusReadFail     = 401001
	StatusGzipReadFail = 401002

	// 服务端引起的错误
	StatusWriteKafkaFail = 501001
)

var statusText = map[int]string{
	StatusOK: "Ok",

	StatusReadFail:     "Read Fail",
	StatusGzipReadFail: "Gzip Read Fail",

	StatusWriteKafkaFail: "Write Kafka Fail",
}

func StatusText(code int) string {
	return statusText[code]
}
