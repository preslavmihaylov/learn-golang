// shippy-service-consignment/main.go
package main

import (
	"context"
	"log"
	"os"

	"github.com/micro/go-micro"

	// Import the generated protobuf code
	proto "github.com/preslavmihaylov/learn-golang/go-micro-tutorial/consignment-service/proto/consignment"
	vesselProto "github.com/preslavmihaylov/learn-golang/go-micro-tutorial/vessel-service/proto/vessel"
)

const (
	port        = ":50051"
	defaultHost = "datastore:27017"
)

func main() {
	srv := micro.NewService(
		micro.Name("shippy.consignment.service"),
	)

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

	consignmentCollection := client.Database("shippy").Collection("consignments")

	repo := consignmentRepository{consignmentCollection}
	vesselClient := vesselProto.NewVesselServiceClient("shippy.vessel.service", srv.Client())

	proto.RegisterShippingServiceHandler(srv.Server(), &service{&repo, vesselClient})

	if err := srv.Run(); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
