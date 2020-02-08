package storage_test

import (
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"sort"
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/vikin91/bid-tracker-go/internal/testutils"
	"github.com/vikin91/bid-tracker-go/pkg/config"
	"github.com/vikin91/bid-tracker-go/pkg/models"
	"github.com/vikin91/bid-tracker-go/pkg/storage"
)

var h = storage.NewMapBiddingSystem()

//scale is related to the maximum number of elements used in benchmarks - 2^scale
const scale = 8

func Test_AllItems(t *testing.T) {
	tests := []struct {
		name           string
		numItemsCreate int
		wantNumItems   int
		wantErr        bool
	}{
		{"Should get empty set of items", 0, 0, false},
		{"Should get 2 items", 2, 2, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h.Reset()
			testutils.CreateTestItems(h, tt.numItemsCreate)

			got, err := h.AllItems()
			assert.Equal(t, tt.wantNumItems, len(got), "Got Wrong number of items")
			if (err != nil) != tt.wantErr {
				t.Errorf(".AllItems() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_AllUsers(t *testing.T) {

	tests := []struct {
		name           string
		numUsersCreate int
		wantNumUsers   int
		wantErr        bool
	}{
		{"Should get empty set of items", 0, 0, false},
		{"Should get 2 items", 2, 2, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h.Reset()
			testutils.CreateTestUsers(h, tt.numUsersCreate)

			got, err := h.AllUsers()
			assert.Equal(t, tt.wantNumUsers, len(got), "Got Wrong number of items")
			if (err != nil) != tt.wantErr {
				t.Errorf(".AllUsers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_AllBids(t *testing.T) {
	tests := []struct {
		name          string
		numBidsCreate int
		wantNumBids   int
		wantErr       bool
	}{
		{"Should get empty set of items", 0, 0, false},
		{"Should get 2 items", 2, 2, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h.Reset()
			testutils.CreateTestBids(h, tt.numBidsCreate, testutils.GenerateSliceOfRandomFloat64(tt.numBidsCreate))

			got, err := h.AllBids()
			assert.Equal(t, tt.wantNumBids, len(got), "Got Wrong number of items")
			if (err != nil) != tt.wantErr {
				t.Errorf(".AllBids() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_CreateItem(t *testing.T) {

	tests := []struct {
		name     string
		itemName string
		wantErr  bool
	}{
		{"Should create an item", "A thing", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h.Reset()
			item := &models.Item{Name: tt.itemName}
			if err := h.CreateItem(item); (err != nil) != tt.wantErr {
				t.Errorf(".CreateItem() error = %v, wantErr %v", err, tt.wantErr)
			}
			//CreateItem should set required fields
			assert.NotNil(t, item.ID)
			assert.NotNil(t, item.CreatedAt)
			assert.NotEqualf(t, config.ZeroUUID, item.ID, "UUID must mot be empty. Is: %s", item.ID.String())
			assert.NotEqualf(t, time.Time{}, item.CreatedAt, "CreatedAt must mot be empty. Is: %s", item.CreatedAt.String())
		})
	}
}

func Test_CreateUser(t *testing.T) {

	tests := []struct {
		name     string
		userName string
		wantErr  bool
	}{
		{"test1", "James Bond", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h.Reset()
			user := &models.User{Name: tt.userName}
			if err := h.CreateUser(user); (err != nil) != tt.wantErr {
				t.Errorf(".CreateUser() error = %v, wantErr %v", err, tt.wantErr)
			}
			//CreateUser should set required fields
			assert.NotNil(t, user.ID)
			assert.NotNil(t, user.CreatedAt)
			assert.NotEqualf(t, config.ZeroUUID, user.ID, "UUID must mot be empty. Is: %s", user.ID.String())
			assert.NotEqualf(t, time.Time{}, user.CreatedAt, "CreatedAt must mot be empty. Is: %s", user.CreatedAt.String())
		})
	}
}

func Test_PlaceBid(t *testing.T) {

	h.Reset()
	users := testutils.CreateTestUsers(h, 10)
	items := testutils.CreateTestItems(h, 10)

	tests := []struct {
		name    string
		userIdx int
		itemIdx int
		amount  float64
		wantErr bool
	}{
		{"Bid should be added", 0, 0, 9.99, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := users[tt.userIdx]
			item := items[tt.itemIdx]
			bid := models.NewBid(item.ID, user.ID, tt.amount)

			if err := h.PlaceBid(bid); (err != nil) != tt.wantErr {
				t.Errorf(".PlaceBid() error = %v, wantErr %v", err, tt.wantErr)
			}
			//PlaceBid should set required fields
			assert.NotNil(t, bid.ID)
			assert.NotNil(t, bid.UserID)
			assert.NotNil(t, bid.ItemID)
			assert.NotNil(t, bid.CreatedAt)
			assert.NotEqualf(t, config.ZeroUUID, bid.ID, "Bid ID must mot be empty. Is: %s", bid.ID.String())
			assert.NotEqualf(t, config.ZeroUUID, bid.UserID, "Bid User ID must mot be empty. Is: %s", bid.UserID.String())
			assert.NotEqualf(t, config.ZeroUUID, bid.ItemID, "Bid Item ID must mot be empty. Is: %s", bid.ItemID.String())
			assert.Falsef(t, bid.CreatedAt.IsZero(), "CreatedAt must mot be empty. Is: %s", bid.CreatedAt.String())
		})
	}
}

func Test_GetUser(t *testing.T) {

	h.Reset()
	users := testutils.CreateTestUsers(h, 10)

	tests := []struct {
		name     string
		ID       uuid.UUID
		expected *models.User
		wantErr  bool
	}{
		{"User 0 should exist", users[0].ID, users[0], false},
		{"User 1 should exist", users[1].ID, users[1], false},
		{"User 9 should exist", users[9].ID, users[9], false},
		{"User should not exist", uuid.NewV4(), &models.User{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := h.GetUser(tt.ID)
			if (err != nil) != tt.wantErr {
				t.Errorf(".GetUser() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(user, tt.expected) {
				t.Errorf(".GetUser() = \n%+v\n, want \n%+v", user, tt.expected)
			}
		})
	}
}

func Test_GetItem(t *testing.T) {

	h.Reset()
	items := testutils.CreateTestItems(h, 10)

	tests := []struct {
		name     string
		ID       uuid.UUID
		expected *models.Item
		wantErr  bool
	}{
		{"Item 0 should exist", items[0].ID, items[0], false},
		{"Item 1 should exist", items[1].ID, items[1], false},
		{"Item 9 should exist", items[9].ID, items[9], false},
		{"Item should not exist", uuid.NewV4(), &models.Item{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := h.GetItem(tt.ID)
			if (err != nil) != tt.wantErr {
				t.Errorf(".GetItem() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(user, tt.expected) {
				t.Errorf(".GetItem() = \n%+v\n, want \n%+v", user, tt.expected)
			}
		})
	}
}

func Test_GetUserBids(t *testing.T) {

	h.Reset()
	num := 10
	amounts := testutils.GenerateSliceOfRandomFloat64(num)
	bids, _, users := testutils.CreateTestBids(h, num, amounts)

	tests := []struct {
		name    string
		userID  uuid.UUID
		wantBid []*models.Bid
		wantErr bool
	}{
		{"Should find the bid", users[0].ID, []*models.Bid{bids[0]}, false},
		{"Should not find the bid", uuid.NewV4(), []*models.Bid{bids[6]}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := h.GetUserBids(tt.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf(".GetUserBids() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				// got is slice, so must be compared against slice
				if !reflect.DeepEqual(got, tt.wantBid) {
					t.Errorf(".GetUserBids() got = \n%+v\n, want \n%+v\n", got, tt.wantBid)
				}
			}
		})
	}
}

func Test_GetBidsOnItem(t *testing.T) {

	h.Reset()
	num := 10
	amounts := testutils.GenerateSliceOfRandomFloat64(num)
	bids, items, _ := testutils.CreateTestBids(h, num, amounts)

	tests := []struct {
		name    string
		itemID  uuid.UUID
		wantBid []*models.Bid
		wantErr bool
	}{
		{"Should find the bid", items[0].ID, []*models.Bid{bids[0]}, false},
		{"Should not find the bid", uuid.NewV4(), []*models.Bid{bids[6]}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := h.GetBidsOnItem(tt.itemID)
			if (err != nil) != tt.wantErr {
				t.Errorf(".GetBidsOnItem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if !reflect.DeepEqual(got, tt.wantBid) {
					t.Errorf(".GetBidsOnItem() got = \n%+v\n, want \n%+v\n", got, tt.wantBid)
				}
			}
		})
	}
}

func Test_GetWinningBid(t *testing.T) {

	h.Reset()
	numItems := 10
	amountsMatrix, maxAmounts := testutils.GenerateAmountsMatrix(numItems, 3*numItems)
	_, items, _ := testutils.CreateTestBidsManyOnItem(h, numItems, amountsMatrix)

	tests := []struct {
		name          string
		itemID        uuid.UUID
		wantMaxAmount float64
		wantErr       bool
	}{
		{"Should find winning bid for 0", items[0].ID, maxAmounts[0], false},
		{"Should find winning bid for 4", items[4].ID, maxAmounts[4], false},
		{"Should not find the bid", uuid.NewV4(), 0.0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			item, _ := h.GetItem(tt.itemID)
			if !tt.wantErr {
				assert.NotNil(t, item.WinningBid)
				assert.Equal(t, tt.wantMaxAmount, item.MaxBidAmount)
			} else {
				assert.Nil(t, item.WinningBid)
				assert.Equal(t, tt.wantMaxAmount, item.MaxBidAmount)
			}

			gotBid, err := h.GetWinningBid(tt.itemID)
			if (err != nil) != tt.wantErr {
				t.Errorf(".GetWinningBid() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.Equal(t, tt.wantMaxAmount, gotBid.Amount)
			}

		})
	}
}

func Test_GetItemsUserHasBid(t *testing.T) {

	h.Reset()
	numItems := 10
	amountsMatrix, _ := testutils.GenerateAmountsMatrix(numItems, 3*numItems)
	_, items, users := testutils.CreateTestBidsManyOnItem(h, numItems, amountsMatrix) // each users bids multiple times but on exactly one item

	tests := []struct {
		name         string
		userID       uuid.UUID
		wantNumItems int
		wantItems    []*models.Item
		wantErr      bool
	}{
		{"User 0 should bid on 0-th item", users[0].ID, 1, []*models.Item{items[0]}, false},
		{"User 1 should bid on 1-st item", users[1].ID, 1, []*models.Item{items[1]}, false},
		{"Unknown user should not be found", uuid.NewV4(), 0, []*models.Item{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			gotItems, err := h.GetItemsUserHasBid(tt.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetItemsUserHasBid() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.wantNumItems, len(gotItems))
			if !tt.wantErr && !reflect.DeepEqual(gotItems, tt.wantItems) {
				t.Errorf("GetItemsUserHasBid() gotItems = \n%+v\n, wantItems \n%+v\n", gotItems, tt.wantItems)
			}
		})
	}
}

func Test_GetItemsUserHasBid_TwoUsers(t *testing.T) {

	h.Reset()
	numItems := 10
	amountsMatrix, _ := testutils.GenerateAmountsMatrix(numItems, 3*numItems)
	_, items, users := testutils.CreateTestTwoUsersBidOnManyItems(h, numItems, amountsMatrix) // each users bids multiple times but on exactly one item

	evenItems := []*models.Item{items[0], items[2], items[4], items[6], items[8]}
	oddItems := []*models.Item{items[1], items[3], items[5], items[7], items[9]}

	tests := []struct {
		name         string
		userID       uuid.UUID
		wantNumItems int
		wantItems    []*models.Item
		wantErr      bool
	}{
		{"User 0 should bid on even items", users[0].ID, 5, evenItems, false},
		{"User 1 should bid on odd items", users[1].ID, 5, oddItems, false},
		{"Unknown user should not be found", uuid.NewV4(), 0, []*models.Item{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			gotItems, err := h.GetItemsUserHasBid(tt.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetItemsUserHasBid() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.wantNumItems, len(gotItems))

			//sorting the results by Name, so that compare is possible
			sort.Slice(gotItems, func(i, j int) bool {
				return gotItems[i].Name > gotItems[j].Name
			})
			sort.Slice(tt.wantItems, func(i, j int) bool {
				return tt.wantItems[i].Name > tt.wantItems[j].Name
			})
			if !tt.wantErr && !reflect.DeepEqual(gotItems, tt.wantItems) {
				t.Errorf("GetItemsUserHasBid() gotItems = \n%+v\n, wantItems \n%+v\n", gotItems, tt.wantItems)
			}
		})
	}
}

//// BENCHMARKS

func Benchmark_PlaceBid_OneUser_OneItem(b *testing.B) {

	user := models.NewUser("James Bond")
	item := models.NewItem("A thing")

	for n := 0; n < b.N; n++ {
		bid := models.NewBid(item.ID, user.ID, 3.1415)
		h.PlaceBid(bid)
	}
}

func Benchmark_PlaceBid_ManyUsers_ManyItems(b *testing.B) {
	for k := 0.; k <= scale; k++ {
		n := int(math.Pow(2, k))
		b.Run(fmt.Sprintf("%d", n), func(b *testing.B) {
			b.StopTimer()
			h.Reset()
			items := testutils.CreateTestItems(h, n)
			users := testutils.CreateTestUsers(h, n)

			randomItemIdx := rand.Int31n(int32(n))
			randomUserIdx := rand.Int31n(int32(n))
			b.StartTimer()
			for i := 0; i < b.N; i++ {
				bid := models.NewBid(items[randomItemIdx].ID, users[randomUserIdx].ID, 3.1415)
				h.PlaceBid(bid)
			}
		})
	}
}

func Benchmark_GetWinningBid(b *testing.B) {
	for k := 0.; k <= scale; k++ {
		n := int(math.Pow(2, k))
		b.Run(fmt.Sprintf("%d", n), func(b *testing.B) {
			b.StopTimer()

			h.Reset()
			numItems := n
			amountsMatrix, _ := testutils.GenerateAmountsMatrix(numItems, 3*numItems)
			_, items, _ := testutils.CreateTestBidsManyOnItem(h, numItems, amountsMatrix)

			b.StartTimer()
			for i := 0; i < b.N; i++ {
				randomItemIdx := rand.Int31n(int32(numItems))
				h.GetWinningBid(items[randomItemIdx].ID)
			}
		})
	}
}

func Benchmark_GetBidsOnItem(b *testing.B) {
	for k := 0.; k <= scale; k++ {
		n := int(math.Pow(2, k))
		b.Run(fmt.Sprintf("%d", n), func(b *testing.B) {
			b.StopTimer()

			h.Reset()
			numItems := n
			amountsMatrix, _ := testutils.GenerateAmountsMatrix(numItems, 3*numItems)
			_, items, _ := testutils.CreateTestBidsManyOnItem(h, numItems, amountsMatrix)

			b.StartTimer()
			for i := 0; i < b.N; i++ {
				randomItemIdx := rand.Int31n(int32(numItems))
				h.GetBidsOnItem(items[randomItemIdx].ID)
			}
		})
	}
}

func Benchmark_GetItemsUserHasBid(b *testing.B) {
	for k := 0.; k <= scale; k++ {
		n := int(math.Pow(2, k))
		b.Run(fmt.Sprintf("%d", n), func(b *testing.B) {
			b.StopTimer()

			h.Reset()
			numItems := n //also numUsers
			amountsMatrix, _ := testutils.GenerateAmountsMatrix(numItems, 3*numItems)
			_, _, users := testutils.CreateTestBidsManyOnItem(h, numItems, amountsMatrix)

			b.StartTimer()
			for i := 0; i < b.N; i++ {
				randomUserIdx := rand.Int31n(int32(numItems))
				h.GetItemsUserHasBid(users[randomUserIdx].ID)
			}
		})
	}
}
