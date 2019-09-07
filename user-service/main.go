package main

import (
	"fmt"
	"github.com/micro/go-micro"
	"log"
	pb "shippy/user-service/proto/user"
)

func main() {

	db, err := CreateConnection()
	defer db.Close()

	if err != nil {
		log.Fatalf("Could not connect to DB: %v", err)
	}

	db.AutoMigrate(&pb.User{})

	repo := &UserRepository{db}

	tokenService := &TokenService{repo}

	server := micro.NewService(
		micro.Name("go.micro.srv.user"),
		micro.Version("latest"),
	)

	server.Init()

	publisher := micro.NewPublisher(topic, server.Client())

	pb.RegisterUserServiceHandler(server.Server(), &handler{repo, tokenService, publisher})

	if err := server.Run(); err != nil {
		fmt.Println(err)
	}
}