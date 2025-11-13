package models

import (
	"time"
)

type Task struct {
	Title       string    `bson:"title"       json:"title"`
	Status      string    `bson:"status"      json:"status"`      
	CompletedAt *time.Time `bson:"completedAt" json:"completedAt"`  
	CreatedAt   time.Time `bson:"createdAt"   json:"createdAt"`
	UpdatedAt   time.Time `bson:"updatedAt"   json:"updatedAt"`
}
