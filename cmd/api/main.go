package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/spf13/viper"
	"github.com/vikin91/bid-tracker-go/internal/testutils"
	"github.com/vikin91/bid-tracker-go/pkg/config"
	"github.com/vikin91/bid-tracker-go/pkg/logging"
	"github.com/vikin91/bid-tracker-go/pkg/server"
	"github.com/vikin91/bid-tracker-go/pkg/storage"
)

func main() {
	config.SetupEnv()
	port := viper.GetString("PORT")
	demo := flag.Bool("demo", false, "Pre-fill with demo data")
	flag.Parse()

	quitServerCh := make(chan struct{})
	errorsCh := make(chan config.ErrorMessage)
	termSignal := make(chan os.Signal, 1)
	signal.Notify(termSignal, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)

	db := storage.NewMapBiddingSystem()

	if *demo {
		numItems := 50
		amountsMatrix, _ := testutils.GenerateAmountsMatrix(numItems, 3*numItems)
		testutils.CreateTestTwoUsersBidOnManyItems(db, numItems, amountsMatrix)
		logging.LogInfo("System populated with demo data")
	}

	server := server.NewServer()
	server.SetupRoutes(db)

	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		route = strings.Replace(route, "/*/", "/", -1)
		fmt.Printf("%s %s\n", method, route)
		return nil
	}

	if err := chi.Walk(server.Mux(), walkFunc); err != nil {
		fmt.Printf("Logging err: %s\n", err.Error())
	}

	go server.ListenAndServe(quitServerCh, errorsCh, port)

	terminateFunc := func(quitServerCh chan struct{}) {
		quitServerCh <- struct{}{}
		close(quitServerCh)
		time.Sleep(time.Second)
	}

	select {
	case event := <-errorsCh:
		logging.LogInfo("Received error message. Waiting for server to stop...")
		logging.LogError(event.Message, event.Err)
		terminateFunc(quitServerCh)
	case <-termSignal:
		logging.LogInfo("Received terminate signal. Waiting for server to stop...")
		terminateFunc(quitServerCh)
	}
}
