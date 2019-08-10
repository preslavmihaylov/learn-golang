package main

import (
	"context"
	"errors"
	"log"

	"github.com/micro/go-micro"
	proto "github.com/preslavmihaylov/learn-golang/go-micro-tutorial/vessel-service/proto/vessel"
)

type vesselRepository struct {
	vessels []*proto.Vessel
}

func (vr *vesselRepository) FindAvailable(spec *proto.Specification) (*proto.Vessel, error) {
	for _, v := range vr.vessels {
		if v.Capacity >= spec.Capacity && v.MaxWeight >= spec.MaxWeight {
			return v, nil
		}
	}
	return nil, errors.New("Available vessel not found")
}

type service struct {
	repo *vesselRepository
}

func (s *service) FindAvailable(ctx context.Context, spec *proto.Specification, res *proto.Response) error {
	vessel, err := s.repo.FindAvailable(spec)
	if err != nil {
		return err
	}

	log.Printf("Found vessel %v for specs %v", vessel, spec)
	res.Vessel = vessel
	return nil
}

func main() {
	repo := vesselRepository{vessels: []*proto.Vessel{
		&proto.Vessel{Id: "vessel001", Name: "Boaty McBoatface", MaxWeight: 200000, Capacity: 500},
	}}
	srv := micro.NewService(
		micro.Name("shippy.vessel.service"),
	)

	srv.Init()

	// Register our service with the gRPC server, this will tie our
	// implementation into the auto-generated interface code for our
	// protobuf definition.
	proto.RegisterVesselServiceHandler(srv.Server(), &service{&repo})

	if err := srv.Run(); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
