package models

import (
	"errors"
	"sync"
)

// define error messages
const (
	ItemCreationFailure = "Failed creating item"
	ItemListFailure     = "Failed listing item"
)

// Item model
type Item struct {
	BaseModel
	Name string `json:"name"`

	mutexBids sync.RWMutex
	bids      []*Bid

	mutexBestBid sync.RWMutex
	WinningBid   *Bid    `json:"-"`
	MaxBidAmount float64 `json:"-"`
}

//NewItem creates an Item
func NewItem(name string) *Item {
	return &Item{
		BaseModel:    NewBaseModel(),
		Name:         name,
		bids:         make([]*Bid, 0),
		MaxBidAmount: float64(0.0),
		WinningBid:   nil,
	}
}

//PlaceNewBid handles the bid placement
func (i *Item) PlaceNewBid(bid *Bid) {
	i.UpdateBestBid(bid)

	i.mutexBids.Lock()
	defer i.mutexBids.Unlock()
	i.bids = append(i.bids, bid)
}

//GetBids handles the bid placement
func (i *Item) GetBids() []*Bid {
	i.mutexBids.RLock()
	defer i.mutexBids.RUnlock()
	return i.bids
}

//UpdateBestBid updates information about currently best bid
func (i *Item) UpdateBestBid(bid *Bid) {
	i.mutexBestBid.Lock()
	defer i.mutexBestBid.Unlock()

	if bid.Amount > i.MaxBidAmount {
		i.MaxBidAmount = bid.Amount
		i.WinningBid = bid
	}
}

//GetWinningBid updates information about currently best bid
func (i *Item) GetWinningBid() (*Bid, error) {
	i.mutexBestBid.RLock()
	defer i.mutexBestBid.RUnlock()

	if i.WinningBid == nil {
		return nil, errors.New("Cannot find valid bids on this item")
	}
	return i.WinningBid, nil
}
