package main

import (
	"net/http"

	"github.com/k0marov/avencia-backend/lib/di"
	"github.com/k0marov/avencia-backend/lib/di/external"
)

func main() {
	http.ListenAndServe(":4244", di.InitializeHandler(di.InitializeBusiness(external.InitializeExternal())))
}
