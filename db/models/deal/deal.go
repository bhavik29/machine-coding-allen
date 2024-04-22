package deal

import (
	"allen-machine-coding/db"
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DealSchema struct {
	Name             string
	MaxNumberOfItems int       `bson:"max_number_of_items"`
	DealEndTime      time.Time `bson:"deal_end_time"`
	IsActive         bool      `bson:"is_active"`
}

func InsertOne(ctx context.Context, d *DealSchema) (string, error) {
	db := db.GetInstance()

	coll := db.Client.Database("allen-machine-coding").Collection("deals")
	result, err := coll.InsertOne(ctx, d)
	if err != nil {
		return "", fmt.Errorf("error in saving to database: %v", err)
	}

	oid, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", errors.New("cannot decode insertID")
	}

	return oid.Hex(), nil
}

func FindOne(ctx context.Context, dealID string) *DealSchema {
	db := db.GetInstance()

	coll := db.Client.Database("allen-machine-coding").Collection("deals")
	filters := bson.D{{Key: "_id", Value: dealID}}
	singleResult := coll.FindOne(ctx, filters)

	var result DealSchema
	singleResult.Decode(&result)

	return &result
}

func UpdateOne(ctx context.Context, dealID string, doc *DealSchema) error {
	db := db.GetInstance()

	coll := db.Client.Database("allen-machine-coding").Collection("deals")
	filter := bson.D{{Key: "_id", Value: dealID}}
	update := bson.D{{Key: "$set",
		Value: bson.D{
			{Key: "name", Value: doc.Name},
			{Key: "max_number_of_items", Value: doc.MaxNumberOfItems},
			{Key: "deal_end_time", Value: doc.DealEndTime},
			{Key: "is_active", Value: doc.IsActive},
		},
	}}
	_, err := coll.UpdateOne(ctx, filter, update)

	if err != nil {
		return fmt.Errorf("error in updating the deal: %v", err)
	}

	return nil
}
