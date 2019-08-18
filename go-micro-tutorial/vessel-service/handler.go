package main

import (
	"context"
	"fmt"
	"log"

	proto "github.com/preslavmihaylov/learn-golang/go-micro-tutorial/vessel-service/proto/vessel"
)

type service struct {
	repo *vesselRepository
}

func (s *service) Create(ctx context.Context, vessel *proto.Vessel, resp *proto.Response) error {
	err := s.repo.Create(vessel)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("failed to create vessel in repo: %s", err)
	}

	resp.Vessel = vessel
	return nil
}

func (s *service) FindAvailable(ctx context.Context, spec *proto.Specification, res *proto.Response) error {
	vessel, err := s.repo.FindAvailable(spec)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("failed to find available vessel in repo: %s", err)
	}

	log.Printf("found vessel %v for specs %v", vessel, spec)
	res.Vessel = vessel
	return nil
}
