package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	uuid "github.com/satori/go.uuid"
	"github.com/vikin91/bid-tracker-go/pkg/logging"
	"github.com/vikin91/bid-tracker-go/pkg/models"
	"github.com/vikin91/bid-tracker-go/pkg/storage"
)

// define error messages
const (
	ItemAccountForbidden  = "Not allowed to get Item Account"
	ItemCreationForbidden = "Not allowed to create Item"
	ItemCreationFailure   = "Failed to create Item"
	BidDecodeFailure      = "Failed to decode a bid"
	UnknownUserBids       = "Cannot find user that places this bid"
	BidPlacementFailure   = "Failed place a bid"
	ItemListForbidden     = "Not allowed to get all Items"
	ResourceNotFound      = "Resource not found"
	ItemNotFound          = "Item not found"
)

//NewItemHandler initializes a new handler
func NewItemHandler(db storage.Storage) *ItemHandler {
	return &ItemHandler{db: db}
}

//ItemHandler is the handler responsible for Item operations
type ItemHandler struct {
	db storage.Storage
}

//Routes returns the routes for the ItemHandler
func (e *ItemHandler) Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Use(render.SetContentType(render.ContentTypeJSON))
	router.Get("/", e.GetItems)
	router.Post("/", e.CreateItem)

	router.Get("/{itemID}/bids", e.GetBids)
	router.Post("/{itemID}/bids", e.PlaceBid)
	router.Get("/{itemID}/winner", e.GetWinner)
	return router
}

// GetItems returns list of items
func (e *ItemHandler) GetItems(w http.ResponseWriter, r *http.Request) {

	items, err := e.db.AllItems()
	if err != nil {
		WriteHTTPErrorCode(w, err, http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, items)
}

// CreateItem creates new item (only admin or item owner)
func (e *ItemHandler) CreateItem(w http.ResponseWriter, r *http.Request) {

	// Item has to be valid
	item := &models.Item{}
	err := json.NewDecoder(r.Body).Decode(item)
	if err != nil {
		logging.LogError("Error decoding item creation request payload", err)
		WriteHTTPErrorCode(w, err, http.StatusBadRequest)
		return
	}
	err = e.db.CreateItem(item)
	if err != nil {
		WriteHTTPErrorCode(w, err, http.StatusInternalServerError)
		return
	}

	WriteHTTPCode(w, http.StatusCreated)
}

// GetBids returns list of bids on item
func (e *ItemHandler) GetBids(w http.ResponseWriter, r *http.Request) {
	item, err := e.findItem(w, r)
	if err != nil {
		return
	}
	bids, err := e.db.GetBidsOnItem(item.ID)
	if err != nil {
		logging.LogError("Cannot get bids on item", err)
		WriteHTTPErrorCode(w, err, http.StatusInternalServerError)
		return
	}
	render.JSON(w, r, bids)
}

// PlaceBid returns list of bids on item
func (e *ItemHandler) PlaceBid(w http.ResponseWriter, r *http.Request) {
	item, err := e.findItem(w, r)
	if err != nil {
		return
	}

	bid := &models.Bid{}
	err = json.NewDecoder(r.Body).Decode(bid)
	if err != nil || (*bid == models.Bid{}) {
		logging.LogError("Error decoding bid", err)
		WriteHTTPErrorCode(w, errors.New(BidDecodeFailure), http.StatusBadRequest)
		return
	}
	bid.ItemID = item.ID
	_, err = e.db.GetUser(bid.UserID)
	if err != nil {
		logging.LogError(UnknownUserBids, err)
		WriteHTTPErrorCode(w, errors.New(UnknownUserBids), http.StatusInternalServerError)
		return
	}
	err = e.db.PlaceBid(bid)
	if err != nil {
		logging.LogError(BidPlacementFailure, err)
		WriteHTTPErrorCode(w, errors.New(BidPlacementFailure), http.StatusInternalServerError)
		return
	}
	WriteHTTPCode(w, http.StatusCreated)
}

// GetWinner returns single winning bid
func (e *ItemHandler) GetWinner(w http.ResponseWriter, r *http.Request) {
	item, err := e.findItem(w, r)
	if err != nil {
		return
	}
	bid, err := e.db.GetWinningBid(item.ID)
	if err != nil {
		logging.LogError("Cannot get winning bid on item", err)
		WriteHTTPErrorCode(w, err, http.StatusInternalServerError)
		return
	}
	render.JSON(w, r, bid)
}

func (e *ItemHandler) findItem(w http.ResponseWriter, r *http.Request) (*models.Item, error) {
	itemID, err := ParseItemID(w, r)
	if err != nil {
		return nil, err
	}
	item, err := e.db.GetItem(itemID)
	if err != nil {
		logging.LogError("Cannot find item", err)
		WriteHTTPErrorCode(w, err, http.StatusNotFound)
		return nil, err
	}
	return item, nil
}

// ParseItemID parses the URLParam and sends the HTTPError Response on failure
func ParseItemID(w http.ResponseWriter, r *http.Request) (uuid.UUID, error) {
	itemID, err := uuid.FromString(chi.URLParam(r, "itemID"))
	if err != nil {
		logging.LogError("Error parsing URL parameter to UUID", err)
		WriteHTTPErrorCode(w, errors.New("Malformed URL Parameter"), http.StatusBadRequest)
		return uuid.FromStringOrNil(""), err
	}
	return itemID, nil
}
