
## 准备客户端

[Kafka 官网 Clients 列表](https://cwiki.apache.org/confluence/display/KAFKA/Clients#Clients-Go(AKAgolang)) 推荐使用 [开源 sarama GO客户端](https://github.com/Shopify/sarama)进行接入。

## 安装环境
1. 请确保安装了 Go 环境，详情请前往 [golang 官网](https://golang.org "golang 官网") 参考；
2. 执行命令：

   ```
   git clone https://github.com/WALL-E/http2kafka
   ```
3. 执行命令：

    ```
    cd http2kafka
    ```

4. 执行命令：

   ```
   export GOPATH=`pwd`
   ```
5. 执行命令安装依赖(请保证联网，需要几分钟，请耐心等待): 

   ```
   go get github.com/Shopify/sarama/ ; go get github.com/bsm/sarama-cluster
   ```

6. 执行命令：

   ```
    go install services
   ```

7. 执行命令：

    ```
    go install services/producer
    ```

8. 执行命令：
    ```
    go install services/consumer
    ```

9. 按照本页下面配置说明配置 conf/kafka.json
10. 生产(没报错说明运行成功):
    ```
    ./bin/producer
    ```

11. 消费(没报错说明运行成功)：

    ```
    ./bin/consumer
    ```


## 准备配置

| Demo 库中配置文件 |配置项| 说明 |
| --- | --- | --- |
| conf/kafka.json | topics | 请修改为控制台上申请的 Topic |
| conf/kafka.json | servers | 请修改为控制台获取的接入点 |
| conf/kafka.json  | consumerGroup | 请修改为控制要申请的 Consumer Group |

## TODO

1. 增加命令行参数，可以指定配置文件
2. 重构代码
3. 增加版本号
4. 从文件内容中判断文件类型，使用http.DetectContentType()
