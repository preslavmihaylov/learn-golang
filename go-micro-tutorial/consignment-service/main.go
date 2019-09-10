// shippy-service-consignment/main.go
package main

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/micro/go-micro"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/metadata"
	"github.com/micro/go-micro/server"

	// Import the generated protobuf code
	proto "github.com/preslavmihaylov/learn-golang/go-micro-tutorial/consignment-service/proto/consignment"
	userProto "github.com/preslavmihaylov/learn-golang/go-micro-tutorial/user-service/proto/user"
	vesselProto "github.com/preslavmihaylov/learn-golang/go-micro-tutorial/vessel-service/proto/vessel"
)

const (
	port        = ":50051"
	defaultHost = "datastore:27017"
)

func AuthWrapper(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, resp interface{}) error {
		if os.Getenv("DISABLE_AUTH") == "true" {
			return fn(ctx, req, resp)
		}

		meta, ok := metadata.FromContext(ctx)
		if !ok {
			return errors.New("no auth meta-data found in request")
		}

		// Note this is now uppercase (not entirely sure why this is...)
		token := meta["Token"]
		log.Println("Authenticating with token: ", token)

		// Auth here
		authClient := userProto.NewUserService("shippy.user.service", client.DefaultClient)
		_, err := authClient.ValidateToken(context.Background(), &userProto.Token{
			Token: token,
		})
		if err != nil {
			return err
		}

		return fn(ctx, req, resp)
	}
}

func main() {
	srv := micro.NewService(
		micro.Name("shippy.consignment.service"),
		micro.Version("latest"),
		micro.WrapHandler(AuthWrapper),
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
	vesselClient := vesselProto.NewVesselService("shippy.vessel.service", srv.Client())

	proto.RegisterShippingServiceHandler(srv.Server(), &service{&repo, vesselClient})

	if err := srv.Run(); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
