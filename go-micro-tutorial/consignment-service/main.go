// shippy-service-consignment/main.go
package main

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/micro/go-micro"

	// Import the generated protobuf code
	proto "github.com/preslavmihaylov/learn-golang/go-micro-tutorial/consignment-service/proto/consignment"
	vesselProto "github.com/preslavmihaylov/learn-golang/go-micro-tutorial/vessel-service/proto/vessel"
)

// consignmentRepository - Dummy repository, this simulates the use of a datastore
// of some kind. We'll replace this with a real implementation later on.
type consignmentRepository struct {
	mu           sync.RWMutex
	consignments []*proto.Consignment
}

// Create a new consignment
func (repo *consignmentRepository) Create(consignment *proto.Consignment) (*proto.Consignment, error) {
	repo.mu.Lock()
	repo.consignments = append(repo.consignments, consignment)
	repo.mu.Unlock()

	return consignment, nil
}

func (repo *consignmentRepository) GetAll() []*proto.Consignment {
	return repo.consignments
}

// Service should implement all of the methods to satisfy the service
// we defined in our protobuf definition. You can check the interface
// in the generated code itself for the exact method signatures etc
// to give you a better idea.
type service struct {
	repo         *consignmentRepository
	vesselClient vesselProto.VesselServiceClient
}

// CreateConsignment - we created just one method on our service,
// which is a create method, which takes a context and a request as an
// argument, these are handled by the gRPC server.
func (s *service) CreateConsignment(ctx context.Context, req *proto.Consignment, resp *proto.Response) error {
	vesselResp, err := s.vesselClient.FindAvailable(ctx, &vesselProto.Specification{
		MaxWeight: req.Weight,
		Capacity:  int32(len(req.Containers)),
	})
	if err != nil {
		return fmt.Errorf("failed to create consignment: %s", err)
	}

	log.Printf("Found vessel for consignment: %s\n", vesselResp.Vessel.Name)
	req.VesselId = vesselResp.GetVessel().Id
	consignment, err := s.repo.Create(req)
	if err != nil {
		return fmt.Errorf("failed to create consignment: %s", err)
	}

	resp.Consignment = consignment
	resp.Created = true

	return nil
}

func (s *service) GetConsignments(ctx context.Context, req *proto.GetRequest, res *proto.Response) error {
	consignments := s.repo.GetAll()
	res.Consignments = consignments

	return nil
}

func main() {
	srv := micro.NewService(
		micro.Name("shippy.consignment.service"),
	)

	srv.Init()
	vesselClient := vesselProto.NewVesselServiceClient("shippy.vessel.service", srv.Client())

	// Register our service with the gRPC server, this will tie our
	// implementation into the auto-generated interface code for our
	// protobuf definition.
	proto.RegisterShippingServiceHandler(srv.Server(), &service{
		repo:         &consignmentRepository{},
		vesselClient: vesselClient,
	})

	if err := srv.Run(); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
