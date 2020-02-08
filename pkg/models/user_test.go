package models_test

import (
	"testing"

	"github.com/vikin91/bid-tracker-go/pkg/models"
)

func Benchmark_Reference_NewUser(b *testing.B) {
	for n := 0; n < b.N; n++ {
		models.NewUser("A name")
	}
}
