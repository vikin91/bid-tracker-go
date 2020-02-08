package models_test

import (
	"testing"

	"github.com/vikin91/bid-tracker-go/pkg/models"
)

func Benchmark_Reference_NewBid(b *testing.B) {
	user := models.NewUser("James Bond")
	item := models.NewItem("A thing")

	for n := 0; n < b.N; n++ {
		models.NewBid(item.ID, user.ID, 0.007)
	}
}
