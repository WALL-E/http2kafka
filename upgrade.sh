#!/bin/bash

file bin/linux_amd64/producer | grep ELF
if test $? -ne 0 ;then
    exit 1
else
    echo "Deploy to the target host"
fi


# kong pro
echo ""
echo "Pro:"
echo "scp -P 4239 bin/linux_amd64/producer gmu@172.16.213.183:/apps/http2kafka"
echo "scp -P 4239 bin/linux_amd64/producer gmu@172.16.213.190:/apps/http2kafka"
