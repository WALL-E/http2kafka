#!/bin/bash

set -x
set -e

PWD=$(cd "$(dirname "$0")"; pwd)

cd $PWD
curl -XPOST http://127.0.0.1:8080/kafka/topoc4/upload -H "Content-Type: multipart/form-data" -F "file=@qos.sol"
