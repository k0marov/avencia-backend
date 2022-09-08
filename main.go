package main

import (
	"net/http"
	"log"

	"github.com/k0marov/avencia-backend/lib/di"
	"github.com/k0marov/avencia-backend/lib/di/external"
)

func main() {
	log.Fatalf("while running handler: %v", http.ListenAndServe(":4244", di.InitializeHandler(di.InitializeBusiness(external.InitializeExternal()))))
}
