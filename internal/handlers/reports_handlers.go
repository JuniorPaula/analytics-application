package handlers

import (
	"c2d-reports/internal/config"
	"c2d-reports/internal/usecases"
	"net/http"
)

var companyToken = config.CompanyToken

func LoadTMR_handler(w http.ResponseWriter, r *http.Request) {

	uc := usecases.ReportTmrUsecase{
		CompanyToken: companyToken,
	}
	uc.LoadTMR()

	defer r.Body.Close()
	w.Write([]byte("Hello form LoadTMR"))
}
