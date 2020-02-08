package handlers_test

import (
	"os"
	"testing"

	"github.com/vikin91/bid-tracker-go/pkg/config"
)

// Executed before test runs in this package (fails otherwise)
func TestMain(m *testing.M) {
	config.SetupEnv()
	os.Exit(m.Run())
}
