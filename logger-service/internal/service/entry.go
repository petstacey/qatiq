package service

import (
	"context"
	"log"
	"time"

	"github.com/pso-dev/qatiq/backend/logger-service/internal/validator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type LogEntry struct {
	ID        string    `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string    `bson:"name" json:"name"`
	Data      string    `bson:"data" json:"data"`
	CreatedAt time.Time `bson:"created_at" json:"createdAt"`
	UpdatedAt time.Time `bson:"updated_at" json:"updatedAt"`
}

func ValidateName(v *validator.Validator, name string) {
	v.Check(name != "", "name", "must be provided")
}

func ValidateData(v *validator.Validator, data string) {
	v.Check(data != "", "data", "must be provided")
}

func ValidateUpdated(v *validator.Validator, created, updated time.Time) {
	v.Check(updated.Sub(created) >= 0, "updatedAt", "cannot be before createdAt")
}

func ValidateEntry(v *validator.Validator, entry LogEntry) {
	ValidateName(v, entry.Name)
	ValidateData(v, entry.Data)
	ValidateUpdated(v, entry.CreatedAt, entry.UpdatedAt)
}

func (s *Service) LogItem(entry LogEntry) error {
	collection := s.DB.Database("logs").Collection("logs")
	_, err := collection.InsertOne(context.TODO(), LogEntry{
		Name:      entry.Name,
		Data:      entry.Data,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		log.Println("error inserting into logs:", err)
		return err
	}
	return nil
}

func (s *Service) GetAll() ([]*LogEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	collection := s.DB.Database("logs").Collection("logs")
	opts := options.Find()
	opts.SetSort(bson.D{{Key: "created_at", Value: -1}})
	cursor, err := collection.Find(context.TODO(), bson.D{}, opts)
	if err != nil {
		log.Println("finding all docs error:", err)
		return nil, err
	}
	defer cursor.Close(ctx)
	var logs []*LogEntry
	for cursor.Next(ctx) {
		var item LogEntry

		err := cursor.Decode(item)
		if err != nil {
			log.Println("error decoding log into logs:", err)
			return nil, err
		} else {
			logs = append(logs, &item)
		}
	}
	return logs, nil
}

func (s *Service) GetOne(id string) (*LogEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	collection := s.DB.Database("logs").Collection("logs")
	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var entry LogEntry
	err = collection.FindOne(ctx, bson.M{"_id": docID}).Decode(&entry)
	if err != nil {
		return nil, err
	}
	return &entry, nil
}

func (s *Service) DropCollection() error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	collection := s.DB.Database("logs").Collection("logs")
	if err := collection.Drop(ctx); err != nil {
		return err
	}
	return nil
}

func (s *Service) Update(entry *LogEntry) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	collection := s.DB.Database("logs").Collection("logs")
	docID, err := primitive.ObjectIDFromHex(entry.ID)
	if err != nil {
		return nil, err
	}
	result, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": docID},
		bson.D{
			{Key: "$set", Value: bson.D{
				{Key: "name", Value: entry.Name},
				{Key: "data", Value: entry.Data},
				{Key: "updated_at", Value: time.Now()},
			}},
		},
	)
	if err != nil {
		return nil, err
	}
	return result, nil
}
