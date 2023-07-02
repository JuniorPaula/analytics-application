package handlers

import (
	"c2d-reports/internal/config"
	"c2d-reports/internal/usecases"
	"net/http"
)

// FindAllReports_handler handles the request to find all reports
// This handler is used for testing purposes, its not used in production
func LoadTMR_handler(w http.ResponseWriter, r *http.Request) {

	uc := usecases.ReportTmrUsecase{
		CompanyToken: config.CompanyToken,
	}
	uc.LoadTMR()

	defer r.Body.Close()
	w.Write([]byte("Hello form LoadTMR"))
}

// FindAllReports_handler handles the request to find all reports
// This handler is used for testing purposes, its not used in production
func DeleteReports_handler(w http.ResponseWriter, r *http.Request) {
	uc := usecases.ReportTmrUsecase{
		CompanyToken: config.CompanyToken,
	}
	uc.DeleteReport()

	defer r.Body.Close()
	w.Write([]byte("Hello form FindAllReports"))
}
