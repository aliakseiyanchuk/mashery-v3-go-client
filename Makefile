TEST?=$$(go list ./... | grep -v 'vendor')
HOSTNAME=github.com
NAMESPACE=aliakseiyanchuk
NAME=mashery-v3-go-client

TRGT_GOOS?=windows
TRGT_GOARCH?=amd64
BINARY=mash-query
VERSION=0.9

default: build

build_linux:
	GOOS=linux GOARCH=arm64 		go build -o ../bin/${BINARY}_${VERSION}_linux_arm			./cmd/mash-query
	GOOS=linux GOARCH=arm GOARM=6 	go build -o ../bin/${BINARY}_${VERSION}_linux_armv6 		./cmd/mash-query
	GOOS=linux GOARCH=amd64 		go build -o ../bin/${BINARY}_${VERSION}_linux_amd64			./cmd/mash-query
	GOOS=linux GOARCH=386 			go build -o ../bin/${BINARY}_${VERSION}_linux_386			./cmd/mash-query
build_mac:
	GOOS=darwin GOARCH=arm64 go build -o ./bin/${BINARY}_${VERSION}_darwin_arm64 		./cmd/mash-query

build: build_mac build_linux
	echo "All done"

test:
	go test ./v3client
