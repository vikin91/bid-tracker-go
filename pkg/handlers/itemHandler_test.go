package handlers_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gavv/httpexpect"
	"github.com/stretchr/testify/assert"

	"github.com/vikin91/bid-tracker-go/internal/testutils"
	"github.com/vikin91/bid-tracker-go/pkg/config"
	"github.com/vikin91/bid-tracker-go/pkg/handlers"
	"github.com/vikin91/bid-tracker-go/pkg/models"
	"github.com/vikin91/bid-tracker-go/pkg/storage"
)

func TestItemHandler_GetItemsEmpty(t *testing.T) {
	db := storage.NewMapBiddingSystem()
	handler := handlers.NewItemHandler(db)

	server := httptest.NewServer(handler.Routes())
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	e.GET("/").
		Expect().
		Status(http.StatusOK).JSON().Array().Empty()
}

func TestItemHandler_GetItems(t *testing.T) {
	db := storage.NewMapBiddingSystem()
	numItems := 1
	amountsMatrix, _ := testutils.GenerateAmountsMatrix(numItems, 3*numItems)
	_, items, _ := testutils.CreateTestTwoUsersBidOnManyItems(db, numItems, amountsMatrix)

	handler := handlers.NewItemHandler(db)

	server := httptest.NewServer(handler.Routes())
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	e.GET("/").
		Expect().
		Status(http.StatusOK).JSON().Array().Contains(items[0])
}

func TestItemHandler_CreateItem(t *testing.T) {
	db := storage.NewMapBiddingSystem()
	handler := handlers.NewItemHandler(db)

	server := httptest.NewServer(handler.Routes())
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	e.GET("/").
		Expect().
		Status(http.StatusOK).JSON().Array().Empty()

	item := map[string]interface{}{
		"name": "A pen",
	}

	e.POST("/").WithJSON(item).
		Expect().
		Status(http.StatusCreated)

	e.GET("/").
		Expect().
		Status(http.StatusOK).JSON().Array().NotEmpty()
}

func TestItemHandler_GetBids(t *testing.T) {
	db := storage.NewMapBiddingSystem()
	numItems := 3
	amountsMatrix, _ := testutils.GenerateAmountsMatrix(numItems, 3*numItems)
	_, items, _ := testutils.CreateTestTwoUsersBidOnManyItems(db, numItems, amountsMatrix)
	handler := handlers.NewItemHandler(db)

	server := httptest.NewServer(handler.Routes())
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	e.GET(fmt.Sprintf("/%s/bids", items[0].ID.String())).
		Expect().
		Status(http.StatusOK).JSON().Array().NotEmpty()
}

func TestItemHandler_PlaceBid(t *testing.T) {
	db := storage.NewMapBiddingSystem()
	items := testutils.CreateTestItems(db, 1)
	users := testutils.CreateTestUsers(db, 1)
	handler := handlers.NewItemHandler(db)

	server := httptest.NewServer(handler.Routes())
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	bid := models.NewBid(config.ZeroUUID, users[0].ID, 99.55)
	bidAfterSaving := models.NewBid(items[0].ID, users[0].ID, 99.55)
	bidAfterSaving.ID = bid.ID
	bidAfterSaving.CreatedAt = bid.CreatedAt

	e.GET(fmt.Sprintf("/%s/bids", items[0].ID.String())).
		Expect().
		Status(http.StatusOK).JSON().Array().Empty()

	e.POST(fmt.Sprintf("/%s/bids", items[0].ID.String())).
		WithJSON(bid).
		Expect().
		Status(http.StatusCreated).NoContent()

	e.GET(fmt.Sprintf("/%s/bids", items[0].ID.String())).
		Expect().
		Status(http.StatusOK).JSON().Array().NotEmpty().Contains(bidAfterSaving)

	e.POST(fmt.Sprintf("/%s/bids", "xxx-trash")).
		WithJSON(bid).
		Expect().
		Status(http.StatusBadRequest).Body().Contains("Malformed URL Parameter")

	fakeBid := map[string]interface{}{
		"money": 1000.000,
	}
	e.POST(fmt.Sprintf("/%s/bids", items[0].ID.String())).
		WithJSON(fakeBid).
		Expect().
		Body().Contains(handlers.BidDecodeFailure)

	//place bid on non existing item
	e.POST(fmt.Sprintf("/%s/bids", config.ZeroUUID.String())).
		WithJSON(bid).
		Expect().
		Status(http.StatusNotFound)

	bidFakeUser := models.NewBid(config.ZeroUUID, config.ZeroUUID, 99.55)
	e.POST(fmt.Sprintf("/%s/bids", items[0].ID.String())).
		WithJSON(bidFakeUser).
		Expect().
		Status(http.StatusInternalServerError).Body().Contains(handlers.UnknownUserBids)
}

func TestItemHandler_GetWinner(t *testing.T) {
	db := storage.NewMapBiddingSystem()
	items := testutils.CreateTestItems(db, 2)
	users := testutils.CreateTestUsers(db, 2)
	handler := handlers.NewItemHandler(db)

	server := httptest.NewServer(handler.Routes())
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	bid1 := models.NewBid(items[0].ID, users[0].ID, 10.1)
	bid2 := models.NewBid(items[0].ID, users[0].ID, 15.00)

	oneSecondLater := time.Now().Local().Add(time.Second * time.Duration(1))
	bid3 := models.NewBid(items[0].ID, users[0].ID, 15.00) //same amount but placed later
	bid3.CreatedAt = oneSecondLater

	assert.True(t, bid2.CreatedAt.Before(bid3.CreatedAt))

	db.PlaceBid(bid1)
	db.PlaceBid(bid2)
	db.PlaceBid(bid3)

	e.GET(fmt.Sprintf("/%s/winner", items[0].ID.String())).
		Expect().
		Status(http.StatusOK).JSON().Object().Equal(bid2)
}
