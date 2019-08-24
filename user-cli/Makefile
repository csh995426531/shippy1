build:
	# dep 工具暂不可用，直接手动编译
	GOOS=linux GOARCH=amd64 go build
	docker build -t user-cli .

run:
	docker run -e MICRO_REGISTRY=mdns user-cli
