package companies

import (
	"c2d-reports/internal/database"
	"c2d-reports/internal/entity"
	"c2d-reports/internal/repositories"
	"errors"
)

func CreateCompanyUsecase(company entity.Company) (entity.Company, error) {
	if err := company.Validate(); err != nil {
		return entity.Company{}, err
	}

	db, err := database.Connect()
	if err != nil {
		return entity.Company{}, err
	}

	repository := repositories.NewCompanyRepository(db)
	company, err = repository.CreateCompany(company)
	if err != nil {
		return entity.Company{}, err
	}

	return company, nil
}

func GetAllCompaniesUsecase() ([]entity.Company, error) {
	db, err := database.Connect()
	if err != nil {
		return nil, err
	}

	repository := repositories.NewCompanyRepository(db)
	companies, err := repository.GetAllCompanies()
	if err != nil {
		return nil, err
	}

	return companies, nil
}

func GetCompanyByIDUsecase(ID int64) (entity.Company, error) {
	db, err := database.Connect()
	if err != nil {
		return entity.Company{}, err
	}

	repository := repositories.NewCompanyRepository(db)
	company, err := repository.GetCompanyByID(ID)
	if err != nil {
		return entity.Company{}, err
	}

	return company, nil
}

func UpdateCompanyUsecase(ID int64, company entity.Company) (entity.Company, error) {
	companyInDB, err := GetCompanyByIDUsecase(ID)
	if err != nil {
		return entity.Company{}, errors.New("company not found")
	}

	if company.CompanyID == 0 {
		company.CompanyID = companyInDB.CompanyID
	}

	if company.CompanyToken == "" {
		company.CompanyToken = companyInDB.CompanyToken
	}

	if company.CompanyName == "" {
		company.CompanyName = companyInDB.CompanyName
	}

	if company.EmailAdmin == "" {
		company.EmailAdmin = companyInDB.EmailAdmin
	}

	db, err := database.Connect()
	if err != nil {
		return entity.Company{}, err
	}

	repository := repositories.NewCompanyRepository(db)
	company, err = repository.UpdateCompany(ID, company)
	if err != nil {
		return entity.Company{}, err
	}

	return company, nil
}
