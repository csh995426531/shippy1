package main

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/micro/go-micro"
	"io/ioutil"
	"log"
	"os"
	pb "shippy/consignment-service/proto/consignment"
)

const (
	ADDRESS				= "localhost:50051"
	DEFAULT_INFO_FILE	= "consignment.json"
)

// 读取 consignment.json 中记录的货物信息
func parseFile(fileName string) (*pb.Consignment, error) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	var consignment *pb.Consignment
	err = json.Unmarshal(data, &consignment)
	if err != nil {
		return nil, errors.New("consignment.json file content error")
	}
	return consignment, nil
}

func main() {
	service := micro.NewService(micro.Name("go.micro.srv.consignment.cli"))
	service.Init()
	// 创建微服务的客户端，简化了手动 Dial 链接服务端的步骤
	client := pb.NewShippingServiceClient("go.micro.srv.consignment", service.Client())

	//// 链接到 gRPC 服务器
	//conn,err := grpc.Dial(ADDRESS, grpc.WithInsecure())
	//if err != nil {
	//	log.Fatalf("connect error : %v\n", err)
	//}
	//defer conn.Close()
	//
	//// 初始化 gRPC 客户端
	//client := pb.NewShippingServiceClient(conn)

	// 在命令行中指定新的货物信息 json 文件
	infoFile := DEFAULT_INFO_FILE
	if len(os.Args) > 1 {
		infoFile = os.Args[1]
	}

	// 解析货物信息
	consignment, err := parseFile(infoFile)
	if err != nil {
		log.Fatalf("parse info file error:%v\n", err)
	}

	// 调用 RPC
	// 将货物存储到我们自己的仓库里
	resp, err := client.CreateConsignment(context.TODO(), consignment)
	if err != nil {
		log.Fatalf("create consignment error:%v\n", err)
	}

	// 新货物是否托运成功
	log.Printf("created:%t\n", resp.Created)
	// 列出目前所有托运的货物
	resp, err = client.GetConsignments(context.Background(), &pb.GetRequest{})
	if err != nil {
		log.Fatalf("failed to list consignments:%v\n", err)
	}

	for _, c := range resp.Consignments {
		log.Printf("%+v", c)
	}
}