package routes

import (
	"c2d-reports/internal/handlers"
	"net/http"
)

var companiesRoutes = []Route{
	{
		URI:    "/companies",
		Method: http.MethodPost,
		Func:   handlers.CreateCompany_hanlder,
	},
	{
		URI:    "/companies",
		Method: http.MethodGet,
		Func:   handlers.GetAllCompanies_handler,
	},
}
