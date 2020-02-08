package handlers

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	uuid "github.com/satori/go.uuid"
	"github.com/vikin91/bid-tracker-go/pkg/logging"
	"github.com/vikin91/bid-tracker-go/pkg/storage"
)

// define error messages
const (
	UserGetForbidden         = "Not allowed to get User"
	UserListForbidden        = "Not allowed to list Users"
	UserUpdateForbidden      = "Not allowed to update User"
	UserDeleteForbidden      = "Not allowed to delete User"
	UserPermissionsForbidden = "Not allowed to get User permissions"
	InvalidResetToken        = "Invalid Password Reset Token"
	MismatchedUserIDs        = "Request User IDs do not match"
)

//NewUserHandler initializes a new handler
func NewUserHandler(db storage.Storage) *UserHandler {
	return &UserHandler{db: db}
}

//UserHandler is the handler responsible for User operations
type UserHandler struct {
	db storage.Storage
}

//Routes returns the routes for the UserHandler
func (e *UserHandler) Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/", e.GetUsers)
	router.Get("/{userID}", e.GetUserByID)
	router.Get("/{userID}/bids", e.GetUserBids)
	router.Get("/{userID}/items", e.GetItemsUserHasBid) //TODO: Check swagger!
	return router
}

// GetUsers returns lists of Users
func (e *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := e.db.AllUsers()
	if err != nil {
		WriteHTTPErrorCode(w, err, http.StatusInternalServerError)
		return
	}
	render.JSON(w, r, users)
}

// GetUserByID returns User for the given user id
func (e *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	userID, err := ParseUserID(w, r)
	if err != nil {
		return
	}

	user, err := e.db.GetUser(userID)
	if err != nil {
		WriteHTTPErrorCode(w, err, http.StatusNotFound)
		return
	}
	render.JSON(w, r, user)
}

// GetUserBids returns User bids
func (e *UserHandler) GetUserBids(w http.ResponseWriter, r *http.Request) {
	userID, err := ParseUserID(w, r)
	if err != nil {
		WriteHTTPErrorCode(w, err, http.StatusInternalServerError)
		return
	}

	bids, err := e.db.GetUserBids(userID)
	if err != nil {
		WriteHTTPErrorCode(w, err, http.StatusNotFound)
		return
	}
	render.JSON(w, r, bids)
}

// GetItemsUserHasBid returns User bids
func (e *UserHandler) GetItemsUserHasBid(w http.ResponseWriter, r *http.Request) {
	userID, err := ParseUserID(w, r)
	if err != nil {
		WriteHTTPErrorCode(w, err, http.StatusInternalServerError)
		return
	}

	items, err := e.db.GetItemsUserHasBid(userID)
	if err != nil {
		WriteHTTPErrorCode(w, err, http.StatusNotFound)
		return
	}
	render.JSON(w, r, items)
}

// ParseUserID parses the URLParam or return an error if there is none
func ParseUserID(w http.ResponseWriter, r *http.Request) (uuid.UUID, error) {
	userID, err := uuid.FromString(chi.URLParam(r, "userID"))
	if err != nil {
		logging.LogError("Error parsing URL parameter to UUID", err)
		WriteHTTPErrorCode(w, errors.New("Malformed URL Parameter"), http.StatusBadRequest)
		return uuid.FromStringOrNil(""), err
	}
	return userID, nil
}
