build:
	protoc -I. --go_out=plugins=micro:. ./proto/user/user.proto
	# dep 工具暂不可用，直接手动编译
	GOOS=linux GOARCH=amd64 go build
	docker build -t user-service .

run:
	docker run --net="host" \
	-p 50053:50051 \
	-e MICRO_SERVER_ADDRESS=:50051 \
	-e MICRO_REGISTRY=mdns \
	-e DB_HOST=172.17.0.6 \
	-e DB_PORT=5432 \
	-e DB_USER=postgres \
	-e DB_PASSWORD=postgres \
	-e DB_NAME=postgres \
	-e MICRO_BROKER=nats \
    -e MICRO_BROKER_ADDRESS=0.0.0.0:4222 \
	user-service
