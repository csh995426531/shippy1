package main

import (
	"context"
	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/config/cmd"
	"log"
	"os"
	pb "shippy/user-service/proto/user"
	microclient "github.com/micro/go-micro/client"
)
func main() {

	cmd.Init()
	// 创建 user-service 微服务的客户端
	client := pb.NewUserServiceClient("go.micro.srv.user", microclient.DefaultClient)

	// 设置命令行参数
	service := micro.NewService(
		micro.Flags(
			cli.StringFlag{
				Name:        "name",
				Usage:       "You full name",
			},
			cli.StringFlag{
				Name:  "email",
				Usage: "Your email",
			},
			cli.StringFlag{
				Name:  "password",
				Usage: "Your password",
			},
			cli.StringFlag{
				Name: "company",
				Usage: "Your company",
			},
		),
	)

	service.Init(
		micro.Action(func(c *cli.Context) {
			//name := c.String("name")
			//email := c.String("email")
			//password := c.String("password")
			//company := c.String("company")
			name := "Ewan Valentine"
			email := "ewan.valentine89@gmail.com"
			password := "Testing123"
			company := "BBC"
			user := pb.User{
				Name:                 name,
				Company:              company,
				Email:                email,
				Password:             password,
			}

			r, err := client.Create(context.TODO(), &user)
			log.Printf("user %T, %#v", user, user)

			if err != nil {
				log.Fatalf("Could not create: %v", err)
			}
			log.Printf("Created: %v", r.User.Id)

			getAll, err := client.GetAll(context.Background(), &pb.Request{})
			if err != nil {
				log.Fatalf("Could not list users: %v", err)
			}
			for _, v := range getAll.Users {
				log.Println(v)
			}

			authResp, err := client.Auth(context.TODO(), &pb.User{
				Email: email,
				Password: password,
			})

			if err != nil {
				log.Fatalf("auth failed: %v", err)
			}
			log.Println("token:", authResp.Token)
			os.Exit(0)
		}),
	)

	if err := service.Run(); err != nil {
		log.Println(err)
	}
}