package storage

import (
	"errors"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/vikin91/bid-tracker-go/pkg/config"
	"github.com/vikin91/bid-tracker-go/pkg/models"
)

/*
 * Current state:
 * Place Bid = O(1)
 * item.AllBids = O(1)
 * item.WinningBid = O(1)
 * user.ItemsBided = O(1)
 */

//MapBiddingSystem is a data structure ...
type MapBiddingSystem struct {
	Items map[uuid.UUID]*models.Item
	Users map[uuid.UUID]*models.User
}

//NewMapBiddingSystem creates empty BiddingSystem
func NewMapBiddingSystem() *MapBiddingSystem {
	return &MapBiddingSystem{
		Items: make(map[uuid.UUID]*models.Item),
		Users: make(map[uuid.UUID]*models.User),
	}
}

//AllItems ...
func (h *MapBiddingSystem) AllItems() ([]*models.Item, error) {
	var values []*models.Item = make([]*models.Item, 0)
	for _, v := range h.Items {
		values = append(values, v)
	}
	return values, nil
}

//CreateItem ...
func (h *MapBiddingSystem) CreateItem(item *models.Item) error {
	if item.ID == config.ZeroUUID {
		item.ID = uuid.NewV4()
		item.CreatedAt = time.Now()
	}
	h.Items[item.ID] = item
	return nil
}

//GetItem ...
func (h *MapBiddingSystem) GetItem(id uuid.UUID) (*models.Item, error) {
	if itm, ok := h.Items[id]; ok {
		return itm, nil
	}
	return &models.Item{}, errors.New("Cannot find item")
}

//AllUsers ...
func (h *MapBiddingSystem) AllUsers() ([]*models.User, error) {
	var values []*models.User = make([]*models.User, 0)
	for _, v := range h.Users {
		values = append(values, v)
	}
	return values, nil
}

//CreateUser ...
func (h *MapBiddingSystem) CreateUser(user *models.User) error {
	if user.ID == config.ZeroUUID {
		user.ID = uuid.NewV4()
		user.CreatedAt = time.Now()
	}
	h.Users[user.ID] = user
	return nil
}

//GetUser ...
func (h *MapBiddingSystem) GetUser(id uuid.UUID) (*models.User, error) {
	if usr, ok := h.Users[id]; ok {
		return usr, nil
	}
	return &models.User{}, errors.New("Cannot find user")
}

//AllBids ...
func (h *MapBiddingSystem) AllBids() ([]*models.Bid, error) {
	var values []*models.Bid = make([]*models.Bid, 0)
	for _, item := range h.Items {
		for _, bidsOnItem := range item.GetBids() {
			values = append(values, bidsOnItem)
		}
	}
	return values, nil
}

//GetUserBids ...
func (h *MapBiddingSystem) GetUserBids(userID uuid.UUID) ([]*models.Bid, error) {
	user, err := h.GetUser(userID)
	if err != nil {
		return make([]*models.Bid, 0), errors.New("User not found")
	}
	return user.GetBids(), nil
}

//PlaceBid (ASSIGNMENT FUNCTION)
func (h *MapBiddingSystem) PlaceBid(bid *models.Bid) error {
	if bid.ID == config.ZeroUUID {
		bid.ID = uuid.NewV4()
	}
	if bid.CreatedAt.IsZero() {
		bid.CreatedAt = time.Now()
	}
	item, err := h.GetItem(bid.ItemID)
	if err != nil {
		return err
	}
	user, err := h.GetUser(bid.UserID)
	if err != nil {
		return err
	}

	item.PlaceNewBid(bid)
	user.PlaceNewBidOnItem(bid, item)
	return nil
}

//GetItemsUserHasBid (ASSIGNMENT FUNCTION) returns a slice of items no which user has placed at least one bid
func (h *MapBiddingSystem) GetItemsUserHasBid(userID uuid.UUID) ([]*models.Item, error) {
	user, err := h.GetUser(userID)
	if err != nil {
		return []*models.Item{}, errors.New("User not found")
	}
	return user.ItemsBid, nil
}

//GetBidsOnItem (ASSIGNMENT FUNCTION) return slice of Bids for an Item
func (h *MapBiddingSystem) GetBidsOnItem(itemID uuid.UUID) ([]*models.Bid, error) {
	item, err := h.GetItem(itemID)
	if err != nil {
		return make([]*models.Bid, 0), errors.New("Item not found")
	}
	return item.GetBids(), nil
}

//GetWinningBid (ASSIGNMENT FUNCTION) returns bid with the highest amount for an item
func (h *MapBiddingSystem) GetWinningBid(itemID uuid.UUID) (*models.Bid, error) {
	item, err := h.GetItem(itemID)
	if err != nil {
		return nil, err
	}
	return item.GetWinningBid()
}

//Reset empties the data structure
func (h *MapBiddingSystem) Reset() {
	h = NewMapBiddingSystem()
}
