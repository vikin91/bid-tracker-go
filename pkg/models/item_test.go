package models_test

import (
	"testing"

	"github.com/vikin91/bid-tracker-go/pkg/models"
)

func Benchmark_Reference_NewItem(b *testing.B) {
	for n := 0; n < b.N; n++ {
		models.NewItem("A name")
	}
}
