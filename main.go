package main

import (
	"log"
	"net/http"

	"github.com/AvenciaLab/avencia-backend/lib/setup/di"
	"github.com/AvenciaLab/avencia-backend/lib/setup/di/external"
)

func main() {
	handler := di.InitializeHandler(di.InitializeBusiness(external.InitializeExternal()))
	log.Fatalf("while running handler: %v", http.ListenAndServe(":4244", handler))
}
