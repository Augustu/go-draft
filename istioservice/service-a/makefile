all: build clean docker

build:
	CGO_ENABLED=0 GOOS=linux go build -ldflags '-w -s' -a -installsuffix cgo -o main .

clean:
	docker image rm -f service-a:v1

docker:
	docker build -t service-a:v1 .
