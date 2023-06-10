package routes

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Route struct {
	URI    string
	Method string
	Func   func(w http.ResponseWriter, r *http.Request)
}

// BootstrapRoutes is a function that receives a mux.Router and returns a mux.Router
// with all the routes of the application
func BootstrapRoutes(router *mux.Router) *mux.Router {
	routes := homeRoutes
	// routes = append(routes, authRoutes...)

	router.HandleFunc(routes.URI, routes.Func).Methods(routes.Method)
	// for _, route := range routes {
	// }

	return router
}
