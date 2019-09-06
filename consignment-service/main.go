package main

import (
	"context"
	"errors"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/metadata"
	"github.com/micro/go-micro/server"
	"log"

	pb "shippy/consignment-service/proto/consignment"
	vesselProto "shippy/vessel-service/proto/vessel"
	userPb "shippy/user-service/proto/user"
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
		micro.WrapHandler(AuthWrapper),
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

func AuthWrapper(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, resp interface{}) error {
		// consignment-service 独立测试时不进行认证, 直接处理
		if os.Getenv("DISABLE_AUTH") == "true" {
			return fn(ctx, req, resp)
		}

		meta, ok := metadata.FromContext(ctx)
		if !ok {
			return errors.New("no auth meta-data found in request")
		}

		token := meta["Token"]

		authClient := userPb.NewUserServiceClient("go.micro.srv.user", client.DefaultClient)
		authResp, err := authClient.ValidateToken(context.Background(), &userPb.Token{
			Token: token,
		})
		log.Println("Auth Resp:", authResp)
		if err != nil {
			return err
		}
		err = fn(ctx, req, resp)
		return err
	}
}
