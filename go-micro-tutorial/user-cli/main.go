package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/micro/go-micro"
	"github.com/micro/go-micro/config/cmd"
	proto "github.com/preslavmihaylov/learn-golang/go-micro-tutorial/user-service/proto/user"
	"golang.org/x/net/context"
)

func main() {
	name := flag.String("name", "pm", "Your full name")
	email := flag.String("email", "pm@pm.com", "Your Email")
	password := flag.String("password", "passsecret", "Your Password")
	company := flag.String("company", "BBC", "Your Company")
	flag.Parse()

	cmd.Init()

	srv := micro.NewService(micro.Name("shippy.user.cli"))
	srv.Init()

	client := proto.NewUserServiceClient("shippy.user.service", srv.Client())

	fmt.Printf("Creating user: %s %s %s %s\n", *name, *email, *password, *company)
	r, err := client.Create(context.TODO(), &proto.User{
		Name:     *name,
		Email:    *email,
		Password: *password,
		Company:  *company,
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
}
