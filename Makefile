TEST?=$$(go list ./... | grep -v 'vendor')
HOSTNAME=github.com
NAMESPACE=aliakseiyanchuk
NAME=mashery-v3-go-client

TRGT_GOOS?=windows
TRGT_GOARCH?=amd64
CONNECT_BIN?=mash-connect.exe
QUERY_BIN?=mash-query.exe

default: build

build:
	GOOS=${TRGT_GOOS} GOARCH=${TRGT_GOARCH} go build -o bin/${CONNECT_BIN} cmd/mash-connect/main.go
	GOOS=${TRGT_GOOS} GOARCH=${TRGT_GOARCH} go build -o bin/${QUERY_BIN} cmd/mash-query/main.go

release:
	GOOS=darwin GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_darwin_amd64
	GOOS=freebsd GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_freebsd_386
	GOOS=freebsd GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_freebsd_amd64
	GOOS=freebsd GOARCH=arm go build -o ./bin/${BINARY}_${VERSION}_freebsd_arm
	GOOS=linux GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_linux_386
	GOOS=linux GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_linux_amd64
	GOOS=linux GOARCH=arm go build -o ./bin/${BINARY}_${VERSION}_linux_arm
	GOOS=openbsd GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_openbsd_386
	GOOS=openbsd GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_openbsd_amd64
	GOOS=solaris GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_solaris_amd64
	GOOS=windows GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_windows_386
	GOOS=windows GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_windows_amd64

test:
	go test ./v3client
#	go test -i $(TEST) || exit 1
#	echo $(TEST) | xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4
