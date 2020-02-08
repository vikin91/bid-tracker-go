package handlers_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect"
	"github.com/stretchr/testify/assert"
	"github.com/vikin91/bid-tracker-go/internal/testutils"
	"github.com/vikin91/bid-tracker-go/pkg/handlers"
	"github.com/vikin91/bid-tracker-go/pkg/storage"
)

func TestUserHandler_GetUsersEmpty(t *testing.T) {
	db := storage.NewMapBiddingSystem()
	handler := handlers.NewUserHandler(db)

	server := httptest.NewServer(handler.Routes())
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	e.GET("/").
		Expect().
		Status(http.StatusOK).JSON().Array().Empty()
}

func TestUserHandler_GetUsers(t *testing.T) {
	db := storage.NewMapBiddingSystem()
	users := testutils.CreateTestUsers(db, 2)
	handler := handlers.NewUserHandler(db)

	server := httptest.NewServer(handler.Routes())
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	e.GET("/").
		Expect().
		Status(http.StatusOK).JSON().Array().NotEmpty().Contains(users[0]).Contains(users[1])
}

func TestUserHandler_GetUserByID(t *testing.T) {
	db := storage.NewMapBiddingSystem()
	users := testutils.CreateTestUsers(db, 2)
	handler := handlers.NewUserHandler(db)

	server := httptest.NewServer(handler.Routes())
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	e.GET(fmt.Sprintf("/%s", users[0].ID.String())).
		Expect().
		Status(http.StatusOK).JSON().Object().Equal(users[0]).NotEqual(users[1])

	e.GET(fmt.Sprintf("/%s", "xxx-not-uuid")).
		Expect().
		Status(http.StatusBadRequest).Body().Contains("Malformed URL Parameter")
}

func TestUserHandler_GetUserBids(t *testing.T) {
	db := storage.NewMapBiddingSystem()
	numItems := 1
	amountsMatrix, _ := testutils.GenerateAmountsMatrix(numItems, 3*numItems)
	bids, _, users := testutils.CreateTestTwoUsersBidOnManyItems(db, numItems, amountsMatrix)

	handler := handlers.NewUserHandler(db)

	server := httptest.NewServer(handler.Routes())
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	assert.NotNil(t, bids[0])
	assert.NotNil(t, bids[1])
	assert.NotNil(t, bids[2])
	assert.Panics(t, func() { fmt.Println(bids[3]) }, "Only 3 bids should be created")

	e.GET(fmt.Sprintf("/%s/bids", users[0].ID.String())).
		Expect().
		Status(http.StatusOK).JSON().Array().Contains(bids[0]).Contains(bids[1]).Contains(bids[2])
}

func TestUserHandler_GetItemsUserHasBid(t *testing.T) {
	db := storage.NewMapBiddingSystem()
	numItems := 4
	amountsMatrix, _ := testutils.GenerateAmountsMatrix(numItems, 3*numItems)
	_, items, users := testutils.CreateTestTwoUsersBidOnManyItems(db, numItems, amountsMatrix)

	handler := handlers.NewUserHandler(db)

	server := httptest.NewServer(handler.Routes())
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	assert.NotNil(t, items[0])
	assert.NotNil(t, items[1])
	assert.NotNil(t, items[2])
	assert.NotNil(t, items[3])
	assert.Panics(t, func() { fmt.Println(items[4]) }, "Only 4 items should be created")

	//user bids even
	e.GET(fmt.Sprintf("/%s/items", users[0].ID.String())).
		Expect().
		Status(http.StatusOK).JSON().Array().Contains(items[0], items[2]).NotContains(items[1], items[3])

	//user bids odd
	e.GET(fmt.Sprintf("/%s/items", users[1].ID.String())).
		Expect().
		Status(http.StatusOK).JSON().Array().Contains(items[1], items[3]).NotContains(items[0], items[2])
}
