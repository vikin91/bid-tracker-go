package models

import (
	uuid "github.com/satori/go.uuid"
)

// define error messages
const (
	BidCreationFailure = "Failed creating bid"
	BidListFailure     = "Failed listing bid"
)

// Bid model
type Bid struct {
	BaseModel
	ItemID uuid.UUID `json:"itemID"`
	UserID uuid.UUID `json:"userID"`
	Amount float64   `json:"amount"`
}

//NewBid creates an Item
func NewBid(itemID uuid.UUID, userID uuid.UUID, amount float64) *Bid {
	return &Bid{
		BaseModel: NewBaseModel(),
		UserID:    userID,
		ItemID:    itemID,
		Amount:    amount,
	}
}
