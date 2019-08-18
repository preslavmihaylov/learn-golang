package main

import (
	"context"
	"errors"
	"fmt"

	proto "github.com/preslavmihaylov/learn-golang/go-micro-tutorial/vessel-service/proto/vessel"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type vesselRepository struct {
	collection *mongo.Collection
}

func (repo *vesselRepository) Create(vessel *proto.Vessel) error {
	_, err := repo.collection.InsertOne(context.Background(), vessel)
	if err != nil {
		return fmt.Errorf("failed to create vessel: %s", err)
	}

	return nil
}

func (repo *vesselRepository) FindAvailable(spec *proto.Specification) (*proto.Vessel, error) {
	cur, err := repo.collection.Find(context.Background(), bson.D{})
	if err != nil {
		return nil, fmt.Errorf("mongo collection failed to get cursor for collection: %s", err)
	}

	for cur.Next(context.Background()) {
		var v *proto.Vessel
		if err := cur.Decode(&v); err != nil {
			return nil, fmt.Errorf("failed to decode vessel: %s", err)
		}

		if v.Capacity >= spec.Capacity && v.MaxWeight >= spec.MaxWeight {
			return v, nil
		}
	}

	return nil, errors.New("available vessel not found")
}
