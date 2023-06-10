package router

import (
	"c2d-reports/internal/router/routes"

	"github.com/gorilla/mux"
)

func Handler() *mux.Router {
	r := mux.NewRouter()

	return routes.BootstrapRoutes(r)
}
