package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

// BaseModel defines the basic fields for each other model
type BaseModel struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
}

//NewBaseModel creates a new BaseModel object with random UUID and CreatedAt set to now
func NewBaseModel() BaseModel {
	return BaseModel{ID: uuid.NewV4(), CreatedAt: time.Now()}
}
