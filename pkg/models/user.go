package models

import (
	"sync"

	uuid "github.com/satori/go.uuid"
)

// define error messages
const (
	UserCreationFailure = "Failed creating user"
	UserListFailure     = "Failed listing users"
)

// User model
type User struct {
	BaseModel
	Name       string `json:"name"`
	mutexBids  sync.RWMutex
	bids       map[uuid.UUID]*Bid //TODO: Could be a simple slice + mutex - not optimizing this, as not required in assignment
	mutexItems sync.Mutex
	//itemsBidFlag guards uniqueness of ItemsBid
	itemsBidFlag map[uuid.UUID]struct{}
	ItemsBid     []*Item `json:"-"`
}

//NewUser creates an User
func NewUser(name string) *User {
	return &User{
		BaseModel:    NewBaseModel(),
		Name:         name,
		bids:         make(map[uuid.UUID]*Bid),
		itemsBidFlag: make(map[uuid.UUID]struct{}),
		ItemsBid:     make([]*Item, 0),
	}
}

//GetBids returns all bids the user has placed
func (u *User) GetBids() []*Bid {
	values := make([]*Bid, 0)

	u.mutexBids.RLock()
	defer u.mutexBids.RUnlock()

	for _, userBids := range u.bids {
		values = append(values, userBids)
	}
	return values
}

//PlaceNewBidOnItem handles the bid placement for user - adds items to a memo
func (u *User) PlaceNewBidOnItem(bid *Bid, item *Item) {
	u.registerBid(bid)

	u.mutexItems.Lock()
	defer u.mutexItems.Unlock()

	if _, ok := u.itemsBidFlag[item.ID]; !ok {
		u.ItemsBid = append(u.ItemsBid, item)
		u.itemsBidFlag[item.ID] = struct{}{}
	}
}

func (u *User) registerBid(bid *Bid) {
	u.mutexBids.Lock()
	defer u.mutexBids.Unlock()
	u.bids[bid.ID] = bid
}
