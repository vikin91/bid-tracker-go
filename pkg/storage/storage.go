package storage

import (
	uuid "github.com/satori/go.uuid"
	"github.com/vikin91/bid-tracker-go/pkg/models"
)

//Storage is an interface for underlying data structure storing a state - useful when implementing multiple storage backends
type Storage interface {
	//CR methods for item
	AllItems() ([]*models.Item, error)
	CreateItem(*models.Item) error
	GetItem(id uuid.UUID) (*models.Item, error)

	//CR methods for user
	AllUsers() ([]*models.User, error)
	CreateUser(*models.User) error
	GetUser(id uuid.UUID) (*models.User, error)

	//CR methods for item
	AllBids() ([]*models.Bid, error)
	GetUserBids(userID uuid.UUID) ([]*models.Bid, error)

	//Important functions required in the assignment
	PlaceBid(*models.Bid) error
	GetBidsOnItem(itemID uuid.UUID) ([]*models.Bid, error)
	GetItemsUserHasBid(userID uuid.UUID) ([]*models.Item, error)
	GetWinningBid(itemID uuid.UUID) (*models.Bid, error)

	Reset()
}
