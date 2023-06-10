package routes

import (
	"c2d-reports/internal/handlers"
	"net/http"
)

var reportsRoutes = []Route{
	{
		URI:    "/tmr/load",
		Method: http.MethodGet,
		Func:   handlers.LoadTMR,
	},
}
