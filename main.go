package main

import (
	"log"
	"net/http"

	"github.com/AvenciaLab/avencia-backend/lib/di"
	"github.com/AvenciaLab/avencia-backend/lib/di/external"
)

func main() {
	log.Fatalf("while running handler: %v", http.ListenAndServe(":4244", di.InitializeHandler(di.InitializeBusiness(external.InitializeExternal()))))
}
