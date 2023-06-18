package handlers

import (
	"c2d-reports/internal/config"
	"c2d-reports/internal/usecases"
	"net/http"
)

func LoadTMR_handler(w http.ResponseWriter, r *http.Request) {

	uc := usecases.ReportTmrUsecase{
		CompanyToken: config.CompanyToken,
	}
	uc.LoadTMR()

	defer r.Body.Close()
	w.Write([]byte("Hello form LoadTMR"))
}
