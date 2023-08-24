package routes

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Route is the model for the routes
type Route struct {
	URI    string
	Method string
	Func   func(w http.ResponseWriter, r *http.Request)
}

// BootstrapRoutes is a function that receives a mux.Router and returns a mux.Router
// with all the routes of the application
func BootstrapRoutes(router *mux.Router) *mux.Router {
	var routes []Route
	routes = append(routes, homeRoutes)
	routes = append(routes, companiesRoutes...)

	for _, route := range routes {
		router.HandleFunc(route.URI, route.Func).Methods(route.Method)

	}

	return router
}
