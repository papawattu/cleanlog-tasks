.PHONY: all build clean client run run-client watch docker-build

all: build client

build:
	@go build -o bin/tasks ./*.go

client:
	@go build -o bin/client client/main.go

run: build
	@./bin/tasks

run-client: client
	@./bin/client

clean:
	@rm -rf bin

watch:
	@air -c .air.toml

docker-build:
	@docker buildx build --platform linux/amd64,linux/arm64 -t ghcr.io/papawattu/cleanlog-tasks:latest .

docker-push: docker-build
	@docker push ghcr.io/papawattu/cleanlog-tasks:latest