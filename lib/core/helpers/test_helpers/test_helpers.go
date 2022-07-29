package test_helpers

import (
	"cloud.google.com/go/firestore"
	"github.com/k0marov/avencia-backend/lib/core/fs_facade"
	"log"
	"math"
	"math/rand"
	"time"
)

var randGen *rand.Rand

func init() {
	seed := time.Now().Unix()
	log.Printf("running tests with random seed: %v", seed)
	randGen = rand.New(rand.NewSource(seed))
}

var StubRunBatch = func(f func(fs_facade.BatchUpdater) error) error {
	return f(func(*firestore.DocumentRef, map[string]any) error {
		return nil
	})
}

func TimeAlmostEqual(t1, t2 time.Time) bool {
	return math.Abs(t1.Sub(t2).Minutes()) < 1
}
