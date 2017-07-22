# Raspberry Pi Build Light Indicator

## Cross Compiling on Windows

### GNU Toolchain

http://gnutoolchains.com/raspberry/tutorial/

### Go Environment Variables

GOOS=linux
GOARM=7
GOARCH=arm
CC=arm-linux-gnueabihf-gcc
CGO_ENABLED=1