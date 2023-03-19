package service

import (
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
)

var ErrorFailedValidation = errors.New("failed validation")

type Service struct {
	DB *mongo.Client
}
