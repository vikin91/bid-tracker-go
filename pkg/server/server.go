package server

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/vikin91/bid-tracker-go/pkg/config"
	"github.com/vikin91/bid-tracker-go/pkg/handlers"
	"github.com/vikin91/bid-tracker-go/pkg/logging"
	"github.com/vikin91/bid-tracker-go/pkg/storage"
)

//Server wraps a chi router (chi.Mux)
type Server struct {
	mux *chi.Mux
}

func newMux() *chi.Mux {
	mux := chi.NewRouter()
	mux.Use(
		render.SetContentType(render.ContentTypeJSON), // Set content-Type headers as application/json
		middleware.RequestID,
		middleware.Logger,          // Log API request calls
		middleware.DefaultCompress, // Compress results, mostly gzipping assets and json
		middleware.Recoverer,       // Recover from panics without crashing server
		middleware.StripSlashes,
	)
	return mux
}

// NewServer creates a router with routes setup
func NewServer() *Server {
	server := Server{mux: newMux()}
	return &server
}

//Mux returns the chi router
func (s *Server) Mux() *chi.Mux {
	return s.mux
}

//SetupRoutes adds all routes that the server should listen to
func (s *Server) SetupRoutes(db storage.Storage) {
	userHandler := handlers.NewUserHandler(db)
	itemHandler := handlers.NewItemHandler(db)

	s.Mux().Route(config.APIPrefixV1, func(r chi.Router) {
		r.Mount("/user", userHandler.Routes())
		r.Mount("/item", itemHandler.Routes())
	})
}

//ListenAndServe starts the server
func (s *Server) ListenAndServe(quit chan struct{}, errors chan config.ErrorMessage, port string) {
	go func() {
		listenAddress := net.JoinHostPort("", port)
		log.Printf("Listening on %s\n", listenAddress)
		if err := http.ListenAndServe(listenAddress, s.mux); err != nil {
			msg := config.ErrorMessage{Message: fmt.Sprintf("Could not listen on port %s", port), Err: err}
			select {
			case errors <- msg:
				logging.LogInfo("Sent on errors channel")
			default:
				logging.LogInfo("Failed to send on errors channel")
			}
		}
	}()

	<-quit
	log.Printf("Server has been shutdown")
}
