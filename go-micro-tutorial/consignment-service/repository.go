package main

import (
	"context"
	"fmt"

	proto "github.com/preslavmihaylov/learn-golang/go-micro-tutorial/consignment-service/proto/consignment"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type consignmentRepository struct {
	collection *mongo.Collection
}

func (repo *consignmentRepository) Create(consignment *proto.Consignment) error {
	_, err := repo.collection.InsertOne(context.Background(), consignment)
	if err != nil {
		return fmt.Errorf("failed to insert consignment in mongodb: %s", err)
	}

	return nil
}

func (repo *consignmentRepository) GetAll() ([]*proto.Consignment, error) {
	cur, err := repo.collection.Find(context.Background(), bson.D{})
	if err != nil {
		return nil, fmt.Errorf("failed to get cursor for collection: %s", err)
	}

	var consignments []*proto.Consignment
	for cur.Next(context.Background()) {
		var consignment *proto.Consignment
		if err := cur.Decode(&consignment); err != nil {
			return nil, fmt.Errorf("failed to decode consignment: %s", err)
		}

		consignments = append(consignments, consignment)
	}

	return consignments, nil
}
