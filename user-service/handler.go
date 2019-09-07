package main

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/micro/go-micro/broker"
	"golang.org/x/crypto/bcrypt"
	"log"
	pb "shippy/user-service/proto/user"
)

const topic = "user.created"

type handler struct {
	repo Repository
	tokenService Authable
	PubSub broker.Broker
}

func (h *handler) Create(ctx context.Context, req *pb.User, resp *pb.Response) error {

	// 哈希处理用户输入的密码
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	//log.Printf("%v\n", req)
	req.Password = string(hashedPwd)
	if err := h.repo.Create(req); err != nil {
		return err
	}
	resp.User = req

	// 发布带有用户所有信息的消息
	if err := h.publishEvent(req); err != nil {
		return err
	}
	return nil
}

// 发送消息通知
func (h *handler) publishEvent(user *pb.User) error {
	body, err := json.Marshal(user)
	if err != nil {
		return err
	}

	msg := &broker.Message{
		Header: map[string]string{
			"id": user.Id,
		},
		Body: body,
	}

	// 发布 user.created topic 消息
	if err := h.PubSub.Publish(topic, msg); err != nil {
		log.Fatalf("[pub] failed: %v\n", err)
	}
	return nil
}

func (h *handler) Get(ctx context.Context, req *pb.User, resp *pb.Response) error {
	u, err := h.repo.Get(req.Id)
	if err != nil {
		return err
	}

	resp.User = u
	return nil
}

func (h *handler) GetAll(ctx context.Context, req *pb.Request, resp *pb.Response) error {
	users, err := h.repo.GetAll()
	if err != nil {
		return err
	}
	resp.Users = users
	return nil
}

func (h *handler) Auth(ctx context.Context, req *pb.User, resp *pb.Token) error {
	u, err := h.repo.GetByEmail(req.Email)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)); err != nil {
		return err
	}

	token,err := h.tokenService.Encode(u)
	if err != nil {
		return err
	}
	resp.Token = token
	return nil
}

func (h *handler) ValidateToken(ctx context.Context, req *pb.Token, resp *pb.Token) error {

	claims, err := h.tokenService.Decode(req.Token)
	if err != nil {
		return err
	}

	log.Println(claims)

	if claims.User.Id == "" {
		return errors.New("invalid user")
	}

	resp.Valid = true

	return nil
}