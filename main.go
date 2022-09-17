package main

import (
	"log"
	"net/http"

	"github.com/AvenciaLab/avencia-backend/lib/setup/di"
	"github.com/AvenciaLab/avencia-backend/lib/setup/di/external"
)
// TODO 5: consider adding more context info to every core_err.Rethrow()
// TODO 5: maybe use context.Context instead of db.DB for the services



func main() {
	handler := di.InitializeHandler(di.InitializeBusiness(external.InitializeExternal()))
	log.Fatalf("while running handler: %v", http.ListenAndServe(":4244", handler))
}
