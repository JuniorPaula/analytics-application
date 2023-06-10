package handlers

import (
	"c2d-reports/internal/usecases"
	"net/http"
)

var companyToken = "203a34841c77f66c9e94524d38d79d"

func LoadTMR_handler(w http.ResponseWriter, r *http.Request) {

	uc := usecases.ReportTmrUsecase{
		CompanyToken: companyToken,
	}
	uc.LoadTMR()

	defer r.Body.Close()
	w.Write([]byte("Hello form LoadTMR"))
}
