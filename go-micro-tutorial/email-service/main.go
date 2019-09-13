package main

import (
	"context"
	"fmt"
	"log"

	"github.com/micro/go-micro"
	proto "github.com/preslavmihaylov/learn-golang/go-micro-tutorial/user-service/proto/user"
)

const topic = "user.created"

type Subscriber struct{}

func (sub *Subscriber) Process(ctx context.Context, usr *proto.User) error {
	log.Println("Picked up a new message")
	log.Println("sending email to:", usr.Name)
	return nil
}

func main() {
	fmt.Println("vim-go")
	srv := micro.NewService(
		micro.Name("shippy-email-service"),
		micro.Version("latest"),
	)

	srv.Init()

	pubsub := srv.Server().Options().Broker
	err := pubsub.Connect()
	if err != nil {
		log.Fatal(err)
	}

	micro.RegisterSubscriber(topic, srv.Server(), &Subscriber{})

	err = srv.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func sendEmail(usr *proto.User) error {
	log.Println("sending email to:", usr.Name)

	return nil
}
