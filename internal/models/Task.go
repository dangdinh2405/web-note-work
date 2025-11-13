package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Task struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Title       string    `bson:"title"       json:"title"`
	Status      string    `bson:"status"      json:"status"`      
	CompletedAt *time.Time `bson:"completedAt" json:"completedAt"`  
	CreatedAt   time.Time `bson:"createdAt"   json:"createdAt"`
	UpdatedAt   time.Time `bson:"updatedAt"   json:"updatedAt"`
}
