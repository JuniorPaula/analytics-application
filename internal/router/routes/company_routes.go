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
	{
		URI:    "/companies/{id}",
		Method: http.MethodGet,
		Func:   handlers.GetCompanyByID_handler,
	},
	{
		URI:    "/companies/{id}",
		Method: http.MethodPut,
		Func:   handlers.UpdateCompany_handler,
	},
	{
		URI:    "/companies/{id}",
		Method: http.MethodDelete,
		Func:   handlers.DeleteCompany_handler,
	},
}
