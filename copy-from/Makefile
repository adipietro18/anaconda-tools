build:
	go build -o ./bin/copy-from

build-linux: build
	GOOS=linux GOARCH=amd64 go build -o bin/linux-64/copy-from

build-all: build build-linux

test:
	go test -coverprofile=coverage.out

coverage: test
	go tool cover -html=coverage.out

clean:
	@rm -f coverage.out
	@rm -rf ./bin

all: build
