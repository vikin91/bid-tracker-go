package testutils

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"sort"
	"strings"

	"github.com/vikin91/bid-tracker-go/pkg/models"
	"github.com/vikin91/bid-tracker-go/pkg/storage"
)

// GetRequestPayload converts a given object into a reader of that object as json payload
func GetRequestPayload(payload interface{}) io.Reader {
	bytes, _ := json.Marshal(payload)
	return strings.NewReader(string(bytes))
}

//GenerateAmountsMatrix generates a matrix of size (numItems x numBidOnItem) and returns it along with an array holding the winning bid amount for each item
func GenerateAmountsMatrix(numItems int, numBidsOnItem int) ([][]float64, []float64) {
	amountsMatrix := make([][]float64, numItems)
	// maxAmounts holds winning bid amount for each item
	maxAmounts := make([]float64, numItems)
	for i := 0; i < numItems; i++ {
		amountsMatrix[i] = make([]float64, numBidsOnItem)
		amountsMatrix[i] = GenerateSliceOfRandomFloat64(numBidsOnItem)

		tmp := make([]float64, numBidsOnItem)
		copy(tmp, amountsMatrix[i])
		sort.Float64s(tmp)
		maxAmounts[i] = tmp[numBidsOnItem-1]
	}
	return amountsMatrix, maxAmounts
}

//EmptyTestDB initializes DB-like data structure
func EmptyTestDB() storage.Storage {
	return EmptyTestSimpleMap()
}

//EmptyTestSimpleMap initializes DB-like data structure
func EmptyTestSimpleMap() storage.Storage {
	return storage.NewMapBiddingSystem()
}

//CreateTestUsers creates some users for test
func CreateTestUsers(h storage.Storage, num int) []*models.User {
	var users []*models.User
	for i := 0; i < num; i++ {
		user := models.NewUser(fmt.Sprintf("James Bond 007-%03d", i))
		users = append(users, user)
		h.CreateUser(user)
	}
	return users
}

//CreateTestItems creates some items for test
func CreateTestItems(h storage.Storage, num int) []*models.Item {

	var items []*models.Item

	for i := 0; i < num; i++ {
		item := models.NewItem(fmt.Sprintf("A-thing-%03d", i))
		items = append(items, item)
		h.CreateItem(item)
	}
	return items
}

//GenerateSliceOfRandomFloat64 generates a slice of length len filled with random float64s
func GenerateSliceOfRandomFloat64(len int) []float64 {
	s := make([]float64, len)
	for idx := range s {
		s[idx] = rand.ExpFloat64()
	}
	return s
}

//CreateTestBids creates some bids for test - on each item, exactly one bid
func CreateTestBids(h storage.Storage, num int, amounts []float64) ([]*models.Bid, []*models.Item, []*models.User) {

	var bids []*models.Bid
	users := CreateTestUsers(h, num)
	items := CreateTestItems(h, num)

	for i := 0; i < num; i++ {
		bid := models.NewBid(items[i].ID, users[i].ID, amounts[i])
		h.PlaceBid(bid)
		itemBids, _ := h.GetBidsOnItem(items[i].ID)
		bids = append(bids, itemBids...)
	}
	return bids, items, users
}

//CreateTestBidsManyOnItem creates some bids for test - multiple bids for single item
func CreateTestBidsManyOnItem(h storage.Storage, num int, amounts [][]float64) ([]*models.Bid, []*models.Item, []*models.User) {

	var bids []*models.Bid
	users := CreateTestUsers(h, num)
	items := CreateTestItems(h, num)

	for itemIdx := 0; itemIdx < len(amounts); itemIdx++ {
		bidAmounts := amounts[itemIdx]
		for bidIdx := 0; bidIdx < len(bidAmounts); bidIdx++ {
			bid := models.NewBid(items[itemIdx].ID, users[itemIdx].ID, bidAmounts[bidIdx])
			h.PlaceBid(bid)
		}
		itemBids, _ := h.GetBidsOnItem(items[itemIdx].ID)
		bids = append(bids, itemBids...)
	}
	return bids, items, users
}

//CreateTestTwoUsersBidOnManyItems creates some bids for test - multiple bids for single item, but only two users
func CreateTestTwoUsersBidOnManyItems(h storage.Storage, numItems int, amounts [][]float64) ([]*models.Bid, []*models.Item, []*models.User) {

	var bids []*models.Bid
	users := CreateTestUsers(h, 2)
	items := CreateTestItems(h, numItems)

	for itemIdx := 0; itemIdx < len(amounts); itemIdx++ {
		bidAmounts := amounts[itemIdx]
		for bidIdx := 0; bidIdx < len(bidAmounts); bidIdx++ {
			userIdx := itemIdx % 2
			bid := models.NewBid(items[itemIdx].ID, users[userIdx].ID, bidAmounts[bidIdx])
			h.PlaceBid(bid)
		}
		itemBids, _ := h.GetBidsOnItem(items[itemIdx].ID)
		bids = append(bids, itemBids...)
	}
	return bids, items, users
}
