package main

import (
	"log"
	"net/http"

	"github.com/AvenciaLab/avencia-backend/lib/setup/di"
	"github.com/AvenciaLab/avencia-backend/lib/setup/di/external"
)
// TODO 1: get rid of using type x = func() and replace it with type x func() everywhere
// TODO 2: maybe use less packages (i.e have only one namespace for every feature)
// TODO 5: consider adding more context info to every core_err.Rethrow()
// TODO 2: maybe rename all package names to camel case since it is idiomatic
// TODO 3: limit available currencies 
// TODO 5: maybe use context.Context instead of db.DB for the services



func main() {
	handler := di.InitializeHandler(di.InitializeBusiness(external.InitializeExternal()))
	log.Fatalf("while running handler: %v", http.ListenAndServe(":4244", handler))
}
