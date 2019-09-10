package main

import (
	"context"
	"fmt"
	"log"

	proto "github.com/preslavmihaylov/learn-golang/go-micro-tutorial/consignment-service/proto/consignment"
	vesselProto "github.com/preslavmihaylov/learn-golang/go-micro-tutorial/vessel-service/proto/vessel"
)

// Service should implement all of the methods to satisfy the service
// we defined in our protobuf definition. You can check the interface
// in the generated code itself for the exact method signatures etc
// to give you a better idea.
type service struct {
	repo         *consignmentRepository
	vesselClient vesselProto.VesselService
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
		log.Println(err)
		return fmt.Errorf("failed to create consignment: %s", err)
	}

	log.Printf("Found vessel for consignment: %s\n", vesselResp.Vessel.Name)
	req.VesselId = vesselResp.GetVessel().Id
	err = s.repo.Create(req)
	if err != nil {
		return fmt.Errorf("failed to create consignment: %s", err)
	}

	resp.Consignment = req
	resp.Created = true
	return nil
}

func (s *service) GetConsignments(ctx context.Context, req *proto.GetRequest, resp *proto.Response) error {
	consignments, err := s.repo.GetAll()
	if err != nil {
		log.Println(err)
		return fmt.Errorf("failed to get consignment from repo: %s", err)
	}

	resp.Consignments = consignments
	return nil
}
