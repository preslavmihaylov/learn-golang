package main

import (
	"context"
	"log"
	"os"

	"github.com/micro/go-micro"
	proto "github.com/preslavmihaylov/learn-golang/go-micro-tutorial/vessel-service/proto/vessel"
)

const (
	port        = ":50051"
	defaultHost = "datastore:27017"
)

func main() {
	srv := micro.NewService(micro.Name("shippy.vessel.service"))

	srv.Init()
	uri := os.Getenv("DB_HOST")
	if uri == "" {
		uri = defaultHost
	}

	client, err := CreateClient(uri)
	if err != nil {
		log.Panic(err)
	}
	defer client.Disconnect(context.TODO())

	repo := vesselRepository{client.Database("shippy").Collection("vessels")}
	err = repo.Create(&proto.Vessel{
		Id:        "vessel001",
		Name:      "Kane's Salty Secret",
		MaxWeight: 200000,
		Capacity:  500,
	})
	if err != nil {
		log.Panic(err)
	}

	proto.RegisterVesselServiceHandler(srv.Server(), &service{&repo})

	if err := srv.Run(); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
