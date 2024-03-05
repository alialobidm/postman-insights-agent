.PHONY: clean build test mock

export GO111MODULE = on

build: clean
	go build -o bin/postman-insights-agent .

clean:
	go clean

mock:
	mkdir -p ./rest/mock
	mockgen -source ./rest/interface.go -destination ./rest/mock/interface.mock.go -package mock

test: mock
	go test ./...
