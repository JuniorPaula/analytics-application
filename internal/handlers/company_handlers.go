package handlers

import (
	"c2d-reports/internal/entity"
	"c2d-reports/internal/usecases/companies"
	"c2d-reports/internal/utils"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
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

func GetCompanyByID_handler(w http.ResponseWriter, r *http.Request) {
	paramsID := mux.Vars(r)["id"]
	ID, err := strconv.ParseInt(paramsID, 10, 64)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	company, err := companies.GetCompanyByIDUsecase(ID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, company)
}

func UpdateCompany_handler(w http.ResponseWriter, r *http.Request) {
	paramsID := mux.Vars(r)["id"]
	ID, err := strconv.ParseInt(paramsID, 10, 64)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

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

	company, err = companies.UpdateCompanyUsecase(ID, company)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, company)
}

func DeleteCompany_handler(w http.ResponseWriter, r *http.Request) {
	paramsID := mux.Vars(r)["id"]
	ID, err := strconv.ParseInt(paramsID, 10, 64)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	err = companies.DeleteCompanyUsecase(ID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, nil)
}
