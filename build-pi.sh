#!/bin/bash
ARCH=arm
echo "build for raspberry PI"
env GOARCH=$ARCH GOOS=linux go build  -buildmode exe -o emonP1-pi.$ARCH emonP1.go
