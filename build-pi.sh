#!/bin/bash
ARCH=arm
echo "build for raspberry PI"
env GOARCH=$ARCH GOOS=linux go build  -buildmode exe -o emonP1-pi.$ARCH emonP1.go
file emonP1-pi.$ARCH

scp -P 2222 emonP1-pi.$ARCH root@81.165.73.9:/tmp
scp         emonP1-pi.$ARCH root@192.168.0.251:/tmp

