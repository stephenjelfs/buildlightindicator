set GOOS=linux
set GOARCH=arm
set GOARM=7
set CC=arm-linux-gnueabihf-gcc
set CGO_ENABLED=1

go build -v
