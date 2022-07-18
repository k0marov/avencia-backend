package main

import (
	"net/http"

	"github.com/k0marov/avencia-backend/lib"
)

func main() {
	http.ListenAndServe(":4244", lib.Initialize())
}
