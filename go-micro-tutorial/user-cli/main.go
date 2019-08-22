package main

import (
	"fmt"
	"log"

	"github.com/micro/go-micro"
	"github.com/micro/go-micro/config/cmd"
	proto "github.com/preslavmihaylov/learn-golang/go-micro-tutorial/user-service/proto/user"
	"golang.org/x/net/context"
)

func main() {
	name := "pm"
	email := "pm@pm.com"
	password := "abcd"
	company := "BBC"

	cmd.Init()

	srv := micro.NewService(micro.Name("shippy.user.cli"))
	srv.Init()

	client := proto.NewUserServiceClient("shippy.user.service", srv.Client())

	fmt.Printf("Creating user: %s %s %s %s\n", name, email, password, company)
	r, err := client.Create(context.TODO(), &proto.User{
		Name:     name,
		Email:    email,
		Password: password,
		Company:  company,
	})
	if err != nil {
		log.Fatalf("Could not create: %v", err)
	}
	log.Printf("Created: %s", r.User.Id)

	getAll, err := client.GetAll(context.Background(), &proto.Request{})
	if err != nil {
		log.Fatalf("Could not list users: %v", err)
	}
	for _, v := range getAll.Users {
		log.Println(v)
	}

	authResp, err := client.Auth(context.TODO(), &proto.User{
		Email:    email,
		Password: password,
	})
	if err != nil {
		log.Fatalf("couldn't authenticate user: %s", err)
	}

	fmt.Println("Access Token:", authResp)

	authResp, err = client.Auth(context.TODO(), &proto.User{
		Email:    email,
		Password: "wrong_pass",
	})
	if err == nil {
		log.Fatalf("[ERR] email with wrong password authenticated successfully!")
	}

	log.Printf("[EXPECTED] couldn't authenticate user %s (%s): %s", email, "wrong_pass", err)
}
