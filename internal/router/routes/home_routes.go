package routes

import "net/http"

// Route is the model for the routes
var homeRoutes = Route{
	URI:    "/",
	Method: http.MethodGet,
	Func: func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	},
}
