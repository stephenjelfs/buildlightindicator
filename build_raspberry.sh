#!/bin/sh
GOOS=linux GOARCH=arm GOARM=7 CC=arm-linux-gnueabihf-gcc CGO_ENABLED=1 go build -v
