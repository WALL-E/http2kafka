#!/bin/bash

# set -x
set -e

PWD=$(cd "$(dirname "$0")"; pwd)

cd $PWD

curl -XPOST http://127.0.0.1:8080/kafka/topic1/upload -H "Content-Type: multipart/form-data" -F "file=@qos.sol"
echo ""

curl -XPOST http://127.0.0.1:8080/kafka/topic2/upload -H "Content-Type: multipart/form-data" -F "file=@qos.sol.gz"
echo ""

curl -XPOST http://127.0.0.1:8080/kafka/topic2/upload -H "Content-Type: multipart/form-data" -F "file=@t.log.gz"
echo ""

curl -XPOST http://127.0.0.1:8080/kafka/topic2/upload -H "Content-Type: multipart/form-data" -F "file=@t.log"
echo ""
