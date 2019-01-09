#!/bin/bash

# set -x
set -e

PWD=$(cd "$(dirname "$0")"; pwd)

cd $PWD

curl -XPOST http://127.0.0.1:8080/http2kafka/v1/http2mq/upload -H "Content-Type: multipart/form-data" -F "file=@qos.sol"
echo ""

curl -XPOST http://127.0.0.1:8080/http2kafka/v1/http2mq/upload -H "Content-Type: multipart/form-data" -F "file=@qos.sol.gz"
echo ""

curl -XPOST http://127.0.0.1:8080/http2kafka/v1/http2mq/upload -H "Content-Type: multipart/form-data" -F "file=@t.log.gz"
echo ""

curl -XPOST http://127.0.0.1:8080/http2kafka/http2mq/upload -H "Content-Type: multipart/form-data" -F "file=@t.log"
echo ""

curl -XPOST --data '{"name": "zhangzheng"}'  http://127.0.0.1:8080/http2kafka/v1/http2mq/send
echo ""

curl -XPOST --data '{"name": "zhangzheng"}'  http://127.0.0.1:8080/http2kafka/v1/http2mq1/send
echo ""
