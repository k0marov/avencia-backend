package test_helpers

import (
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

func TimeAlmostEqual(t1, t2 time.Time) bool {
	return math.Abs(t1.Sub(t2).Minutes()) < 1
}
