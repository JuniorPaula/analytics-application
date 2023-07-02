package router

import (
	"c2d-reports/internal/router/routes"

	"github.com/gorilla/mux"
)

// Handler is the function that returns a mux.Router with all the routes of the application
// It is called from cmd/main.go
func Handler() *mux.Router {
	r := mux.NewRouter()

	return routes.BootstrapRoutes(r)
}
