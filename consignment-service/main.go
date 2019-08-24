package main

import (
	"log"

	pb "shippy/consignment-service/proto/consignment"
	vesselProto "shippy/vessel-service/proto/vessel"
	"github.com/micro/go-micro"
	"os"
)
const (
	DEFAULT_HOST = "datastore:27017"
)

func main() {

	// 获取容器设置的数据库地址环境变量的值
	host := os.Getenv("DB_HOST")

	if host == "" {
		host = DEFAULT_HOST
	}
	session, err := CreateSession(host)
	// 创建于 MongoDb 的主会话，需在退出 main() 时候手动释放链接
	defer session.Close()
	if err != nil {
		log.Fatalf("create session error:%v : %v\n", host, err)
	}

	server := micro.NewService(
		// 必须和 consignment.proto 中的 package 一致
		micro.Name("go.micro.srv.consignment"),
		micro.Version("latest"),
	)

	// 解析命令行参数
	server.Init()
	// 作为 vessel-service 的客户端
	vClient := vesselProto.NewVesselServiceClient("go.micro.srv.vessel", server.Client())
	// 将 server 作为微服务的服务端
	pb.RegisterShippingServiceHandler(server.Server(), &handler{session, vClient})

	if err := server.Run(); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
