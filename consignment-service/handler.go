// 实现微服务的服务端，处理业务逻辑
package main

import (
	"context"
	"gopkg.in/mgo.v2"
	"log"
	pb "shippy/consignment-service/proto/consignment"
	vesselProto "shippy/vessel-service/proto/vessel"
)

// 微服务服务端 struct handler 必须实现 protobuf 中定义的 rpc 方法
// 实现方法的传参等可参考生成的 consignment.pb.go
type handler struct {
	session *mgo.Session
	vesselClient vesselProto.VesselServiceClient
}

// 从主会话中 Close() 出新会话处理查询
func (h *handler) GetRepo() Repository {
	return &ConsignmentRepository{h.session.Clone()}
}

func (h *handler) CreateConsignment (ctx context.Context, req *pb.Consignment, resp *pb.Response) error {
	defer h.GetRepo().Close()

	// 检查是否有适合的货轮
	vReq := &vesselProto.Specification{
		Capacity:             int32(len(req.Containers)),
		MaxWeight:            req.Weight,
	}
	vResp, err := h.vesselClient.FindAvailable(context.Background(), vReq)
	if err != nil {
		return err
	}

	// 货物被承运
	log.Printf("found vessel: %s\n", vResp.Vessel.Name)
	req.VesselId = vResp.Vessel.Id

	err = h.GetRepo().Create(req)
	if err != nil {
		return err
	}
	resp.Created = true
	resp.Consignment = req
	return nil
}

func (h *handler) GetConsignments (ctx context.Context, req *pb.GetRequest, resp *pb.Response) error{
	defer h.GetRepo().Close()
	consignments, err := h.GetRepo().GetAll()
	if err != nil {
		return err
	}
	resp.Consignments = consignments
	return nil
}

