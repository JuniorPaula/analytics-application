package handlers

import (
	"c2d-reports/internal/entity"
	"c2d-reports/internal/usecases/companies"
	"c2d-reports/internal/utils"
	"encoding/json"
	"io"
	"net/http"
)

func CreateCompany_hanlder(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	var company entity.Company
	err = json.Unmarshal(body, &company)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	company, err = companies.CreateCompanyUsecase(company)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, company)
}

func GetAllCompanies_handler(w http.ResponseWriter, r *http.Request) {
	companies, err := companies.GetAllCompaniesUsecase()
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, companies)
}
