DOCKER_IMAGE_NAME ?= status-bots

docker-image:
	docker build --file _assets/Dockerfile . -t $(DOCKER_IMAGE_NAME):latest

allbots:
	go build -i -o ./pinger     -v  ./cmd/pinger
	go build -i -o ./chanreader -v ./cmd/chanreader
	go build -i -o ./pingerweb -v ./cmd/pingerweb
	
